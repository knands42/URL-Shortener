package otel

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type OpenTelemetry struct {
	exporter sdktrace.SpanExporter
}

func NewOpenTelemetry(exporter sdktrace.SpanExporter) *OpenTelemetry {
	return &OpenTelemetry{
		exporter: exporter,
	}
}

func (ot *OpenTelemetry) GetTracer() trace.Tracer {
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSyncer(ot.exporter),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	otel.Meter("url-shortener-app")
	tracer := otel.Tracer("url-shortener-app")
	return tracer
}
