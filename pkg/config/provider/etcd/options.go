package etcd

import (
  "time"

  "github.com/ushakovn/boiler/pkg/env"
  v3 "go.etcd.io/etcd/client/v3"
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
  const appName = "boiler"

  endpoints := env.Get(env.EtcdEndpointsKey).
    OrDefault(env.EtcdEndpointsDefault).
    String()

  return func(o *calledOptions) {
    o.config = config{
      // Etcd client config
      client: v3.Config{
        Username:  appName,
        Endpoints: []string{endpoints},
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
