package env

import "os"

const (
  PrometheusEndpointKey     Key = "BOILER_PROMETHEUS_ENDPOINT"
  PrometheusEndpointDefault Env = "localhost:9090"

  GrafanaEndpointKey     Key = "BOILER_GRAFANA_ENDPOINT"
  GrafanaEndpointDefault Env = "localhost:3000"

  JaegerEndpointKey     Key = "BOILER_JAEGER_ENDPOINT"
  JaegerEndpointDefault Env = "localhost:4318"

  EtcdEndpointsKey     Key = "BOILER_ETCD_ENDPOINTS"
  EtcdEndpointsDefault Env = "localhost:2379"
)

type (
  Key string
  Env string
)

func (k Key) String() string {
  return string(k)
}

func (e Env) String() string {
  return string(e)
}

func (e Env) OrDefault(env Env) Env {
  if e != "" {
    return e
  }
  return env
}

func Get(key Key) Env {
  return Env(os.Getenv(key.String()))
}
