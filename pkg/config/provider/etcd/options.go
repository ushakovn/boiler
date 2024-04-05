package etcd

import (
  "time"
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
  return func(o *calledOptions) {
    o.config = config{
      appName:  "boiler",
      cacheTTL: 15 * time.Second,
    }
  }
}

func callOptions(calls ...Option) *calledOptions {
  calls = append([]Option{WithDefaultConfig()}, calls...)
  o := new(calledOptions)

  for _, call := range calls {
    call(o)
  }
  return o
}
