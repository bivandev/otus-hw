package main

import (
	"fmt"

	"github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/config"
)

type Config struct {
	LogLvl      string `json:"logLevel"`
	RabbitMQURL string `json:"rabbitMqUrl"`
}

func ValidateConfig(cfg *Config) error {
	if cfg.RabbitMQURL == "" {
		return fmt.Errorf("%w: rabbitmq URL is required", config.ErrInvalidConfig)
	}

	return nil
}
