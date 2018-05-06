package main

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// A buffered channel that we can send work requests on.
var WorkQueue = make(chan WorkRequest, 100)

var id int

type CollectorHandler struct{}

// Collects requests from the frontend, and place workrequest in workQueue
func (c *CollectorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Make sure we can only be called with an HTTP POST request.
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
	}
	myJson := string(body)

	tree, err := parseJsonArray(myJson)
	if err != nil {
		log.Error(err)
	}

	root := parseBinaryTree(tree)
	work := WorkRequest{ID: id, Tree: &root}
	log.Info("Workrequest created")

	// TODO: get last ID from database, use that one + 1
	id = id + 1

	// Push the work onto the queue.
	WorkQueue <- work
	log.Info("Work request queued")

	// And let the user know their work request was created.
	w.WriteHeader(http.StatusCreated)
}
