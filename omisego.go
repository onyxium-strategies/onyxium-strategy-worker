package main

import (
	"bitbucket.org/onyxium/onyxium-strategy-worker/models"
	// "flag"
	"fmt"
	omg "github.com/Alainy/OmiseGo-Go-SDK"
	"github.com/joho/godotenv"
	// log "github.com/sirupsen/logrus"
	"net/url"
	"os"
)

var (
	baseTokenId string
	// SeedLedger  = flag.Bool("seed", false, "Seed the database with all available  tokens for users.")
	currencies = map[string]string{
		"BTC": "Bitcoin",
		"ETH": "Ethereum",
		"NEO": "NEO",
		"OMG": "Omisego",
		"LTC": "Litecoin",
	}
)

const subUnitToUnit = 100000000

func initOmisego() error {
	err := seedTokens()
	if err != nil {
		return err
	}

	err = getBaseTokenId()
	if err != nil {
		return err
	}

	err = getPrimaryWalletAddress()
	return err
}

func initOMGClient() (*omg.AdminAPI, error) {
	// Get authentication and connection to the eWallet Admin API.
	adminURL := &url.URL{
		Scheme: "http",
		Host:   "localhost:4000",
		Path:   "/api/admin",
	}

	client, err := omg.NewClient(os.Getenv("accessKey"), os.Getenv("secretKey"), adminURL)
	if err != nil {
		return nil, err
	}
	OMGProvider := omg.AdminAPI{
		Client: client,
	}
	return &OMGProvider, nil
}

func seedTokens() error {
	tokenList, err := env.Ledger.TokenAll(omg.ListParams{})
	if err != nil {
		return err
	}
	for symbol, name := range currencies {
		if !symbolInSlice(symbol, tokenList.Data) {
			// Mint tokens for the master account
			body := omg.TokenCreateParams{
				Name:          name,
				Symbol:        symbol,
				Description:   name,
				SubunitToUnit: subUnitToUnit,
				Amount:        2100000000000000,
			}
			_, err := env.Ledger.TokenCreate(body)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getBaseTokenId() error {
	if baseTokenId = os.Getenv("baseTokenId"); baseTokenId == "" {
		tokenList, err := env.Ledger.TokenAll(omg.ListParams{})
		if err != nil {
			return err
		}
		for _, token := range tokenList.Data {
			if token.Symbol == "BTC" {
				baseTokenId = token.Id
				myEnv, err := godotenv.Read()
				if err != nil {
					return err
				}
				myEnv["baseTokenId"] = token.Id
				err = godotenv.Write(myEnv, "./.env")
				return err
			}
		}
	}
	return nil
}

func getPrimaryWalletAddress() error {
	body := omg.ListByIdParams{
		os.Getenv("accountId"),
		omg.ListParams{},
	}
	walletList, err := env.Ledger.AccountGetWallets(body)
	if err != nil {
		return err
	}
	for _, wallet := range walletList.Data {
		if wallet.Identifier == "primary" {
			myEnv, err := godotenv.Read()
			if err != nil {
				return err
			}
			myEnv["primaryWalletAddress"] = wallet.Address
			err = godotenv.Write(myEnv, "./.env")
			return err
		}
	}
	return nil
}

func NewOMGUser(user *models.User) error {
	// create user on ewallet
	userCreateBody := omg.UserParams{
		ProviderUserId: user.Id.Hex(),
		Username:       user.Email,
	}
	_, err := env.Ledger.UserCreate(userCreateBody)
	if err != nil {
		return err
	}

	// get primary wallet of user
	getWalletBody := omg.ListByProviderUserIdParams{
		user.Id.Hex(),
		omg.ListParams{},
	}
	walletList, err := env.Ledger.UserGetWalletsByProviderUserId(getWalletBody)
	if err != nil {
		return err
	}
	for _, wallet := range walletList.Data {
		if wallet.Identifier == "primary" {
			// send initial funds of 10 BTC
			transactionBody := omg.TransactionCreateParams{
				IdempotencyToken: omg.NewIdempotencyToken(),
				FromAddress:      os.Getenv("primaryWalletAddress"),
				ToAddress:        wallet.Address,
				TokenId:          os.Getenv("baseTokenId"),
				Amount:           10 * subUnitToUnit,
			}
			_, err = env.Ledger.TransactionCreate(transactionBody)
			return err
		}
	}
	return fmt.Errorf("Unable to send initial funds to new user.")
}

func symbolInSlice(a string, list []omg.Token) bool {
	for _, b := range list {
		if b.Symbol == a {
			return true
		}
	}
	return false
}
