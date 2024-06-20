package app

import (
  "context"
  "fmt"

  "github.com/99designs/gqlgen/graphql"
  "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
  log "github.com/sirupsen/logrus"
  "google.golang.org/grpc"
)

type ServiceType uint32

const (
  UnknownServiceType ServiceType = 0
  GrpcServiceType    ServiceType = 1
  GqlgenServiceType  ServiceType = 2
  HttpServiceType    ServiceType = 3
)

var knownServiceTypes = map[ServiceType]struct{}{
  GrpcServiceType:   {},
  GqlgenServiceType: {},
  HttpServiceType:   {},
}

type Service interface {
  RegisterService(params *RegisterParams) error
}

type RegisterParams struct {
  // App context
  appCtx context.Context
  // Services types
  serviceTypes map[ServiceType]struct{}
  // gRPC params
  grpcParams *GrpcParams
  // GraphQL params
  gqlgenParams *GqlgenParams
}

type GrpcParams struct {
  // gRPC server
  grpcServer *grpc.Server
  // gRPC port
  grpcServerPort int
  // gRPC HTTP proxy dial options
  grpcClientOptions []grpc.DialOption
  // gRPC HTTP proxy serve mux
  grpcHttpProxyServeMux *runtime.ServeMux
  // Open API specification doc
  docOpenAPI []byte
}

func (p *RegisterParams) Context() context.Context {
  // Boiler app context
  return p.appCtx
}

func (p *RegisterParams) Grpc() *GrpcParams {
  return p.grpcParams
}

func (p *GrpcParams) GrpcServiceRegistrar() grpc.ServiceRegistrar {
  return p.grpcServer
}

func (p *GrpcParams) GrpcHttpProxyServeMux() *runtime.ServeMux {
  return p.grpcHttpProxyServeMux
}

func (p *GrpcParams) GrpcServerEndpoint() string {
  return fmt.Sprint("localhost", ":", p.grpcServerPort)
}

func (p *GrpcParams) GrpcClientOptions() []grpc.DialOption {
  return p.grpcClientOptions
}

func (p *GrpcParams) SetOpenAPIDoc(doc []byte) {
  p.docOpenAPI = doc
}

type GqlgenParams struct {
  // GraphQL schema
  gqlgenSchema graphql.ExecutableSchema
}

func (p *RegisterParams) Gqlgen() *GqlgenParams {
  return p.gqlgenParams
}

func (p *GqlgenParams) GqlgenSchema() graphql.ExecutableSchema {
  return p.gqlgenSchema
}

func (p *GqlgenParams) SetGqlgenSchema(schema graphql.ExecutableSchema) {
  if schema == nil {
    // GraphQL executable schema must be set
    log.Fatalf("boiler: gqlgen schema is a nil")
  }
  p.gqlgenSchema = schema
}

func (p *RegisterParams) SetServiceType(serviceType ServiceType) {
  if _, ok := knownServiceTypes[serviceType]; !ok {
    // Unknown service types not allowed
    log.Fatalf("boiler: unknown service type")
  }
  p.serviceTypes[serviceType] = struct{}{}
}
