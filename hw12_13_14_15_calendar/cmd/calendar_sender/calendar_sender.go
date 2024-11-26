package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/rabbitmq"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config_sender.json", "Path to configuration file")
}

func main() {
	lvl := new(slog.LevelVar)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
	slog.SetDefault(logger)

	flag.Parse()

	var cfg Config
	if err := config.LoadConfig[Config](configFile, &cfg); err != nil {
		slog.Error("error while loading configuration", "error", err)
		os.Exit(1)
	}

	if err := ValidateConfig(&cfg); err != nil {
		slog.Error("error while validating configuration", "error", err)
		os.Exit(1)
	}

	if v, ok := config.LogLevelNames[cfg.LogLvl]; ok {
		lvl.Set(v)
	}

	if err := run(cfg); err != nil {
		slog.Error("application error", "error", err)
		os.Exit(1)
	}
}

func run(cfg Config) error {
	rmq, err := rabbitmq.NewQueue(cfg.RabbitMQURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer rmq.Close()

	messages, err := rmq.Consume()
	if err != nil {
		return fmt.Errorf("failed to consume messages: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		slog.Info("Sender started")
		for {
			select {
			case <-ctx.Done():
				slog.Info("Service shutting down...")
				return
			case msg, ok := <-messages:
				if !ok {
					slog.Info("Message channel closed")
					return
				}

				if err = processMessage(msg); err != nil {
					slog.Error("failed to process message", "error", err)
					errCh <- err
				}
			}
		}
	}()

	select {
	case err = <-errCh:
		return err
	case <-ctx.Done():
		slog.Info("Service stopped by user")
		return nil
	}
}

func processMessage(body []byte) error {
	var notification rabbitmq.Notification
	if err := json.Unmarshal(body, &notification); err != nil {
		return fmt.Errorf("failed to deserialize message: %w", err)
	}

	slog.Info("Sending notification", "notification", notification)
	return nil
}
