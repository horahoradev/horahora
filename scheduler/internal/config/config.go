package config

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env"
	"github.com/go-redsync/redsync"
	stomp "github.com/go-stomp/stomp/v3"
	"github.com/gomodule/redigo/redis"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	proto "github.com/horahoradev/horahora/video_service/protocol"
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

type RedisInfo struct {
	Hostname string `env:"redis_host,required"`
	Port     int    `env:"redis_port,required"`
	Password string `env:"redis_pass,required"`
}

type RabbitmqInfo struct {
	Hostname string `env:"rabbit_host,required"`
	Port     int    `env:"rabbit_port,required"`
	Username string `env:"rabbit_user,required"`
	Password string `env:"rabbit_pass,required"`
}

type config struct {
	PostgresInfo
	RedisInfo
	RabbitmqInfo
	RabbitConn              *stomp.Conn
	RedisPool               *redis.Pool
	Redlock                 *redsync.Redsync
	VideoOutputLoc          string
	VideoServiceGRPCAddress string `env:"VideoServiceGRPCAddress,required"`
	NumberOfRetries         int    `env:"NumberOfRetries,required"`
	Conn                    *sqlx.DB
	GRPCConn                *grpc.ClientConn
	Client                  proto.VideoServiceClient
	SocksConnStr            string        `env:"SocksConn,required"`
	SyncPollDelay           time.Duration `env:"SyncPollDelay,required"`
	MaxFS                   uint64        `env:"MaxDLFileSize,required"`
	AcceptLanguage          string        `env:"AcceptLanguage"`
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
	config.VideoOutputLoc = "/tmp"

	// I'm putting this here because it makes it easier to do integration tests
	// https://www.calhoun.io/connecting-to-a-postgresql-database-with-gos-database-sql-package/
	config.Conn, err = sqlx.Connect("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		config.PostgresInfo.Hostname, config.PostgresInfo.Username, config.PostgresInfo.Password, config.PostgresInfo.Db))
	if err != nil {
		log.Fatalf("Could not connect to postgres. Err: %s", err)
	}

	config.Conn.SetMaxOpenConns(50)

	// Rabbitmq
	var options []func(*stomp.Conn) error = []func(*stomp.Conn) error{
		stomp.ConnOpt.Login(config.RabbitmqInfo.Username, config.RabbitmqInfo.Password),
		stomp.ConnOpt.Host("/"),
	}
	var serverAddr = flag.String("server", fmt.Sprintf("%s:%d", config.RabbitmqInfo.Hostname, config.RabbitmqInfo.Port),
		"STOMP server endpoint")

	conn, err := stomp.Dial("tcp", *serverAddr, options...)
	if err != nil {
		return nil, err
	}
	config.RabbitConn = conn

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
