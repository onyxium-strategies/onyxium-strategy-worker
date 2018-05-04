package main

import (
	log "github.com/sirupsen/logrus"
)

var WorkerQueue chan chan WorkRequest

func StartDispatcher(nworkers int) {
	// First, initialize the channel we are going to put the workers' work channels into.
	WorkerQueue = make(chan chan WorkRequest, nworkers)

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
				log.Info("Received work requeust")
				go func() {
					worker := <-WorkerQueue

					log.Info("Dispatching work request")
					worker <- work
				}()
			}
		}
	}()
}