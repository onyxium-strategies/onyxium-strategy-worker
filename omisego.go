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
	OMGProvider omg.AdminAPI
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

func initOmisego() error {
	err := authenticateClient()
	if err != nil {
		return err
	}

	err = seedTokens()
	if err != nil {
		return err
	}

	err = getBaseTokenId()
	if err != nil {
		return err
	}

	err = getPrimaryWalletAddress()
	if err != nil {
		return err
	}

	return nil
}

func authenticateClient() error {
	// Get authentication and connection to the eWallet Admin API.
	adminURL := &url.URL{
		Scheme: "http",
		Host:   "localhost:4000",
		Path:   "/api/admin",
	}

	client, err := omg.NewClient(os.Getenv("accessKey"), os.Getenv("secretKey"), adminURL)
	if err != nil {
		return err
	}
	OMGProvider = omg.AdminAPI{
		Client: client,
	}
	return nil
}

func seedTokens() error {
	tokenList, err := OMGProvider.TokenAll(omg.ListParams{})
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
				SubunitToUnit: 1,
				Amount:        21000000,
			}
			_, err := OMGProvider.TokenCreate(body)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getBaseTokenId() error {
	if baseTokenId = os.Getenv("baseTokenId"); baseTokenId == "" {
		tokenList, err := OMGProvider.TokenAll(omg.ListParams{})
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
	walletList, err := OMGProvider.AccountGetWallets(body)
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

func NewUser(user *models.User) error {
	// create user on ewallet
	userCreateBody := omg.UserParams{
		ProviderUserId: user.Id.Hex(),
		Username:       user.Email,
	}
	_, err := OMGProvider.UserCreate(userCreateBody)
	if err != nil {
		return err
	}

	// create primary wallet for user
	getWalletBody := omg.ProviderUserIdParam{
		ProviderUserId: user.Id.Hex(),
	}
	walletList, err := OMGProvider.UserGetWalletsByProviderUserId(getWalletBody)
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
				Amount:           10,
			}
			_, err = OMGProvider.TransactionCreate(transactionBody)
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
