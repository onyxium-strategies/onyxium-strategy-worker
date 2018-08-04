package main

import (
	"bitbucket.org/onyxium/onyxium-strategy-worker/models"
	log "github.com/sirupsen/logrus"
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

	// Checks if there is any work and then dispatches it to a free worker
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

}
