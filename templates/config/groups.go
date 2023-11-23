package config

import "github.com/ushakovn/boiler/pkg/config"

var (
  PgConn = pgConn_ConfigGroup{}
  Kafka  = kafka_ConfigGroup{}
)

var (
  _ = PgConn
  _ = Kafka
)

type (
  pgConn_ConfigGroup struct{}
  kafka_ConfigGroup  struct{}
)

type (
  pgConn_ConfigKey string
  kafka_ConfigKey  string
)

func configValue(configKey string) config.Value {
  return config.ClientConfig().GetValue(configKey)
}
