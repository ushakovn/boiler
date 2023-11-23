package config

type genConfigDesc struct {
  ConfigGroups []*groupDesc
}

type groupDesc struct {
  GroupName string
  GroupKeys []*groupKeyDesc
}

type groupKeyDesc struct {
  KeyName    string
  ValueType  string
  KeyComment string
}

type goPackageDesc struct {
  CustomName  string
  ImportLine  string
  ImportAlias string
  IsBuiltin   bool
  IsInstall   bool
}

func (g *GenConfig) loadGenConfigDesc() (*genConfigDesc, error) {
  // TODO
  return nil, nil
}

var (
  // TODO
  _ = configPackage
  _ = groupsPackages
)

var configPackage = []*goPackageDesc{
  {
    CustomName: "go/context",
    ImportLine: "context",
    IsBuiltin:  true,
  },
  {
    CustomName: "go/time",
    ImportLine: "time",
    IsBuiltin:  true,
  },
}

var groupsPackages = []*goPackageDesc{
  {
    CustomName: "boiler/config",
    ImportLine: "github.com/ushakovn/boiler/pkg/config",
    IsInstall:  true,
  },
}
