package tracing

import (
  "context"
  "sync"

  log "github.com/sirupsen/logrus"
  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
  "go.opentelemetry.io/otel/sdk/resource"
  sdktrace "go.opentelemetry.io/otel/sdk/trace"
  semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
  "go.opentelemetry.io/otel/trace"
)

var (
  once   sync.Once
  tracer trace.Tracer
)

func InitTracer(ctx context.Context, serviceName, serviceVer string) (shutdowns []func(ctx context.Context) error) {
  once.Do(func() {
    exporter, err := otlptracehttp.New(ctx,
      otlptracehttp.WithInsecure(),
    )
    if err != nil {
      log.Fatalf("tracing: otlptracehttp.New: %v", err)
    }
    shutdowns = append(shutdowns, exporter.Shutdown)

    res, err := resource.Merge(
      resource.Default(),
      resource.NewWithAttributes(
        semconv.SchemaURL,
        semconv.ServiceName(serviceName),
        semconv.ServiceVersion(serviceVer),
      ),
    )
    if err != nil {
      log.Fatalf("tracing: resource.Merge: %v", err)
    }

    tr := sdktrace.NewTracerProvider(
      sdktrace.WithBatcher(exporter),
      sdktrace.WithResource(res),
    )
    shutdowns = append(shutdowns, exporter.Shutdown)

    otel.SetTracerProvider(tr)
    tracer = otel.Tracer(serviceName)
  })

  return shutdowns
}

func Tracer() trace.Tracer {
  if tracer == nil {
    panic("tracer: tracer not initialized")
  }
  return tracer
}

func StartSpanFromContext(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
  return Tracer().Start(ctx, name, opts...)
}
