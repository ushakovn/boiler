package app

import (
  "github.com/gin-gonic/gin"
  "github.com/ushakovn/boiler/internal/pkg/aggr"
  "google.golang.org/grpc"
  "google.golang.org/grpc/stats"
  "google.golang.org/grpc/tap"
)

type Option func(o *calledAppOptions)

type calledAppOptions struct {
  grpcServePort     int
  grpcServerOptions []grpc.ServerOption

  gqlgenServePort int
  gqlgenHandlers  []gin.HandlerFunc
}

func callAppOptions(calls ...Option) *calledAppOptions {
  const (
    defaultGrpcPort   = 8082
    defaultGqlgenPort = 8080
  )
  o := &calledAppOptions{
    grpcServePort:   defaultGrpcPort,
    gqlgenServePort: defaultGqlgenPort,
  }
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
    serverOptions := aggr.Map(interceptors, func(interceptor grpc.UnaryServerInterceptor) grpc.ServerOption {
      return grpc.UnaryInterceptor(interceptor)
    })
    o.grpcServerOptions = append(o.grpcServerOptions, serverOptions...)
  }
}

func WithGrpcStatsHandlers(handlers ...stats.Handler) Option {
  return func(o *calledAppOptions) {
    serverOptions := aggr.Map(handlers, func(handler stats.Handler) grpc.ServerOption {
      return grpc.StatsHandler(handler)
    })
    o.grpcServerOptions = append(o.grpcServerOptions, serverOptions...)
  }
}

func WithGrpcTapHandlers(handlers ...tap.ServerInHandle) Option {
  return func(o *calledAppOptions) {
    serverOptions := aggr.Map(handlers, func(handler tap.ServerInHandle) grpc.ServerOption {
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

func WithGqlgenHandles(handlers ...gin.HandlerFunc) Option {
  return func(o *calledAppOptions) {
    o.gqlgenHandlers = append(o.gqlgenHandlers, handlers...)
  }
}
