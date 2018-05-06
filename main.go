package main

import (
	"bitbucket.org/visa-startups/coinflow-strategy-worker/models"
	"flag"
	"github.com/johntdyer/slackrus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

type Env struct {
	db models.Datastore
}

var (
	NWorkers = flag.Int("n", 4, "The number of workers to start")
	HTTPAddr = flag.String("http", "127.0.0.1:8000", "Address to listen for HTTP requests on")
	Verbose  = flag.Int("v", 1, "The level of verbosity of the log statements")
	env      = &Env{}
)

func main() {
	// Parse the command-line flags.
	flag.Parse()

	initLogger(*Verbose)

	db, err := models.InitDB("localhost")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	env.db = db

	// Start the dispatcher.
	log.Infof("Starting the dispatcher with %d workers", *NWorkers)
	StartDispatcher(*NWorkers)

	// Register our collector as an HTTP handler function.
	collectorUrl := "/api/work"
	log.Infof("Registering the collector on %s", collectorUrl)
	http.Handle(collectorUrl, &CollectorHandler{})

	// Start the HTTP server!
	log.Infof("HTTP server listening on %s", *HTTPAddr)
	if err := http.ListenAndServe(*HTTPAddr, nil); err != nil {
		log.Fatal(err)
	}
}

func initLogger(level int) {
	log.SetOutput(os.Stdout)

	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)

	switch level {
	case 0:
		log.SetLevel(log.DebugLevel)
	case 1:
		log.SetLevel(log.InfoLevel)
	case 2:
		log.SetLevel(log.WarnLevel)
	case 3:
		log.SetLevel(log.ErrorLevel)
	case 4:
		log.SetLevel(log.FatalLevel)
	case 5:
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.WarnLevel)
	}

	log.AddHook(&slackrus.SlackrusHook{
		HookURL:        "https://hooks.slack.com/services/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",
		AcceptedLevels: slackrus.LevelThreshold(log.PanicLevel),
		Channel:        "#development",
		IconEmoji:      ":robot_face:",
		Username:       "workerbot",
	})
}
