package main

import (
	"fmt"
	// "log"
	"net/http"

	// "golang.org/x/net/websocket"
	"encoding/json"
	"github.com/gorilla/websocket"
	"time"
)

func echoHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		time.Sleep(time.Second)
		market, err = getMarketSummary()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Got market summary")

		myJson, err := json.Marshal(market)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = conn.WriteMessage(websocket.TextMessage, myJson)
		if err != nil {
			fmt.Println(err)
			break
		}
	}

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {

	http.HandleFunc("/echo", echoHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
