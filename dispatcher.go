package main

import (
	"bitbucket.org/visa-startups/coinflow-strategy-worker/models"
	log "github.com/sirupsen/logrus"
	"time"
)

var WorkerQueue chan chan *models.Strategy

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
				log.Info("Received work requeust")
				go func(work *models.Strategy) {
					worker := <-WorkerQueue

					log.Info("Dispatching work request")
					worker <- work
				}(work)
			}
		}
	}()

	go func() {
		for {
			strategies, err := env.DataStore.GetPausedStrategies()
			log.Info("Dispatching paused strategies")
			if err != nil {
				log.Fatal(err)
			}
			for _, strategy := range strategies {
				WorkQueue <- &strategy
			}
			time.Sleep(time.Second)
		}
	}()
}
