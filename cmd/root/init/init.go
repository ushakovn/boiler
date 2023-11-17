package init

import (
  "context"
  "fmt"

  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
  cmdGqlgen "github.com/ushakovn/boiler/cmd/root/init/gqlgen"
  cmdGrpc "github.com/ushakovn/boiler/cmd/root/init/grpc"
  cmdProtoDeps "github.com/ushakovn/boiler/cmd/root/init/protodeps"
  cmdStorage "github.com/ushakovn/boiler/cmd/root/init/storage"
  "github.com/ushakovn/boiler/internal/boiler/gen"
  "github.com/ushakovn/boiler/internal/pkg/gens/project"
)

var (
  flagProjectConfigPath string
  flagGoModVersion      string
)

var CmdInit = &cobra.Command{
  Use: "init",

  Short: "Init a template for a microservice application in the Go language",
  Long:  `Init a template for a microservice application in the Go language`,

  PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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
  CmdInit.AddCommand(cmdGrpc.CmdGrpc, cmdGqlgen.CmdGqlgen, cmdProtoDeps.CmdProtoDeps, cmdStorage.CmdStorage)

  CmdInit.PersistentFlags().StringVar(&flagProjectConfigPath, "project-config-path", "", "path to project directories config in json/yaml")
  CmdInit.PersistentFlags().StringVar(&flagGoModVersion, "go-mod-version", "", "go mod version for project")
}
