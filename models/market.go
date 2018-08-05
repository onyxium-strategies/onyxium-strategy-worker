package models

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func (db *MGO) GetLatestMarket() (map[string]Market, error) {
	record := &MarketRecord{}
	err := db.DB("onyxium").C("market").Find(nil).Sort("-$natural").One(record)
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
	err := db.DB("onyxium").C("market").Find(bson.M{"_id": bson.M{"$gte": fromId, "$lt": toId}}).Sort("$natural").One(fromRecord)
	if err != nil {
		return nil, fmt.Errorf("Failed to get history market record with TimeframeInMS %d.", TimeframeInMS)
	}
	return fromRecord.Market, nil
}
