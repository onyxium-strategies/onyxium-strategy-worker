package main

import (
	"bitbucket.org/onyxium/onyxium-strategy-worker/models"
	log "github.com/sirupsen/logrus"
	"time"
)

var WorkerQueue chan chan *models.Strategy

// A buffered channel that we can send work requests on.
var WorkQueue = make(chan models.Strategy, 100)

func StartDispatcher(nworkers int) {
	// First, initialize the channel we are going to put the workers' work channels into.
	WorkerQueue = make(chan chan *models.Strategy, nworkers)

	// Now, create all of our workers.
	for i := 0; i < nworkers; i++ {
		log.Info("Starting worker", i+1)
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-WorkQueue:
				log.Infof("Received work with id: %s", work.Id.Hex())
				go func(work models.Strategy) {
					worker := <-WorkerQueue
					log.Info("Dispatching work request %s", work.Id.Hex())
					worker <- &work
				}(work)
			}
		}
	}()

	go func() {
		for {
			strategies, err := env.DataStore.GetPausedStrategies()
			if len(strategies) > 0 {
				log.Info("Dispatching paused strategies")
			}
			if err != nil {
				log.Fatal(err)
			}
			for _, strategy := range strategies {
				log.Infof("strategy: %s", strategy.Id.Hex())
				WorkQueue <- strategy
			}
			time.Sleep(time.Second)
		}
	}()
}
