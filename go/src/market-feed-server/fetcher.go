package main

import (
	"encoding/json"
	// "fmt"
	"log"
	"net/http"
	// "net/url"
)

type MarketSummaryResponse struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	MarketSummary `json:"result"`
}

type MarketSummary []struct {
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

type Market struct {
	MarketName        string
	High              float64
	Low               float64
	Volume            float64
	Last              float64
	BaseVolume        float64
	TimeStamp         string
	Bid               float64
	Ask               float64
	OpenBuyOrders     int
	OpenSellOrders    int
	PrevDay           float64
	Created           string
	DisplayMarketName interface{}
}

/*
url: https://bittrex.com/api/v1.1/public/getmarketsummary?market=btc-ltc
response:
{   "success":true,
    "message":"",
    "result":[{
        "MarketName":"BTC-LTC",
        "High":0.02129998,
        "Low":0.01951000,
        "Volume":56074.33572962,
        "Last":0.02019990,
        "BaseVolume":1138.16360733,
        "TimeStamp":"2018-02-21T20:07:38.94",
        "Bid":0.02011701,
        "Ask":0.02019986,
        "OpenBuyOrders":2126,
        "OpenSellOrders":4863,
        "PrevDay":0.02064007,
        "Created":"2014-02-13T00:00:00"
    }]
}
*/

var market map[string]Market

// Return
func getMarketSummary() {
	url := "https://bittrex.com/api/v1.1/public/getmarketsummaries"

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	// Fill the record with the data from the JSON
	var record MarketSummaryResponse

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}

	market = make(map[string]Market)

	for _, i := range record.MarketSummary {
		market[i.MarketName] = Market{
			MarketName:        i.MarketName,
			High:              i.High,
			Low:               i.Low,
			Volume:            i.Volume,
			Last:              i.Last,
			BaseVolume:        i.BaseVolume,
			TimeStamp:         i.TimeStamp,
			Bid:               i.Bid,
			Ask:               i.Ask,
			OpenBuyOrders:     i.OpenBuyOrders,
			OpenSellOrders:    i.OpenSellOrders,
			PrevDay:           i.PrevDay,
			Created:           i.Created,
			DisplayMarketName: i.DisplayMarketName,
		}
	}
	return
}

// func main() {
// 	getMarketSummary()
// 	fmt.Println(market["BTC-LTC"])
// }
