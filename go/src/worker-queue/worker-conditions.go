package main

import (
	"fmt"
	// "gopkg.in/mgo.v2/bson"
)

type Market struct {
	MarketName        string      `json:"MarketName"`
	High              float64     `json:"High"`
	Low               float64     `json:"Low"`
	Volume            float64     `json:"Volume"`
	Last              float64     `json:"Last"`
	BaseVolume        float64     `json:"BaseVolume"`
	TimeStamp         string      `json:"TimeStamp"`
	Bid               float64     `json:"Bid"`
	Ask               float64     `json:"Ask"`
	OpenBuyOrders     int         `json:"OpenBuyOrders"`
	OpenSellOrders    int         `json:"OpenSellOrders"`
	PrevDay           float64     `json:"PrevDay"`
	Created           string      `json:"Created"`
	DisplayMarketName interface{} `json:"DisplayMarketName"`
}

// To get the current value of a metric
func getMetricValue(baseMetric string, market Market) float64 {
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
		markets := getCurrentMarket()

		if *Verbose >= 3 {
			fmt.Printf("\n%x condition iterations after last action", i)
			fmt.Println("\nCURRENT ROOT: ", root)
			fmt.Println("CURRENT NODE: ", tree)
		}

		// if all conditions true do action then Tree.Left
		// else go to next sibling Tree.Right
		doAction := true

		for _, condition := range tree.Conditions {
			market := markets[condition.BaseCurrency+"-"+condition.QuoteCurrency]
			if *Verbose >= 3 {
				fmt.Println("MARKET last:", market.Last)
			}

			switch condition.ConditionType {
			case "absolute-above":
				currentValue := getMetricValue(condition.BaseMetric, market)

				if currentValue < condition.Value {
					doAction = false

					if *Verbose >= 3 {
						fmt.Println("COMPARISON Market ", condition.BaseMetric, currentValue, " <= than condition value ", condition.Value)
					}
				}
			case "absolute-below":
				currentValue := getMetricValue(condition.BaseMetric, market)

				if currentValue > condition.Value {
					doAction = false

					if *Verbose >= 3 {
						fmt.Println("COMPARISON Market ", currentValue, " >= than condition value ", condition.Value)
					}
				}
			case "percentage-increase":
				pastMarkets := getHistoryRecord(TimeframeInMS)
				pastMarket := pastMarkets[condition.BaseCurrency+"-"+condition.QuoteCurrency]

				newValue := getMetricValue(condition.BaseMetric, market)
				oldValue := getMetricValue(condition.BaseMetric, pastMarket)

				percentage := (newValue - oldValue) / oldValue

				if percentage < condition.Value {
					doAction = false
				}

			case "percentage-decrease":
				pastMarkets := getHistoryRecord(TimeframeInMS)
				pastMarket := pastMarkets[condition.BaseCurrency+"-"+condition.QuoteCurrency]

				newValue := getMetricValue(condition.BaseMetric, market)
				oldValue := getMetricValue(condition.BaseMetric, pastMarket)

				percentage := (newValue - oldValue) / oldValue

				if percentage > -condition.Value {
					doAction = false
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
