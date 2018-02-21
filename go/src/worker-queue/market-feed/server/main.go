package main

import (
	"fmt"
	// "log"
	"net/http"

	// "golang.org/x/net/websocket"
	"github.com/gorilla/websocket"
	// "encoding/json"
	"time"
)

func echoHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		time.Sleep(2 * time.Second)
		err = conn.WriteMessage(websocket.TextMessage, "hoi")
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
