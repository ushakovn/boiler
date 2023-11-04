package app

import "google.golang.org/grpc"

type Service interface {
  Register(params *RegisterParams)
}

type RegisterParams struct {
  GrpcServer *grpc.Server
}
