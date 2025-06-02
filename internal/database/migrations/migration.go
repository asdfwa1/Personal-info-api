package migrations

import (
	"Test_Task_EffMob/internal/config"
	"Test_Task_EffMob/internal/database"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log/slog"
)

func RunMigrations(cfg *config.Config, flag string) error {
	if flag != "yes" {
		slog.Info("Migrations are disabled")
		return nil
	}
	connStr := database.DSN(cfg)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("Failed to close database connection", "error", err)
		}
	}()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		slog.Error("Failed to create DB driver", "error", err)
		return err
	}

	migrationsPath := cfg.MigrationPath
	if migrationsPath == "" {
		migrationsPath = "file://internal/database/migrations"
		slog.Debug("MigrationsPath env not set, using default path")
	}
	slog.Info("Running migrations from", "path", migrationsPath)

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		slog.Error("Failed to create migrate instance", "error", err)
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("Failed to apply migrations", "error", err)
		return err
	}
	version, _, err := m.Version()
	if err != nil {
		slog.Error("Failed to get migration version", "error", err)
		return err
	}
	slog.Info("Migrations applied successfully", "version", version)

	return nil
}
