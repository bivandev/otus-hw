package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

var LogLevelNames = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"error": slog.LevelError,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
}

var (
	ErrInvalidConfig  = errors.New("error invalid config")
	ErrConfigNotFound = errors.New("error config not found")
	ErrUnsupportedExt = errors.New("unsupported file extension")
)

// LoadConfig загружает и парсит конфигурацию из JSON файла.
func LoadConfig[T any](path string, cfg *T) error {
	if filepath.Ext(path) != ".json" {
		return fmt.Errorf("%w: %s", ErrUnsupportedExt, filepath.Ext(path))
	}

	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrConfigNotFound
		}

		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			slog.Error("failed to close config file", "error", cerr)
		}
	}()

	if err = json.NewDecoder(file).Decode(&cfg); err != nil {
		return fmt.Errorf("%w: failed to decode config: %w", ErrInvalidConfig, err)
	}

	return nil
}
