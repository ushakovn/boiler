package config

import (
  "fmt"
  "os"
  "path/filepath"

  "github.com/ushakovn/boiler/internal/pkg/filer"
  "github.com/ushakovn/boiler/pkg/config/types"
  "gopkg.in/yaml.v3"
)

func findConfig(configPath string) bool {
  return filer.IsExistedFile(configPath)
}

func ParseConfig(configPath string) (*Parsed, error) {
  absConfigPath, err := filepath.Abs(configPath)
  if err != nil {
    return nil, fmt.Errorf("config file not found for path: %s", configPath)
  }
  configBuf, err := os.ReadFile(absConfigPath)
  if err != nil {
    return nil, fmt.Errorf("config file read failed: %w", err)
  }
  parsed := &Parsed{}

  if err = yaml.Unmarshal(configBuf, parsed); err != nil {
    return nil, fmt.Errorf("yaml unmarshal failed: %w", err)
  }
  return parsed, err
}

func collectAppInfo(as *AppSection) AppInfo {
  return AppInfo{
    Name:        as.Name,
    Version:     as.Version,
    Description: as.Description,
  }
}

func collectConfigValues(cs CustomSection) (configValues, error) {
  values := configValues{}

  for csKey, csVal := range cs {
    value, err := types.CastValue(csVal.Type, csVal.Value)
    if err != nil {
      return nil, fmt.Errorf("cast config value failed: %w", err)
    }
    strKey := csKey.String()

    values[strKey] = value
  }
  return values, nil
}
