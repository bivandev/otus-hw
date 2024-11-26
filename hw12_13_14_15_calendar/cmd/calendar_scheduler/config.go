package main

import (
	"fmt"
	"time"

	"github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/config"
)

type Config struct {
	LogLvl      string        `json:"logLevel"`
	UsePostgres bool          `json:"postgres"`
	PostgresDSN string        `json:"postgresDsn"`
	RabbitMQURL string        `json:"rabbitMqUrl"`
	Interval    time.Duration `json:"interval"`
}

func ValidateConfig(cfg *Config) error {
	if cfg.PostgresDSN == "" && cfg.UsePostgres {
		return fmt.Errorf("%w: postgres DSN is required", config.ErrInvalidConfig)
	}

	if cfg.RabbitMQURL == "" {
		return fmt.Errorf("%w: rabbitmq URL is required", config.ErrInvalidConfig)
	}

	return nil
}
