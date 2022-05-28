package main

import (
	"log"
	"net/http"

	"io/ioutil"
	"net/http"

	"github.com/horahoradev/horahora/ws_kafka_proxy/config"

	"golang.org/x/net/websocket"
)

// Echo the data received on the WebSocket.
func EchoServer(ws *websocket.Conn) {
	buf, err := ioutil.ReadAll(ws)
	if err != nil {
		ws.
	}
}

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Could not initialize config. Err: %s", err)
	}

	http.Handle("/", websocket.Handler(EchoServer))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
