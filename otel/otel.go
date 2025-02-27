package otel

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type Otel struct {
	tracer *trace.TracerProvider
}

func New(ctx context.Context, url, serviceName string) (*Otel, error) {
	tp, err := newTraceProvider(ctx, url, serviceName)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tp)

	return &Otel{
		tracer: tp,
	}, nil
}

func (o *Otel) SetLogger(name string) {
	provider := log.NewLoggerProvider()
	slog.SetDefault(otelslog.NewLogger(name, otelslog.WithLoggerProvider(provider)))
}

func (o *Otel) Close(ctx context.Context) error {
	return o.tracer.Shutdown(ctx)
}

func newTraceProvider(ctx context.Context, url, serviceName string) (*trace.TracerProvider, error) {
	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint(url), otlptracehttp.WithInsecure())
	if err != nil {
		return nil, err
	}

	return trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName)))), nil
}
