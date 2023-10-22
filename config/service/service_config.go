package service

import (
	"os"
)

type ServiceConfig struct {
	Database string
	Env      string
	Host     string
	Name     string
	Port     string
	Version  string
}

func Get() ServiceConfig {
	return ServiceConfig{
		Database: os.Getenv("SERVICE_DATABASE"),
		Env:      os.Getenv("SERVICE_ENV"),
		Host:     os.Getenv("SERVICE_HOST"),
		Name:     os.Getenv("SERVICE_NAME"),
		Port:     os.Getenv("SERVICE_PORT"),
		Version:  os.Getenv("SERVICE_VERSION"),
	}
}
