package provider

import (
	"context"
	"fmt"

	"github.com/rss3-network/global-indexer/internal/config"
	"github.com/rss3-network/global-indexer/internal/constant"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	tracer "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

func ProvideOpenTelemetryTracer(configFile *config.File) (trace.TracerProvider, error) {
	if configFile.Telemetry == nil {
		return otel.GetTracerProvider(), nil
	}

	options := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(configFile.Telemetry.Endpoint),
	}

	if configFile.Telemetry.Insecure {
		options = append(options, otlptracehttp.WithInsecure())
	}

	exporter, err := otlptrace.New(context.TODO(), otlptracehttp.NewClient(options...))
	if err != nil {
		return nil, fmt.Errorf("new exporter: %w", err)
	}

	var (
		serviceName    = constant.BuildServiceName()
		serviceVersion = constant.BuildServiceVersion()
	)

	tracerProvider := tracer.NewTracerProvider(
		tracer.WithBatcher(exporter),
		tracer.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(serviceVersion),
		)),
	)

	return tracerProvider, nil
}
