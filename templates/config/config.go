package config

import (
  "context"
  "time"
)

const (
  PgConn_Timeout pgConn_ConfigKey = "pg_conn_timeout"
)

const (
  Kafka_ErrorsPerTopic kafka_ConfigKey = "kafka_errors_per_topic"
)

func (c pgConn_ConfigGroup) Timeout() time.Duration {
  return configValue(string(PgConn_Timeout)).Duration()
}

func (c pgConn_ConfigGroup) Context(parent context.Context) context.Context {
  return context.WithValue(parent, PgConn_Timeout, c.Timeout())
}

func (c kafka_ConfigGroup) ErrorsPerTopic() string {
  return configValue(string(Kafka_ErrorsPerTopic)).String()
}

func (c kafka_ConfigGroup) Context(parent context.Context) context.Context {
  return context.WithValue(parent, Kafka_ErrorsPerTopic, c.ErrorsPerTopic())
}

var (
  _ = time.Time{}
)
