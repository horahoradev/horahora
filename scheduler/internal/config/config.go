package config

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	proto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"time"
)

type PostgresInfo struct {
	Hostname string `env:"pgs_host,required"`
	Port     int    `env:"pgs_port,required"`
	Username string `env:"pgs_user"`
	Password string `env:"pgs_pass"`
	Db       string `env:"pgs_db,required"`
}

type RedisInfo struct {
	Hostname string `env:"redis_host,required"`
	Port     int    `env:"redis_port,required"`
	Password string `env:"redis_pass,required"`
}

type config struct {
	PostgresInfo
	RedisInfo
	RedisPool               *redis.Pool
	Redlock                 *redsync.Redsync
	VideoOutputLoc          string
	VideoServiceGRPCAddress string `env:"VideoServiceGRPCAddress,required"`
	NumberOfRetries         int    `env:"NumberOfRetries,required"`
	Conn                    *sqlx.DB
	GRPCConn                *grpc.ClientConn
	Client                  proto.VideoServiceClient
	SocksConnStr            string `env:"SocksConn,required"`
}

func New() (*config, error) {
	config := config{}
	err := env.Parse(&config.PostgresInfo)
	if err != nil {
		return nil, err
	}

	err = env.Parse(&config)
	if err != nil {
		return nil, err
	}
	config.VideoOutputLoc = "./videos"

	// I'm putting this here because it makes it easier to do integration tests
	// https://www.calhoun.io/connecting-to-a-postgresql-database-with-gos-database-sql-package/
	config.Conn, err = sqlx.Connect("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		config.PostgresInfo.Hostname, config.PostgresInfo.Username, config.PostgresInfo.Password, config.PostgresInfo.Db))
	if err != nil {
		log.Fatalf("Could not connect to postgres. Err: %s", err)
	}

	config.Conn.SetMaxOpenConns(50)

	err = env.Parse(&config.RedisInfo)
	if err != nil {
		return nil, err
	}

	config.RedisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%d", config.RedisInfo.Hostname, config.RedisInfo.Port),
				redis.DialPassword(config.RedisInfo.Password))
		},
		TestOnBorrow:    nil,
		MaxIdle:         10,
		MaxActive:       10,
		IdleTimeout:     2 * time.Minute,
		Wait:            false,
		MaxConnLifetime: 0,
	}

	config.Redlock = redsync.New([]redsync.Pool{config.RedisPool})

	opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(100 * time.Millisecond)),
		grpc_retry.WithMax(5),
	}

	config.GRPCConn, err = grpc.Dial(config.VideoServiceGRPCAddress, grpc.WithInsecure(),
		//grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)))
	if err != nil {
		log.Fatal(err)
	}

	config.Client = proto.NewVideoServiceClient(config.GRPCConn)

	return &config, err
}
