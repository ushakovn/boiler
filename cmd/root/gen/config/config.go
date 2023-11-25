package config

import (
  "context"
  "fmt"

  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
  "github.com/ushakovn/boiler/internal/boiler/gen"
  "github.com/ushakovn/boiler/internal/pkg/gens/config"
)

var CmdConfig = &cobra.Command{
  Use: "config",

  Short: "Generate a boiler app config",
  Long:  `Generate a boiler app config`,

  RunE: func(cmd *cobra.Command, args []string) error {
    ctx := context.Background()

    generator, err := gen.NewGenerator(config.Config{})
    if err != nil {
      return fmt.Errorf("boiler: failed to create generator: %w", err)
    }
    if err = generator.Generate(ctx); err != nil {
      return fmt.Errorf("boiler: generator failed: %w", err)
    }
    log.Infof("boiler: app config generated")

    return nil
  },
}
