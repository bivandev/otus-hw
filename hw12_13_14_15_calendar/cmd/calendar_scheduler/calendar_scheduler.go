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
	"time"

	"github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/rabbitmq"
	memorystorage "github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config_scheduler.json", "Path to configuration file")
}

func main() {
	flag.Parse()

	lvl := new(slog.LevelVar)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
	slog.SetDefault(logger)

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

	if err = rmq.DeclareQueue(); err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	storage, err := initStorage(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	defer storage.Close()

	calendar := app.New(storage)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	errCh := make(chan error, 1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err = processEvents(ctx, calendar, rmq); err != nil {
					errCh <- err
					return
				}
			}
		}
	}()

	select {
	case err = <-errCh:
		return fmt.Errorf("service encountered an error: %w", err)
	case <-ctx.Done():
		slog.Info("Service stopped by user")
		return nil
	}
}

func processEvents(ctx context.Context, app *app.App, rmq rabbitmq.Queue) error {
	events, err := app.GetNotification(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch events: %w", err)
	}

	for _, event := range events {
		body, err := json.Marshal(event)
		if err != nil {
			slog.Error("failed to serialize notification", "error", err)
			continue
		}

		if err = rmq.Publish(body); err != nil {
			slog.Error("failed to publish message", "error", err)
			continue
		}

		slog.Info("Notification sent", "EventID", event.EventID)
	}

	if err = app.CleanOldEvents(ctx); err != nil {
		slog.Error("failed to clean old events", "error", err)
	}

	return nil
}

func initStorage(cfg Config) (app.Storage, error) {
	if !cfg.UsePostgres {
		return memorystorage.New(), nil
	}

	dbStorage, err := sqlstorage.New(cfg.PostgresDSN)
	if err != nil {
		return nil, fmt.Errorf("error connecting to postgres: %w", err)
	}

	return dbStorage, nil
}
