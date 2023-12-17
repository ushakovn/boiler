package tracer

import (
  "go.opentelemetry.io/otel/trace"
  "go.opentelemetry.io/otel/trace/noop"
)

var (
  noopTracer = newNoopTracer()
)

func newNoopTracer() trace.Tracer {
  const tracerName = "noopTracer"

  return noop.NewTracerProvider().
    Tracer(tracerName)
}
