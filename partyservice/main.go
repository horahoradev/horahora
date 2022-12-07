package main

import (
	"log"

	"github.com/horahoradev/horahora/partyservice/internal/config"
	"github.com/horahoradev/horahora/partyservice/internal/protocol"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config fail: %v", err)
	}

	server, err := protocol.New(cfg.SqlClient, cfg.VideoClient)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Run(cfg.GRPCPort)
	if err != nil {
		log.Fatalf("grpc server exited with err: %v", err)
	}

	log.Print("Video service finished executing")
}
