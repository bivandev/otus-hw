package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("successfully run command with environment", func(t *testing.T) {
		env := Environment{
			"TEST_ENV_VAR": {Value: "test_value", NeedRemove: false},
		}
		cmd := []string{"env"}

		exitCode, err := RunCmd(cmd, env)
		require.NoError(t, err)
		require.Equal(t, 0, exitCode)
	})

	t.Run("command not provided", func(t *testing.T) {
		env := Environment{}
		cmd := []string{}

		exitCode, err := RunCmd(cmd, env)
		require.Error(t, err)
		require.Equal(t, 1, exitCode)
	})

	t.Run("remove environment variable", func(t *testing.T) {
		os.Setenv("REMOVE_ME", "should_be_removed")
		env := Environment{
			"REMOVE_ME": {Value: "", NeedRemove: true},
		}
		cmd := []string{"env"}

		exitCode, err := RunCmd(cmd, env)
		require.NoError(t, err)
		require.Equal(t, 0, exitCode)

		_, exists := os.LookupEnv("REMOVE_ME")
		require.False(t, exists)
	})

	t.Run("command fails with error", func(t *testing.T) {
		env := Environment{}
		cmd := []string{"false"}

		exitCode, err := RunCmd(cmd, env)
		require.NoError(t, err)
		require.Equal(t, 1, exitCode)
	})
}
