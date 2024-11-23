package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "Connection timeout")

	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go-telnet --timeout=10s host port")
		os.Exit(1)
	}

	host, port := flag.Arg(0), flag.Arg(1)
	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to %s: %v\n", address, err)

		stop()

		os.Exit(1) //nolint:gocritic
	}
	defer client.Close()

	go func() {
		if err := client.Receive(); err != nil && !errors.Is(err, io.EOF) {
			fmt.Fprintf(os.Stderr, "Receive error: %v\n", err)
		}

		stop()
	}()

	go func() {
		if err := client.Send(); err != nil {
			fmt.Fprintf(os.Stderr, "Send error: %v\n", err)
		}
		stop()
	}()

	<-ctx.Done()

	fmt.Fprintln(os.Stderr, "...Connection closed, exiting")
}
