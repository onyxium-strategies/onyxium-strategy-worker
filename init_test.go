package main

import (
	"bitbucket.org/onyxium/onyxium-strategy-worker/models"
	"fmt"
	omg "github.com/Alainy/OmiseGo-Go-SDK"
	"github.com/icrowley/fake"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

type FakeDataStore struct {
	models.DataStore // This way we only have to implement the methods we want to test.
}

func (FakeDataStore) GetLatestMarket() (map[string]models.Market, error) {
	record := map[string]models.Market{
		"BTC-ETH": models.Market{
			Volume: 10000,
			Last:   0.07,
			Bid:    0.068,
			Ask:    0.072,
		},
	}
	return record, nil
}

func (FakeDataStore) GetHistoryMarket(TimeframeInMS int) (map[string]models.Market, error) {
	// Since we want to test percentage-increase and percentage-decrease we need a record
	// that is lower and higher than GetLatestMarket record
	// use TimeframeInMS == 1 for increase and TimeframeInMS == 0 for decrease
	record := make(map[string]models.Market)
	if TimeframeInMS == 1 {
		record["BTC-ETH"] = models.Market{
			Volume: 7000,  //+42.85%
			Last:   0.065, //+7.69%
			Bid:    0.063, //+7.93%
			Ask:    0.067, //+7.46%
		}
	} else { // TimeframeInMs == 0
		record["BTC-ETH"] = models.Market{
			Volume: 12000, //-16.66%
			Last:   0.075, //-6.66%
			Bid:    0.073, //-6.84%
			Ask:    0.077, //-6.49%
		}
	}
	return record, nil
}

func (FakeDataStore) UserActivate(id string, token string) error {
	ok := bson.IsObjectIdHex(id)
	if !ok {
		return fmt.Errorf("Incorrect id hex received: %s", id)
	}
	user := models.User{
		Email: "example@gmail.com",
	}
	if ok, err := models.ComparePasswords(token, []byte(user.Email)); ok && err == nil {
		user.IsActivated = true
		user.ActivatedAt = time.Now()
		return nil
	} else {
		return err
	}
}

func (FakeDataStore) UserAll() ([]models.User, error) {
	users := []models.User{
		{Email: "test@gmail.com", Password: "pwd"},
		{Email: "test2@gmail.com", Password: "pwd2"},
	}
	return users, nil
}

func (FakeDataStore) UserCreate(user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Id = bson.NewObjectId()
	pwd, err := models.HashAndSalt(user.Password)
	if err != nil {
		return err
	}
	user.Password = pwd
	return nil
}

func (FakeDataStore) StrategyCreate(strategy *models.Strategy) error {
	return nil
}

func (FakeDataStore) StrategiesGetPaused() ([]models.Strategy, error) {
	var strategies []models.Strategy
	return strategies, nil
}

func (FakeDataStore) StrategyUpdate(strategy *models.Strategy) error {
	return nil
}

type FakeLedger struct {
	omg.EwalletAdminAPI
}

func (FakeLedger) UserCreate(reqBody omg.UserParams) (*omg.User, error) {
	log.Info("hoi")
	user := omg.User{
		Object:         "user",
		Id:             fake.CharactersN(8),
		Username:       fake.UserName(),
		ProviderUserId: fake.CharactersN(8),
		Email:          fake.EmailAddress(),
		CreatedAt:      time.Now().String(),
		UpdatedAt:      time.Now().String(),
	}
	return &user, nil
}

func (FakeLedger) UserGetWalletsByProviderUserId(reqBody omg.ProviderUserIdParam) (*omg.WalletList, error) {
	walletList := omg.WalletList{
		Object: "walletlist",
		Data: []omg.Wallet{{
			Object:     "wallet",
			Address:    fake.CharactersN(16),
			Name:       "primary",
			Identifier: "primary",
		}},
	}
	return &walletList, nil
}

func (FakeLedger) TokenAll(reqBody omg.ListParams) (*omg.TokenList, error) {
	tokenList := omg.TokenList{
		Object: "tokenlist",
		Data: []omg.Token{{
			Object:        "token",
			Id:            fake.CharactersN(16),
			Symbol:        fake.CharactersN(3),
			SubunitToUnit: 100000000,
		}},
	}
	return &tokenList, nil
}

func (FakeLedger) TokenCreate(reqBody omg.TokenCreateParams) (*omg.Token, error) {
	token := omg.Token{
		Object:        "token",
		Id:            fake.CharactersN(16),
		Symbol:        fake.CharactersN(3),
		SubunitToUnit: 100000000,
	}
	return &token, nil
}

func (FakeLedger) AccountGetWallets(reqBody omg.ListByIdParams) (*omg.WalletList, error) {
	walletList := omg.WalletList{
		Object: "walletlist",
		Data: []omg.Wallet{{
			Object:     "wallet",
			Address:    fake.CharactersN(16),
			Name:       "primary",
			Identifier: "primary",
		}},
	}
	return &walletList, nil
}

func (FakeLedger) TransactionCreate(reqBody omg.TransactionCreateParams) (*omg.Transaction, error) {
	transaction := omg.Transaction{
		Object: "transaction",
		Id:     fake.CharactersN(16),
		From: omg.TransactionSource{
			Object:  "transactionsource",
			Address: reqBody.FromAddress,
			Amount:  float64(reqBody.Amount),
		},
		To: omg.TransactionSource{
			Object:  "transactionsource",
			Address: reqBody.ToAddress,
			Amount:  float64(reqBody.Amount),
		},
	}
	return &transaction, nil
}

func TestMain(m *testing.M) {
	// log.SetLevel(log.DebugLevel)
	log.SetOutput(ioutil.Discard)
	env.DataStore = FakeDataStore{}
	env.Ledger = FakeLedger{}
	os.Exit(m.Run())
}
