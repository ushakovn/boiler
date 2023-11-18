package sql

import (
  "context"
  "encoding/json"
  "fmt"
  "os"

  "github.com/ushakovn/boiler/internal/pkg/executor"
  "github.com/ushakovn/boiler/internal/pkg/filer"
  "github.com/ushakovn/boiler/internal/pkg/validator"
  "gopkg.in/yaml.v3"
)

type pgDumpOptionTyp string

const (
  WithPgDumpFile   pgDumpOptionTyp = "pg_dump"
  WithPgConfigFile pgDumpOptionTyp = "pg_config"
)

type PgDumpOption interface {
  Call(ctx context.Context) (PgDumpBuf []byte, err error)
}

func NewPgDumpOption(typ pgDumpOptionTyp, filePath string) PgDumpOption {
  var option PgDumpOption
  switch typ {
  case WithPgDumpFile:
    option = &withPgDumpFile{filePath: filePath}
  case WithPgConfigFile:
    option = &withPgConfigFile{filePath: filePath}
  }
  return option
}

type withPgDumpFile struct {
  filePath string
}

func (option *withPgDumpFile) Call(context.Context) ([]byte, error) {
  pgDumpBuf, err := os.ReadFile(option.filePath)
  if err != nil {
    return nil, fmt.Errorf("os.ReadFile filePath: %w", err)
  }
  return pgDumpBuf, nil
}

type withPgConfigFile struct {
  filePath string
}

func (option *withPgConfigFile) Call(ctx context.Context) ([]byte, error) {
  buf, err := os.ReadFile(option.filePath)
  if err != nil {
    return nil, fmt.Errorf("os.ReadFile filePath: %w", err)
  }
  fileExtension := filer.ExtractFileExtension(option.filePath)

  config, err := parsePgConfig(fileExtension, buf)
  if err != nil {
    return nil, fmt.Errorf("parsePgConfig: %w", err)
  }
  pgDumpBuf, err := execPgDump(ctx, *config)
  if err != nil {
    return nil, fmt.Errorf("execPgDump: %w", err)
  }
  return pgDumpBuf, nil
}

func parsePgConfig(fileExtension string, buf []byte) (*PgConfig, error) {
  type wrappedConfig struct {
    PgConfig *PgConfig `json:"pg_config" yaml:"pg_config"`
  }
  var (
    config *wrappedConfig
    err    error
  )
  switch fileExtension {
  case "yml", "yaml", "YML", "YAML":
    err = yaml.Unmarshal(buf, &config)
  case "json", "JSON":
    err = json.Unmarshal(buf, &config)
  default:
    err = fmt.Errorf("unsupported file extension: %s", fileExtension)
  }
  if err != nil {
    return nil, err
  }
  if err = validator.ValidateStructWithTags(config, "json", "yaml"); err != nil {
    return nil, err
  }
  return config.PgConfig, err
}

type PgConfig struct {
  Host     string `json:"host" yaml:"host"`
  Port     string `json:"port" yaml:"port"`
  User     string `json:"user" yaml:"user"`
  DBName   string `json:"db_name" yaml:"db_name"`
  Password string `json:"password" yaml:"password"`
}

func (c *PgConfig) pgDumpCmd() (name string, args []string, err error) {
  if err = os.Setenv("PGPASSWORD", c.Password); err != nil {
    return "", nil, fmt.Errorf("os.Setenv: %w", err)
  }
  args = []string{
    "--no-owner",
    "-h", c.Host,
    "-p", c.Port,
    "-U", c.User,
    c.DBName,
  }
  name = "pg_dump"

  return name, args, nil
}

func execPgDump(ctx context.Context, config PgConfig) ([]byte, error) {
  name, args, err := config.pgDumpCmd()
  if err != nil {
    return nil, fmt.Errorf("conn.pgDump: %w", err)
  }
  buf, err := executor.ExecCmdCtxWithOut(ctx, name, args...)
  if err != nil {
    return nil, fmt.Errorf("executor.ExecCmdCtxWithOut: %w", err)
  }
  return buf, nil
}
