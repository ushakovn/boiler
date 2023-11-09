package gqlgen

type gqlgenServiceDesc struct {
  ServicePackages []*goPackageDesc
}

type goPackageDesc struct {
  CustomName  string
  ImportLine  string
  ImportAlias string
  IsBuiltin   bool
  IsInstall   bool
}

var gqlgenServicePackages = []*goPackageDesc{
  {
    CustomName:  "boiler/app",
    ImportLine:  "github.com/ushakovn/boiler/pkg/app",
    ImportAlias: "",
    IsInstall:   true,
  },
}

func buildGqlgenServiceDesc() *gqlgenServiceDesc {
  return &gqlgenServiceDesc{
    ServicePackages: gqlgenServicePackages,
  }
}
