package init

import (
  "context"
  "fmt"

  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
  cmdGqlgen "github.com/ushakovn/boiler/cmd/root/init/gqlgen"
  cmdGrpc "github.com/ushakovn/boiler/cmd/root/init/grpc"
  "github.com/ushakovn/boiler/internal/boiler/gen"
  "github.com/ushakovn/boiler/internal/pkg/gens/project"
)

var flagProjectConfigPath string

var CmdInit = &cobra.Command{
  Use: "init",

  SuggestFor: []string{
    "initialize",
  },

  Short: "Init a template for a microservice application in the Go language",
  Long:  `Init a template for a microservice application in the Go language`,

  PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
    ctx := context.Background()

    initor, err := gen.NewInitor(project.Config{
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
  CmdInit.AddCommand(cmdGrpc.CmdGrpc, cmdGqlgen.CmdGqlgen)

  CmdInit.PersistentFlags().StringVar(&flagProjectConfigPath, "proj-conf", "", "path to project directories config in json/yaml")
}
