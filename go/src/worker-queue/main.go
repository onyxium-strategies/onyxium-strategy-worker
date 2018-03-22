package main

import (
	"../database"
	"flag"
	"fmt"
	"gopkg.in/mgo.v2"
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

	fmt.Println(Verbose)

	// Start and share database connection
	// https://stackoverflow.com/questions/31218008/sharing-a-globally-defined-db-conn-with-multiple-packages-in-golang
	// Possible interesting: https://hackernoon.com/how-to-work-with-databases-in-golang-33b002aa8c47
	var err error
	database.DBCon, err = mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer database.DBCon.Close()
	database.DBCon.SetMode(mgo.Monotonic, true)
	fmt.Println(database.DBCon)

	// Start the dispatcher.
	fmt.Println("Starting the dispatcher")
	StartDispatcher(*NWorkers)

	// Register our collector as an HTTP handler function.
	fmt.Println("Registering the collector")
	http.HandleFunc("/work", Collector)

	// Start the HTTP server!
	fmt.Println("HTTP server listening on", *HTTPAddr)
	if err := http.ListenAndServe(*HTTPAddr, nil); err != nil {
		fmt.Println(err.Error())
	}
}
