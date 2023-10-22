package mongo

import (
	"context"
	"log"

	"sync"
	"time"

	"github.com/mrangelba/go-toolkit/config"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

var once sync.Once
var instance *mongo.Database

func Get() *mongo.Database {
	once.Do(func() {
		instance = new()
	})

	return instance
}

func new() *mongo.Database {
	cfg := config.Get()

	credential := options.Credential{
		Username: cfg.Mongo.User,
		Password: cfg.Mongo.Password,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.Mongo.Uri).SetAuth(credential)

	if cfg.IsDev() {
		cmdMonitor := &event.CommandMonitor{
			Started: func(_ context.Context, evt *event.CommandStartedEvent) {
				log.Printf("MongoDB Command: %v", evt.Command)
			},
		}
		clientOptions.SetMonitor(cmdMonitor)
	} else {
		clientOptions.SetMonitor(otelmongo.NewMonitor())
	}

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatalf("MongoDB Error: %v", err)
	}

	return client.Database(cfg.Mongo.Database)
}
