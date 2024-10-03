package otel

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/trace"
)

type JaegerExporter struct {
	Endpoint string
}

func NewJaegerExporter(ctx context.Context, endpoint string) (trace.SpanExporter, error) {
	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithInsecure(), otlptracehttp.WithEndpoint(endpoint))
	if err != nil {
		return nil, err
	}

	return exporter, nil
}
