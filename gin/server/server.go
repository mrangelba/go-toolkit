package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-contrib/requestid"
	"github.com/mrangelba/go-toolkit/config"
	"github.com/mrangelba/go-toolkit/logger"
	"github.com/mrangelba/go-toolkit/module"

	"github.com/gin-gonic/gin"
)

var once sync.Once
var instance *http.Server

func Get(modules ...module.Module) *http.Server {
	once.Do(func() {
		instance = newServer(modules...)
	})

	return instance
}

func Run(modules ...module.Module) {
	server := Get(modules...)

	log := logger.Get()

	log.Info().Msgf("Starting HTTP Server REST on port %s ...", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Msgf("Server HTTP Server REST failed to start: %s", err)
	}
}

func newServer(modules ...module.Module) *http.Server {
	cfg := config.Get()
	log := logger.Get()

	r := gin.Default()

	if cfg.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Use(requestid.New())

	registerModules(r, modules...)

	server := &http.Server{
		Addr:    cfg.Service.Port,
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-quit
		log.Info().Msg("Shutdown HTTP Server REST ...")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Errorf("Server Shutdown: %s", err)
		}
		log.Info().Msg("Server HTTP Server REST exit")
	}()

	return server
}

func registerModules(r *gin.Engine, modules ...module.Module) {
	for _, module := range modules {
		module.Register(r)
	}
}
