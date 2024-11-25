package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

var levelNames = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"error": slog.LevelError,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
}

type Config struct {
	LogLvl      string `json:"logLevel"`
	GRPC        string `json:"grpcAddr"`
	REST        string `json:"restAddr"`
	UsePostgres bool   `json:"postgres"`
	PostgresDSN string `json:"postgresDsn"`

	LogLevel slog.Level
}

var (
	ErrInvalidConfig  = errors.New("error invalid config")
	ErrConfigNotFound = errors.New("error config not found")
)

// LoadConfig загружает и парсит конфигурацию из JSON файла.
func LoadConfig(path string) (*Config, error) {
	var cfg Config

	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if path == "config.json" {
				DefaultConfig(&cfg)

				return &cfg, nil
			}

			return nil, ErrConfigNotFound
		}

		return &cfg, fmt.Errorf("failed to open config file: %w", err)
	}

	if filepath.Ext(path) != ".json" {
		return nil, fmt.Errorf("unsupported file extension: %s, expected a .json file", filepath.Ext(path))
	}

	defer func() {
		if err = file.Close(); err != nil {
			slog.Error("failed to close config file", "error", err)
		}
	}()

	if err = json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	err = ValidateConfig(&cfg)

	cfg.LogLevel = levelNames[cfg.LogLvl]

	return &cfg, err
}

const (
	defaultPortREST = "8080"
	defaultPortGRPC = "9000"
	defaultLogLevel = "info"
)

func ValidateConfig(cfg *Config) error {
	if cfg.PostgresDSN == "" && cfg.UsePostgres {
		return fmt.Errorf("%w: postgres DSN is required", ErrInvalidConfig)
	}

	return nil
}

func DefaultConfig(cfg *Config) {
	if _, ok := levelNames[cfg.LogLvl]; !ok {
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
