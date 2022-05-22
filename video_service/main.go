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
		log.Fatalf("Failed to initialize config. Err: %s", err)
	}

	err = grpcserver.NewGRPCServer(conf.BucketName, conf.SqlClient, conf.GRPCPort, conf.OriginFQDN, conf.Local,
		conf.RedisConn, conf.UserClient, conf.Tracer, conf.StorageBackend, conf.StorageAPIID, conf.StorageAPIKey,
		conf.ApprovalThreshold, conf.StorageEndpoint, conf.MaxDLFileSize)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Video service finished executing")
}
