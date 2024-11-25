package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"slices"

	migrate "github.com/rubenv/sql-migrate"
)

const (
	argMigrate = "migrate"
	argDown    = "down"

	dialect = "postgres"
)

//go:embed *.sql
var migrations embed.FS

func Migrate(ctx context.Context, conn *sql.DB, args ...string) error {
	if !slices.Contains(args, argMigrate) {
		return nil
	}

	slog.InfoContext(ctx, "migrating")

	var (
		operation = migrate.Up
		message   = "Applied"
	)

	if slices.Contains(args, argDown) {
		operation = migrate.Down
		message = "Reverted"
	}

	mg := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: migrations,
		Root:       ".",
	}

	num, err := migrate.ExecContext(ctx, conn, dialect, mg, operation)
	if err != nil {
		return fmt.Errorf("error %s migrations: %w", message, err)
	}

	slog.InfoContext(ctx, fmt.Sprintf("%s %d migrations", message, num))

	return nil
}
