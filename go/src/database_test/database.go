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

type MarketRecords struct {
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

	// for key, _ := range m {
	// 	fmt.Println(key)
	// }
	// fmt.Println(m["_id"].ObjectId)

	// Get multiple marketrecords
	markets := make([]MarketRecords, 0)
	err = database.DBCon.DB("coinflow").C("market").Find(nil).Sort("-$natural").All(&markets)

	// Get one marketrecord
	// markets := MarketRecords{}
	// err = database.DBCon.DB("coinflow").C("market").Find(nil).Sort("-$natural").One(&markets)
	// err = database.DBCon.DB("coinflow").C("market").Find(nil).Limit(1).Sort("-$natural").One(&markets)

	fmt.Printf("size of collection: %d\n", len(markets))
	// for _, market := range markets {
	// 	fmt.Println(market.Id.Time())
	// 	fmt.Println(market.Market["BTC-LTC"])
	// }

	// Get records between certain time range
	seconds := 5
	toDate := bson.Now()
	toId := bson.NewObjectIdWithTime(toDate)
	fromDate := toDate.Add(-time.Duration(seconds) * time.Second)
	fromId := bson.NewObjectIdWithTime(fromDate)
	err = database.DBCon.DB("coinflow").C("market").Find(bson.M{"_id": bson.M{"$gte": fromId, "$lt": toId}}).Sort("-$natural").All(&markets)

	fmt.Printf("size of range: %d\n", len(markets))
	fmt.Printf("Now: %s\n", time.Now())
	var sumAsk, sumBid, sumLast float64
	for _, market := range markets {
		fmt.Println(market.Id.Time())
		sumAsk = sumAsk + market.Market["BTC-LTC"].Ask
		sumBid = sumBid + market.Market["BTC-LTC"].Bid
		sumLast = sumLast + market.Market["BTC-LTC"].Last
	}
	denominator := float64(len(markets))
	avgAsk, avgBid, avgLast := sumAsk/denominator, sumBid/denominator, sumLast/denominator
	fmt.Printf("The last %d seconds the average price-bid: %.8f, price-last: %.8f, price-ask: %.8f\n", seconds, avgBid, avgLast, avgAsk)

}
