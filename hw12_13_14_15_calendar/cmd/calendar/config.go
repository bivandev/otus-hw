package main

import (
	"fmt"

	"github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/config"
)

type Config struct {
	LogLvl      string `json:"logLevel"`
	GRPC        string `json:"grpcAddr"`
	REST        string `json:"restAddr"`
	UsePostgres bool   `json:"postgres"`
	PostgresDSN string `json:"postgresDsn"`
}

const (
	defaultPortREST = "8080"
	defaultPortGRPC = "9000"
	defaultLogLevel = "info"
)

func ValidateConfig(cfg *Config) error {
	if cfg.PostgresDSN == "" && cfg.UsePostgres {
		return fmt.Errorf("%w: postgres DSN is required", config.ErrInvalidConfig)
	}

	return nil
}

func DefaultConfig(cfg *Config) {
	if _, ok := config.LogLevelNames[cfg.LogLvl]; !ok {
		cfg.LogLvl = defaultLogLevel
	}

	if cfg.REST == "" {
		cfg.REST = defaultPortREST
	}

	if cfg.GRPC == "" {
		cfg.GRPC = defaultPortGRPC
	}

	if cfg.PostgresDSN == "" && cfg.UsePostgres {
		cfg.UsePostgres = false
	}
}
