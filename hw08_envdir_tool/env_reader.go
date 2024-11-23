package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if strings.Contains(name, "=") {
			return nil, fmt.Errorf("invalid file name: %s", name)
		}

		filePath := filepath.Join(dir, name)
		f, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		var value string
		if scanner.Scan() {
			value = processEnvValue([]byte(scanner.Text()))
		}

		if err = scanner.Err(); err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
		}

		env[name] = EnvValue{
			Value:      value,
			NeedRemove: value == "",
		}
	}

	return env, nil
}

// processEnvValue trims whitespace, replaces null bytes with newlines, and cleans up the environment value.
func processEnvValue(data []byte) string {
	data = bytes.ReplaceAll(data, []byte("\x00"), []byte("\n"))
	return strings.TrimRight(string(data), " \t\r\n")
}
