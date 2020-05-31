package config

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/go-redis/redis"
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
	RedisConn              *redis.Client
	GRPCPort               int    `env:"GRPCPort,required"`
	UserServiceGRPCAddress string `env:"UserServiceGRPCAddress,required"`
	BucketName             string `env:"BucketName,required"`
	Local                  bool   `env:"Local,required"` // If running locally, no s3 uploads
	// (this is a workaround for getting IAM permissions into pods running on minikube)
}

func New() (*config, error) {
	config := config{}
	err := env.Parse(&config.PostgresInfo)
	if err != nil {
		return nil, err
	}

	err = env.Parse(&config.RedisInfo)
	if err != nil {
		return nil, err
	}

	config.RedisConn = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisInfo.Hostname, config.RedisInfo.Port),
		Password: config.RedisInfo.Password, // no password set
		DB:       0,                         // use default DB
	})

	err = env.Parse(&config)
	return &config, err
}
