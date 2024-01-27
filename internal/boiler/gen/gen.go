package gen

import (
  "context"
  "fmt"

  "github.com/ushakovn/boiler/internal/pkg/gens/config"
  "github.com/ushakovn/boiler/internal/pkg/gens/gqlgen"
  "github.com/ushakovn/boiler/internal/pkg/gens/grpc"
  "github.com/ushakovn/boiler/internal/pkg/gens/protodeps"
  "github.com/ushakovn/boiler/internal/pkg/gens/rpc"
  "github.com/ushakovn/boiler/internal/pkg/gens/storage"
)

type Generator interface {
  Generate(ctx context.Context) error
}

func NewGenerator(cfg any) (Generator, error) {
  var (
    g   Generator
    err error
  )
  switch c := cfg.(type) {
  case rpc.Config:
    g, err = rpc.NewRpc(c)
  case storage.ConfigPath:
    g, err = storage.NewStorage(c)
  case grpc.Config:
    g, err = grpc.NewGrpc(c)
  case gqlgen.Config:
    g, err = gqlgen.NewGqlgen(c)
  case protodeps.Config:
    g, err = protodeps.NewProtoDeps(c)
  case config.Config:
    g, err = config.NewGenConfig(c)
  default:
    err = fmt.Errorf("unsupported generator type")
  }
  return g, err
}
