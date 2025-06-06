package database

import (
	"Test_Task_EffMob/internal/config"
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"log/slog"
	"time"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	connStr := DSNWithPostgresHost(cfg)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		if err = db.Close(); err != nil {
			slog.Error("Failed to close DB", slog.String("error", err.Error()))
		}
		connStr = DSNWithDBHost(cfg)
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			return nil, err
		}
		if err = db.PingContext(ctx); err != nil {
			return nil, err
		}
	}

	slog.InfoContext(context.Background(), "Successfully connected to database", "host", cfg.PostgresHost, "port", cfg.PostgresPort, "dbname", cfg.PostgresDbName)

	return &Database{DB: db}, nil
}

func DSNWithDBHost(cfg *config.Config) string {
	return "postgres://" + cfg.PostgresUser + ":" + cfg.PostgresPassword + "@" +
		cfg.DBhost + ":" + cfg.PostgresPort + "/" + cfg.PostgresDbName + "?sslmode=disable"
}

func DSNWithPostgresHost(cfg *config.Config) string {
	return "postgres://" + cfg.PostgresUser + ":" + cfg.PostgresPassword + "@" +
		cfg.PostgresHost + ":" + cfg.PostgresPort + "/" + cfg.PostgresDbName + "?sslmode=disable"
}

func (d *Database) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}
	return nil
}
