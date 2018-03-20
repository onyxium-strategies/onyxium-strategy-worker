package main

import (
	"fmt"
	// "gopkg.in/mgo.v2/bson"
	"../database"
	"log"
	"time"
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
			fmt.Println("Add ourselves into the worker queue.")
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				// Receive a work request.
				fmt.Println("worker", w.ID, ": Received work request ", work.ID, work.Tree.Conditions[0].ConditionType)

				// Close the connection after a duration
				timer := time.NewTimer(10 * time.Second)
				go func() {
					<-timer.C
					w.Stop()
				}()

				markets := make(map[string]Market)
				err := database.DBCon.DB("coinflow").C("market").Find(nil).Limit(1).Sort("-$natural").One(&markets)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("Timestamp: %s, Bid: %f, Ask: %f, Last: %f \n", markets["BTC-LTC"].TimeStamp, markets["BTC-LTC"].Bid, markets["BTC-LTC"].Ask, markets["BTC-LTC"].Last)

				walk(work.Tree.Left, work.Tree.Left, markets)

			case <-w.QuitChan:
				// We have been asked to stop.
				fmt.Printf("worker%d stopping\n", w.ID)
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

// Maybe make a Method from this function
func walk(tree *Tree, root *Tree, markets map[string]Market) {
	i := 0
	for tree != nil {

		// get latest market update
		err := database.DBCon.DB("coinflow").C("market").Find(nil).Limit(1).Sort("-$natural").One(&markets)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("\n%x condition iterations after last action", i)
		fmt.Println("\nCURRENT ROOT: ", root)
		fmt.Println("CURRENT NODE: ", tree)

		// if all conditions true do action then Tree.Left
		// else go to next sibling Tree.Right
		doAction := true

		for _, condition := range tree.Conditions {
			market := markets[condition.BaseCurrency+"-"+condition.QuoteCurrency]
			fmt.Println("MARKET last:", market.Last)
			switch condition.ConditionType {
			case "absolute-above":
				switch condition.BaseMetric {
				case "price":
					if market.Last < condition.Value {
						fmt.Println("COMPARISON Market", market.Last, "< than condition value", condition.Value)
						doAction = false
					} else {
						fmt.Println("COMPARISON Market", market.Last, "> than condition value", condition.Value)
					}
				}
			case "absolute-below":
				switch condition.BaseMetric {
				case "price":
					if market.Last > condition.Value {
						doAction = false
					}
				}
			case "absolute-increase":
				fmt.Println("COMPARISON absolute-increase not supported yet, so doAction set to false")
				doAction = false
			case "absolute-decrease":
				fmt.Println("COMPARISON absolute-decrease not supported yet, so doAction set to false")
				doAction = false
			}
		}

		// In next section, new root and or tree will be set for next while loop.

		if doAction {
			fmt.Println("ACTION:", tree.Action)
			if tree.Left == nil {
				fmt.Println("\nNO MORE STATEMENT AFTER THIS ACTION STATEMENT, I'M DONE")
				tree = nil
			} else {
				tree = tree.Left
				root = root.Left
				fmt.Println("JUMPING to left")
			}
			i = 0
		} else {
			if tree.Right == nil {
				fmt.Println("JUMPING to root: ", root)
				time.Sleep(3 * time.Second)
				tree = root
				i += 1
			} else {
				fmt.Println("JUMPING to right")
				tree = tree.Right
			}
		}

	}

}
