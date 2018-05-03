package main

import (
	"flag"
	"fmt"
	// "log"
	"bitbucket.org/visa-startups/coinflow-strategy-worker/models"
	"net/http"
)

var (
	NWorkers = flag.Int("n", 4, "The number of workers to start")
	HTTPAddr = flag.String("http", "127.0.0.1:8000", "Address to listen for HTTP requests on")
	Verbose  = flag.Int("v", 0, "The level of verbosity of the print statements") //
)

func main() {
	// Parse the command-line flags.

	flag.Parse()

	fmt.Println(*Verbose)

	models.InitDB("localhost")
	defer models.DBCon.Close()

	// use case example
	// market, err := models.GetLatestMarket()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(market.Market["BTC-LTC"])
	// market, err = models.GetHistoryMarket(100000)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(market.Market["BTC-LTC"])

	// Start the dispatcher.
	fmt.Println("Starting the dispatcher")
	StartDispatcher(*NWorkers)

	// Register our collector as an HTTP handler function.
	fmt.Println("Registering the collector")
	http.HandleFunc("/api/work", Collector)

	// Start the HTTP server!
	fmt.Println("HTTP server listening on", *HTTPAddr)
	if err := http.ListenAndServe(*HTTPAddr, nil); err != nil {
		fmt.Println(err.Error())
	}
}
