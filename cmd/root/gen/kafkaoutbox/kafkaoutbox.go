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

  Short: "Generate a Kafka outbox components",
  Long:  `Generate a Kafka outbox components`,

  RunE: func(cmd *cobra.Command, args []string) error {
    ctx := context.Background()

    generator, err := gen.NewGenerator(kafkaoutbox.Config{})
    if err != nil {
      return fmt.Errorf("boiler: failed to create generator: %w", err)
    }
    if err = generator.Generate(ctx); err != nil {
      return fmt.Errorf("boiler: generator failed: %w", err)
    }
    log.Infof("boiler: kafka outbox components generated")

    return nil
  },
}
