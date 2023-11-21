package tracing

import (
  "context"
  "sync"
  "time"

  "github.com/99designs/gqlgen/graphql"
  log "github.com/sirupsen/logrus"
  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/attribute"
  otelCodes "go.opentelemetry.io/otel/codes"
  "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
  "go.opentelemetry.io/otel/sdk/resource"
  sdktrace "go.opentelemetry.io/otel/sdk/trace"
  semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "google.golang.org/grpc/metadata"
  "google.golang.org/grpc/status"
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

func SpanFromContext(ctx context.Context) trace.Span {
  // Wrap span from context
  return trace.SpanFromContext(ctx)
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
  spanCtx = metadata.AppendToOutgoingContext(spanCtx, "trace-id", span.SpanContext().TraceID().String())

  // Handle request
  if resp, err = handler(spanCtx, req); err != nil {
    // Set span error status
    errString := err.Error()
    span.SetStatus(otelCodes.Error, errString)
    // Set gRPC error attributes
    span.SetAttributes(attribute.String("grpcError", errString))
    span.SetAttributes(attribute.String("grpcStatusCode", status.Code(err).String()))
  }
  return resp, err
}

func GqlgenOperationMiddleware(ctx context.Context, handler graphql.OperationHandler) graphql.ResponseHandler {
  // Tracing middleware
  if graphql.HasOperationContext(ctx) {
    // Extract operation context
    opCtx := graphql.GetOperationContext(ctx)
    var (
      opTyp  string
      opName string
    )
    if operation := opCtx.Operation; operation != nil {
      opTyp = opCtx.Operation.Name
      opName = string(opCtx.Operation.Operation)
    }
    var span trace.Span

    ctx, span = StartContextWithSpan(ctx, opCtx.OperationName,
      // Start span options
      trace.WithSpanKind(trace.SpanKindServer),
      trace.WithTimestamp(time.Now().UTC()),

      // Operation context info
      trace.WithAttributes(
        attribute.String("operationType", opTyp),
        attribute.String("operationName", opName),

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
  }
  return handler(ctx)
}

func GqlgenResponseMiddleware(ctx context.Context, handler graphql.ResponseHandler) *graphql.Response {
  span := SpanFromContext(ctx)
  graphql.RegisterExtension(ctx, "traceID", span.SpanContext().TraceID().String())
  // Handle errors
  errors := graphql.GetErrors(ctx)

  if len(errors) != 0 {
    // Set span error status
    errString := errors.Error()
    span.SetStatus(otelCodes.Error, errString)
    // Set GraphQL error attributes
    span.SetAttributes(attribute.String("graphqlError", errString))
  }

  // Handle request
  return handler(ctx)
}
