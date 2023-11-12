package grpc

import (
  "context"
  "fmt"

  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
  "github.com/ushakovn/boiler/internal/boiler/gen"
  "github.com/ushakovn/boiler/internal/pkg/gens/grpc"
)

var CmdGrpc = &cobra.Command{
  Use: "grpc",

  Short: "Init a gRPC template components",
  Long:  `Init a gRPC template components`,

  RunE: func(cmd *cobra.Command, args []string) error {
    ctx := context.Background()

    initor, err := gen.NewInitor(grpc.Config{})
    if err != nil {
      return fmt.Errorf("boiler: failed to create initor: %w", err)
    }
    if err = initor.Init(ctx); err != nil {
      return fmt.Errorf("boiler: initor failed: %w", err)
    }
    log.Infof("boiler: grpc components initialized")

    return nil
  },
}
