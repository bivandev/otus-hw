package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Success(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.json")

	cfgData := map[string]interface{}{
		"logLevel": "debug",
	}

	file, err := os.Create(configFile)
	assert.NoError(t, err)

	defer file.Close() //nolint:errcheck

	err = json.NewEncoder(file).Encode(cfgData)
	assert.NoError(t, err)

	cfg, err := LoadConfig(configFile)
	assert.NoError(t, err)

	assert.Equal(t, "debug", cfg.LogLvl)
	assert.Equal(t, levelNames["debug"], cfg.LogLevel)
}

func TestLoadConfig_UnsupportedExtension(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "txt")

	_, err := LoadConfig(configFile)

	assert.Error(t, err)
	assert.Equal(t, err, ErrConfigNotFound)
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.json")

	file, err := os.Create(configFile)
	assert.NoError(t, err)

	defer file.Close() //nolint:errcheck

	_, err = file.WriteString(`{invalid json}`)
	assert.NoError(t, err)

	_, err = LoadConfig(configFile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode config")
}
