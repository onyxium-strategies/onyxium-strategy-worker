package main

import (
	"bitbucket.org/onyxium/onyxium-strategy-worker/models"
	log "github.com/sirupsen/logrus"
	// "time"
	omg "github.com/Alainy/OmiseGo-Go-SDK"
	"os"
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

				root, err := work.BsonTree.Search(work.State)
				if err != nil {
					log.Error(err)
					w.Stop()
					continue
				}
				if root.Order != nil {
					err = w.CheckOrderFilled(root, work)
					if err != nil {
						log.Error(err)
						w.Stop()
						continue
					}
				} else {
					_, err = w.WalkSiblings(root, work)
					if err != nil {
						log.Error(err)
						w.Stop()
						continue
					}
				}
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
		doAction, currentValue := CheckConditions(tree.Conditions, latestMarkets)
		if doAction {
			strategy.Status = "executing"
			strategy.State = tree.Id
			err = env.DataStore.StrategyUpdate(strategy)
			if err != nil {
				return nil, err
			}
			ExecuteAction(tree, strategy, currentValue)
			strategy.Status = "paused"
			err = env.DataStore.StrategyUpdate(strategy)
			if err != nil {
				return nil, err
			}
			// return strategy, nil
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

// TODO in the future just get the data from the exchange, for now manually check it
func (w *Worker) CheckOrderFilled(tree *models.Tree, strategy *models.Strategy) error {
	latestMarkets, err := env.DataStore.GetLatestMarket()
	if err != nil {
		return err
	}
	latestMarket := latestMarkets[tree.Action.BaseCurrency+"-"+tree.Action.QuoteCurrency]

	var filled bool
	switch tree.Action.OrderType {
	case "limit-buy":
		if tree.Order.Rate > latestMarket.Last {
			filled = true
		}
	case "limit-sell":
		if tree.Order.Rate < latestMarket.Last {
			filled = true
		}
	}
	if filled {
		tree.Order.Status = "filled"
		err = env.DataStore.StrategyUpdate(strategy)
		if err != nil {
			return err
		}

		// Trade currency
		err = ExchangeTokens(tree, strategy)
		if err != nil {
			return err
		}

		if tree.Left == nil {
			strategy.Status = "stopped"
			err = env.DataStore.StrategyUpdate(strategy)
			if err != nil {
				return err
			}
			log.Info("NO MORE STATEMENT AFTER THIS ACTION STATEMENT, I'M DONE")
			return nil
		} else {
			strategy.Status = "paused"
			strategy.State = tree.Left.Id
			err = env.DataStore.StrategyUpdate(strategy)
			if err != nil {
				return err
			}
			log.Info("JUMPING to left")
			return nil
		}
	}
	log.Info("Order is still pending.")
	strategy.Status = "paused"
	err = env.DataStore.StrategyUpdate(strategy)
	return nil
}

func ExchangeTokens(tree *models.Tree, strategy *models.Strategy) error {
	var baseTokenId, quoteTokenId string

	listBody := omg.ListParams{}
	tokens, err := env.Ledger.TokenAll(listBody)
	if err != nil {
		return err
	}
	for _, token := range tokens.Data {
		if token.Symbol == tree.Action.BaseCurrency {
			baseTokenId = token.Id
		}
		if token.Symbol == tree.Action.QuoteCurrency {
			quoteTokenId = token.Id
		}
	}

	// get primary wallet of user
	getWalletBody := omg.ProviderUserIdParam{
		ProviderUserId: strategy.UserId.Hex(),
	}
	walletList, err := env.Ledger.UserGetWalletsByProviderUserId(getWalletBody)
	if err != nil {
		return err
	}

	// BUY: -base, +quote
	// SELL: +base, -quote
	for _, wallet := range walletList.Data {
		if wallet.Identifier == "primary" {
			log.Infof("Amount to transfer: %d", int(tree.Action.Value*tree.Action.Quantity*subUnitToUnit))
			// Note that it is important that the user first sends us money so we can confirm that he has enough
			// before we send him back tokens
			if tree.Action.OrderType == "limit-buy" {
				transactionBody := omg.TransactionCreateParams{
					IdempotencyToken: omg.NewIdempotencyToken(),
					FromAddress:      wallet.Address,
					ToAddress:        os.Getenv("primaryWalletAddress"),
					TokenId:          baseTokenId,
					Amount:           int(tree.Action.Value * tree.Action.Quantity * subUnitToUnit),
				}
				_, err = env.Ledger.TransactionCreate(transactionBody)
				if err != nil {
					return err
				}
				transactionBody = omg.TransactionCreateParams{
					IdempotencyToken: omg.NewIdempotencyToken(),
					FromAddress:      os.Getenv("primaryWalletAddress"),
					ToAddress:        wallet.Address,
					TokenId:          quoteTokenId,
					Amount:           int(tree.Action.Value * tree.Action.Quantity * subUnitToUnit),
				}
				_, err = env.Ledger.TransactionCreate(transactionBody)
				if err != nil {
					return err
				}
			}
			if tree.Action.OrderType == "limit-sell" {
				transactionBody := omg.TransactionCreateParams{
					IdempotencyToken: omg.NewIdempotencyToken(),
					FromAddress:      wallet.Address,
					ToAddress:        os.Getenv("primaryWalletAddress"),
					TokenId:          quoteTokenId,
					Amount:           int(tree.Action.Value * tree.Action.Quantity * subUnitToUnit),
				}
				_, err = env.Ledger.TransactionCreate(transactionBody)
				if err != nil {
					return err
				}
				transactionBody = omg.TransactionCreateParams{
					IdempotencyToken: omg.NewIdempotencyToken(),
					FromAddress:      os.Getenv("primaryWalletAddress"),
					ToAddress:        wallet.Address,
					TokenId:          baseTokenId,
					Amount:           int(tree.Action.Value * tree.Action.Quantity * subUnitToUnit),
				}
				_, err = env.Ledger.TransactionCreate(transactionBody)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func CheckConditions(conditions []models.Condition, latestMarkets map[string]models.Market) (bool, float64) {
	var currentValue float64
	for _, condition := range conditions {
		if condition.ConditionType == "greater-than-or-equal-to" || condition.ConditionType == "less-than-or-equal-to" {
			log.Infof("If the %s on the market %s/%s is %s than %.8f.", condition.BaseMetric, condition.BaseCurrency, condition.QuoteCurrency, condition.ConditionType, condition.Value)
		} else {
			log.Infof("If the %s on the market %s/%s has %s with %.3f percentage within %d minutes.", condition.BaseMetric, condition.BaseCurrency, condition.QuoteCurrency, condition.ConditionType, condition.Value, condition.TimeframeInMS/60000)
		}

		latestMarket := latestMarkets[condition.BaseCurrency+"-"+condition.QuoteCurrency]

		switch condition.ConditionType {
		case "greater-than-or-equal-to":
			currentValue = getMetricValue(condition.BaseMetric, latestMarket)
			log.Debugf("MARKET %s: %.8f", condition.BaseMetric, currentValue)
			if currentValue < condition.Value {
				return false, currentValue
				log.Debugf("doAction FALSE: Market %s with value %.8f is < than condition value %.8f", condition.BaseMetric, currentValue, condition.Value)
			}
		case "less-than-or-equal-to":
			currentValue = getMetricValue(condition.BaseMetric, latestMarket)
			log.Debugf("MARKET %s: %.8f", condition.BaseMetric, currentValue)
			if currentValue > condition.Value {
				return false, currentValue
				log.Debugf("doAction FALSE: Market %s with value %.8f is > than condition value %.8f", condition.BaseMetric, currentValue, condition.Value)
			}
		case "percentage-increase":
			historyMarkets, err := env.DataStore.GetHistoryMarket(condition.TimeframeInMS)
			if err != nil {
				log.Fatal(err)
			}
			historyMarket := historyMarkets[condition.BaseCurrency+"-"+condition.QuoteCurrency]

			currentValue = getMetricValue(condition.BaseMetric, latestMarket)
			pastValue := getMetricValue(condition.BaseMetric, historyMarket)

			percentage := (currentValue - pastValue) / pastValue
			log.Debugf("MARKET %s changed with %.3f", condition.BaseMetric, percentage)
			if percentage < condition.Value {
				return false, currentValue
				log.Debugf("doAction FALSE: Market %s with percentage difference of %.3f is < than condition value %.3f", condition.BaseMetric, percentage, condition.Value)
			}

		case "percentage-decrease":
			historyMarkets, err := env.DataStore.GetHistoryMarket(condition.TimeframeInMS)
			if err != nil {
				log.Fatal(err)
			}
			historyMarket := historyMarkets[condition.BaseCurrency+"-"+condition.QuoteCurrency]

			currentValue = getMetricValue(condition.BaseMetric, latestMarket)
			pastValue := getMetricValue(condition.BaseMetric, historyMarket)

			percentage := (currentValue - pastValue) / pastValue
			log.Debugf("MARKET %s changed with %.3f", condition.BaseMetric, percentage)
			if percentage > -condition.Value {
				return false, currentValue
				log.Debugf("COMPARISON Market %s with percentage difference of %.3f is > than condition value -%.3f", condition.BaseMetric, percentage, condition.Value)
			}
		default:
			return false, currentValue
			log.Warningf("Unknown ConditionType %s", condition.ConditionType)
		}
	}
	return true, currentValue
}

func ExecuteAction(tree *models.Tree, strategy *models.Strategy, currentValue float64) error {
	var rate float64
	action := tree.Action

	switch action.ValueType {
	case "absolute":
		log.Infof("Set a %s order for %.8f %s at %.8f %s/%s for a total of %.8f %s.", action.OrderType, action.Quantity, action.QuoteCurrency, action.Value, action.BaseCurrency, action.QuoteCurrency, action.Value*action.Quantity, action.BaseCurrency)
		rate = action.Value
	case "relative-above":
		log.Infof("Set a %s order for %.8f %s at the rate of %s + %.8f %s/%s per unit.", action.OrderType, action.Quantity, action.QuoteCurrency, action.ValueQuoteMetric, action.Value, action.BaseCurrency, action.QuoteCurrency)
		rate = currentValue + action.Value
	case "relative-below":
		log.Infof("Set a %s order for %.8f %s at the rate of %s - %.8f %s/%s per unit.", action.OrderType, action.Quantity, action.QuoteCurrency, action.ValueQuoteMetric, action.Value, action.BaseCurrency, action.QuoteCurrency)
		rate = currentValue - action.Value
	case "percentage-above":
		log.Infof("Set a %s order for %.8f %s at the  rate of %s * (1 + %.8f) %s/%s per unit.", action.OrderType, action.Quantity, action.QuoteCurrency, action.ValueQuoteMetric, action.Value, action.BaseCurrency, action.QuoteCurrency)
		rate = currentValue * (1 + action.Value)
	case "percentage-below":
		log.Infof("Set a %s order for %.8f %s at the  rate of %s * (1 - %.8f) %s/%s per unit.", action.OrderType, action.Quantity, action.QuoteCurrency, action.ValueQuoteMetric, action.Value, action.BaseCurrency, action.QuoteCurrency)
		rate = currentValue * (1 - action.Value)
	}
	// TODO set out order on exchange and get orderUUID
	order, err := models.NewOrder("example-remote-order-id", rate)
	if err != nil {
		return err
	}
	tree.Order = order
	err = env.DataStore.StrategyUpdate(strategy)
	return err
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
