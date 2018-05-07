package main

import (
	"bitbucket.org/visa-startups/coinflow-strategy-worker/models"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
)

type FakeDataStore struct{}

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
			Volume: 7000,
			Last:   0.065,
			Bid:    0.063,
			Ask:    0.067,
		}
	} else { // TimeframeInMs == 0
		record["BTC-ETH"] = models.Market{
			Volume: 12000,
			Last:   0.075,
			Bid:    0.073,
			Ask:    0.077,
		}
	}
	return record, nil
}

func TestMain(m *testing.M) {
	log.SetLevel(log.DebugLevel)
	env.DataStore = FakeDataStore{}
	os.Exit(m.Run())
}

// This func just serves to play around with the mock db
func TestDBB(t *testing.T) {
	m, _ := env.DataStore.GetLatestMarket()
	t.Log(m["BTC-ETH"])
	m, _ = env.DataStore.GetHistoryMarket(0)
	t.Log(m["BTC-ETH"])
	m, _ = env.DataStore.GetHistoryMarket(1)
	t.Log(m["BTC-ETH"])
}
