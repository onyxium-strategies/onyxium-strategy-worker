package database

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var (
	DBCon *mgo.Session
)

type Market struct {
	MarketName        string      `json:"MarketName" bson:"MarketName"`
	High              float64     `json:"High" bson:"High"`
	Low               float64     `json:"Low" bson:"Low"`
	Volume            float64     `json:"Volume" bson:"Volume"`
	Last              float64     `json:"Last" bson:"Last"`
	BaseVolume        float64     `json:"BaseVolume" bson:"BaseVolume"`
	TimeStamp         string      `json:"TimeStamp" bson:"TimeStamp"`
	Bid               float64     `json:"Bid" bson:"Bid"`
	Ask               float64     `json:"Ask" bson:"Ask"`
	OpenBuyOrders     int         `json:"OpenBuyOrders" bson:"OpenBuyOrders"`
	OpenSellOrders    int         `json:"OpenSellOrders" bson:"OpenSellOrders"`
	PrevDay           float64     `json:"PrevDay" bson:"PrevDay"`
	Created           string      `json:"Created" bson:"Created"`
	DisplayMarketName interface{} `json:"DisplayMarketName" bson:"DisplayMarketName"`
}

type MarketRecord struct {
	Id     bson.ObjectId     `json:"id" bson:"_id,omitempty"`
	Market map[string]Market `bson:",inline"`
}

func getLatestRecord() MarketRecord {
	record := MarketRecord{}
	err := DBCon.DB("coinflow").C("market").Find(nil).Sort("-$natural").One(&record)
	if err != nil {
		panic(err)
	}
	return record
}

func getHistoryRecord(TimeframeInMS int) MarketRecord {
	toDate := bson.Now()
	toId := bson.NewObjectIdWithTime(toDate)
	fromDate := toDate.Add(-time.Duration(TimeframeInMS) * time.Millisecond)
	fromId := bson.NewObjectIdWithTime(fromDate)

	// Get the record that is x millisecond old
	fromRecord := MarketRecord{}
	err := DBCon.DB("coinflow").C("market").Find(bson.M{"_id": bson.M{"$gte": fromId, "$lt": toId}}).Sort("$natural").One(&fromRecord)
	if err != nil {
		panic(err)
	}
	return fromRecord
}
