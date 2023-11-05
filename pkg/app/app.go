package app

import (
  "net"
  "syscall"

  "github.com/gin-gonic/gin"
  log "github.com/sirupsen/logrus"
  "github.com/ushakovn/boiler/pkg/closer"
  "google.golang.org/grpc"
  "google.golang.org/grpc/reflection"
)

type App struct {
  grpcPort int
  httpPort int

  grpcServer *grpc.Server
  httpServer *gin.Engine

  appCloser closer.Closer
}

func NewApp(calls ...Option) *App {
  options := callAppOptions(calls...)

  grpcServer := grpc.NewServer(options.grpcOptions...)

  httpServer := gin.New()
  httpServer.Use(options.ginHandles...)

  appCloser := closer.NewCloser(syscall.SIGTERM, syscall.SIGKILL)

  return &App{
    grpcServer: grpcServer,
    httpServer: httpServer,
    appCloser:  appCloser,
  }
}

func (a *App) Run(services ...Service) {
  a.registerServices(services...)
  a.registerGrpc()
}

func (a *App) registerGrpc() {
  const (
    defaultNetwork = "tcp"
    defaultHost    = "localhost:82"
  )
  if a.grpcServer != nil {
    lister, err := net.Listen(defaultNetwork, defaultHost)
    if err != nil {
      log.Fatalf("boiler: register gprc failed: %v", err)
    }
    reflection.Register(a.grpcServer)

    if err = a.grpcServer.Serve(lister); err != nil {
      log.Fatalf("Boiler: grpc server serve failed: %v", err)
    }
  }
}

func (a *App) registerServices(services ...Service) {
  params := &RegisterParams{
    GrpcServer: a.grpcServer,
  }
  for _, service := range services {
    service.Register(params)
  }
}
