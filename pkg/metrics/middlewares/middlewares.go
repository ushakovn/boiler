package middlewares

import (
  "context"
  "sync"
  "time"

  "github.com/99designs/gqlgen/graphql"
  "github.com/prometheus/client_golang/prometheus"
  log "github.com/sirupsen/logrus"
  "github.com/ushakovn/boiler/pkg/metrics"
  "google.golang.org/grpc"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"
)

type mwMetrics struct {
  grpcReqDur     *prometheus.HistogramVec
  grpcReqCount   *prometheus.CounterVec
  gqlgenReqDur   *prometheus.HistogramVec
  gqlgenReqCount *prometheus.CounterVec
}

var (
  m    *mwMetrics
  mu   sync.Mutex
  once sync.Once
)

// InitMetrics USE ONLY AFTER CALL config.InitClientConfig
func InitMetrics() {
  once.Do(func() {
    // Latency metric for gRPC
    grpcRequestDurationHistogram := metrics.NewHistogramVec(
      "grpc_request_duration_seconds_histogram",
      "Histogram of gRPC request duration in seconds",
      prometheus.LinearBuckets(0, 0.300, 4),
      []string{"method", "code"},
    )
    // RPS metric for gRPC
    grpcRequestCounter := metrics.NewCounterVec(
      "grpc_request_counter",
      "Counter of gRPC requests",
      []string{"method", "code"},
    )

    // Latency metric for GraphQL
    gqlgenRequestDurationHistogram := metrics.NewHistogramVec(
      "gqlgen_request_duration_seconds_histogram",
      "Histogram of GraphQL request duration in seconds",
      prometheus.LinearBuckets(0, 0.300, 4),
      []string{"method"},
    )

    // RPS metric for GraphQL
    gqlgenRequestCounter := metrics.NewCounterVec(
      "gqlgen_request_counter",
      "Counter of GraphQL requests",
      []string{"method"},
    )

    mu.Lock()
    defer mu.Unlock()

    m = &mwMetrics{
      grpcReqDur:   grpcRequestDurationHistogram,
      grpcReqCount: grpcRequestCounter,

      gqlgenReqDur:   gqlgenRequestDurationHistogram,
      gqlgenReqCount: gqlgenRequestCounter,
    }
  })
}

// GrpcServerUnaryInterceptor USE ONLY AFTER CALL InitMetrics
func GrpcServerUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
  reqStartTime := time.Now()

  statusCode := codes.OK
  methodName := info.FullMethod

  // Handle request
  resp, respErr := handler(ctx, req)
  // Evaluate duration
  reqDurSec := time.Since(reqStartTime).Seconds()

  if respErr != nil {
    statusCode = status.Code(respErr)
  }

  // Try to observe duration
  if durHist, err := m.grpcReqDur.GetMetricWithLabelValues(methodName, statusCode.String()); err != nil {
    log.Errorf("metrics: grpc duration histogram error: %v", err)
  } else {
    durHist.Observe(reqDurSec)
  }

  // Try to increment counter
  reqCounter, err := m.grpcReqCount.GetMetricWithLabelValues(methodName, statusCode.String())
  if err != nil {
    log.Errorf("metrics: grpc request counter error: %v", err)
  } else {
    reqCounter.Inc()
  }

  return resp, respErr
}

// GqlgenOperationMiddleware USE ONLY AFTER CALL InitMetrics
func GqlgenOperationMiddleware(ctx context.Context, handler graphql.OperationHandler) graphql.ResponseHandler {
  reqStartTime := time.Now()

  defer func() {
    // Evaluate duration
    reqDurSec := time.Since(reqStartTime).Seconds()

    if graphql.HasOperationContext(ctx) {
      // Extract operation context
      opCtx := graphql.GetOperationContext(ctx)
      var opName string

      if operation := opCtx.Operation; operation != nil {
        opName = operation.Name
      }

      // Try to observe request duration
      if durHist, err := m.gqlgenReqDur.GetMetricWithLabelValues(opName); err != nil {
        log.Errorf("metrics: gqlgen duration histogram error: %v", err)
      } else {
        durHist.Observe(reqDurSec)
      }

      // Try to increment counter
      if reqCounter, err := m.gqlgenReqCount.GetMetricWithLabelValues(opName); err != nil {
        log.Errorf("metrics: gqlgen request counter error: %v", err)
      } else {
        reqCounter.Inc()
      }
    }
  }()

  // Handle request
  return handler(ctx)
}
