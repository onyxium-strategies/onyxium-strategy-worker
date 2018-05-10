package models

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Market struct {
	MarketName        string      `json:"MarketName" bson:"marketName"`
	High              float64     `json:"High" bson:"high"`
	Low               float64     `json:"Low" bson:"low"`
	Volume            float64     `json:"Volume" bson:"volume"`
	Last              float64     `json:"Last" bson:"last"`
	BaseVolume        float64     `json:"BaseVolume" bson:"baseVolume"`
	TimeStamp         string      `json:"TimeStamp" bson:"timeStamp"`
	Bid               float64     `json:"Bid" bson:"bid"`
	Ask               float64     `json:"Ask" bson:"ask"`
	OpenBuyOrders     int         `json:"OpenBuyOrders" bson:"openBuyOrders"`
	OpenSellOrders    int         `json:"OpenSellOrders" bson:"openSellOrders"`
	PrevDay           float64     `json:"PrevDay" bson:"prevDay"`
	Created           string      `json:"Created" bson:"created"`
	DisplayMarketName interface{} `json:"DisplayMarketName" bson:"displayMarketName"`
}

type MarketRecord struct {
	Id     bson.ObjectId     `json:"id" bson:"_id,omitempty"`
	Market map[string]Market `bson:",inline"`
}

func (db *MGO) GetLatestMarket() (map[string]Market, error) {
	record := &MarketRecord{}
	err := db.DB("coinflow").C("market").Find(nil).Sort("-$natural").One(record)
	if err != nil {
		return nil, fmt.Errorf("Failed to get latest market record.")
	}
	return record.Market, nil
}

func (db *MGO) GetHistoryMarket(TimeframeInMS int) (map[string]Market, error) {
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
