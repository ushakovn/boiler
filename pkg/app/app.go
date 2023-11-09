package app

import (
  "fmt"
  "net"
  "sync"
  "syscall"

  "github.com/gin-gonic/gin"
  log "github.com/sirupsen/logrus"
  "github.com/ushakovn/boiler/pkg/closer"
  "github.com/ushakovn/boiler/pkg/gqlgen"
  "google.golang.org/grpc"
  "google.golang.org/grpc/reflection"
)

type App struct {
  once sync.Once
  mu   sync.Mutex

  grpcPort   int
  gqlgenPort int

  grpcServer   *grpc.Server
  gqlgenServer *gin.Engine

  appCloser closer.Closer
}

func NewApp(calls ...Option) *App {
  // Call all options
  options := callAppOptions(calls...)

  // Create grpc server with options
  grpcServer := grpc.NewServer(options.grpcServerOptions...)

  gin.SetMode(gin.ReleaseMode)
  // Create gqlgen server with handlers
  gqlgenServer := gin.New()

  gqlgenServer.Use(gin.Recovery())
  gqlgenServer.Use(options.gqlgenHandlers...)

  // Create app closer
  appCloser := closer.NewCloser(syscall.SIGTERM, syscall.SIGKILL)

  // Return app
  return &App{
    grpcPort:   options.grpcServePort,
    gqlgenPort: options.gqlgenServePort,

    grpcServer:   grpcServer,
    gqlgenServer: gqlgenServer,

    appCloser: appCloser,
  }
}

func (a *App) Run(services ...Service) {
  a.once.Do(func() {
    a.registerServices(services...)
    log.Infof("boiler: app bootstrapped")
    a.waitServicesShutdown()
  })
}

func (a *App) waitServicesShutdown() {
  a.appCloser.WaitAll()
}

func (a *App) registerServices(services ...Service) {
  params := &RegisterParams{
    grpcServer: a.grpcServer,
  }
  for _, service := range services {
    service.RegisterService(params)
  }
  log.Infof("boiler: app services registered")

  serviceTypes := params.serviceTypesValues()

  if _, ok := serviceTypes[GrpcServiceTyp]; ok {
    a.registerGrpc()
  }
  if _, ok := serviceTypes[GqlgenServiceTyp]; ok {
    a.registerGqlgen()
  }
  log.Infof("boiler: app servers registered")
}

func (a *App) registerGrpc() {
  address := fmt.Sprint("localhost", ":", a.grpcPort)

  lister, err := net.Listen("tcp", address)
  if err != nil {
    log.Fatalf("boiler: register gprc failed: %v", err)
  }
  reflection.Register(a.grpcServer)

  go func() {
    log.Infof("boiler: grpc server running on port: %d", a.grpcPort)

    if err = a.grpcServer.Serve(lister); err != nil {
      log.Fatalf("boiler: grpc server run failed: %v", err)
      a.appCloser.CloseAll()
    }
  }()
}

func (a *App) registerGqlgen() {
  gqlgen.SandboxHandler("boiler", "/")
  address := fmt.Sprint("localhost", ":", a.gqlgenPort)

  go func() {
    log.Infof("boiler: gqlgen server running on port: %d", a.gqlgenPort)

    if err := a.gqlgenServer.Run(address); err != nil {
      log.Errorf("boiler: gqlgen server run failed: %v", err)
      a.appCloser.CloseAll()
    }
  }()
}

func (a *App) GqlgenServer() *gin.Engine {
  a.mu.Lock()
  defer a.mu.Unlock()

  return a.gqlgenServer
}
