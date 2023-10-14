package sql

import (
  "encoding/json"
  "fmt"
  "os"
  "os/exec"

  "github.com/ushakovn/boiler/pkg/utils"
  "gopkg.in/yaml.v3"
)

type pgDumpOptionTyp string

const (
  WithPgDumpFile   pgDumpOptionTyp = "pg_dump"
  WithPgConfigFile pgDumpOptionTyp = "pg_config"
)

type PgDumpOption interface {
  Call() (PgDumpBuf []byte, err error)
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

func (option *withPgDumpFile) Call() ([]byte, error) {
  pgDumpBuf, err := os.ReadFile(option.filePath)
  if err != nil {
    return nil, fmt.Errorf("os.ReadFile filePath: %w", err)
  }
  return pgDumpBuf, nil
}

type withPgConfigFile struct {
  filePath string
}

func (option *withPgConfigFile) Call() ([]byte, error) {
  buf, err := os.ReadFile(option.filePath)
  if err != nil {
    return nil, fmt.Errorf("os.ReadFile filePath: %w", err)
  }
  fileExtension := utils.ExtractFileExtension(option.filePath)

  config, err := parsePgConfig(fileExtension, buf)
  if err != nil {
    return nil, fmt.Errorf("parsePgConfig: %w", err)
  }
  pgDumpBuf, err := execPgDump(config)
  if err != nil {
    return nil, fmt.Errorf("execPgDump: %w", err)
  }
  return pgDumpBuf, nil
}

func parsePgConfig(fileExtension string, buf []byte) (PgConfig, error) {
  type wrappedConfig struct {
    PgConfig PgConfig `json:"pg_config" yaml:"pg_config"`
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
  return config.PgConfig, err
}

type PgConfig struct {
  DB       string `json:"db" yaml:"db"`
  Host     string `json:"host" yaml:"host"`
  Port     string `json:"port" yaml:"port"`
  User     string `json:"user" yaml:"user"`
  Password string `json:"password" yaml:"password"`
}

func (c *PgConfig) pgDumpCmd() (name string, args []string, err error) {
  if err := os.Setenv("PGPASSWORD", c.Password); err != nil {
    return "", nil, fmt.Errorf("os.Setenv: %w", err)
  }
  args = []string{
    "--no-owner",
    "-h", c.Host,
    "-p", c.Port,
    "-U", c.User,
    c.DB,
  }
  name = "pg_dump"

  return name, args, nil
}

func execPgDump(config PgConfig) ([]byte, error) {
  name, args, err := config.pgDumpCmd()
  if err != nil {
    return nil, fmt.Errorf("conn.pgDump: %w", err)
  }
  cmd := exec.Command(name, args...)
  if cmd.Err != nil {
    return nil, fmt.Errorf("exec.CommandContext %w", cmd.Err)
  }
  buf, err := cmd.Output()
  if err != nil {
    return nil, fmt.Errorf("cmd.Output: %w", err)
  }
  return buf, nil
}