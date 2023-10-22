package sqlite

import (
	"log"
	"sync"

	"github.com/mrangelba/go-toolkit/config"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"gorm.io/driver/sqlite"
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

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})

	if err != nil {
		log.Fatalf("sqlite error: %v", err)
	}

	err = db.Use(tracing.NewPlugin(
		tracing.WithAttributes(semconv.DBSystemSqlite),
		tracing.WithDBName(cfg.Postgres.Database),
	))

	if err != nil {
		log.Fatalf("sqlite error: %v", err)
	}

	return db
}
