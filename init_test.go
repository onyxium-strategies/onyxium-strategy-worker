package main

import (
	"bitbucket.org/visa-startups/coinflow-strategy-worker/models"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
	// "io/ioutil"
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

func (FakeDataStore) UserCreate(user *models.User) (*models.User, error) {
	pwd, err := models.HashAndSalt(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = pwd
	return user, nil
}

func (FakeDataStore) StrategyCreate(name string, jsonTree string, bsonTree *models.Tree) (models.Strategy, error) {
	strategy := models.Strategy{
		Id:       bson.NewObjectId(),
		Name:     name,
		JsonTree: jsonTree,
		BsonTree: bsonTree,
		Status:   "paused",
		State:    bsonTree.Id,
	}
	return strategy, nil
}

// func (FakeDataStore) GetPausedStrategies() ([]models.Strategy, error) {
// 	var strategies []models.Strategy
// 	return strategies, nil
// }

func TestMain(m *testing.M) {
	log.SetLevel(log.DebugLevel)
	// log.SetOutput(ioutil.Discard)
	env.DataStore = FakeDataStore{}
	os.Exit(m.Run())
}
