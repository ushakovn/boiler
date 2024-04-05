package etcd

import (
  "sync"

  log "github.com/sirupsen/logrus"
  v3 "go.etcd.io/etcd/client/v3"
)

var (
  once   sync.Once
  client *v3.Client
)

func Client() *v3.Client {
  if client == nil {
    log.Fatalf("boiler: etcd client not initialized")
  }
  return client
}
