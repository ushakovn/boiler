package app

import (
  "net/http"

  "github.com/99designs/gqlgen/graphql"
  "github.com/samber/lo"
  "github.com/ushakovn/boiler/pkg/tracing"
  "google.golang.org/grpc"
  "google.golang.org/grpc/stats"
  "google.golang.org/grpc/tap"
)

type Option func(o *calledAppOptions)

type calledAppOptions struct {
  // gRPC
  grpcServePort     int
  grpcServerOptions []grpc.ServerOption

  // GraphQL
  gqlgenServePort int
  gqlgenMWs       []func(http.Handler) http.Handler

  gqlgenFieldMWs     []graphql.FieldMiddleware
  gqlgenOperationMWs []graphql.OperationMiddleware
  gqlgenResponseMWs  []graphql.ResponseMiddleware
}

func defaultOptions() []Option {
  const (
    defaultGrpcPort   = 8082
    defaultGqlgenPort = 8080
  )
  options := []Option{
    // Port options
    WithGrpcServePort(defaultGrpcPort),
    WithGqlgenServePort(defaultGqlgenPort),

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

func WithGrpcServePort(port int) Option {
  return func(o *calledAppOptions) {
    o.grpcServePort = port
  }
}

func WithGrpcServerOptions(options ...grpc.ServerOption) Option {
  return func(o *calledAppOptions) {
    o.grpcServerOptions = append(o.grpcServerOptions, options...)
  }
}

func WithGrpcUnaryServerInterceptors(interceptors ...grpc.UnaryServerInterceptor) Option {
  return func(o *calledAppOptions) {
    serverOptions := lo.Map(interceptors, func(interceptor grpc.UnaryServerInterceptor, _ int) grpc.ServerOption {
      return grpc.UnaryInterceptor(interceptor)
    })
    o.grpcServerOptions = append(o.grpcServerOptions, serverOptions...)
  }
}

func WithGrpcStatsHandlers(handlers ...stats.Handler) Option {
  return func(o *calledAppOptions) {
    serverOptions := lo.Map(handlers, func(handler stats.Handler, _ int) grpc.ServerOption {
      return grpc.StatsHandler(handler)
    })
    o.grpcServerOptions = append(o.grpcServerOptions, serverOptions...)
  }
}

func WithGrpcTapHandlers(handlers ...tap.ServerInHandle) Option {
  return func(o *calledAppOptions) {
    serverOptions := lo.Map(handlers, func(handler tap.ServerInHandle, _ int) grpc.ServerOption {
      return grpc.InTapHandle(handler)
    })
    o.grpcServerOptions = append(o.grpcServerOptions, serverOptions...)
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
