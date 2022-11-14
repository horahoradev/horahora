package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

type PostgresInfo struct {
	Hostname string `env:"pgs_host,required"`
	Port     int    `env:"pgs_port,required"`
	Username string `env:"pgs_user"`
	Password string `env:"pgs_pass"`
	Db       string `env:"pgs_db,required"`
}

type config struct {
	GRPCPort int `env:"GRPCPort,required"`

	VideoserviceGRPCAddress string `env:"VideoServiceGRPCAddress,required"`
	VideoClient             videoproto.VideoServiceClient

	PostgresInfo
	SqlClient *sqlx.DB
}

func New() (*config, error) {
	config := config{}
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	err = env.Parse(&config.PostgresInfo)
	if err != nil {
		return nil, err
	}

	config.SqlClient, err = sqlx.Connect("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable connect_timeout=180", config.PostgresInfo.Hostname, config.PostgresInfo.Username, config.PostgresInfo.Password, config.PostgresInfo.Db))
	if err != nil {
		return nil, fmt.Errorf("Could not connect to postgres. Err: %s", err)
	}

	opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(100 * time.Millisecond)),
		grpc_retry.WithMax(5),
	}

	conn, err := grpc.Dial(config.VideoserviceGRPCAddress, grpc.WithInsecure(),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpc_retry.UnaryClientInterceptor(opts...))))
	if err != nil {
		return nil, err
	}

	config.VideoClient = videoproto.NewVideoServiceClient(conn)

	return &config, err
}
