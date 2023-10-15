package factory

import (
  "fmt"

  "github.com/ushakovn/boiler/internal/boiler/gen"
  "github.com/ushakovn/boiler/internal/pkg/gens/project"
  "github.com/ushakovn/boiler/internal/pkg/gens/rpc"
  "github.com/ushakovn/boiler/internal/pkg/gens/storage"
)

type Typ string

const (
  ProjectTyp Typ = "project"
  RpcType    Typ = "rpc"
  StorageTyp Typ = "storage"
)

type Generators []gen.Generator

type CommonConfig struct {
  Project project.Config
  Rpc     rpc.Config
  Storage storage.Config
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
  default:
    err = fmt.Errorf("unsupported generator type")
  }
  return g, err
}
