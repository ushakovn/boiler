package gen

import (
  "context"
  "fmt"

  "github.com/ushakovn/boiler/internal/pkg/gens/clients"
  "github.com/ushakovn/boiler/internal/pkg/gens/gqlgen"
  "github.com/ushakovn/boiler/internal/pkg/gens/grpc"
  "github.com/ushakovn/boiler/internal/pkg/gens/rpc"
  "github.com/ushakovn/boiler/internal/pkg/gens/storage"
)

type Generator interface {
  Generate(ctx context.Context) error
}

func NewGenerator(config any) (Generator, error) {
  var (
    g   Generator
    err error
  )
  switch c := config.(type) {
  case rpc.Config:
    g, err = rpc.NewRpc(c)
  case storage.Config:
    g, err = storage.NewStorage(c)
  case grpc.Config:
    g, err = grpc.NewGrpc(c)
  case gqlgen.Config:
    g, err = gqlgen.NewGqlgen(c)
  case clients.Config:
    g, err = clients.NewClients(c)
  default:
    err = fmt.Errorf("unsupported generator type")
  }
  return g, err
}
