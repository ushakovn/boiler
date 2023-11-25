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

  Short: "Init a boiler app config",
  Long:  `Init a boiler app config`,

  RunE: func(cmd *cobra.Command, args []string) error {
    ctx := context.Background()

    initor, err := gen.NewInitor(config.Config{})
    if err != nil {
      return fmt.Errorf("boiler: failed to create initor: %w", err)
    }
    if err = initor.Init(ctx); err != nil {
      return fmt.Errorf("boiler: initor failed: %w", err)
    }
    log.Infof("boiler: app config initialized")

    return nil
  },
}
