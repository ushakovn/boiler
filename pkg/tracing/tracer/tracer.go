package tracer

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

// Suppress unused variable
var _ = tracer

var (
  mu     sync.Mutex
  once   sync.Once
  tracer trace.Tracer
)

func Tracer() trace.Tracer {
  // Lock before returning
  mu.Lock()
  defer mu.Unlock()

  if tracer == nil {
    panic("tracer: tracer not initialized")
  }
  return tracer
}

func StartContextWithSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
  // Context and span
  return Tracer().Start(ctx, name, opts...)
}

func StartContextSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) context.Context {
  // Only context
  ctx, _ = StartContextWithSpan(ctx, name, opts...)
  return ctx
}

func SpanFromContext(ctx context.Context) trace.Span {
  // Wrap span from context
  return trace.SpanFromContext(ctx)
}

func InitTracer(ctx context.Context, serviceName, serviceVer string) (shutdowns []func(ctx context.Context) error) {
  once.Do(func() {
    exporter, err := otlptracehttp.New(ctx,
      otlptracehttp.WithInsecure(),
    )
    if err != nil {
      log.Fatalf("tracer: otlptracehttp.New: %v", err)
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
      log.Fatalf("tracer: resource.Merge: %v", err)
    }

    tr := sdktrace.NewTracerProvider(
      sdktrace.WithBatcher(exporter),
      sdktrace.WithResource(res),
    )
    shutdowns = append(shutdowns, exporter.Shutdown)

    otel.SetTracerProvider(tr)

    mu.Lock()
    defer mu.Unlock()

    tracer = otel.Tracer(serviceName)
  })

  return shutdowns
}
