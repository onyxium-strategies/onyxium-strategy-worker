package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	// Enable closing connection with ctrl-c
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	ticker := time.NewTicker(2 * time.Second).C
	done := make(chan bool)

	// Close the connection after a duration
	timer := time.NewTimer(20 * time.Second).C
	// go func() {
	// 	time.Sleep(10 * time.Second)
	// 	done <- true
	// }()

	for {
		select {
		case <-ticker:
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			var objmap map[string]*json.RawMessage
			err = json.Unmarshal(message, &objmap)
			fmt.Printf("recv: %s", objmap["BTC-LTC"])
			fmt.Println("\n")

		case <-timer:
			log.Println("timer ended")
			return
		case <-done:
			log.Println("job done")
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
