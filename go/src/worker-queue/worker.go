package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
	"worker-queue/database"
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
			fmt.Println("Worker", w.ID, "Add ourselves into the worker queue.")
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				// Receive a work request.
				fmt.Println("Worker", w.ID, "Received work request ", work.ID, work.Tree.Conditions[0].ConditionType)

				walk(work.Tree.Left, work.Tree.Left)
				fmt.Printf("Worker %d's work is done\n", w.ID)

			case <-w.QuitChan:
				// We have been asked to stop.
				fmt.Printf("Worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
//
// Note that the worker will only stop *after* it has finished its work.
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

// To get the current value of a metric
func getCurrentValue(baseMetric string, timeframeInMS int, market Market) float64 {
	switch baseMetric {
	case "price-last":
		currentValue := market.Last
	case "price-ask":
		currentValue := market.Ask
	case "price-bid":
		currentValue := market.Bid
	case "volume":
		currentValue := market.Volume
	}

	return currentValue
}

// Walk: https://en.m.wikipedia.org/wiki/Left-child_right-sibling_binary_tree
// Maybe make a Method from this function
func walk(tree *Tree, root *Tree) {
	i := 0
	for tree != nil {

		// get latest market update
		markets := make(map[string]Market)
		err := database.DBCon.DB("coinflow").C("market").Find(nil).Limit(1).Sort("-$natural").One(&markets)
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
			market := markets[condition.BaseCurrency+"-"+condition.QuoteCurrency]
			if *Verbose >= 3 {
				fmt.Println("MARKET last:", market.Last)
			}
			switch condition.ConditionType {
			case "absolute-above":
				var currentValue float64
				switch condition.BaseMetric {
				case "price-last":
					currentValue := market.Last
					if currentValue <= condition.Value {
						doAction = false
					}
				case "price-ask":
					currentValue := market.Ask
					if currentValue <= condition.Value {
						doAction = false
					}
				case "price-bid":
					currentValue := market.Bid
					if currentValue <= condition.Value {
						doAction = false
					}
				case "volume":
					currentValue := market.Volume
					if currentValue <= condition.Value {
						doAction = false
					}
				default:
					fmt.Errorf("Condition BaseMetric %s does not exist", condition.BaseMetric)
				}
				if *Verbose >= 3 {
					fmt.Println("COMPARISON Market ", currentValue, " <= than condition value ", condition.Value)
				}
			case "absolute-below":
				var currentValue float64
				switch condition.BaseMetric {
				case "price-last":
					currentValue := market.Last
				case "price-ask":
					currentValue := market.Ask
				case "price-bid":
					currentValue := market.Bid
				case "volume":
					currentValue := market.Volume
				default:
					fmt.Errorf("Condition BaseMetric %s does not exist", condition.BaseMetric)
				}
				if *Verbose >= 3 {
					fmt.Println("COMPARISON Market ", currentValue, " >= than condition value ", condition.Value)
				}
				if currentValue >= condition.Value {
					doAction = false
				}
			case "percentage-increase":
				currentTime := time.Now()
				pastTime := currentTime.Add(-time.Duration(condition.TimeframeInMS) * time.Millisecond)
				fmt.Println(pastTime)
				pastMarkets := make(map[string]Market)
				err := database.DBCon.DB("coinflow").C("market").Find(nil).Limit(1).Sort("$natural").One(&pastMarkets)
				// err := database.DBCon.DB("coinflow").C("market").Find(bson.M{"_id": {"$lt": pastTime}}).Limit(1).Sort("-$natural").One(&pastMarkets)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(pastMarkets[condition.BaseCurrency+"-"+condition.QuoteCurrency].Last)

				// if *Verbose >= 3 {
				// 	fmt.Println("COMPARISON Market ", currentValue, " >= than condition value ", condition.Value)
				// }
				// if currentValue >= condition.Value {
				// 	doAction = false
				// }
			case "percentage-decrease":
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
