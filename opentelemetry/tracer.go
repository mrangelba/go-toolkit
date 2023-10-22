package opentelemetry

import (
	"context"

	"github.com/mrangelba/go-toolkit/config"
	"github.com/mrangelba/go-toolkit/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func InitTracer() func(context.Context) error {
	cfg := config.Get()
	oplog := logger.Get()

	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if cfg.OpenTelemetry.InsecureNode {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(cfg.OpenTelemetry.ExporterOTLPEndPoint),
		),
	)

	if err != nil {
		oplog.Err(err).Msg("Could not initialize exporter")
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", cfg.Service.Name),
			attribute.String("service.version", cfg.Service.Version),
			attribute.String("service.environment", cfg.Service.Env),
			attribute.String("library.language", "go"),
		),
	)

	if err != nil {
		oplog.Err(err).Msg("Could not set resources")
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)

	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		oplog.Err(err).Msg("Error on OpenTelemetry")
	}))

	return exporter.Shutdown
}
