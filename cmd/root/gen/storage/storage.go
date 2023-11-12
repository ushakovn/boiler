package storage

import (
  "context"
  "fmt"

  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
  "github.com/ushakovn/boiler/internal/boiler/gen"
  "github.com/ushakovn/boiler/internal/pkg/gens/storage"
)

var (
  flagPgConfigPath string
  flagPgDumpPath   string
)

var CmdStorage = &cobra.Command{
  Use: "storage",

  Short: "Generate a storage template components",
  Long:  `Generate a storage template components`,

  RunE: func(cmd *cobra.Command, args []string) error {
    ctx := context.Background()

    generator, err := gen.NewGenerator(storage.Config{
      PgConfigPath: flagPgConfigPath,
      PgDumpPath:   flagPgDumpPath,
    })
    if err != nil {
      return fmt.Errorf("boiler: failed to create generator: %w", err)
    }
    if err = generator.Generate(ctx); err != nil {
      return fmt.Errorf("boiler: generator failed: %w", err)
    }
    log.Infof("boiler: storage components generated")

    return nil
  },
}

func init() {
  CmdStorage.Flags().StringVar(&flagPgConfigPath, "pg-conf", "", "path to postgres connection config in json/yaml")
  CmdStorage.Flags().StringVar(&flagPgDumpPath, "pg-dump", "", "path to postgres dump in sql")
}
