package gen

import (
  "context"
  "fmt"

  "github.com/ushakovn/boiler/internal/pkg/gens/grpc"
  "github.com/ushakovn/boiler/internal/pkg/gens/project"
)

type Initor interface {
  Init(ctx context.Context) error
}

func NewInitor(config any) (Initor, error) {
  var (
    g   Initor
    err error
  )
  switch c := config.(type) {
  case project.Config:
    g, err = project.NewProject(c)
  case grpc.Config:
    g, err = grpc.NewGrpc(c)
  default:
    err = fmt.Errorf("unsupported initor type")
  }
  return g, err
}
