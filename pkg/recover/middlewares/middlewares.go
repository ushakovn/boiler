package middlewares

import (
  "context"

  "github.com/99designs/gqlgen/graphql"
  log "github.com/sirupsen/logrus"
  "google.golang.org/grpc"
)

func GrpcServerUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
  defer func() {
    if rec := recover(); rec != nil {
      // Log gRPC call info
      log.Errorf("boiler: panic during grpc call handling: %s", info.FullMethod)
      // Log panic info
      log.Errorf("boiler: panic recovered: %v", rec)
    }
  }()

  return handler(ctx, req)
}

func GqlgenOperationMiddleware(ctx context.Context, handler graphql.OperationHandler) graphql.ResponseHandler {
  defer func() {
    if rec := recover(); rec != nil {
      // Log GraphQL operation info
      if graphql.HasOperationContext(ctx) {
        opCtx := graphql.GetOperationContext(ctx)
        log.Errorf("boiler: panic during graphql operation handling: %s", opCtx.OperationName)
      }
      // Log panic info
      log.Errorf("boiler: panic recovered: %v", rec)
    }
  }()

  return handler(ctx)
}