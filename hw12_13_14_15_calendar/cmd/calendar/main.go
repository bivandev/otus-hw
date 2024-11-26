package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/server/grpc"
	memorystorage "github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/devv4n/otus-hw/hw12_13_14_15_calendar/migrations"
	"github.com/jackc/pgx/v5/stdlib"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config.json", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	lvl := new(slog.LevelVar)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))

	slog.SetDefault(logger)

	var cfg Config
	if err := config.LoadConfig[Config](configFile, &cfg); err != nil {
		if errors.Is(err, config.ErrConfigNotFound) && configFile == "config.json" {
			DefaultConfig(&cfg)
		}

		slog.Error("error while loading configuration", "error", err)
	}

	if err := ValidateConfig(&cfg); err != nil {
		slog.Error("error while validate configuration", "error", err)

		os.Exit(1)
	}

	lvl.Set(config.LogLevelNames[cfg.LogLvl])

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage, err := initStorage(ctx, cfg)
	if err != nil {
		slog.Error("error initializing storage", "error", err)

		cancel()
		os.Exit(1) //nolint:gocritic
	}

	calendar := app.New(storage)

	srv := grpc.New(
		calendar,
		grpc.Config{
			GRPC: cfg.GRPC,
			REST: cfg.REST,
		},
	)

	errCh := make(chan error, 2)

	go srv.ServeUserAPI(errCh)
	go srv.ServeGatewayAPI(ctx, errCh)

	slog.Info("calendar is running...")

	select {
	case err = <-errCh:
		slog.Error("Service encountered an error", "error", err)

		os.Exit(1)
	case <-ctx.Done():
		slog.Info("Service stopped by user")
	}
}

func initStorage(ctx context.Context, cfg Config) (app.Storage, error) {
	if !cfg.UsePostgres {
		return memorystorage.New(), nil
	}

	dbStorage, err := sqlstorage.New(cfg.PostgresDSN)
	if err != nil {
		return nil, fmt.Errorf("error connecting to postgres: %w", err)
	}

	if args := flag.Args(); len(args) > 0 {
		if err = migrations.Migrate(ctx, stdlib.OpenDBFromPool(dbStorage.Pool), args...); err != nil {
			slog.Error("error while migrate", "error", err)
		}
	}

	return dbStorage, nil
}
