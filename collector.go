package main

import (
	"bitbucket.org/visa-startups/coinflow-strategy-worker/models"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// A buffered channel that we can send work requests on.
var WorkQueue = make(chan *models.Strategy, 100)

// Collects requests from the frontend, and place strategy in workQueue
func Collector(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
	}
	jsonString := string(body)

	jsonTree, err := parseJsonArray(jsonString)
	if err != nil {
		respondWithError(w, 400, err.Error())
		log.Info("Bad request parseJsonArray, responded with error")
		return
	}

	binaryTree, err := parseBinaryTree(jsonTree)
	if err != nil {
		respondWithError(w, 400, err.Error())
		log.Info("Bad request parseBinaryTree, responded with error")
		return
	}
	binaryTree.SetIdsForBinarySearch()
	strategy, err := models.NewStrategy("teststrategy", jsonString, binaryTree)
	if err != nil {
		respondWithError(w, 400, err.Error())
		log.Info("Bad request StrategyCreate, responded with error")
	}

	err = env.DataStore.StrategyCreate(strategy)
	if err != nil {
		respondWithError(w, 400, err.Error())
	}

	log.Info("Workrequest created")

	// Push the work onto the queue.
	WorkQueue <- strategy
	log.Info("Work request queued")

	// And let the user know their work request was created.
	w.WriteHeader(http.StatusCreated)
}
