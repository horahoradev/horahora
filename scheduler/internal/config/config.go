package config

import (
	"fmt"
	proto "github.com/horahoradev/horahora/video_service/protocol"
	"google.golang.org/grpc"
	"log"

	"github.com/caarlos0/env"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresInfo struct {
	Hostname string `env:"pgs_host,required"`
	Port     int    `env:"pgs_port,required"`
	Username string `env:"pgs_user"`
	Password string `env:"pgs_pass"`
	Db       string `env:"pgs_db,required"`
}

type config struct {
	PostgresInfo
	VideoOutputLoc          string
	VideoServiceGRPCAddress string `env:"VideoServiceGRPCAddress",required"`
	NumberOfRetries         int    `env:"NumberOfRetries",required"`
	Conn                    *sqlx.DB
	GRPCConn                *grpc.ClientConn
	Client                  proto.VideoServiceClient
}

func New() (*config, error) {
	config := config{}
	err := env.Parse(&config.PostgresInfo)
	if err != nil {
		return nil, err
	}

	err = env.Parse(&config)
	config.VideoOutputLoc = "./videos"

	// I'm putting this here because it makes it easier to do integration tests
	// https://www.calhoun.io/connecting-to-a-postgresql-database-with-gos-database-sql-package/
	config.Conn, err = sqlx.Connect("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", config.Hostname, config.Username, config.Password, config.Db))
	if err != nil {
		log.Fatalf("Could not connect to postgres. Err: %s", err)
	}

	config.GRPCConn, err = grpc.Dial(config.VideoServiceGRPCAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	config.Client = proto.NewVideoServiceClient(config.GRPCConn)

	return &config, err
}
