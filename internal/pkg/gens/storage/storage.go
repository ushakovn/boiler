package storage

import (
  "context"
  "fmt"

  "github.com/ushakovn/boiler/internal/boiler/gen"
  "github.com/ushakovn/boiler/internal/pkg/sql"
)

type storage struct {
  dumpSQL    *sql.DumpSQL
  schemaDesc *schemaDesc
}

type Config struct {
  PgConfigPath string
  PgDumpPath   string
}

func (c *Config) Validate() error {
  if (c.PgDumpPath == "" && c.PgConfigPath == "") || (c.PgDumpPath != "" && c.PgConfigPath != "") {
    return fmt.Errorf("pg dump path OR pg config path must be specified")
  }
  return nil
}

func NewStorage(config Config) (gen.Generator, error) {
  if err := config.Validate(); err != nil {
    return nil, err
  }
  var (
    filePath string
    option   sql.PgDumpOption
  )
  if filePath = config.PgConfigPath; filePath != "" {
    option = sql.NewPgDumpOption(sql.WithPgConfigFile, filePath)
  }
  if filePath = config.PgDumpPath; filePath != "" {
    option = sql.NewPgDumpOption(sql.WithPgDumpFile, filePath)
  }
  dumpSQL, err := sql.DumpSchemaSQL(option)
  if err != nil {
    return nil, fmt.Errorf("sql.DumpSchemaSQL: %w", err)
  }
  return &storage{
    dumpSQL: dumpSQL,
  }, nil
}

func (g *storage) Generate(ctx context.Context) error {
  return nil
}
