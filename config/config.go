package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/mrangelba/go-toolkit/config/log"
	"github.com/mrangelba/go-toolkit/config/mongo"
	"github.com/mrangelba/go-toolkit/config/opentelemetry"
	"github.com/mrangelba/go-toolkit/config/postgres"
	"github.com/mrangelba/go-toolkit/config/rabbitmq"
	"github.com/mrangelba/go-toolkit/config/service"
)

type Config struct {
	Service       service.ServiceConfig
	Log           log.LogConfig
	Mongo         mongo.MongoConfig
	OpenTelemetry opentelemetry.OTelConfig
	Postgres      postgres.PostgresConfig
	RabbitMQ      rabbitmq.RabbitMQConfig
}

func Run() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Error loading .env file %v", err)
	}
}

func Get() Config {
	return Config{
		Service:       service.Get(),
		Log:           log.Get(),
		Mongo:         mongo.Get(),
		OpenTelemetry: opentelemetry.Get(),
		Postgres:      postgres.Get(),
		RabbitMQ:      rabbitmq.Get(),
	}
}

func (c *Config) IsProd() bool {
	return c.Service.Env == "production"
}

func (c *Config) IsDev() bool {
	return c.Service.Env == "development"
}

func (c *Config) IsTest() bool {
	return c.Service.Env == "test"
}
