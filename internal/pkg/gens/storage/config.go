package storage

import (
  "fmt"
  "os"

  "github.com/samber/lo"
  "github.com/ushakovn/boiler/internal/pkg/filer"
  "gopkg.in/yaml.v3"
)

// ConfigPath USE IT FOR GENERATORS FACTORY
type ConfigPath string

// Config USE ConfigPath FOR GENERATORS FACTORY INSTEAD OF IT
type Config struct {
  PgConfig      *PgConfig                `yaml:"pg_config"`
  PgDumpPath    string                   `yaml:"pg_dump_path"`
  PgTableConfig *PgTableConfig           `yaml:"pg_table_config"`
  PgTypeConfig  map[string]*PgTypeConfig `yaml:"pg_type_config"`
}

type PgConfig struct {
  Host     string `yaml:"host"`
  Port     string `yaml:"port"`
  User     string `yaml:"user"`
  DBName   string `yaml:"db_name"`
  Password string `yaml:"password"`
}

type PgTableConfig struct {
  PgColumnFilter       *PgColFilter `yaml:"pg_column_filter"`
  PgSkipTables         []string     `yaml:"pg_skip_tables"`
  PgSkipCustomStorages []string     `yaml:"pg_skip_custom_storages"`
}

type PgColFilter struct {
  AllByDefault bool                           `yaml:"all_by_default"`
  String       []string                       `yaml:"string"`
  Numeric      []string                       `yaml:"numeric"`
  Overrides    map[string]map[string][]string `yaml:"overrides"`
}

type PgTypeConfig struct {
  GoType     string `yaml:"go_type"`
  GoZeroType string `yaml:"go_zero_type"`
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

func (c *Config) WithInitial() *Config {
  pg := c.PgConfig

  if pg == nil {
    c.PgConfig = &PgConfig{}
  }
  table := c.PgTableConfig

  if table == nil {
    filters := &PgColFilter{
      String:    make([]string, 0),
      Numeric:   make([]string, 0),
      Overrides: make(map[string]map[string][]string),
    }
    table = &PgTableConfig{
      PgColumnFilter: filters,
    }
  }
  filters := table.PgColumnFilter

  if filters == nil {
    filters = &PgColFilter{
      String:    make([]string, 0),
      Numeric:   make([]string, 0),
      Overrides: make(map[string]map[string][]string),
    }
  }
  overrides := filters.Overrides

  if overrides == nil {
    overrides = make(map[string]map[string][]string)
  }
  types := c.PgTypeConfig

  if types == nil {
    types = make(map[string]*PgTypeConfig)
  }

  return c
}

func (c *Config) skipStorage(modelName string) bool {
  config := c.PgTableConfig.PgSkipCustomStorages
  return lo.Contains(config, modelName)
}

func (c *Config) skipTable(tableName string) bool {
  config := c.PgTableConfig.PgSkipTables
  return lo.Contains(config, tableName)
}

func (c ConfigPath) String() string {
  return string(c)
}
