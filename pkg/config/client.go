package config

import (
  "context"

  "github.com/ushakovn/boiler/pkg/config/provider"
  "github.com/ushakovn/boiler/pkg/config/types"
)

type Client interface {
  GetAppInfo() AppInfo
  GetValue(ctx context.Context, key string) types.Value
  WatchValue(ctx context.Context, key string, action func(types.Value))
}

type configClient struct {
  app       AppInfo
  providers []provider.Values
}

func newClient(app AppInfo, providers ...provider.Values) Client {
  return &configClient{
    app:       app,
    providers: providers,
  }
}

func (c *configClient) GetAppInfo() AppInfo {
  return c.app
}

func (c *configClient) GetValue(ctx context.Context, key string) types.Value {
  for _, p := range c.providers {
    if value := p.Get(ctx, key); !value.IsNil() {
      return value
    }
  }
  return types.NewNilValue()
}

func (c *configClient) WatchValue(ctx context.Context, key string, action func(types.Value)) {
  go func(ctx context.Context) {
    for _, p := range c.providers {
      p.Watch(ctx, key, action)
    }
  }(ctx)
}

type noopConfigClient struct{}

func newNoopClient() *noopConfigClient {
  return &noopConfigClient{}
}

func (c *noopConfigClient) GetAppInfo() AppInfo {
  const (
    name    = "app"
    version = "v0.0.0"
  )
  return AppInfo{
    Name:    name,
    Version: version,
  }
}

func (c *noopConfigClient) GetValue(_ context.Context, _ string) types.Value {
  return types.NewNilValue()
}

func (c *noopConfigClient) WatchValue(_ context.Context, _ string, action func(types.Value)) {
  action(types.NewNilValue())
}
