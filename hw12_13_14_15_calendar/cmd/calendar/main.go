package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/config"
	internalhttp "github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/server/http"
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

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		slog.Error("error while loading configuration", "error", err)
	}

	lvl.Set(cfg.LogLevel)

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

	server := internalhttp.NewServer(calendar, cfg)

	go func() {
		<-ctx.Done()

		ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		slog.Info("stopping service")

		storage.Close()

		if err = server.Stop(ctx); err != nil {
			slog.Error("failed to stop http server", "error", err)
		}
	}()

	slog.Info("calendar is running...")

	if err = server.Start(ctx); err != nil {
		slog.Error("failed to start http server", "error", err)
		cancel()
		os.Exit(1)
	}
}

func initStorage(ctx context.Context, cfg *config.Config) (app.Storage, error) {
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
