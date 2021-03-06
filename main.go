package main

import (
	"github.com/onyxium-strategies/onyxium-strategy-worker/models"
	// "encoding/json"
	"flag"
	omg "github.com/Alainy/OmiseGo-Go-SDK"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/johntdyer/slackrus"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

type Env struct {
	DataStore models.DataStore
	Ledger    omg.EwalletAdminAPI
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
	env.DataStore = db

	// TEMPORARY To test ToKaryArray - needs to be a unit test in future
	// strategy, _ := env.DataStore.StrategyGet("5b68598724979051f911631d")
	// array, _ := strategy.Tree.ToKaryArray()
	// log.Infof("%+v", array)
	// response, _ := json.Marshal(array)
	// log.Info(string(response))
	// log.Fatal("stop")

	// Connect to the local ledger
	ledger, err := initOMGClient()
	if err != nil {
		log.Fatal(err)
	}
	env.Ledger = ledger

	err = initOmisego()
	if err != nil {
		log.Fatal(err)
	}

	// Start the dispatcher.
	log.Infof("Starting the dispatcher with %d workers", *NWorkers)
	StartDispatcher(*NWorkers)

	// Start paused strategies collector
	IdleStategyCollector()

	router := mux.NewRouter()
	s := router.PathPrefix("/api").Subrouter()
	// Register our collector as an HTTP handler function.
	s.Path("/login").HandlerFunc(Login).Methods("POST")

	s.Path("/user").HandlerFunc(UserAll).Methods("GET")
	s.Path("/user/{id}").HandlerFunc(UserGet).Methods("GET")
	s.Path("/user/{id}").HandlerFunc(UserUpdate).Methods("PUT")
	s.Path("/user").HandlerFunc(UserCreate).Methods("POST")
	s.Path("/user/{id}").HandlerFunc(UserDelete).Methods("DELETE")
	s.Path("/confirm-email").Queries("token", "{token}").Queries("id", "{id}").HandlerFunc(EmailConfirm).Methods("GET")

	s.Path("/strategy").HandlerFunc(StrategyCreateCollector).Methods("POST")
	s.Path("/strategy").HandlerFunc(StrategyAll).Methods("GET")
	s.Path("/strategy/{id}").HandlerFunc(StrategyGet).Methods("GET")
	s.Path("/strategy/{id}").HandlerFunc(StrategyUpdate).Methods("PUT")
	s.Path("/strategy/{id}").HandlerFunc(StrategyDelete).Methods("DELETE")

	s.Path("/balances/{userId}").HandlerFunc(BalancesGet).Methods("GET")
	s.Path("/transactions/{userId}").HandlerFunc(TransactionsGet).Methods("GET")

	// Start the HTTP server!
	log.Infof("HTTP server listening on %s", *HTTPAddr)
	if err := http.ListenAndServe(*HTTPAddr, handlers.LoggingHandler(os.Stdout, router)); err != nil {
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
		IconEmoji:      ":golang:",
		Username:       "Go",
	})
}
