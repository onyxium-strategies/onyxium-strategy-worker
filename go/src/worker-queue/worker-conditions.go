package main

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

func getCurrentMarket() map[string]Market {
	seconds = 10
	toDate := bson.Now()
	toId := bson.NewObjectIdWithTime(toDate)
	fromDate := toDate.Add(-time.Duration(seconds) * time.Second)
	fromId := bson.NewObjectIdWithTime(fromDate)

	// Get the latest record
	toRecord := MarketRecord{}
	err = database.DBCon.DB("coinflow").C("market").Find(bson.M{"_id": bson.M{"$gte": fromId, "$lt": toId}}).Sort("-$natural").One(&toRecord)
	if err != nil {
		log.Fatal(err)
	}
	return toRecord
}

func getHistoricMarket(seconds int) map[string]Market {
	toDate := bson.Now()
	toId := bson.NewObjectIdWithTime(toDate)
	fromDate := toDate.Add(-time.Duration(seconds) * time.Second)
	fromId := bson.NewObjectIdWithTime(fromDate)

	// Get the record that is x seconds old
	fromRecord := MarketRecord{}
	err = database.DBCon.DB("coinflow").C("market").Find(bson.M{"_id": bson.M{"$gte": fromId, "$lt": toId}}).Sort("$natural").One(&fromRecord)
	if err != nil {
		log.Fatal(err)
	}
	return fromRecord
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

				if currentValue <= condition.Value {
					doAction = false

					if *Verbose >= 3 {
						fmt.Println("COMPARISON Market ", condition.BaseMetric, currentValue, " <= than condition value ", condition.Value)
					}
				}
			case "absolute-below":
				currentValue := getMetricValue(condition.BaseMetric, market)

				if currentValue >= condition.Value {
					doAction = false

					if *Verbose >= 3 {
						fmt.Println("COMPARISON Market ", currentValue, " >= than condition value ", condition.Value)
					}
				}
			case "percentage-increase":
				pastMarkets := getHistoricMarket(TimeframeInMS)
				pastMarket := pastMarkets[condition.BaseCurrency+"-"+condition.QuoteCurrency]

				newValue = getMetricValue(condition.BaseMetric, market)
				oldValue = getMetricValue(condition.BaseMetric, pastMarket)

				percentage = (newValue - oldValue) / oldValue

				if percentage >= condition.Value {
					doAction = true
				}

			case "percentage-decrease":
				pastMarkets := getHistoricMarket(1000)
				fmt.Println("percentage-decrease not supported yet")
				doAction = false
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
