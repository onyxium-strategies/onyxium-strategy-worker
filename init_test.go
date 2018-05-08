package main

import (
	"bitbucket.org/visa-startups/coinflow-strategy-worker/models"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
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

func TestMain(m *testing.M) {
	// log.SetLevel(log.DebugLevel)
	log.SetOutput(ioutil.Discard)
	env.DataStore = FakeDataStore{}
	os.Exit(m.Run())
}
