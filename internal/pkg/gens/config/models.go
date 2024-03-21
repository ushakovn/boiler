package config

import (
  "fmt"
  "path/filepath"
  "strings"

  "github.com/ushakovn/boiler/pkg/config"
  "golang.org/x/text/cases"
  "golang.org/x/text/language"
)

type genConfigDesc struct {
  ConfigGroups     []*groupDesc
  ConfigPackages   []*goPackageDesc
  GroupsPackages   []*goPackageDesc
  ProviderPackages []*goPackageDesc
}

type groupDesc struct {
  GroupName string
  GroupKeys []*groupKeyDesc
}

type groupKeyDesc struct {
  KeyName     string
  KeyNameTrim string
  KeyComment  string
  ValueType   string
  ValueCall   string
}

type goPackageDesc struct {
  CustomName  string
  ImportLine  string
  ImportAlias string
  IsBuiltin   bool
  IsInstall   bool
}

func (g *GenConfig) loadGenConfigDesc() (*genConfigDesc, error) {
  configPath := filepath.Join(".config", "app_config.yaml")

  parsed, err := config.ParseConfig(configPath)
  if err != nil {
    return nil, fmt.Errorf("config parsing failed:\n%v", err)
  }
  if err = parsed.Validate(); err != nil {
    return nil, fmt.Errorf("config validation failed:\n%v", err)
  }
  genConfig := buildGenConfig(parsed.Custom)

  return genConfig, nil
}

func buildGenConfig(customSection config.CustomSection) *genConfigDesc {
  configGroups := buildConfigGroups(customSection)

  return &genConfigDesc{
    ConfigGroups:     configGroups,
    ConfigPackages:   buildPackagesByNames(configPackages...),
    GroupsPackages:   buildPackagesByNames(groupsPackages...),
    ProviderPackages: buildPackagesByNames(providerPackages...),
  }
}

func buildConfigGroups(customSection config.CustomSection) []*groupDesc {
  groupsKeys := map[string][]*groupKeyDesc{}

  for key, val := range customSection {
    keyName := key.String()
    grKeys := groupsKeys[val.Group]

    keyNameTrim := buildGroupKeyNameTrim(keyName, val.Group)
    keyComment := strings.TrimSpace(val.Description)

    valueType := buildGroupKeyValueType(val.Type)
    valueCall := buildGroupKeyValueCall(val.Type)

    grKeys = append(grKeys, &groupKeyDesc{
      KeyName:     keyName,
      KeyNameTrim: keyNameTrim,
      KeyComment:  keyComment,
      ValueType:   valueType,
      ValueCall:   valueCall,
    })
    groupsKeys[val.Group] = grKeys
  }
  configGroups := make([]*groupDesc, 0, len(groupsKeys))

  for groupName, groupKeys := range groupsKeys {
    configGroups = append(configGroups, &groupDesc{
      GroupName: groupName,
      GroupKeys: groupKeys,
    })
  }
  return configGroups
}

func buildGroupKeyNameTrim(key, group string) string {
  trim := strings.TrimPrefix(key, group)
  trim = strings.Trim(trim, "_- ")
  return trim
}

func buildGroupKeyValueType(typ string) string {
  if packageTyp, ok := valueTypeToPackageType[typ]; ok {
    return packageTyp
  }
  return typ
}

func buildGroupKeyValueCall(typ string) string {
  return cases.Title(language.Und, cases.NoLower).String(typ)
}

func buildPackagesByNames(names ...string) []*goPackageDesc {
  packages := make([]*goPackageDesc, 0, len(names))

  for _, name := range names {
    pack, ok := packagesByNames[name]
    if !ok {
      continue
    }
    packages = append(packages, pack)
  }
  return packages
}

var valueTypeToPackageType = map[string]string{
  "time":     "time.Time",
  "duration": "time.Duration",
}

var configPackages = []string{
  contextPackageName,
  configPackageName,
  typesPackageName,
}

var groupsPackages = []string{
  contextPackageName,
  configPackageName,
  typesPackageName,
}

var providerPackages = []string{
  contextPackageName,
  atomicPackageName,
  typesPackageName,
}

const (
  timePackageName    = "time"
  contextPackageName = "context"
  atomicPackageName  = "atomic"
  configPackageName  = "config"
  typesPackageName   = "types"
)

var packagesByNames = map[string]*goPackageDesc{
  timePackageName: {
    CustomName: "go/time",
    ImportLine: "time",
    IsBuiltin:  true,
  },
  contextPackageName: {
    CustomName: "go/context",
    ImportLine: "context",
    IsBuiltin:  true,
  },
  atomicPackageName: {
    CustomName: "go/atomic",
    ImportLine: "sync/atomic",
    IsBuiltin:  true,
  },
  configPackageName: {
    CustomName: "boiler/config",
    ImportLine: "github.com/ushakovn/boiler/pkg/config",
    IsInstall:  true,
  },
  typesPackageName: {
    CustomName: "boiler/config-types",
    ImportLine: "github.com/ushakovn/boiler/pkg/config/types",
    IsInstall:  true,
  },
}
