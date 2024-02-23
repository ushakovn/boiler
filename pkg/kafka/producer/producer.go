package producer

import (
  "fmt"

  "github.com/IBM/sarama"
)

type Config struct {
  Brokers []string
}

func New(config Config) (sarama.SyncProducer, error) {
  producer, err := sarama.NewSyncProducer(config.Brokers, newConfig())
  if err != nil {
    return nil, fmt.Errorf("sarama.NewSyncProducer: %w", err)
  }
  return producer, nil
}

func newConfig() *sarama.Config {
  config := sarama.NewConfig()

  config.Producer.Return.Successes = true
  config.Producer.Return.Errors = true

  config.Producer.Idempotent = true
  config.Producer.RequiredAcks = sarama.WaitForAll

  config.Net.MaxOpenRequests = 1

  return config
}
