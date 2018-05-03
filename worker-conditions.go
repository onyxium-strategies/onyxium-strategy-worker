package main

import (
	"bitbucket.org/visa-startups/coinflow-strategy-worker/models"
	log "github.com/sirupsen/logrus"
	"time"
	// "gopkg.in/mgo.v2/bson"
)

// To get the current value of a metric
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

// Walk: https://en.m.wikipedia.org/wiki/Left-child_right-sibling_binary_tree
// Maybe make a Method from this function
func walk(tree *Tree, root *Tree) {
	i := 0
	for tree != nil {

		// get latest market update
		latestMarkets, err := models.GetLatestMarket()
		if err != nil {
			log.Fatal(err)
		}

		log.Infof("%d condition iterations after last action", i)
		log.Infof("CURRENT ROOT: %+v", root)
		log.Infof("CURRENT NODE: %+v", tree)

		// if all conditions true do action then Tree.Left
		// else go to next sibling Tree.Right
		doAction := true

		for _, condition := range tree.Conditions {
			latestMarket := latestMarkets.Market[condition.BaseCurrency+"-"+condition.QuoteCurrency]

			log.Infof("MARKET last: %8f", latestMarket.Last)

			switch condition.ConditionType {
			case "absolute-above":
				currentValue := getMetricValue(condition.BaseMetric, latestMarket)
				log.Infof("Current value %8f", currentValue)
				if currentValue < condition.Value {
					doAction = false

					log.Infof("COMPARISON Market %s with value %.8f is < than condition value %.8f", condition.BaseMetric, currentValue, condition.Value)
				}
			case "absolute-below":
				currentValue := getMetricValue(condition.BaseMetric, latestMarket)

				if currentValue > condition.Value {
					doAction = false

					log.Infof("COMPARISON Market %s with value %.8f is > than condition value %.8f", condition.BaseMetric, currentValue, condition.Value)
				}
			case "percentage-increase":
				historyMarkets, err := models.GetHistoryMarket(condition.TimeframeInMS)
				if err != nil {
					log.Fatal(err)
				}
				historyMarket := historyMarkets.Market[condition.BaseCurrency+"-"+condition.QuoteCurrency]

				newValue := getMetricValue(condition.BaseMetric, latestMarket)
				oldValue := getMetricValue(condition.BaseMetric, historyMarket)

				percentage := (newValue - oldValue) / oldValue

				if percentage < condition.Value {
					doAction = false

					log.Infof("COMPARISON Market %s with percentage difference of %.3f is < than condition value %.3f", condition.BaseMetric, percentage, condition.Value)
				}

			case "percentage-decrease":
				historyMarkets, err := models.GetHistoryMarket(condition.TimeframeInMS)
				if err != nil {
					log.Fatal(err)
				}
				historyMarket := historyMarkets.Market[condition.BaseCurrency+"-"+condition.QuoteCurrency]

				newValue := getMetricValue(condition.BaseMetric, latestMarket)
				oldValue := getMetricValue(condition.BaseMetric, historyMarket)

				percentage := (newValue - oldValue) / oldValue

				if percentage > -condition.Value {
					doAction = false

					log.Infof("COMPARISON Market %s with percentage difference of %.3f is > than condition value -%.3f", condition.BaseMetric, percentage, condition.Value)
				}
			}

		}

		if doAction {
			log.Info("ACTION:", tree.Action)
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
				log.Infof("JUMPING to root: %+v", root)
				time.Sleep(3 * time.Second)
				tree = root
				i += 1
			} else {
				log.Info("JUMPING to right")
				tree = tree.Right
			}
		}

	}

}
