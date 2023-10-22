package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
	"github.com/mrangelba/go-toolkit/config"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var once sync.Once
var instance *sql.DB

func Get() *sql.DB {
	once.Do(func() {
		instance = new()
	})

	return instance
}

func new() *sql.DB {
	cfg := config.Get()
	datasource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database)

	db, err := otelsql.Open("postgres", datasource,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
		otelsql.WithDBName("db-sqllite"))

	if err != nil {
		log.Fatalf("postgres rrror: %v", err)
	}

	return db
}
