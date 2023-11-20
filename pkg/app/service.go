package app

import (
  "github.com/99designs/gqlgen/graphql"
  "google.golang.org/grpc"
)

type ServiceType uint32

const (
  UnknownServiceTyp ServiceType = 0
  GrpcServiceTyp    ServiceType = 1
  GqlgenServiceTyp  ServiceType = 2
)

var knownServiceTypes = map[ServiceType]struct{}{
  GrpcServiceTyp:   {},
  GqlgenServiceTyp: {},
}

type Service interface {
  RegisterService(params *RegisterParams)
}

type RegisterParams struct {
  // Services types
  serviceTypes []*ServiceType
  // gRPC
  grpcServer *grpc.Server
  // GraphQL
  gqlgenSchema graphql.ExecutableSchema
}

func (p *RegisterParams) SetServiceType(serviceType ServiceType) {
  if _, ok := knownServiceTypes[serviceType]; !ok {
    panic("boiler: unknown service type")
  }
  p.serviceTypes = append(p.serviceTypes, &serviceType)
}

func (p *RegisterParams) GrpcServiceRegistrar() grpc.ServiceRegistrar {
  if p.grpcServer == nil {
    panic("boiler: grpc server is a nil")
  }
  return p.grpcServer
}

func (p *RegisterParams) SetGqlgenExecutableSchema(schema graphql.ExecutableSchema) {
  p.gqlgenSchema = schema
}

func (p *RegisterParams) serviceTypesValues() map[ServiceType]struct{} {
  var typeValue ServiceType
  typesValues := map[ServiceType]struct{}{}

  for _, serviceType := range p.serviceTypes {
    if serviceType == nil {
      typeValue = UnknownServiceTyp
    } else {
      typeValue = *serviceType
    }
    if _, ok := typesValues[typeValue]; ok {
      continue
    }
    typesValues[typeValue] = struct{}{}
  }

  return typesValues
}
