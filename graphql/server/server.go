package server

import (
	"sync"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
)

var once sync.Once
var instance *handler.Server

func Get(schema graphql.ExecutableSchema) *handler.Server {
	once.Do(func() {
		instance = new(schema)
	})

	return instance
}

func new(schema graphql.ExecutableSchema) *handler.Server {
	server := handler.NewDefaultServer(schema)

	return server
}
