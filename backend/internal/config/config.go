package config

import (
	"github.com/hamillka/team25/backend/internal/db"
	"github.com/hamillka/team25/backend/internal/logger"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DB      db.DatabaseConfig `envconfig:"DB"`
	Port    string            `envconfig:"PORT"`
	Timeout int64             `envconfig:"TIMEOUT"`
	Log     logger.LogConfig  `envconfig:"LOG"`
}

func New() (*Config, error) {
	var config Config

	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
