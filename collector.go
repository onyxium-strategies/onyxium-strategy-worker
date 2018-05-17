package main

import (
	"bitbucket.org/onyxium/onyxium-strategy-worker/models"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// A buffered channel that we can send work requests on.
var WorkQueue = make(chan *models.Strategy, 100)

type CollectorBody struct {
	Name   string        `json:"name"`
	UserId string        `json:"userId"`
	Tree   []interface{} `json:"tree"`
}

// Collects requests from the frontend, and place strategy in workQueue
func Collector(w http.ResponseWriter, r *http.Request) {
	var collector CollectorBody
	err := json.NewDecoder(r.Body).Decode(&collector)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	binaryTree, err := parseBinaryTree(collector.Tree)
	if err != nil {
		respondWithError(w, 400, err.Error())
		log.Infof("Bad request parseBinaryTree, responded with error msg: %s", err)
		return
	}
	binaryTree.SetIdsForBinarySearch()
	strategy, err := models.NewStrategy(collector.Name, collector.UserId, collector.Tree, binaryTree)
	if err != nil {
		respondWithError(w, 400, err.Error())
		log.Infof("Bad request StrategyCreate, responded with error msg: %s", err)
	}

	err = env.DataStore.StrategyCreate(strategy)
	if err != nil {
		respondWithError(w, 400, err.Error())
	}
	log.Info("Workrequest created")
	// And let the user know their work request was created.
	w.WriteHeader(http.StatusCreated)
}
