package producer

import (
  "fmt"

  "github.com/IBM/sarama"
)

type Config struct {
  Brokers []string
}

func New(config Config) (sarama.SyncProducer, error) {
  producer, err := sarama.NewSyncProducer(config.Brokers, sarama.NewConfig())
  if err != nil {
    return nil, fmt.Errorf("sarama.NewSyncProducer: %w", err)
  }
  return producer, nil
}
