package main

import (
	"fmt"

	"github.com/horahoradev/horahora/video_service/internal/config"
	"github.com/horahoradev/horahora/video_service/internal/grpcserver"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func main() {
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", conf.Hostname, conf.Username, conf.Password, conf.Db))
	if err != nil {
		log.Fatalf("Could not connect to postgres. Err: %s", err)
	}

	err = grpcserver.NewGRPCServer(conf.BucketName, conn, conf.GRPCPort, conf.UserServiceGRPCAddress, conf.Local)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Video service finished executing")
}
