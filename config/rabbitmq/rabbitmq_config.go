package rabbitmq

import (
	"os"
)

type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

func Get() RabbitMQConfig {
	return RabbitMQConfig{
		Host:     os.Getenv("RABBITMQ_HOST"),
		Port:     os.Getenv("RABBITMQ_PORT"),
		User:     os.Getenv("RABBITMQ_USER"),
		Password: os.Getenv("RABBITMQ_PASSWORD"),
	}
}
