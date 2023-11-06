package app

import "google.golang.org/grpc"

type ServiceType int

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
  Register(params *RegisterParams)
}

type RegisterParams struct {
  serviceTypes []*ServiceType
  grpcServer   *grpc.Server
}

func (p *RegisterParams) SetServiceType(serviceType ServiceType) {
  if _, ok := knownServiceTypes[serviceType]; !ok {
    panic("unknown service type")
  }
  p.serviceTypes = append(p.serviceTypes, &serviceType)
}

func (p *RegisterParams) GrpcServiceRegistrar() grpc.ServiceRegistrar {
  if p.grpcServer == nil {
    panic("grpc server is a nil")
  }
  return p.grpcServer
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
