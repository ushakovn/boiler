package storage

import (
  "fmt"
  "os"

  validation "github.com/go-ozzo/ozzo-validation"
  "github.com/ushakovn/boiler/internal/pkg/filer"
  "github.com/ushakovn/boiler/internal/pkg/validator"
  "gopkg.in/yaml.v3"
)

// ConfigPath USE IT FOR GENERATORS FACTORY
type ConfigPath string

// Config USE ConfigPath FOR GENERATORS FACTORY INSTEAD OF IT
type Config struct {
  PgConfig     *PgConfig                `yaml:"pg_config"`
  PgDumpPath   string                   `yaml:"pg_dump_path"`
  PgTypeConfig map[string]*PgTypeConfig `yaml:"pg_type_config"`
}

type PgConfig struct {
  Host     string `yaml:"host"`
  Port     string `yaml:"port"`
  User     string `yaml:"user"`
  DBName   string `yaml:"db_name"`
  Password string `yaml:"password"`
}

type PgTypeConfig struct {
  GoType     string `yaml:"go_type"`
  GoZeroType string `yaml:"go_zero_type"`
}

func (c ConfigPath) String() string {
  return string(c)
}

func (c ConfigPath) Parse() (*Config, error) {
  path := c.String()

  if !filer.IsExistedFile(path) {
    return nil, fmt.Errorf("file not found")
  }
  buf, err := os.ReadFile(path)
  if err != nil {
    return nil, fmt.Errorf("os.ReadFile: %w", err)
  }
  type wrapped struct {
    Config *Config `yaml:"storage_config"`
  }
  wrap := &wrapped{}

  if err = yaml.Unmarshal(buf, wrap); err != nil {
    return nil, fmt.Errorf("yaml.Unmarshal: %w", err)
  }
  return wrap.Config, nil
}

func (c *Config) Validate() error {
  if !validator.OneOfNotZero(c.PgConfig, c.PgDumpPath) {
    return fmt.Errorf("only one of fields must be specified: pg_config, pg_dump_path")
  }

  if c.PgConfig != nil {
    if err := c.PgConfig.Validate(); err != nil {
      return fmt.Errorf("pg_config invalid: %w", err)
    }
  }
  if c.PgDumpPath != "" {
    if !filer.IsExistedFile(c.PgDumpPath) {
      return fmt.Errorf("pg_dump_path: file not found")
    }
  }

  if err := validation.ValidateStruct(c,
    validation.Field(&c.PgTypeConfig, validation.Each(validation.Required)),
  ); err != nil {
    return fmt.Errorf("pg_type_config invalid: %w", err)
  }

  return nil
}

func (c *PgConfig) Validate() error {
  return validation.ValidateStruct(c,
    validation.Field(&c.Host, validation.Required),
    validation.Field(&c.Port, validation.Required),
    validation.Field(&c.User, validation.Required),
    validation.Field(&c.DBName, validation.Required),
    // Password is optionally field
  )
}

func (c *PgTypeConfig) Validate() error {
  return validation.ValidateStruct(c,
    // Prefix is optionally field
    validation.Field(&c.GoType, validation.Required),
    validation.Field(&c.GoZeroType, validation.Required),
  )
}
