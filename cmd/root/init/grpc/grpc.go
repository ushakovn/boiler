package grpc

import (
  "context"
  "fmt"

  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
  "github.com/ushakovn/boiler/internal/boiler/gen"
  "github.com/ushakovn/boiler/internal/pkg/gens/grpc"
  "github.com/ushakovn/boiler/internal/pkg/gens/project"
)

var (
  flagProjectConfigPath string
  flagGoModVersion      string
)

var CmdGrpc = &cobra.Command{
  Use: "grpc",

  Short: "Init a gRPC components",
  Long:  `Init a gRPC components`,

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

  PreRunE: func(cmd *cobra.Command, args []string) error {
    ctx := context.Background()

    initor, err := gen.NewInitor(project.Config{
      GoModVersion:    flagGoModVersion,
      ProjectDescPath: flagProjectConfigPath,
    })
    if err != nil {
      return fmt.Errorf("boiler: failed to create initor: %w", err)
    }
    if err = initor.Init(ctx); err != nil {
      return fmt.Errorf("boiler: initor failed: %w", err)
    }
    log.Infof("boiler: project components initialized")

    return nil
  },
}

func init() {
  CmdGrpc.Flags().StringVar(&flagProjectConfigPath, "project-config-path", "", "path to project directories config in json/yaml")
  CmdGrpc.Flags().StringVar(&flagGoModVersion, "go-mod-version", "", "go mod version for project")
}
