package storage

import (
  "context"
  "fmt"

  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
  "github.com/ushakovn/boiler/internal/boiler/gen"
  "github.com/ushakovn/boiler/internal/pkg/gens/storage"
)

var (
  flagConfigPath string
)

var CmdStorage = &cobra.Command{
  Use: "storage",

  Short: "Generate a storage components",
  Long:  `Generate a storage components`,

  RunE: func(cmd *cobra.Command, args []string) error {
    ctx := context.Background()

    generator, err := gen.NewGenerator(storage.ConfigPath(flagConfigPath))
    if err != nil {
      return fmt.Errorf("boiler: failed to create generator: %w", err)
    }
    if err = generator.Generate(ctx); err != nil {
      return fmt.Errorf("boiler: generator failed: %w", err)
    }
    log.Infof("boiler: storage components generated")

    return nil
  },
}

func init() {
  CmdStorage.Flags().StringVar(&flagConfigPath, "config-path", "", "path to storage generator config")
  _ = CmdStorage.MarkFlagRequired("config-path")
}
