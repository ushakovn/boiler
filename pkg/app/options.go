package app

import (
  "sync"

  "github.com/gin-gonic/gin"
  "google.golang.org/grpc"
)

type Option func(o *calledAppOptions)

type calledAppOptions struct {
  mu sync.Mutex

  grpcServePort int
  httpServePort int

  grpcOptions []grpc.ServerOption
  ginHandles  []gin.HandlerFunc
}

func callAppOptions(calls ...Option) *calledAppOptions {
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

func WithHttpServePort(port int) Option {
  return func(o *calledAppOptions) {
    o.httpServePort = port
  }
}

func WithGrpcServerOptions(grpcOptions ...grpc.ServerOption) Option {
  return func(o *calledAppOptions) {
    o.mu.Lock()
    defer o.mu.Unlock()
    o.grpcOptions = append(o.grpcOptions, grpcOptions...)
  }
}

func WithGinHandles(ginHandles ...gin.HandlerFunc) Option {
  return func(o *calledAppOptions) {
    o.mu.Lock()
    defer o.mu.Unlock()
    o.ginHandles = append(o.ginHandles, ginHandles...)
  }
}
