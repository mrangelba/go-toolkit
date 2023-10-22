package opentelemetry

import (
	"os"
)

type OTelConfig struct {
	ExporterOTLPEndPoint string
	InsecureNode         bool
	Enabled              bool
}

func Get() OTelConfig {
	return OTelConfig{
		ExporterOTLPEndPoint: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		InsecureNode:         os.Getenv("OTEL_INSECURE_MODE") == "true",
		Enabled:              os.Getenv("OTEL_ENABLED") == "true",
	}
}
