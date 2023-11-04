package jwtauth

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

// Option for queue system
type Option func(*config)

// WithHandler set handler function for request id with context
func WithHandler(handler gin.HandlerFunc) Option {
	return func(cfg *config) {
		cfg.handler = handler
	}
}

func WithFirebaseAuth(firebaseAuth *auth.Client) Option {
	return func(cfg *config) {
		cfg.firebaseAuth = firebaseAuth
	}
}
