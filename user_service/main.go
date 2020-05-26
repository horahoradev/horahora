package main

import (
	"fmt"
	"log"

	"git.horahora.org/otoman/user-service.git/internal/auth"

	"github.com/jmoiron/sqlx"

	"git.horahora.org/otoman/user-service.git/internal/grpcserver"

	"git.horahora.org/otoman/user-service.git/internal/config"
	_ "github.com/lib/pq"
)

func main() {
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	// https://www.calhoun.io/connecting-to-a-postgresql-database-with-gos-database-sql-package/
	conn, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", conf.Hostname, conf.Username, conf.Password, conf.Db))
	if err != nil {
		log.Fatalf("Could not connect to postgres. Err: %s", err)
	}

	privateKey, err := auth.ParsePrivateKey(conf.RSAKeypair)
	if err != nil {
		log.Fatalf("Could not parse RSA keypair. Err: %s", err)
	}

	err = grpcserver.NewGRPCServer(conn, privateKey, conf.GRPCPort)
	if err != nil {
		log.Fatalf("gRPC server terminated with error: %s", err)
	}

	log.Print("User service finished executing")
}
