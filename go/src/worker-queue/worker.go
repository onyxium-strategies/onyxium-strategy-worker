package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
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

var addr = flag.String("addr", "localhost:8080", "http service address")

// This function "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop.
func (w *Worker) Start() {

	go func() {
		for {
			// Add ourselves into the worker queue.
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				// Receive a work request.
				fmt.Println("worker", w.ID, ": Received work request ", work.ID)

				// Connect with websocket url
				u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
				log.Printf("connecting to %s", u.String())

				conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
				if err != nil {
					log.Fatal("dial:", err)
				}
				defer conn.Close()

				ticker := time.NewTicker(1 * time.Second).C
				done := make(chan bool)

				// Close the connection after a duration
				timer := time.NewTimer(20 * time.Second).C

				for {
					select {
					case <-ticker:
						_, message, err := conn.ReadMessage()
						if err != nil {
							log.Println("read:", err)
							return
						}
						var markets map[string]Market
						if err = json.Unmarshal(message, &markets); err != nil {
							fmt.Println("error:", err)
						}
						fmt.Println(markets["BTC-LTC"])
						fmt.Printf("Bid: %f, Ask: %f, Last: %f", markets["BTC-LTC"].Bid, markets["BTC-LTC"].Ask, markets["BTC-LTC"].Last)
						fmt.Println("\n")

					case <-timer:
						log.Println("timer ended")
						return
					case <-done:
						log.Println("job done")
						return

						// Cleanly close the connection by sending a close message and then
						// waiting (with timeout) for the server to close the connection.
						err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
						if err != nil {
							log.Println("write close:", err)
							return
						}
						select {
						case <-done:
						case <-time.After(time.Second):
						}
						return
					}
				}

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
