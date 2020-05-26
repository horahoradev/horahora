package config

import (
	"github.com/caarlos0/env"
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

	err = env.Parse(&config)
	return &config, err
}
