package app

import (
  "context"
  "fmt"
  "net"
  "net/http"
  "sync"
  "syscall"

  "github.com/go-chi/chi/v5"
  log "github.com/sirupsen/logrus"
  "github.com/ushakovn/boiler/pkg/closer"
  "github.com/ushakovn/boiler/pkg/gqlgen"
  "github.com/ushakovn/boiler/pkg/tracing"
  "google.golang.org/grpc"
  "google.golang.org/grpc/reflection"
)

type App struct {
  once sync.Once
  mu   sync.Mutex

  grpcPort   int
  gqlgenPort int

  grpcServer   *grpc.Server
  gqlgenRouter chi.Router

  appCtx    context.Context
  appCloser closer.Closer
}

func NewApp(calls ...Option) *App {
  // Call all options
  options := callAppOptions(calls...)

  // Create grpc server with options
  grpcServer := grpc.NewServer(options.grpcServerOptions...)

  // Create gqlgen router with middlewares
  gqlgenRouter := chi.NewRouter().With(options.gqlgenMiddlewares...)

  // Create app context
  appCtx := context.Background()

  // Create app closer
  appCloser := closer.NewCloser(syscall.SIGTERM, syscall.SIGKILL)

  // Return app
  return &App{
    grpcPort:   options.grpcServePort,
    gqlgenPort: options.gqlgenServePort,

    grpcServer:   grpcServer,
    gqlgenRouter: gqlgenRouter,

    appCtx:    appCtx,
    appCloser: appCloser,
  }
}

func (a *App) Run(services ...Service) {
  a.once.Do(func() {
    a.registerApp(a.registerParams(), services...)

    log.Infof("boiler: app bootstrapped")

    a.waitAppShutdown()
  })
}

func (a *App) waitAppShutdown() {
  a.appCloser.WaitAll()
}

func (a *App) registerApp(params *RegisterParams, services ...Service) {
  a.registerServices(params, services...)
  a.registerServicesComponents(params, services...)
}

func (a *App) registerParams() *RegisterParams {
  return &RegisterParams{
    grpcServer: a.grpcServer,
  }
}

func (a *App) registerServices(params *RegisterParams, services ...Service) {
  for _, service := range services {
    service.RegisterService(params)
  }

  log.Infof("boiler: app services registered")
}

func (a *App) registerServicesComponents(params *RegisterParams, _ ...Service) {
  // Collect service types
  serviceTypes := params.serviceTypesValues()

  // gRPC components
  if _, ok := serviceTypes[GrpcServiceTyp]; ok {
    a.registerGrpcServer()
  }
  // GraphQL components
  if _, ok := serviceTypes[GqlgenServiceTyp]; ok {
    a.registerGqlgenSchemaServer(params)
    a.registerGqlgenSandbox()
    a.registerGqlgenServer()
  }
  // Tracing components
  a.registerTracer()

  log.Infof("boiler: app services components registered")
}

func (a *App) registerGrpcServer() {
  address := fmt.Sprint("localhost", ":", a.grpcPort)

  lister, err := net.Listen("tcp", address)
  if err != nil {
    log.Fatalf("boiler: register gprc failed: %v", err)
  }
  reflection.Register(a.grpcServer)

  log.Infof("boiler: grpc server running on port: %d", a.grpcPort)

  go func() {
    if err = a.grpcServer.Serve(lister); err != nil {
      log.Fatalf("boiler: grpc server run failed: %v", err)
      a.appCloser.CloseAll()
    }
  }()
}

func (a *App) registerGqlgenServer() {
  address := fmt.Sprint("localhost", ":", a.gqlgenPort)

  log.Infof("boiler: gqlgen server running on port: %d", a.gqlgenPort)

  go func() {
    if err := http.ListenAndServe(address, a.gqlgenRouter); err != nil {
      log.Errorf("boiler: gqlgen server run failed: %v", err)
      a.appCloser.CloseAll()
    }
  }()
}

func (a *App) registerGqlgenSchemaServer(params *RegisterParams) {
  if params.gqlgenSchemaServer == nil {
    panic("boiler: gqlgen schema server is a nil")
  }
  a.gqlgenRouter.Handle("/query", *params.gqlgenSchemaServer)
}

func (a *App) registerGqlgenSandbox() {
  const (
    title    = "Boiler"
    endpoint = "/query"
  )
  sandbox := gqlgen.SandboxHandler(title, endpoint)
  a.gqlgenRouter.Handle("/", sandbox)

  log.Infof("boiler: gqlgen sanbox registered for: %s endpoint", endpoint)
}

// GqlgenRouter MUTATE APP ROUTER IN YOUR OWN RISK
func (a *App) GqlgenRouter() chi.Router {
  a.mu.Lock()
  defer a.mu.Unlock()

  return a.gqlgenRouter
}

func (a *App) registerTracer() {
  const (
    serviceName = "Boiler"
    serviceVer  = "v0.0.1"
  )
  shutdowns := tracing.InitTracer(a.appCtx, serviceName, serviceVer)

  log.Infof("boiler: tracing registered")

  a.appCloser.Add(shutdowns...)
}
