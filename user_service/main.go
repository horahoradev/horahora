package main

import (
	"log"

	"github.com/horahoradev/horahora/user_service/internal/auth"

	"github.com/horahoradev/horahora/user_service/internal/grpcserver"

	"github.com/horahoradev/horahora/user_service/internal/config"
	_ "github.com/lib/pq"
)

func main() {
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := auth.ParsePrivateKey(conf.RSAKeypair)
	if err != nil {
		log.Fatalf("Could not parse RSA keypair. Err: %s", err)
	}

	err = grpcserver.NewGRPCServer(conf.DbConn, privateKey, conf.GRPCPort)
	if err != nil {
		log.Fatalf("gRPC server terminated with error: %s", err)
	}

	log.Print("User service finished executing")
}
