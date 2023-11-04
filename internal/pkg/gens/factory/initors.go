package factory

import (
  "fmt"

  "github.com/ushakovn/boiler/internal/boiler/gen"
  "github.com/ushakovn/boiler/internal/pkg/gens/grpc"
)

func NewInitor(config CommonConfig, typ Typ) (gen.Initor, error) {
  var (
    i   gen.Initor
    err error
  )
  switch typ {
  case GrpcTyp:
    i, err = grpc.NewGrpc(config.Grpc)
  default:
    err = fmt.Errorf("unsupported initor type")
  }
  return i, err
}
