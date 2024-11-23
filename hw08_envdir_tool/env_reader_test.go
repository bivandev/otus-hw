package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	tempDir := t.TempDir()

	createFile(t, filepath.Join(tempDir, "FOO"), "123")
	createFile(t, filepath.Join(tempDir, "BAR"), "value \nignored line")
	createFile(t, filepath.Join(tempDir, "EMPTY"), "")
	createFile(t, filepath.Join(tempDir, "NULL_BYTES"), "value\x00with\x00nulls\n")

	env, err := ReadDir(tempDir)
	require.NoError(t, err)

	require.Equal(t, Environment{
		"FOO":        {Value: "123", NeedRemove: false},
		"BAR":        {Value: "value", NeedRemove: false},
		"EMPTY":      {Value: "", NeedRemove: true},
		"NULL_BYTES": {Value: "value\nwith\nnulls", NeedRemove: false},
	}, env)

	createFile(t, filepath.Join(tempDir, "INVALID=NAME"), "invalid")
	_, err = ReadDir(tempDir)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid file name")

	emptyDir := t.TempDir()
	env, err = ReadDir(emptyDir)
	require.NoError(t, err)
	require.Empty(t, env)
}

func createFile(t *testing.T, path, content string) {
	t.Helper()

	err := os.WriteFile(path, []byte(content), 0o644)
	require.NoError(t, err)
}
