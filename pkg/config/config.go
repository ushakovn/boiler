package config

import (
  "path/filepath"
  "sync"

  log "github.com/sirupsen/logrus"
  "github.com/ushakovn/boiler/pkg/config/provider/etcd"
  "github.com/ushakovn/boiler/pkg/config/provider/local"
)

var (
  once   sync.Once
  client Client
)

var (
  noopClient = newNoopClient()
)

func InitClient() {
  once.Do(func() {
    // Default config path
    path := filepath.Join(".config", "app_config.yaml")

    if !findConfig(path) {
      log.Warnf("config: file not found: %s", path)
      // Use noop client if config not found
      client = newNoopClient()
      return
    }
    parsed, err := ParseConfig(path)
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
    // Use config client
    client = newClient(app,
      etcd.New(etcd.WithAppName(app.Name)),
      local.New(values),
    )
  })
}
