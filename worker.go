package main

import (
	"bitbucket.org/onyxium/onyxium-strategy-worker/models"
	log "github.com/sirupsen/logrus"
	// "time"
)

type Worker struct {
	ID          int
	Work        chan *models.Strategy
	WorkerQueue chan chan *models.Strategy
	QuitChan    chan bool
}

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func NewWorker(id int, workerQueue chan chan *models.Strategy) Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:          id,
		Work:        make(chan *models.Strategy),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool)}

	return worker
}

// This function "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop.
func (w *Worker) Start() {
	go func() {
		for {
			log.Infof("Worker %d Add ourselves into the worker queue.", w.ID)
			w.WorkerQueue <- w.Work

			select {
			case <-w.QuitChan:
				// We have been asked to stop.
				log.Infof("Worker %d stopping", w.ID)
				return
			case work := <-w.Work:
				// Receive a work request.
				work.Status = "running"
				err := env.DataStore.StrategyUpdate(work)
				if err != nil {
					log.Error(err)
					w.Stop()
					continue
				}
				log.Infof("Worker %d Received work request %s", w.ID, work.Id.Hex())

				root, _ := work.BsonTree.Search(work.State)
				w.WalkSiblings(root, work)
				log.Infof("Worker %d work is done", w.ID)
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
// Note that the worker will only stop *after* it has finished its work.
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

func (w *Worker) WalkSiblings(tree *models.Tree, strategy *models.Strategy) (*models.Strategy, error) {
	// We could update state also between siblings
	if tree != nil {
		latestMarkets, err := env.DataStore.GetLatestMarket()
		if err != nil {
			return nil, err
		}
		log.Infof("Worker %d is in CURRENT NODE: %+v", w.ID, tree)
		doAction := CheckConditions(tree.Conditions, latestMarkets)
		if doAction {
			strategy.Status = "executing"
			strategy.State = tree.Id
			err = env.DataStore.StrategyUpdate(strategy)
			if err != nil {
				return nil, err
			}
			ExecuteAction(tree.Action)

			if tree.Left == nil {
				// state == tree.id
				strategy.Status = "stopped"
				err = env.DataStore.StrategyUpdate(strategy)
				if err != nil {
					return nil, err
				}
				log.Info("NO MORE STATEMENT AFTER THIS ACTION STATEMENT, I'M DONE")
				return strategy, nil
			} else {
				strategy.Status = "paused"
				strategy.State = tree.Left.Id
				err = env.DataStore.StrategyUpdate(strategy)
				if err != nil {
					return nil, err
				}
				log.Info("JUMPING to left")
				return strategy, nil
			}
		} else {
			if tree.Right == nil {
				// state == first sibling id
				strategy.Status = "paused"
				err = env.DataStore.StrategyUpdate(strategy)
				if err != nil {
					return nil, err
				}
				log.Infof("None of the conditions of each child are true. PAUSING strategy: %s", strategy.Id.Hex())
				return strategy, nil
			} else {
				log.Info("None of the conditions are true. JUMPING to right sibling.")
				return w.WalkSiblings(tree.Right, strategy)
			}
		}
	}
	return strategy, nil
}

func CheckConditions(conditions []models.Condition, latestMarkets map[string]models.Market) bool {
	for _, condition := range conditions {
		if condition.ConditionType == "greater-than-or-equal-to" || condition.ConditionType == "less-than-or-equal-to" {
			log.Infof("If the %s on the market %s/%s is %s than %.8f.", condition.BaseMetric, condition.BaseCurrency, condition.QuoteCurrency, condition.ConditionType, condition.Value)
		} else {
			log.Infof("If the %s on the market %s/%s has %s with %.3f percentage within %d minutes.", condition.BaseMetric, condition.BaseCurrency, condition.QuoteCurrency, condition.ConditionType, condition.Value, condition.TimeframeInMS/60000)
		}

		latestMarket := latestMarkets[condition.BaseCurrency+"-"+condition.QuoteCurrency]

		switch condition.ConditionType {
		case "greater-than-or-equal-to":
			currentValue := getMetricValue(condition.BaseMetric, latestMarket)
			log.Debugf("MARKET %s: %.8f", condition.BaseMetric, currentValue)
			if currentValue < condition.Value {
				return false
				log.Debugf("doAction FALSE: Market %s with value %.8f is < than condition value %.8f", condition.BaseMetric, currentValue, condition.Value)
			}
		case "less-than-or-equal-to":
			currentValue := getMetricValue(condition.BaseMetric, latestMarket)
			log.Debugf("MARKET %s: %.8f", condition.BaseMetric, currentValue)
			if currentValue > condition.Value {
				return false
				log.Debugf("doAction FALSE: Market %s with value %.8f is > than condition value %.8f", condition.BaseMetric, currentValue, condition.Value)
			}
		case "percentage-increase":
			historyMarkets, err := env.DataStore.GetHistoryMarket(condition.TimeframeInMS)
			if err != nil {
				log.Fatal(err)
			}
			historyMarket := historyMarkets[condition.BaseCurrency+"-"+condition.QuoteCurrency]

			currentValue := getMetricValue(condition.BaseMetric, latestMarket)
			pastValue := getMetricValue(condition.BaseMetric, historyMarket)

			percentage := (currentValue - pastValue) / pastValue
			log.Debugf("MARKET %s changed with %.3f", condition.BaseMetric, percentage)
			if percentage < condition.Value {
				return false
				log.Debugf("doAction FALSE: Market %s with percentage difference of %.3f is < than condition value %.3f", condition.BaseMetric, percentage, condition.Value)
			}

		case "percentage-decrease":
			historyMarkets, err := env.DataStore.GetHistoryMarket(condition.TimeframeInMS)
			if err != nil {
				log.Fatal(err)
			}
			historyMarket := historyMarkets[condition.BaseCurrency+"-"+condition.QuoteCurrency]

			currentValue := getMetricValue(condition.BaseMetric, latestMarket)
			pastValue := getMetricValue(condition.BaseMetric, historyMarket)

			percentage := (currentValue - pastValue) / pastValue
			log.Debugf("MARKET %s changed with %.3f", condition.BaseMetric, percentage)
			if percentage > -condition.Value {
				return false
				log.Debugf("COMPARISON Market %s with percentage difference of %.3f is > than condition value -%.3f", condition.BaseMetric, percentage, condition.Value)
			}
		default:
			return false
			log.Warningf("Unknown ConditionType %s", condition.ConditionType)
		}
	}
	return true
}

func ExecuteAction(action models.Action) {
	switch action.ValueType {
	case "absolute":
		log.Infof("Set a %s order for %.8f %s at %.8f %s/%s for a total of %.8f %s.", action.OrderType, action.Quantity, action.QuoteCurrency, action.Value, action.BaseCurrency, action.QuoteCurrency, action.Value*action.Quantity, action.BaseCurrency)
	case "relative-above", "relative-below":
		log.Infof("Set a %s order for %.8f %s at the future rate of %s +/- %.8f %s/%s per unit.", action.OrderType, action.Quantity, action.QuoteCurrency, action.ValueQuoteMetric, action.Value, action.BaseCurrency, action.QuoteCurrency)
	case "percentage-above", "percentage-below":
		log.Infof("Set a %s order for %.8f %s at the future rate of %s * (1 +/- %.8f %s/%s per unit.", action.OrderType, action.Quantity, action.QuoteCurrency, action.ValueQuoteMetric, action.Value, action.BaseCurrency, action.QuoteCurrency)
	}
}

// To get the current market value of a metric
func getMetricValue(baseMetric string, market models.Market) float64 {
	var currentValue float64

	switch baseMetric {
	case "price-last":
		currentValue = market.Last
	case "price-ask":
		currentValue = market.Ask
	case "price-bid":
		currentValue = market.Bid
	case "volume":
		currentValue = market.Volume
	default: // This could be dangerous because now the value 0 is returned
		// TODO Need to find a solution that doesn't cause a run-time panic but does stop this worker
		log.Errorf("models.Condition BaseMetric %s does not exist", baseMetric)
	}

	return currentValue
}
