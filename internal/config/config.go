package config

import (
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

type Config struct {
	PostgresUser     string
	PostgresPassword string
	PostgresDbName   string
	PostgresHost     string
	PostgresPort     string
	ServerPort       string
	MigrationFlag    string
	MigrationPath    string
	AgifyURL         string
	GenderizeURL     string
	NationalizeURL   string
}

func LoadCfg() *Config {
	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file", "error", err)
		panic(err)
	}
	cfg := &Config{
		PostgresUser:     os.Getenv("Postgres_User"),
		PostgresPassword: os.Getenv("Postgres_Password"),
		PostgresDbName:   os.Getenv("Postgres_DbName"),
		PostgresHost:     os.Getenv("Postgres_Host"),
		PostgresPort:     os.Getenv("Postgres_Port"),
		ServerPort:       os.Getenv("Server_Port"),
		MigrationFlag:    os.Getenv("Migration_Flag"),
		MigrationPath:    os.Getenv("Migration_Path"),
		AgifyURL:         os.Getenv("Agify_URL"),
		GenderizeURL:     os.Getenv("Genderize_URL"),
		NationalizeURL:   os.Getenv("Nationalize_URL"),
	}
	return cfg
}
