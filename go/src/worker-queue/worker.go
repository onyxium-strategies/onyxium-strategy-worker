package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
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

				// Connect with DB
				session, err := mgo.Dial("localhost")
				if err != nil {
					panic(err)
				}
				defer session.Close()

				// Optional. Switch the session to a monotonic behavior.
				session.SetMode(mgo.Monotonic, true)

				c := session.DB("coinflow").C("market")

				// ticker := time.NewTicker(1 * time.Second).C
				// done := make(chan bool)

				result := make(map[string]Market)
				err = c.Find(nil).Limit(1).Sort("-$natural").One(&result)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("Timestamp: %s, Bid: %f, Ask: %f, Last: %f \n", result["BTC-LTC"].TimeStamp, result["BTC-LTC"].Bid, result["BTC-LTC"].Ask, result["BTC-LTC"].Last)

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

// Walk a tree example
func walk(tree *Tree) {
	if tree != nil {
		fmt.Println(tree)
		walk(tree.Left)
		walk(tree.Right)
	}
}
