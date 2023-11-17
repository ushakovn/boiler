package protodeps

import (
  "context"
  "fmt"

  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
  "github.com/ushakovn/boiler/internal/boiler/gen"
  "github.com/ushakovn/boiler/internal/pkg/gens/protodeps"
)

var CmdProtoDeps = &cobra.Command{
  Use: "proto-deps",

  Short: "Init a Proto dependencies components",
  Long:  `Init a Proto dependencies components`,

  RunE: func(cmd *cobra.Command, args []string) error {
    ctx := context.Background()

    generator, err := gen.NewInitor(protodeps.Config{})
    if err != nil {
      return fmt.Errorf("boiler: failed to create initor: %w", err)
    }
    if err = generator.Init(ctx); err != nil {
      return fmt.Errorf("boiler: initor failed: %w", err)
    }
    log.Infof("boiler: proto dependencies components initialized")

    return nil
  },
}
