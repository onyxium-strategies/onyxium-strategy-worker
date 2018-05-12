package main

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// A buffered channel that we can send work requests on.
var WorkQueue = make(chan WorkRequest, 100)

var id int

// Collects requests from the frontend, and place workrequest in workQueue
func Collector(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
	}
	jsonString := string(body)

	jsonTree, err := parseJsonArray(jsonString)
	if err != nil {
		respondWithError(w, 400, err.Error())
		log.Info("Bad request, responded with error")
		return
	}

	binaryTree, err := parseBinaryTree(jsonTree)
	if err != nil {
		respondWithError(w, 400, err.Error())
		log.Info("Bad request, responded with error")
		return
	}

	work := WorkRequest{ID: id, Tree: &binaryTree}
	log.Info("Workrequest created")

	// TODO: get last ID from database, use that one + 1
	id = id + 1

	// Push the work onto the queue.
	WorkQueue <- work
	log.Info("Work request queued")

	// And let the user know their work request was created.
	w.WriteHeader(http.StatusCreated)
}
