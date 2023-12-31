package grpc

import (
  "context"
  "fmt"

  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
  "github.com/ushakovn/boiler/internal/boiler/gen"
  "github.com/ushakovn/boiler/internal/pkg/executor"
  "github.com/ushakovn/boiler/internal/pkg/gens/grpc"
  "github.com/ushakovn/boiler/internal/pkg/gens/project"
  "github.com/ushakovn/boiler/templates"
)

var (
  flagProjectConfigPath string
  flagGoModVersion      string
  flagGoModName         string
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
      GoModName:       flagGoModName,
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

  PostRunE: func(cmd *cobra.Command, args []string) error {
    return execMakeGrpcBinDeps()
  },
}

func init() {
  CmdGrpc.Flags().StringVar(&flagProjectConfigPath, "project-config-path", "", "path to project directories config in json/yaml")
  CmdGrpc.Flags().StringVar(&flagGoModVersion, "go-version", "", "go version for project")
  CmdGrpc.Flags().StringVar(&flagGoModName, "go-module", "", "go module name for project")
}

func execMakeGrpcBinDeps() error {
  ctx := context.Background()

  if err := executor.ExecCmdCtx(ctx, "make", templates.GrpcMakeMkBinDepsName); err != nil {
    return fmt.Errorf("boiler: failed to exec: make %s", templates.GrpcMakeMkBinDepsName)
  }
  return nil
}
