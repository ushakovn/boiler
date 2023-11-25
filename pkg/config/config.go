package config

import (
  "path/filepath"
  "sync"

  log "github.com/sirupsen/logrus"
)

// Suppress unused variable
var _ = client

var (
  mu     sync.Mutex
  once   sync.Once
  client Client
)

type Client interface {
  GetAppInfo() AppInfo
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
    mu.Lock()
    defer mu.Unlock()

    configPath := filepath.Join(".boiler", "config.yaml")

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

type noopConfigClient struct{}

func newNoopClient() *noopConfigClient {
  return &noopConfigClient{}
}

func (c *noopConfigClient) GetValue(string) Value {
  return &configValue{}
}

func (c *noopConfigClient) GetAppInfo() AppInfo {
  const (
    name    = "Boiler"
    version = "v0.0.1"
  )
  return AppInfo{
    Name:    name,
    Version: version,
  }
}
