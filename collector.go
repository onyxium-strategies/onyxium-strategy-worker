package main

import (
	"encoding/json"
	"github.com/onyxium-strategies/onyxium-strategy-worker/models"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type CollectorBody struct {
	Name   string        `json:"name"`
	UserId string        `json:"userId"`
	Tree   []interface{} `json:"tree"`
}

// Collects requests from the frontend, and place strategy in the database
func StrategyCreateCollector(w http.ResponseWriter, r *http.Request) {
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
	strategy, err := models.NewStrategy(collector.Name, collector.UserId, binaryTree)
	if err != nil {
		respondWithError(w, 400, err.Error())
		log.Infof("Bad request StrategyCreate, responded with error msg: %s", err)
	}

	err = env.DataStore.StrategyCreate(strategy)
	if err != nil {
		respondWithError(w, 400, err.Error())
	}
	log.Info("Strategy created")

	// And let the user know their work request was created.
	payload := map[string]interface{}{
		"id":        strategy.Id,
		"name":      strategy.Name,
		"status":    strategy.Status,
		"state":     strategy.State,
		"createdAt": strategy.CreatedAt.Unix(),
		"updatedAt": strategy.UpdatedAt.Unix(),
		"tree":      strategy.Tree.ToKaryArray(),
	}
	respondWithJSON(w, http.StatusOK, payload)
}

func IdleStategyCollector() {
	// Puts work into the WorkQueue
	go func() {
		for {
			strategies, err := env.DataStore.StrategiesGetIdle()
			if len(strategies) > 0 {
				log.Info("Dispatching idle strategies")
			}
			if err != nil {
				log.Fatal(err)
			}
			for _, strategy := range strategies {
				log.Infof("strategy: %s", strategy.Id.Hex())
				WorkQueue <- strategy
			}
			time.Sleep(time.Second) // TODO: check if a second has passed instead of waiting one second. Of nog beter is eigenlijk gwn een signaaltje krijgen als er nieuwe data is of wijzigingen in de db
		}
	}()
}
