package config

import (
	"github.com/caarlos0/env"
)

type Config struct {
	UserPGSHostname      string `env:"user_pgs_host,required"`
	VideoPGSHostname     string `env:"video_pgs_host,required"`
	SchedulerPGSHostname string `env:"scheduler_pgs_host,required"`

	UserPGSUsername      string `env:"user_pgs_username,required"`
	VideoPGSUsername     string `env:"video_pgs_username,required"`
	SchedulerPGSUsername string `env:"scheduler_pgs_username,required"`

	UserPGSPassword      string `env:"user_pgs_password,required"`
	VideoPGSPassword     string `env:"video_pgs_password,required"`
	SchedulerPGSPassword string `env:"scheduler_pgs_password,required"`

	UserPGSDatabase      string `env:"user_pgs_database,required"`
	VideoPGSDatabase     string `env:"video_pgs_database,required"`
	SchedulerPGSDatabase string `env:"scheduler_pgs_database,required"`

	BackblazeID     string `env:"StorageAPIID,required"`
	BackblazeAPIKey string `env:"StorageAPIKey,required"`
}

func New() (*Config, error) {
	config := Config{}
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	return &config, err
}
