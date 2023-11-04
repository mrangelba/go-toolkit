package firebase_auth

import (
	"context"
	"sync"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/mrangelba/go-toolkit/logger"
)

var once sync.Once
var instance *auth.Client

func GetClient() *auth.Client {
	once.Do(func() {
		instance = new()
	})

	return instance
}

func new() *auth.Client {
	ctx := context.Background()

	log := logger.Get()

	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to create firebase Auth client")
	}

	return authClient
}
