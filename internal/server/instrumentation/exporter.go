package instrumentation

import (
	"context"

	"github.com/marechal-dev/RouteBastion-Broker/internal/utils"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
)

func InitExporter(config utils.AppEnvConfig) (*otlptrace.Exporter, error) {
	headers := map[string]string{
  	"content-type": "application/json",
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(config.OtelEndpoint),
			otlptracehttp.WithHeaders(headers),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		return nil, err
	}

	return exporter, nil
}
