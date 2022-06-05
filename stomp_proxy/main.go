// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"io"
	"net"
	"net/http"

	frame "github.com/go-stomp/stomp/v3/frame"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var addr = flag.String("addr", "0.0.0.0:15674", "http service address")

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true },
	ReadBufferSize:  8192,
	WriteBufferSize: 8192}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("upgrade: %v", err)
		return
	}

	rabbitConn, err := net.Dial("tcp", "rabbitmq:61613")
	if err != nil {
		log.Errorf("amqp dial: %v", err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// client <- rabbitmq
	go func(rabbitConn net.Conn, c *websocket.Conn) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				buf := make([]byte, 4096)
				_, err := rabbitConn.Read(buf)
				if err != nil {
					log.Errorf("read: %v", err)
					return
				}

				writer, err := c.NextWriter(1)
				if err != nil {
					log.Errorf("nextwriter: %v", err)
					return
				}

				_, err = writer.Write(buf)
				if err != nil {
					writer.Close()
					log.Errorf("write: %v", err)
					return
				}
				writer.Close()
			}
		}
	}(rabbitConn, c)

	// client -> rabbitmq
	for {
		_, wsReader, err := c.NextReader()
		if err != nil {
			log.Errorf("ws reader: %v", err)
			return
		}

		frameReader := frame.NewReader(wsReader)
		if err != nil {
			log.Errorf("new reader: %v", err)
			return
		}
		for {
			frameMsg, err := frameReader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Errorf("frame reader: %v", err)
				return
			}

			// What type of frame is this? If relevant, check user permissions (this will happen later)
			// nil is ok, just a heartbeat
			if frameMsg != nil {
				switch frameMsg.Command {
				case frame.SUBSCRIBE:
					f := frameMsg.Header
					log.Infof("Subscribing to %v", f)
					break
				case frame.CONNECT, frame.CONNECTED, frame.STOMP, frame.ACK, frame.NACK, frame.DISCONNECT:
					break
				case frame.BEGIN, frame.COMMIT, frame.ABORT, frame.SEND:
					log.Errorf("Client tried to use forbidden command %v", frameMsg.Command)
					return
				default:
					// anything else? no
					log.Errorf("Client tried to use unhandled command %v", frameMsg.Command)
					return
				}

			}

			var b bytes.Buffer
			foo := bufio.NewWriter(&b)
			frameWriter := frame.NewWriter(foo)
			frameWriter.Write(frameMsg)

			_, err = rabbitConn.Write(b.Bytes())
			if err != nil {
				log.Errorf("frame write: %v", err)
				return
			}
		}
	}

}

func main() {
	http.HandleFunc("/ws", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
