package etcd

import (
  "context"
  "fmt"
  "time"

  "github.com/jellydator/ttlcache/v3"
  log "github.com/sirupsen/logrus"
  "github.com/ushakovn/boiler/internal/pkg/stringer"
  "github.com/ushakovn/boiler/pkg/config/provider"
  "github.com/ushakovn/boiler/pkg/config/types"
  v3 "go.etcd.io/etcd/client/v3"
)

type etcd struct {
  appName    string
  client     *v3.Client
  cachedKeys *ttlcache.Cache[string, types.Value]
}

type config struct {
  client   v3.Config
  appName  string
  cacheTTL time.Duration
}

func New(calls ...Option) provider.Values {
  options := callOptions(calls...)

  client, err := v3.New(options.client)
  if err != nil {
    log.Fatalf("config: failed to create etcd values provider: %v", err)
  }
  cache := ttlcache.New[string, types.Value](
    ttlcache.WithTTL[string, types.Value](
      options.cacheTTL,
    ),
  )
  return &etcd{
    appName:    options.appName,
    client:     client,
    cachedKeys: cache,
  }
}

func (e *etcd) Get(ctx context.Context, key string) types.Value {
  key = e.buildKey(key)

  if value := e.getCached(key); !value.IsNil() {
    return value
  }
  return e.get(ctx, key)
}

func (e *etcd) Watch(ctx context.Context, key string, action func(value types.Value)) {
  key = e.buildKey(key)

  ch := e.client.Watch(ctx, key)
  for {
    select {
    case resp := <-ch:
      for _, event := range resp.Events {
        isCreateOrUpdate := event.IsCreate() || event.IsModify()

        if !isCreateOrUpdate || event.Kv == nil {
          continue
        }
        value := types.NewValue(string(event.Kv.Value))
        e.cachedKeys.Set(key, value, ttlcache.DefaultTTL)

        action(value)
      }
    case <-ctx.Done():
      return
    }
  }
}

func (e *etcd) get(ctx context.Context, key string) types.Value {
  resp, err := e.client.Get(ctx, key,
    v3.WithSort(v3.SortByVersion, v3.SortDescend),
  )
  if err != nil || resp.Count == 0 || len(resp.Kvs) == 0 {
    return types.NewNilValue()
  }
  kv := resp.Kvs[len(resp.Kvs)-1]
  value := types.NewValue(string(kv.Value))

  e.cachedKeys.Set(key, value, ttlcache.DefaultTTL)
  return value
}

func (e *etcd) getCached(key string) types.Value {
  if value := e.cachedKeys.Get(key); value != nil {
    return value.Value()
  }
  return types.NewNilValue()
}

func (e *etcd) buildKey(key string) string {
  key = fmt.Sprintf("%s_%s", e.appName, key)
  key = stringer.StringToSnakeCase(key)
  return key
}
