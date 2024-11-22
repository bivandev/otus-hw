package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var address string
	var timeout time.Duration
	flag.StringVar(&address, "address", "localhost:8080", "Address to connect to (host:port)")
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "Connection timeout")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect: %v\n", err)
		return
	}
	defer client.Close()

	go func() {
		if err := client.Receive(); err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "Failed to receive data: %v\n", err)
		}
		cancel()
	}()

	go func() {
		if err := client.Send(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to send data: %v\n", err)
			cancel()
		}
	}()

	select {
	case <-sigCh:
		fmt.Fprintln(os.Stderr, "Received SIGINT, terminating...")
	case <-ctx.Done():
	}

	fmt.Fprintln(os.Stderr, "...EOF")
}
