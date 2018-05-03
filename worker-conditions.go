package main

import (
	"fmt"
	"log"
	"time"
	"worker-queue/models"
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
		fmt.Errorf("Condition BaseMetric %s does not exist", baseMetric)
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

		if *Verbose >= 3 {
			fmt.Printf("\n%x condition iterations after last action", i)
			fmt.Println("\nCURRENT ROOT: ", root)
			fmt.Println("CURRENT NODE: ", tree)
		}

		// if all conditions true do action then Tree.Left
		// else go to next sibling Tree.Right
		doAction := true

		for _, condition := range tree.Conditions {
			latestMarket := latestMarkets.Market[condition.BaseCurrency+"-"+condition.QuoteCurrency]
			if *Verbose >= 3 {
				fmt.Println("MARKET last:", latestMarket.Last)
			}

			switch condition.ConditionType {
			case "absolute-above":
				currentValue := getMetricValue(condition.BaseMetric, latestMarket)
				fmt.Println("Currenvalue", currentValue)
				if currentValue < condition.Value {
					doAction = false

					if *Verbose >= 3 {
						fmt.Printf("COMPARISON Market %s with value %.8f is < than condition value %.8f\n", condition.BaseMetric, currentValue, condition.Value)
					}
				}
			case "absolute-below":
				currentValue := getMetricValue(condition.BaseMetric, latestMarket)

				if currentValue > condition.Value {
					doAction = false

					if *Verbose >= 3 {
						fmt.Printf("COMPARISON Market %s with value %.8f is > than condition value %.8f\n", condition.BaseMetric, currentValue, condition.Value)
					}
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
					if *Verbose >= 3 {
						fmt.Printf("COMPARISON Market %s with percentage difference of %.3f is < than condition value %.3f\n", condition.BaseMetric, percentage, condition.Value)
					}
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
					if *Verbose >= 3 {
						fmt.Printf("COMPARISON Market %s with percentage difference of %.3f is > than condition value -%.3f\n", condition.BaseMetric, percentage, condition.Value)
					}
				}
			}

		}

		if doAction {
			if *Verbose >= 3 {
				fmt.Println("ACTION:", tree.Action)
			}
			if tree.Left == nil {
				if *Verbose >= 3 {
					fmt.Println("\nNO MORE STATEMENT AFTER THIS ACTION STATEMENT, I'M DONE")
				}
				tree = nil
			} else {
				tree = tree.Left
				root = root.Left
				if *Verbose >= 3 {
					fmt.Println("JUMPING to left")
				}
			}
			i = 0
		} else {
			if tree.Right == nil {
				if *Verbose >= 3 {
					fmt.Println("JUMPING to root: ", root)
				}
				time.Sleep(3 * time.Second)
				tree = root
				i += 1
			} else {
				if *Verbose >= 3 {
					fmt.Println("JUMPING to right")
				}
				tree = tree.Right
			}
		}

	}

}
