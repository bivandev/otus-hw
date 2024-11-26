package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Config struct {
	LogLvl string `json:"logLevel"`
}

func TestLoadConfig_Success(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.json")

	cfgData := Config{LogLvl: "debug"}

	file, err := os.Create(configFile)
	assert.NoError(t, err)

	defer file.Close()

	err = json.NewEncoder(file).Encode(&cfgData)
	assert.NoError(t, err)

	var cfg Config
	err = LoadConfig(configFile, &cfg)
	assert.NoError(t, err)
	assert.Equal(t, "debug", cfg.LogLvl)
}

func TestLoadConfig_UnsupportedExtension(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.txt")

	_, err := os.Create(configFile)
	assert.NoError(t, err)

	var cfg Config
	err = LoadConfig(configFile, &cfg)
	assert.ErrorIs(t, err, ErrUnsupportedExt)
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.json")

	file, err := os.Create(configFile)
	assert.NoError(t, err)

	defer file.Close()

	_, err = file.WriteString(`{invalid json}`)
	assert.NoError(t, err)

	var cfg Config
	err = LoadConfig(configFile, &cfg)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidConfig)
}
