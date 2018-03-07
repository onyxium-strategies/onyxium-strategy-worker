package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// A buffered channel that we can send work requests on.
var WorkQueue = make(chan WorkRequest, 100)

var id int

// Collects requests from the frontend, and place workrequest in workQueue
func Collector(w http.ResponseWriter, r *http.Request) {

	// Make sure we can only be called with an HTTP POST request.
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	myJson := string(body)

	tree, err := parseJson(myJson)
	if err != nil {
		fmt.Println(err)
	}

	root := parseBinaryTree(tree)
	work := WorkRequest{ID: id, Tree: &root}

	// TODO: get last ID from database, use that one + 1
	id = id + 1

	fmt.Println("Workrequest tree created")

	// Push the work onto the queue.
	WorkQueue <- work
	fmt.Println("Work request queued")

	// And let the user know their work request was created.
	w.WriteHeader(http.StatusCreated)

	return
}
