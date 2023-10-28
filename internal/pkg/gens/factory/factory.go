package factory

import (
  "fmt"

  "github.com/ushakovn/boiler/internal/boiler/gen"
  "github.com/ushakovn/boiler/internal/pkg/gens/gqlgen"
  "github.com/ushakovn/boiler/internal/pkg/gens/project"
  "github.com/ushakovn/boiler/internal/pkg/gens/rpc"
  "github.com/ushakovn/boiler/internal/pkg/gens/storage"
)

type Typ string

const (
  ProjectTyp Typ = "project"
  RpcType    Typ = "rpc"
  StorageTyp Typ = "storage"
  GqlgenTyp  Typ = "gqlgen"
)

type Generators []gen.Generator

type CommonConfig struct {
  Project project.Config
  Rpc     rpc.Config
  Storage storage.Config
  Gqlgen  gqlgen.Config
}

func NewGenerator(config CommonConfig, typ Typ) (gen.Generator, error) {
  var (
    g   gen.Generator
    err error
  )
  switch typ {
  case ProjectTyp:
    g, err = project.NewProject(config.Project)
  case RpcType:
    g, err = rpc.NewRpc(config.Rpc)
  case StorageTyp:
    g, err = storage.NewStorage(config.Storage)
  case GqlgenTyp:
    g, err = gqlgen.NewGqlgen(config.Gqlgen)
  default:
    err = fmt.Errorf("unsupported generator type")
  }
  return g, err
}
