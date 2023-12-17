package config

import (
  "context"
  "path/filepath"
  "sync"

  log "github.com/sirupsen/logrus"
)

type (
  ctxKey struct{}
)

var (
  once   sync.Once
  client Client
)

type Client interface {
  GetAppInfo() AppInfo
  GetValue(configKey string) Value
}

func ClientConfig(ctx context.Context) Client {
  if ctxClient, ok := ctx.Value(ctxKey{}).(Client); ok {
    return ctxClient
  }
  if client == nil {
    return noopClient
  }
  return client
}

func ContextWithClientConfig(parent context.Context, client Client) context.Context {
  return context.WithValue(parent, ctxKey{}, client)
}

func InitClientConfig() {
  once.Do(func() {
    configPath := filepath.Join(".config", "app_config.yaml")

    if !findConfig(configPath) {
      log.Warnf("config: file not found: %s", configPath)
      // Use noop client if config not found
      client = newNoopClient()
      return
    }

    parsed, err := ParseConfig(configPath)
    if err != nil {
      log.Fatalf("config: parsing failed:\n%v", err)
    }
    if err = parsed.Validate(); err != nil {
      log.Fatalf("config: validation failed:\n%v", err)
    }
    app := collectAppInfo(parsed.App)

    values, err := collectConfigValues(parsed.Custom)
    if err != nil {
      log.Fatalf("boiler: values collecting failed: %v", err)
    }

    // Use default client
    client = newClient(app, values)
  })
}

type configClient struct {
  app    AppInfo
  values configValues
}

func newClient(app AppInfo, values configValues) *configClient {
  return &configClient{
    app:    app,
    values: values,
  }
}

func (c *configClient) GetValue(configKey string) Value {
  if value, ok := c.values[configKey]; ok {
    return value
  }
  return &configValue{}
}

func (c *configClient) GetAppInfo() AppInfo {
  return c.app
}

