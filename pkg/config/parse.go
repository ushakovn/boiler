package config

import (
  "fmt"
  "os"
  "path/filepath"

  "github.com/spf13/cast"
  "github.com/ushakovn/boiler/internal/pkg/filer"
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

func collectAppInfo(appSection *AppSection) AppInfo {
  return AppInfo{
    Name:        appSection.Name,
    Version:     appSection.Version,
    Description: appSection.Description,
  }
}

func collectConfigValues(customSection CustomSection) (configValues, error) {
  values := configValues{}

  for csKey, csVal := range customSection {
    castedValue, err := castConfigValue(csVal.Type, csVal.Value)
    if err != nil {
      return nil, fmt.Errorf("cast config value failed: %w", err)
    }
    stringKey := csKey.String()

    values[stringKey] = &configValue{
      value: castedValue,
    }
  }
  return values, nil
}

func castConfigValue(typ, rawValue string) (any, error) {
  var (
    value any
    err   error
  )
  switch typ {
  case "int":
    value, err = cast.ToIntE(rawValue)
  case "int32":
    value, err = cast.ToInt32E(rawValue)
  case "int64":
    value, err = cast.ToInt64E(rawValue)
  case "float32":
    value, err = cast.ToFloat32E(rawValue)
  case "float64":
    value, err = cast.ToFloat64E(rawValue)
  case "uint32":
    value, err = cast.ToUint32E(rawValue)
  case "uint64":
    value, err = cast.ToUint64E(rawValue)
  case "string":
    value, err = cast.ToStringE(rawValue)
  case "bool":
    value, err = cast.ToBoolE(rawValue)
  case "time":
    value, err = cast.ToTimeE(rawValue)
  case "duration":
    value, err = cast.ToDurationE(rawValue)
  default:
    err = fmt.Errorf("unexpected config value type: %s", typ)
  }
  return value, err
}
