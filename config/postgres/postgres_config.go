package postgres

import (
	"os"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func Get() PostgresConfig {
	return PostgresConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Database: os.Getenv("POSTGRES_DATABASE"),
	}
}
