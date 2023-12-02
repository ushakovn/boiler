package app

import (
  "net/http"

  "github.com/99designs/gqlgen/graphql"
  grpcMW "github.com/grpc-ecosystem/go-grpc-middleware"
  recover "github.com/ushakovn/boiler/pkg/recover/middlewares"
  tracing "github.com/ushakovn/boiler/pkg/tracing/middlewares"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"
  "google.golang.org/grpc/stats"
  "google.golang.org/grpc/tap"
)

type Option func(o *calledAppOptions)

type calledAppOptions struct {
  // gRPC
  grpcServePort     int
  grpcHttpProxyPort int
  grpcServerOptions []grpc.ServerOption

  grpcStatsHandler       stats.Handler
  grpcTapInHandler       tap.ServerInHandle
  grpcServerInterceptors []grpc.UnaryServerInterceptor

  // GraphQL
  gqlgenServePort int
  gqlgenMWs       []func(http.Handler) http.Handler

  gqlgenFieldMWs     []graphql.FieldMiddleware
  gqlgenOperationMWs []graphql.OperationMiddleware
  gqlgenResponseMWs  []graphql.ResponseMiddleware
}

func defaultOptions() []Option {
  const (
    defaultGrpcPort          = 8082
    defaultGrpcHttpProxyPort = 8084
    defaultGqlgenPort        = 8080
  )
  options := []Option{
    // Port options
    WithGrpcServePort(defaultGrpcPort),
    WithGrpcHttpProxyPort(defaultGrpcHttpProxyPort),
    WithGqlgenServePort(defaultGqlgenPort),

    // Panic recover options
    WithGrpcUnaryServerInterceptors(recover.GrpcServerUnaryInterceptor),
    WithGqlgenOperationMiddlewares(recover.GqlgenOperationMiddleware),

    // Tracing options
    WithGrpcUnaryServerInterceptors(tracing.GrpcServerUnaryInterceptor),
    WithGqlgenOperationMiddlewares(tracing.GqlgenOperationMiddleware),
    WithGqlgenResponseMiddlewares(tracing.GqlgenResponseMiddleware),
  }
  return options
}

func callAppOptions(calls ...Option) *calledAppOptions {
  calls = append(calls, defaultOptions()...)
  o := &calledAppOptions{}

  for _, call := range calls {
    call(o)
  }
  return o
}

func buildGrpcServerOptions(options *calledAppOptions) []grpc.ServerOption {
  // Set stats handler
  if h := options.grpcStatsHandler; h != nil {
    options.grpcServerOptions = append(options.grpcServerOptions, grpc.StatsHandler(h))
  }
  // Set tap in handler
  if h := options.grpcTapInHandler; h != nil {
    options.grpcServerOptions = append(options.grpcServerOptions, grpc.InTapHandle(h))
  }
  // Set interceptors chain
  if len(options.grpcServerInterceptors) > 0 {
    chain := grpcMW.ChainUnaryServer(options.grpcServerInterceptors...)
    options.grpcServerOptions = append(options.grpcServerOptions, grpc.UnaryInterceptor(chain))
  }
  return options.grpcServerOptions
}

func defaultGrpcClientOptions() []grpc.DialOption {
  return []grpc.DialOption{
    // Without TLS/SSL
    grpc.WithTransportCredentials(insecure.NewCredentials()),
  }
}

func WithGrpcServePort(port int) Option {
  return func(o *calledAppOptions) {
    o.grpcServePort = port
  }
}

func WithGrpcHttpProxyPort(port int) Option {
  return func(o *calledAppOptions) {
    o.grpcHttpProxyPort = port
  }
}

func WithGrpcUnaryServerInterceptors(interceptors ...grpc.UnaryServerInterceptor) Option {
  return func(o *calledAppOptions) {
    o.grpcServerInterceptors = append(o.grpcServerInterceptors, interceptors...)
  }
}

func WithGrpcStatsHandler(handler stats.Handler) Option {
  return func(o *calledAppOptions) {
    o.grpcStatsHandler = handler
  }
}

func WithGrpcTapHandler(handler tap.ServerInHandle) Option {
  return func(o *calledAppOptions) {
    o.grpcTapInHandler = handler
  }
}

// WithGrpcServerOptions USE IN YOUR OWN RISK
func WithGrpcServerOptions(options ...grpc.ServerOption) Option {
  return func(o *calledAppOptions) {
    o.grpcServerOptions = append(o.grpcServerOptions, options...)
  }
}

func WithGqlgenServePort(port int) Option {
  return func(o *calledAppOptions) {
    o.gqlgenServePort = port
  }
}

func WithGqlgenMiddlewares(middlewares ...func(http.Handler) http.Handler) Option {
  return func(o *calledAppOptions) {
    o.gqlgenMWs = append(o.gqlgenMWs, middlewares...)
  }
}

func WithGqlgenFieldMiddlewares(middlewares ...graphql.FieldMiddleware) Option {
  return func(o *calledAppOptions) {
    o.gqlgenFieldMWs = append(o.gqlgenFieldMWs, middlewares...)
  }
}

func WithGqlgenOperationMiddlewares(middlewares ...graphql.OperationMiddleware) Option {
  return func(o *calledAppOptions) {
    o.gqlgenOperationMWs = append(o.gqlgenOperationMWs, middlewares...)
  }
}

func WithGqlgenResponseMiddlewares(middlewares ...graphql.ResponseMiddleware) Option {
  return func(o *calledAppOptions) {
    o.gqlgenResponseMWs = append(o.gqlgenResponseMWs, middlewares...)
  }
}
