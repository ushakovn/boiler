package kafkaoutbox

import (
  "context"
  "fmt"

  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
  "github.com/ushakovn/boiler/internal/boiler/gen"
  "github.com/ushakovn/boiler/internal/pkg/gens/kafkaoutbox"
)

var CmdKafkaoutbox = &cobra.Command{
  Use: "kafkaoutbox",

  Short: "Init a Kafka outbox components",
  Long:  `Init a Kafka outbox components`,

  RunE: func(cmd *cobra.Command, args []string) error {
    ctx := context.Background()

    initor, err := gen.NewInitor(kafkaoutbox.Config{})
    if err != nil {
      return fmt.Errorf("boiler: failed to create initor: %w", err)
    }
    if err = initor.Init(ctx); err != nil {
      return fmt.Errorf("boiler: initor failed: %w", err)
    }
    log.Infof("boiler: kafka outbox components initialized")

    return nil
  },
}
