package tracing

import (
  "context"
  "net/http"
  "sync"
  "time"

  "github.com/99designs/gqlgen/graphql"
  log "github.com/sirupsen/logrus"
  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/attribute"
  "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
  "go.opentelemetry.io/otel/sdk/resource"
  sdktrace "go.opentelemetry.io/otel/sdk/trace"
  semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
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

func StartContextWithSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
  // Context and span
  return Tracer().Start(ctx, name, opts...)
}

func StartContextSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) context.Context {
  // Only context
  ctx, _ = StartContextWithSpan(ctx, name, opts...)
  return ctx
}

func GrpcServerUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
  // Tracing interceptor
  spanCtx, span := StartContextWithSpan(ctx, info.FullMethod,
    // Start span options
    trace.WithSpanKind(trace.SpanKindServer),
    trace.WithTimestamp(time.Now().UTC()),
  )
  defer span.End(
    trace.WithStackTrace(true),
    trace.WithTimestamp(time.Now().UTC()),
  )
  return handler(spanCtx, req)
}

func GqlgenMiddleware(next http.Handler) http.Handler {
  // Tracing middleware
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    reqCtx := r.Context()

    if graphql.HasOperationContext(reqCtx) {
      // Extract operation context
      opCtx := graphql.GetOperationContext(reqCtx)

      spanCtx, span := StartContextWithSpan(reqCtx, opCtx.RawQuery,
        // Start span options
        trace.WithSpanKind(trace.SpanKindServer),
        trace.WithTimestamp(time.Now().UTC()),

        // Operation context info
        trace.WithAttributes(
          attribute.String("operationName", opCtx.OperationName),

          attribute.String("statsOperationStart",
            opCtx.Stats.OperationStart.Format(time.RFC3339),
          ),

          attribute.StringSlice("statsRead", []string{
            opCtx.Stats.Read.Start.Format(time.RFC3339),
            opCtx.Stats.Read.End.Format(time.RFC3339),
          }),

          attribute.StringSlice("statsParsing", []string{
            opCtx.Stats.Parsing.Start.Format(time.RFC3339),
            opCtx.Stats.Parsing.End.Format(time.RFC3339),
          }),

          attribute.StringSlice("statsParsing", []string{
            opCtx.Stats.Validation.Start.Format(time.RFC3339),
            opCtx.Stats.Validation.End.Format(time.RFC3339),
          }),
        ),
      )
      defer span.End(
        trace.WithStackTrace(true),
        trace.WithTimestamp(time.Now().UTC()),
      )
      next.ServeHTTP(w, r.WithContext(spanCtx))
    }
  })
}
