package protodeps

import (
  "context"
  "fmt"

  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
  "github.com/ushakovn/boiler/internal/boiler/gen"
  "github.com/ushakovn/boiler/internal/pkg/gens/protodeps"
)

var (
  flagGitHubToken   string
  flagProtoDepsPath string
)

var CmdProtoDeps = &cobra.Command{
  Use: "proto-deps",

  Short: "Generate a Proto dependencies",
  Long:  `Generate a Proto dependencies`,

  RunE: func(cmd *cobra.Command, args []string) error {
    ctx := context.Background()

    generator, err := gen.NewGenerator(protodeps.Config{
      GithubToken:   flagGitHubToken,
      ProtoDepsPath: flagProtoDepsPath,
    })
    if err != nil {
      return fmt.Errorf("boiler: failed to create generator: %w", err)
    }
    if err = generator.Generate(ctx); err != nil {
      return fmt.Errorf("boiler: generator failed: %w", err)
    }
    log.Infof("boiler: proto dependencies generated")

    return nil
  },
}

func init() {
  CmdProtoDeps.Flags().StringVar(&flagGitHubToken, "github-token", "", "access token for github api")
  CmdProtoDeps.Flags().StringVar(&flagProtoDepsPath, "proto-deps-path", "", "path to grpc clients file in json/yaml")

  CmdProtoDeps.MarkFlagsRequiredTogether("github-token", "proto-deps-path")
}
