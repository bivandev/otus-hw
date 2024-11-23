package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (int, error) {
	if len(cmd) == 0 {
		return 1, errors.New("no command provided")
	}

	for key, val := range env {
		if val.NeedRemove {
			if err := os.Unsetenv(key); err != nil {
				return 0, err
			}

			continue
		}

		if err := os.Setenv(key, val.Value); err != nil {
			return 0, err
		}
	}

	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	command.Env = os.Environ()
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode(), nil
		}

		return 1, err
	}

	return 0, nil
}
