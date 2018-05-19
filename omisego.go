package main

import (
	"bitbucket.org/onyxium/onyxium-strategy-worker/models"
	omg "bitbucket.org/onyxium/onyxium-strategy-worker/omisego"
	"flag"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
)

var (
	AdminUser   omg.AdminAPI
	ServerUser  omg.EWalletAPI
	baseTokenId string
	SeedLedger  = flag.Bool("seed", false, "Seed the database with all available minted tokens for users.")
)

func initOmisego() error {
	err := authenticateClient()
	if err != nil {
		return err
	}
	if *SeedLedger {
		err = seedMintedTokens()
		if err != nil {
			return err
		}
	}
	if baseTokenId == "" {
		baseTokenId = os.Getenv("baseTokenId")
	}
	return nil
}

func authenticateClient() error {
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
		return err
	}
	adminClient := omg.AdminAPI{client}
	authToken, err := adminClient.Login(loginBody)
	if err != nil {
		return err
	}

	AdminUser = omg.AdminAPI{
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
	accessKey, err := AdminUser.AccessKeyCreate()
	if err != nil {
		return err
	}
	ServerUser = omg.EWalletAPI{
		Client: &omg.Client{
			BaseURL: ewalletURL,
			Auth: &omg.ServerAuth{
				AccessKey: accessKey.AccessKey,
				SecretKey: accessKey.SecretKey,
			},
			HttpClient: &http.Client{},
		},
	}
	return nil
}

func seedMintedTokens() error {
	// Mint tokens for the master account
	body := omg.MintedTokenCreateParams{
		Name:          "Ethereum",
		Symbol:        "ETH",
		Description:   "Base coin",
		SubunitToUnit: 1,
		Amount:        21000000,
	}
	mintedToken, err := AdminUser.MintedTokenCreate(body)
	if err != nil {
		return err
	}
	baseTokenId = mintedToken.Id
	log.Info(baseTokenId)
	return nil
}

func newUser(user *models.User) error {
	userCreateBody := omg.UserParams{
		ProviderUserId: user.Id.Hex(),
		Username:       user.Email,
	}
	_, err := ServerUser.UserCreate(userCreateBody)
	if err != nil {
		return err
	}
	creditBalanceBody := omg.BalanceAdjustmentParams{
		ProviderUserId: user.Id.Hex(),
		TokenId:        baseTokenId,
		Amount:         100,
	}

	_, err = ServerUser.UserCreditBalance(creditBalanceBody)
	return err
}
