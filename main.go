package main

import (
	"bitbucket.org/onyxium/onyxium-strategy-worker/models"
	omg "bitbucket.org/onyxium/onyxium-strategy-worker/omisego"
	"flag"
	"github.com/gorilla/mux"
	"github.com/johntdyer/slackrus"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
)

type Env struct {
	DataStore  models.DataStore
	AdminUser  omg.AdminAPI
	ServerUser omg.EWalletAPI
}

var (
	NWorkers    = flag.Int("n", 4, "The number of workers to start")
	HTTPAddr    = flag.String("http", "127.0.0.1:8000", "Address to listen for HTTP requests on")
	Verbose     = flag.Int("v", 1, "The level of verbosity of the log statements")
	env         = &Env{}
	baseTokenId string
	SeedLedger  = flag.Bool("seed", false, "Seed the database with all available minted tokens for users.")
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

	// Connect to the local ledger
	initOmisego()

	// Start the dispatcher.
	log.Infof("Starting the dispatcher with %d workers", *NWorkers)
	StartDispatcher(*NWorkers)

	router := mux.NewRouter()
	s := router.PathPrefix("/api").Subrouter()
	// Register our collector as an HTTP handler function.
	s.Path("/work").HandlerFunc(Collector).Methods("POST")
	s.Path("/user").HandlerFunc(UserAll).Methods("GET")
	s.Path("/user/{id}").HandlerFunc(UserGet).Methods("GET")
	s.Path("/user/{id}").HandlerFunc(UserUpdate).Methods("PUT")
	s.Path("/user").HandlerFunc(UserCreate).Methods("POST")
	s.Path("/user/{id}").HandlerFunc(UserDelete).Methods("DELETE")
	s.Path("/confirm-email").Queries("token", "{token}").Queries("id", "{id}").HandlerFunc(EmailConfirm).Methods("GET")

	// Start the HTTP server!
	log.Infof("HTTP server listening on %s", *HTTPAddr)
	if err := http.ListenAndServe(*HTTPAddr, router); err != nil {
		log.Fatal(err)
	}
}

func initOmisego() {
	// Get authentication and connection to the eWallet and Admin API.
	adminURL := &url.URL{
		Scheme: "http",
		Host:   "localhost:4000",
		Path:   "/admin/api",
	}
	ewalletURL := &url.URL{
		Scheme: "http",
		Host:   "localhost:4000",
		Path:   "/api",
	}

	loginBody := omg.LoginParams{
		Email:    os.Getenv("email"),
		Password: os.Getenv("pwd"),
	}
	client, err := omg.NewClient(os.Getenv("apiKeyId"), os.Getenv("apiKey"), adminURL)
	if err != nil {
		log.Fatal(err)
	}
	adminClient := omg.AdminAPI{client}
	authToken, err := adminClient.Login(loginBody)
	if err != nil {
		log.Fatal(err)
	}

	env.AdminUser = omg.AdminAPI{
		Client: &omg.Client{
			BaseURL: adminURL,
			Auth: &omg.AdminUserAuth{
				ApiKeyId:      os.Getenv("apiKeyId"),
				ApiKey:        os.Getenv("apiKey"),
				UserId:        authToken.UserId,
				UserAuthToken: authToken.AuthenticationToken,
			},
			HttpClient: &http.Client{},
		},
	}
	accessKey, err := env.AdminUser.AccessKeyCreate()
	if err != nil {
		log.Fatal(err)
	}
	env.ServerUser = omg.EWalletAPI{
		Client: &omg.Client{
			BaseURL: ewalletURL,
			Auth: &omg.ServerAuth{
				AccessKey: accessKey.AccessKey,
				SecretKey: accessKey.SecretKey,
			},
			HttpClient: &http.Client{},
		},
	}

	// Mint tokens for the master account
	if *SeedLedger {
		body := omg.MintedTokenCreateParams{
			Name:          "Bitcoin",
			Symbol:        "BTC",
			Description:   "Base coin",
			SubunitToUnit: 1,
			Amount:        21000000,
		}
		mintedToken, err := env.AdminUser.MintedTokenCreate(body)
		if err != nil {
			log.Fatal(err)
		}
		baseTokenId = mintedToken.Id
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
