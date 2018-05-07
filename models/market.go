package models

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
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

func (db *DB) GetLatestMarket() (map[string]Market, error) {
	record := &MarketRecord{}
	err := db.DB("coinflow").C("market").Find(nil).Sort("-$natural").One(record)
	if err != nil {
		return nil, fmt.Errorf("Failed to get latest market record.")
	}
	return record.Market, nil
}

func (db *DB) GetHistoryMarket(TimeframeInMS int) (map[string]Market, error) {
	toDate := bson.Now()
	toId := bson.NewObjectIdWithTime(toDate)
	fromDate := toDate.Add(-time.Duration(TimeframeInMS) * time.Millisecond)
	fromId := bson.NewObjectIdWithTime(fromDate)

	// Get the record that is x millisecond old
	fromRecord := &MarketRecord{}
	err := db.DB("coinflow").C("market").Find(bson.M{"_id": bson.M{"$gte": fromId, "$lt": toId}}).Sort("$natural").One(fromRecord)
	if err != nil {
		return nil, fmt.Errorf("Failed to get history market record with TimeframeInMS %d.", TimeframeInMS)
	}
	return fromRecord.Market, nil
}
