package gqlgen

import "path/filepath"

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
    CustomName:  "gqlgen/handler",
    ImportLine:  "github.com/99designs/gqlgen/graphql/handler",
    ImportAlias: "",
    IsInstall:   true,
  },
  {
    CustomName:  "boiler/app",
    ImportLine:  "github.com/ushakovn/boiler/pkg/app",
    ImportAlias: "",
    IsInstall:   true,
  },
}

func (g *Gqlgen) buildGqlgenServiceDesc() *gqlgenServiceDesc {
  return &gqlgenServiceDesc{
    ServicePackages: g.buildGqlgenServicePackages(),
  }
}

func (g *Gqlgen) buildGqlgenServicePackages() []*goPackageDesc {
  servicePackages := make([]*goPackageDesc, 0, len(gqlgenServicePackages)+1)

  servicePackages = append(servicePackages, gqlgenServicePackages...)
  servicePackages = append(servicePackages, g.buildGqlgenGeneratedPackage())

  return servicePackages
}

func (g *Gqlgen) buildGqlgenGeneratedPackage() *goPackageDesc {
  packagePath := filepath.Join(g.goModuleName, "internal", "app", "graph", "generated")

  return &goPackageDesc{
    CustomName:  "gqlgen/generated",
    ImportAlias: "",
    ImportLine:  packagePath,
  }
}
