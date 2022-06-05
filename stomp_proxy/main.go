// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"io"
	"net/http"
	"time"

	frame "github.com/go-stomp/stomp/v3/frame"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var addr = flag.String("addr", "0.0.0.0:15674", "http service address")

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }} // use default options

var cstDialer = websocket.Dialer{
	ReadBufferSize:   4096,
	WriteBufferSize:  4096,
	HandshakeTimeout: 30 * time.Second,
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("upgrade: %v", err)
		return
	}

	rabbitConn, _, err := cstDialer.Dial("ws://rabbitmq:61614/ws", nil)
	if err != nil {
		log.Errorf("amqp dial: %v", err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// client <- rabbitmq
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				mt, msg, err := rabbitConn.ReadMessage()
				if err != nil {
					log.Errorf("reader: %v", err)
					return
				}

				writer, err := c.NextWriter(mt)
				if err != nil {
					log.Errorf("nextwriter: %v", err)
					return
				}

				_, err = writer.Write(msg)
				if err != nil {
					writer.Close()
					log.Errorf("write: %v", err)
					return
				}
				writer.Close()
			}
		}
	}()

	// client -> rabbitmq
	for {
		mt, wsReader, err := c.NextReader()
		if err != nil {
			log.Errorf("ws reader: %v", err)
			return
		}

		wsWriter, err := rabbitConn.NextWriter(mt)
		if err != nil {
			log.Errorf("ws writer: %v", err)
			return
		}

		for {
			frameReader := frame.NewReader(wsReader)
			if err != nil {
				log.Errorf("new reader: %v", err)
				wsWriter.Close()
				return
			}

			frameMsg, err := frameReader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Errorf("frame reader: %v", err)
				wsWriter.Close()
				return
			}

			// What type of frame is this? If relevant, check user permissions (this will happen later)
			// nil is ok, just a heartbeat
			if frameMsg != nil {
				switch frameMsg.Command {
				case frame.CONNECT, frame.CONNECTED, frame.STOMP, frame.SUBSCRIBE, frame.ACK, frame.NACK, frame.DISCONNECT:
					break
				case frame.BEGIN, frame.COMMIT, frame.ABORT, frame.SEND:
					wsWriter.Close()
					log.Errorf("Client tried to use forbidden command %v", frameMsg.Command)
					return
				default:
					// anything else? no
					wsWriter.Close()
					log.Errorf("Client tried to use unhandled command %v", frameMsg.Command)
					return
				}
			}

			frameWriter := frame.NewWriter(wsWriter)
			err = frameWriter.Write(frameMsg)
			if err != nil {
				wsWriter.Close()
				log.Errorf("frame write: %v", err)
				return
			}
		}
		wsWriter.Close()
	}

}

func main() {
	http.HandleFunc("/ws", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
