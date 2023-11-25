package gen

import (
  "context"
  "fmt"

  "github.com/ushakovn/boiler/internal/pkg/gens/config"
  "github.com/ushakovn/boiler/internal/pkg/gens/gqlgen"
  "github.com/ushakovn/boiler/internal/pkg/gens/grpc"
  "github.com/ushakovn/boiler/internal/pkg/gens/project"
  "github.com/ushakovn/boiler/internal/pkg/gens/protodeps"
  "github.com/ushakovn/boiler/internal/pkg/gens/storage"
)

type Initor interface {
  Init(ctx context.Context) error
}

func NewInitor(cfg any) (Initor, error) {
  var (
    g   Initor
    err error
  )
  switch c := cfg.(type) {
  case project.Config:
    g, err = project.NewProject(c)
  case grpc.Config:
    g, err = grpc.NewGrpc(c)
  case gqlgen.Config:
    g, err = gqlgen.NewGqlgen(c)
  case protodeps.Config:
    g, err = protodeps.NewProtoDeps(c)
  case storage.Config:
    g, err = storage.NewStorage(c)
  case config.Config:
    g, err = config.NewGenConfig(c)
  default:
    err = fmt.Errorf("unsupported initor type")
  }
  return g, err
}
