package gqlgen

import (
  "context"
  "fmt"

  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
  "github.com/ushakovn/boiler/internal/boiler/gen"
  "github.com/ushakovn/boiler/internal/pkg/executor"
  "github.com/ushakovn/boiler/internal/pkg/gens/gqlgen"
  "github.com/ushakovn/boiler/internal/pkg/gens/project"
  "github.com/ushakovn/boiler/templates"
)

var (
  flagGqlgenConfigPath string
)

var (
  flagProjectConfigPath string
  flagGoModVersion      string
  flagGoModName         string
)

var CmdGqlgen = &cobra.Command{
  Use: "gqlgen",

  Short: "Init a GraphQL components",
  Long:  `Init a GraphQL components`,

  RunE: func(cmd *cobra.Command, args []string) error {
    ctx := context.Background()

    initor, err := gen.NewInitor(gqlgen.Config{
      GqlgenDescPath: flagGqlgenConfigPath,
    })
    if err != nil {
      return fmt.Errorf("boiler: failed to create initor: %w", err)
    }
    if err = initor.Init(ctx); err != nil {
      return fmt.Errorf("boiler: initor failed: %w", err)
    }
    log.Infof("boiler: gqlgen components initialized")

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
    return execMakeGqlgenBinDeps()
  },
}

func init() {
  CmdGqlgen.Flags().StringVar(&flagGqlgenConfigPath, "gqlgen-config-path", "", "path to gqlgen directories config in json/yaml")

  CmdGqlgen.Flags().StringVar(&flagProjectConfigPath, "project-config-path", "", "path to project directories config in json/yaml")
  CmdGqlgen.Flags().StringVar(&flagGoModVersion, "go-version", "", "go version for project")
  CmdGqlgen.Flags().StringVar(&flagGoModName, "go-module", "", "go module name for project")
}

func execMakeGqlgenBinDeps() error {
  ctx := context.Background()

  if err := executor.ExecCmdCtx(ctx, "make", templates.GqlgenMakeMkBinDepsName); err != nil {
    return fmt.Errorf("boiler: failed to exec: make %s", templates.GqlgenMakeMkBinDepsName)
  }
  return nil
}
