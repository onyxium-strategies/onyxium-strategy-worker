package main

import (
	"bitbucket.org/visa-startups/coinflow-strategy-worker/models"
	log "github.com/sirupsen/logrus"
	"time"
)

type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
}

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func NewWorker(id int, workerQueue chan chan WorkRequest) Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkRequest),
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
			case work := <-w.Work:
				// Receive a work request.
				log.Infof("Worker %d Received work request %d", w.ID, work.ID)

				w.Walk(work.Tree.Left, work.Tree.Left)
				log.Infof("Worker %d work is done", w.ID)

			case <-w.QuitChan:
				// We have been asked to stop.
				log.Infof("Worker %d stopping", w.ID)
				return
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

// The worker will walk over the strategy tree till it reaches a leaf.
func (w *Worker) Walk(tree *Tree, root *Tree) {
	i := 0
	for tree != nil {

		// get latest market update
		latestMarkets, err := env.db.GetLatestMarket()
		if err != nil {
			log.Fatal(err)
		}

		log.Infof("Worker %d is in CURRENT NODE: %+v", w.ID, tree)

		// If all conditions true do action then Tree.Left else go to next sibling Tree.Right
		// If order to check if all conditions are true we check if not any is false
		// All true == not any False
		doAction := true

		for _, condition := range tree.Conditions {
			if condition.ConditionType == "greater-than-or-equal-to" || condition.ConditionType == "less-than-or-equal-to" {
				log.Infof("If the %s on the market %s/%s is %s than %.8f.", condition.BaseMetric, condition.BaseCurrency, condition.QuoteCurrency, condition.ConditionType, condition.Value)
			} else {
				log.Infof("If the %s on the market %s/%s has %s with %.3f percentage within %d minutes.", condition.BaseMetric, condition.BaseCurrency, condition.QuoteCurrency, condition.ConditionType, condition.Value, condition.TimeframeInMS/60000)
			}

			latestMarket := latestMarkets.Market[condition.BaseCurrency+"-"+condition.QuoteCurrency]

			switch condition.ConditionType {
			case "greater-than-or-equal-to":
				currentValue := getMetricValue(condition.BaseMetric, latestMarket)
				log.Debugf("MARKET %s: %.8f", condition.BaseMetric, currentValue)
				if currentValue < condition.Value {
					doAction = false
					log.Debugf("doAction FALSE: Market %s with value %.8f is < than condition value %.8f", condition.BaseMetric, currentValue, condition.Value)
				}
			case "less-than-or-equal-to":
				currentValue := getMetricValue(condition.BaseMetric, latestMarket)
				log.Debugf("MARKET %s: %.8f", condition.BaseMetric, currentValue)
				if currentValue > condition.Value {
					doAction = false
					log.Debugf("doAction FALSE: Market %s with value %.8f is > than condition value %.8f", condition.BaseMetric, currentValue, condition.Value)
				}
			case "percentage-increase":
				historyMarkets, err := env.db.GetHistoryMarket(condition.TimeframeInMS)
				if err != nil {
					log.Fatal(err)
				}
				historyMarket := historyMarkets.Market[condition.BaseCurrency+"-"+condition.QuoteCurrency]

				currentValue := getMetricValue(condition.BaseMetric, latestMarket)
				pastValue := getMetricValue(condition.BaseMetric, historyMarket)

				percentage := (currentValue - pastValue) / pastValue
				log.Debugf("MARKET %s changed with %.3f", condition.BaseMetric, percentage)
				if percentage < condition.Value {
					doAction = false
					log.Debugf("doAction FALSE: Market %s with percentage difference of %.3f is < than condition value %.3f", condition.BaseMetric, percentage, condition.Value)
				}

			case "percentage-decrease":
				historyMarkets, err := env.db.GetHistoryMarket(condition.TimeframeInMS)
				if err != nil {
					log.Fatal(err)
				}
				historyMarket := historyMarkets.Market[condition.BaseCurrency+"-"+condition.QuoteCurrency]

				currentValue := getMetricValue(condition.BaseMetric, latestMarket)
				pastValue := getMetricValue(condition.BaseMetric, historyMarket)

				percentage := (currentValue - pastValue) / pastValue
				log.Debugf("MARKET %s changed with %.3f", condition.BaseMetric, percentage)
				if percentage > -condition.Value {
					doAction = false
					log.Debugf("COMPARISON Market %s with percentage difference of %.3f is > than condition value -%.3f", condition.BaseMetric, percentage, condition.Value)
				}
			default:
				doAction = false
				log.Warningf("Unknown ConditionType %s", condition.ConditionType)
			}

		}

		if doAction {
			switch tree.Action.ValueType {
			case "absolute":
				log.Infof("Set a %s order for %.8f %s at %.8f %s/%s for a total of %.8f %s.", tree.Action.OrderType, tree.Action.Quantity, tree.Action.QuoteCurrency, tree.Action.Value, tree.Action.BaseCurrency, tree.Action.QuoteCurrency, tree.Action.Value*tree.Action.Quantity, tree.Action.BaseCurrency)
			case "relative-above", "relative-below":
				log.Infof("Set a %s order for %.8f %s at the future rate of %s +/- %.8f %s/%s per unit.", tree.Action.OrderType, tree.Action.Quantity, tree.Action.QuoteCurrency, tree.Action.ValueQuoteMetric, tree.Action.Value, tree.Action.BaseCurrency, tree.Action.QuoteCurrency)
			case "percentage-above", "percentage-below":
				log.Infof("Set a %s order for %.8f %s at the future rate of %s * (1 +/- %.8f %s/%s per unit.", tree.Action.OrderType, tree.Action.Quantity, tree.Action.QuoteCurrency, tree.Action.ValueQuoteMetric, tree.Action.Value, tree.Action.BaseCurrency, tree.Action.QuoteCurrency)
			}

			if tree.Left == nil {
				log.Info("NO MORE STATEMENT AFTER THIS ACTION STATEMENT, I'M DONE")
				tree = nil
			} else {
				tree = tree.Left
				root = root.Left
				log.Info("JUMPING to left")
			}
			i = 0
		} else {
			if tree.Right == nil {
				log.Infof("None of the conditions of each child are true. JUMPING to root: %+v", root)
				time.Sleep(3 * time.Second) // Change in production to check each time new market data is available
				tree = root
				i += 1
				log.Infof("%d condition iterations after last action", i)
			} else {
				log.Info("None of the conditions are true. JUMPING to right sibling.")
				tree = tree.Right
			}
		}
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
	default:
		log.Errorf("Condition BaseMetric %s does not exist", baseMetric)
	}

	return currentValue
}
