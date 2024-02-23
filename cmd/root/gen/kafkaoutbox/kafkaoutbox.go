package kafkaoutbox

import (
  "context"
  "fmt"

  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
  "github.com/ushakovn/boiler/internal/boiler/gen"
  "github.com/ushakovn/boiler/internal/pkg/gens/kafkaoutbox"
)

var (
  flagValidateProto     bool
  flagStorageConfigPath string
)

var CmdKafkaoutbox = &cobra.Command{
  Use: "kafkaoutbox",

  Short: "Generate a Kafka outbox components",
  Long:  `Generate a Kafka outbox components`,

  RunE: func(cmd *cobra.Command, args []string) error {
    ctx := context.Background()

    generator, err := gen.NewGenerator(kafkaoutbox.Config{
      ValidateProto:     flagValidateProto,
      StorageConfigPath: flagStorageConfigPath,
    })
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

func init() {
  CmdKafkaoutbox.Flags().BoolVar(&flagValidateProto, "validate-proto", false, "validate kafkaoutbox proto with pg schema")
  CmdKafkaoutbox.Flags().StringVar(&flagStorageConfigPath, "storage-config-path", "", "path to storage generator config")

  CmdKafkaoutbox.MarkFlagsRequiredTogether("validate-proto", "storage-config-path")
}
