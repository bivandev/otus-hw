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
	Port        int    `json:"port"`
	UsePostgres bool   `json:"postgres"`
	PostgresDSN string `json:"postgresDsn"`

	LogLevel slog.Level
}

var ErrInvalidConfig = errors.New("error invalid config")

// LoadConfig загружает и парсит конфигурацию из JSON файла.
func LoadConfig(path string) (*Config, error) {
	var cfg Config

	if filepath.Ext(path) != ".json" {
		return nil, fmt.Errorf("unsupported file extension: %s, expected a .json file", filepath.Ext(path))
	}

	file, err := os.Open(path)
	if err != nil {
		err = DefaultConfig(&cfg)

		return &cfg, fmt.Errorf("failed to open config file: %w", err)
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
	defaultPort     = 8080
	defaultLogLevel = "info"
)

func ValidateConfig(cfg *Config) error {
	if cfg.PostgresDSN == "" && cfg.UsePostgres {
		return fmt.Errorf("%w: postgres DSN is required", ErrInvalidConfig)
	}

	return nil
}

func DefaultConfig(cfg *Config) error {
	if _, ok := levelNames[cfg.LogLvl]; !ok {
		cfg.LogLvl = defaultLogLevel
	}

	if cfg.Port == 0 {
		cfg.Port = defaultPort
	}

	if cfg.PostgresDSN == "" && cfg.UsePostgres {
		cfg.UsePostgres = false
	}

	return nil
}
