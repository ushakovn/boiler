package config

import (
  "log"
  "sync"
)

// Suppress unused variable
var _ = client

var (
  mu     sync.Mutex
  once   sync.Once
  client Client
)

type Client interface {
  GetValue(configKey string) Value
}

func ClientConfig() Client {
  // Lock before returning
  mu.Lock()
  defer mu.Unlock()

  if client == nil {
    panic("config: client not initialized")
  }
  return client
}

func InitClientConfig() {
  once.Do(func() {
    parsed, err := ParseConfig()
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
    mu.Lock()
    defer mu.Unlock()

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
