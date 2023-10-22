package sqlite

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
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
	db, err := otelsql.Open("sqlite3", "file::memory:",
		otelsql.WithAttributes(semconv.DBSystemSqlite),
		otelsql.WithDBName("db-sqllite"))

	if err != nil {
		log.Fatalf("sqlite rrror: %v", err)
	}

	return db
}
