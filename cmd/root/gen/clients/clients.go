package clients

import (
  "context"
  "fmt"

  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
  "github.com/ushakovn/boiler/internal/boiler/gen"
  "github.com/ushakovn/boiler/internal/pkg/gens/clients"
)

var (
  flagGitHubToken      string
  flagProtoClientsPath string
)

var CmdClients = &cobra.Command{
  Use: "clients",

  Short: "Generate a grpc clients for proto import line",
  Long:  `Generate a grpc clients for proto import line`,

  RunE: func(cmd *cobra.Command, args []string) error {
    ctx := context.Background()

    generator, err := gen.NewGenerator(clients.Config{
      GithubToken:     flagGitHubToken,
      ClientsDescPath: flagProtoClientsPath,
    })
    if err != nil {
      return fmt.Errorf("boiler: failed to create generator: %w", err)
    }
    if err = generator.Generate(ctx); err != nil {
      return fmt.Errorf("boiler: generator failed: %w", err)
    }
    log.Infof("boiler: grpc clients for proto import lines generated")

    return nil
  },
}

func init() {
  CmdClients.Flags().StringVar(&flagGitHubToken, "github-token", "", "access token for github api")
  CmdClients.Flags().StringVar(&flagProtoClientsPath, "proto-clients", "", "path to grpc clients file in json/yaml")
}
