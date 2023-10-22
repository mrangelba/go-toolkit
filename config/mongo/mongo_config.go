package mongo

import (
	"os"
)

type MongoConfig struct {
	Uri      string
	User     string
	Password string
	Database string
}

func Get() MongoConfig {
	return MongoConfig{
		Uri:      os.Getenv("MONGO_URI"),
		User:     os.Getenv("MONGO_USER"),
		Password: os.Getenv("MONGO_PASSWORD"),
		Database: os.Getenv("MONGO_DATABASE"),
	}
}
