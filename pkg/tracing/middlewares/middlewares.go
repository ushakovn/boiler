package middlewares

import (
  "context"
  "time"

  "github.com/99designs/gqlgen/graphql"
  "github.com/ushakovn/boiler/pkg/tracing/tracer"
  "go.opentelemetry.io/otel/attribute"
  otelCodes "go.opentelemetry.io/otel/codes"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "google.golang.org/grpc/metadata"
  "google.golang.org/grpc/status"
)

func GrpcServerUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
  // Tracing interceptor
  spanCtx, span := tracer.StartContextWithSpan(ctx, info.FullMethod,
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
      opTyp = string(operation.Operation)
      opName = operation.Name
    }
    var span trace.Span

    ctx, span = tracer.StartContextWithSpan(ctx, opCtx.OperationName,
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
  span := tracer.SpanFromContext(ctx)
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
