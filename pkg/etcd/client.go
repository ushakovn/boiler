package etcd

import (
  log "github.com/sirupsen/logrus"
  "github.com/ushakovn/boiler/pkg/env"
  v3 "go.etcd.io/etcd/client/v3"
)

func NewClient(appName string) *v3.Client {
  if appName == "" {
    log.Fatalf("boiler: app name not specified")
  }
  endpoints := env.Get(env.EtcdEndpointsKey).
    OrDefault(env.EtcdEndpointsDefault).
    String()

  client, err := v3.New(v3.Config{
    Username:  appName,
    Endpoints: []string{endpoints},
  })
  if err != nil {
    log.Fatalf("boiler: failed to create etcd client: %v", err)
  }
  return client
}
