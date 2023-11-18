package gqlgen

import (
  "context"
  "fmt"

  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
  "github.com/ushakovn/boiler/internal/boiler/gen"
  "github.com/ushakovn/boiler/internal/pkg/gens/gqlgen"
  "github.com/ushakovn/boiler/internal/pkg/gens/project"
)

var (
  flagGqlgenConfigPath string
)

var (
  flagProjectConfigPath string
  flagGoModVersion      string
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
  CmdGqlgen.Flags().StringVar(&flagGqlgenConfigPath, "gqlgen-config-path", "", "path to gqlgen directories config in json/yaml")

  CmdGqlgen.Flags().StringVar(&flagProjectConfigPath, "project-config-path", "", "path to project directories config in json/yaml")
  CmdGqlgen.Flags().StringVar(&flagGoModVersion, "go-mod-version", "", "go mod version for project")
}
