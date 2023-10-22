package server

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/mrangelba/go-toolkit/grpc/interceptors"
	"github.com/mrangelba/go-toolkit/logger"

	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/reflection"
)

var once sync.Once
var instance *grpc.Server

func Get(registerServices func(*grpc.Server)) *grpc.Server {
	once.Do(func() {
		instance = new(registerServices)
	})

	return instance
}

func new(registerServices func(*grpc.Server)) *grpc.Server {
	oplog := logger.Get()
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.AuthInterceptor,
			otelgrpc.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			otelgrpc.StreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(),
		),
	)
	oplog.Info().Msg("Server GRPC created")

	registerServices(server)
	reflection.Register(server)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-quit
		oplog.Info().Msg("Shutdown GRPC ...")
		server.GracefulStop()
		oplog.Info().Msg("Server GRPC exit")
	}()

	return server
}
