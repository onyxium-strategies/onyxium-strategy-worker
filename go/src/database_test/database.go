package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"worker-queue/database"
)

type Market struct {
	MarketName        string      `json:"MarketName"`
	High              float64     `json:"High"`
	Low               float64     `json:"Low"`
	Volume            float64     `json:"Volume"`
	Last              float64     `json:"Last"`
	BaseVolume        float64     `json:"BaseVolume"`
	TimeStamp         string      `json:"TimeStamp"`
	Bid               float64     `json:"Bid"`
	Ask               float64     `json:"Ask"`
	OpenBuyOrders     int         `json:"OpenBuyOrders"`
	OpenSellOrders    int         `json:"OpenSellOrders"`
	PrevDay           float64     `json:"PrevDay"`
	Created           string      `json:"Created"`
	DisplayMarketName interface{} `json:"DisplayMarketName"`
}

type MarketRecord struct {
	Id     bson.ObjectId     `json:"id" bson:"_id,omitempty"`
	Market map[string]Market `bson:",inline"`
}

// dit file dient om te oefenen met database queries

// handige links:
// https://stackoverflow.com/questions/31502195/mgo-query-objectid-for-range-of-time-values
func main() {
	var err error
	database.DBCon, err = mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer database.DBCon.Close()
	database.DBCon.SetMode(mgo.Monotonic, true)

	// bson.M kan je gebruiken om gwn een interface te krijgen
	// alleen dat moet je elke waarde type casten om er iets mee te doen
	// var m bson.M
	// err = database.DBCon.DB("coinflow").C("market").Find(nil).One(&m)

	// Get all marketrecord
	// markets := make([]MarketRecord, 0)
	// err = database.DBCon.DB("coinflow").C("market").Find(nil).Sort("-$natural").All(&markets)

	// Get latest marketrecord
	market := MarketRecord{}
	err = database.DBCon.DB("coinflow").C("market").Find(nil).Sort("-$natural").One(&market)

	// Get records between certain time range
	seconds := 10000
	toDate := bson.Now()
	toId := bson.NewObjectIdWithTime(toDate)
	fromDate := toDate.Add(-time.Duration(seconds) * time.Second)
	fromId := bson.NewObjectIdWithTime(fromDate)

	// Get the latest record
	toRecord := MarketRecord{}
	err = database.DBCon.DB("coinflow").C("market").Find(bson.M{"_id": bson.M{"$gte": fromId, "$lt": toId}}).Sort("-$natural").One(&toRecord)

	// Get the record that is x seconds old
	fromRecord := MarketRecord{}
	err = database.DBCon.DB("coinflow").C("market").Find(bson.M{"_id": bson.M{"$gte": fromId, "$lt": toId}}).Sort("$natural").One(&fromRecord)

	fmt.Printf("%s BTC-LTC record with last price: %.8f\n", toRecord.Id.Time(), toRecord.Market["BTC-LTC"].Last)
	fmt.Printf("%s BTC-LTC record with last price: %.8f\n", fromRecord.Id.Time(), fromRecord.Market["BTC-LTC"].Last)

	// Compute the percentage difference between the 2 points
	nieuw := toRecord.Market["BTC-LTC"].Last
	oud := fromRecord.Market["BTC-LTC"].Last
	fmt.Printf("De afgelopen %.2f seconden heeft er een procentuele verandering van %.3f%% plaatsgevonden.\n", toRecord.Id.Time().Sub(fromRecord.Id.Time()).Seconds(), (nieuw-oud)/oud*100)
}
