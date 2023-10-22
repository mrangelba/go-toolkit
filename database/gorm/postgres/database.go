package postgres

import (
	"fmt"
	"log"
	"sync"

	"github.com/mrangelba/go-toolkit/config"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

var once sync.Once
var instance *gorm.DB

func Get() *gorm.DB {
	once.Do(func() {
		instance = new()
	})

	return instance
}

func new() *gorm.DB {
	cfg := config.Get()
	datasource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database)

	db, err := gorm.Open(postgres.Open(datasource), &gorm.Config{})

	if err != nil {
		log.Fatalf("postgres rrror: %v", err)
	}

	err = db.Use(tracing.NewPlugin(
		tracing.WithAttributes(semconv.DBSystemPostgreSQL),
		tracing.WithDBName(cfg.Postgres.Database),
	))

	if err != nil {
		log.Fatalf("postgres rrror: %v", err)
	}

	return db
}
