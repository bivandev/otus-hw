package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage: go-envdir <envdir> <command> [<args>...]")
		os.Exit(1)
	}

	envDir := os.Args[1]
	cmd := os.Args[2:]

	env, err := ReadDir(envDir)
	if err != nil {
		fmt.Printf("error reading envdir: %v\n", err)
		os.Exit(1)
	}

	exitCode, err := RunCmd(cmd, env)
	if err != nil {
		fmt.Printf("error running command: %v\n", err)
		os.Exit(1)
	}

	os.Exit(exitCode)
}
