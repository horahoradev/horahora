package main

import (
	"github.com/horahoradev/horahora/video_service/internal/config"
	"github.com/horahoradev/horahora/video_service/internal/grpcserver"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func main() {
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	err = grpcserver.NewGRPCServer(conf.BucketName, conf.SqlClient, conf.GRPCPort, conf.UserServiceGRPCAddress, conf.Local, conf.RedisConn, conf.UserClient)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Video service finished executing")
}
