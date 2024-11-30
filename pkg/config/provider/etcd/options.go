package etcd

import (
  "time"

  "github.com/ushakovn/boiler/pkg/env"
  v3 "go.etcd.io/etcd/client/v3"
  "go.uber.org/zap"
  "google.golang.org/grpc"
)

type Option func(*calledOptions)

type calledOptions struct {
  // Embed config
  config
}

func WithAppName(appName string) Option {
  return func(o *calledOptions) {
    o.appName = appName
  }
}

func WithCacheTTL(ttl time.Duration) Option {
  return func(o *calledOptions) {
    o.config.cacheTTL = ttl
  }
}

func WithDefaultConfig() Option {
  const (
    appName = "boiler"
    timeout = 100 * time.Millisecond
  )

  endpoints := env.Get(env.EtcdEndpointsKey).
    OrDefault(env.EtcdEndpointsDefault).
    String()

  return func(o *calledOptions) {
    o.config = config{
      // Etcd client config
      client: v3.Config{
        Endpoints:   []string{endpoints},
        DialTimeout: timeout,
        Username:    appName,
        DialOptions: []grpc.DialOption{
          grpc.WithTimeout(timeout),
        },
        Logger: zap.NewNop(),
      },
      // Values provider config
      appName:  appName,
      cacheTTL: 15 * time.Second,
    }
  }
}

func callOptions(calls ...Option) *calledOptions {
  calls = append(defaultOptions(), calls...)
  o := new(calledOptions)

  for _, call := range calls {
    call(o)
  }
  return o
}

func defaultOptions() []Option {
  return []Option{WithDefaultConfig()}
}
