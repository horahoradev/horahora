package config

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
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
	RSAKeypair string `env:"RSA_KEYPAIR,required"`
	GRPCPort   int64  `env:"GRPCPort,required"`
	DbConn     *sqlx.DB
}

func New() (*config, error) {
	config := config{}
	err := env.Parse(&config.PostgresInfo)
	if err != nil {
		return nil, err
	}

	err = env.Parse(&config)

	// https://www.calhoun.io/connecting-to-a-postgresql-database-with-gos-database-sql-package/
	conn, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable connect_timeout=180", config.Hostname, config.Username, config.Password, config.Db))
	if err != nil {
		return nil, fmt.Errorf("could not connect to postgres. Err: %s", err)
	}

	log.Infof("RSAKeypair: %s", config.RSAKeypair)

	config.DbConn = conn
	return &config, err
}
