package main

import (
	"fmt"
	"net/http"
	"time"
)

// A buffered channel that we can send work requests on.
var WorkQueue = make(chan WorkRequest, 100)

func Collector(w http.ResponseWriter, r *http.Request) {
  // Make sure we can only be called with an HTTP POST request.
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	json, err := parseJson(r.FormValue("json"))
	if err != nil {
		http.Error(w, err, http.StatusBadRequest)
		return
	}


  // Now, we information and make a WorkRequest out of them.
	// condition := Condition{ConditionType: "percentage-increase", BaseCurrency: "ETH", QuoteCurrency: "OMG", TimeframeInMS: 3600000, BaseMetric: "price", Value: 20}
	// action := Action{Type0: "order", OrderType: "limit-buy", OrderValueType: "absolute", BaseCurrency: "ETH", QuoteCurrency: "OMG", Quantity: 100, Value: 0.012}
	// leftchild := Tree{Left: nil, Right: nil, Condition: condition, Action: action}
	// root := Tree{Left: &leftchild, Right: nil, Condition: condition, Action: action}



	work := WorkRequest{Id: 1, Tree: &root}
	fmt.Println("Workrequest tree created")


  // Push the work onto the queue.
	WorkQueue <- work
	fmt.Println("Work request queued")

  // And let the user know their work request was created.
	w.WriteHeader(http.StatusCreated)
	return
}
