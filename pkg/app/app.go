package app

import (
  "context"
  "encoding/json"
  "fmt"
  "net"
  "net/http"
  "sync"
  "syscall"
  "time"

  "github.com/99designs/gqlgen/graphql"
  "github.com/99designs/gqlgen/graphql/handler"
  "github.com/go-chi/chi/v5"
  runtimeGrpc "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
  log "github.com/sirupsen/logrus"
  httpswagger "github.com/ushakovn/boiler/internal/pkg/thirdparty/swaggo/http-swagger"
  "github.com/ushakovn/boiler/pkg/closer"
  "github.com/ushakovn/boiler/pkg/config"
  "github.com/ushakovn/boiler/pkg/gqlgen"
  "github.com/ushakovn/boiler/pkg/grpcx/http-gateway"
  "github.com/ushakovn/boiler/pkg/logger"
  mw "github.com/ushakovn/boiler/pkg/metrics/middlewares"
  "github.com/ushakovn/boiler/pkg/tracing/tracer"
  "google.golang.org/grpc"
  "google.golang.org/grpc/reflection"

  metrics "github.com/ushakovn/boiler/pkg/metrics/handler"
)

type App struct {
  // Sync
  once sync.Once
  mu   sync.Mutex

  // gRPC
  grpcPort   int
  grpcServer *grpc.Server

  // gRPC HTTP proxy
  grpcHttpProxyPort int

  // GraphQL
  gqlgenPort   int
  gqlgenRouter chi.Router
  gqlgenServer *handler.Server

  gqlgenFieldMWs     []graphql.FieldMiddleware
  gqlgenOperationMWs []graphql.OperationMiddleware
  gqlgenResponseMWs  []graphql.ResponseMiddleware

  // Duty HTTP router
  dutyHttpPort   int
  dutyHttpRouter chi.Router

  // Shutdown
  appCtx    context.Context
  appCloser closer.Closer
}

func NewApp(calls ...Option) *App {
  // Set log options
  logger.SetDefaultLogOptions()

  // Call all options
  options := callAppOptions(calls...)

  // Build gRPC server options
  grpcServerOptions := buildGrpcServerOptions(options)

  // Create gRPC server with options
  grpcServer := grpc.NewServer(grpcServerOptions...)

  // Create GraphQL router with middlewares
  gqlgenRouter := chi.NewRouter().With(options.gqlgenMWs...)

  // Create duty HTTP router
  dutyHttpRouter := chi.NewRouter()

  // Create app context
  appCtx := context.Background()

  // Create app closer
  appCloser := closer.NewCloser(
    syscall.SIGTERM,
    syscall.SIGKILL,
    syscall.SIGINT,
  )

  // Registering pre run components
  registerPreRunComponents()

  // Return app
  return &App{
    grpcPort:   options.grpcServePort,
    grpcServer: grpcServer,

    grpcHttpProxyPort: options.grpcHttpProxyPort,

    gqlgenPort:   options.gqlgenServePort,
    gqlgenRouter: gqlgenRouter,

    gqlgenFieldMWs:     options.gqlgenFieldMWs,
    gqlgenOperationMWs: options.gqlgenOperationMWs,
    gqlgenResponseMWs:  options.gqlgenResponseMWs,

    dutyHttpPort:   options.dutyHttpServePort,
    dutyHttpRouter: dutyHttpRouter,

    appCtx:    appCtx,
    appCloser: appCloser,
  }
}

func registerPreRunComponents() {
  // Config components
  registerConfigClient()
}

func registerConfigClient() {
  config.InitClientConfig()
  log.Infof("boiler: config client registered")
}

func (a *App) Run(services ...Service) {
  defer func() {
    if rec := recover(); rec != nil {
      log.Errorf("boiler: app panic recovered: %v", rec)
    }
  }()

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
  grpcParams := &GrpcParams{
    grpcServer:            a.grpcServer,
    grpcServerPort:        a.grpcPort,
    grpcClientOptions:     defaultGrpcClientOptions(),
    grpcHttpProxyServeMux: runtimeGrpc.NewServeMux(),
  }
  gqlgenParams := &GqlgenParams{}

  return &RegisterParams{
    appCtx:       a.appCtx,
    grpcParams:   grpcParams,
    gqlgenParams: gqlgenParams,
  }
}

func (a *App) registerServices(params *RegisterParams, services ...Service) {
  for _, service := range services {
    if err := service.RegisterService(params); err != nil {
      log.Fatalf("boiler: app service registration failed: %v", err)
    }
  }
  log.Infof("boiler: app services registered")
}

func (a *App) registerServicesComponents(params *RegisterParams, _ ...Service) {
  // Collect service types
  serviceTypes := params.serviceTypesValues()

  // Confirm service types
  if _, ok := serviceTypes[UnknownServiceTyp]; ok {
    log.Fatalf("boiler: encountered unknown service type")
  }

  // gRPC components
  if _, ok := serviceTypes[GrpcServiceTyp]; ok {
    p := params.Grpc()
    a.registerGrpcServer()
    a.registerGrpcHttpProxyServer(p)
    a.registerGrpcSwagger(p)
  }

  // GraphQL components
  if _, ok := serviceTypes[GqlgenServiceTyp]; ok {
    p := params.Gqlgen()
    a.registerGqlgenSchemaServer(p)
    a.registerGqlgenAroundMWs()
    a.registerGqlgenSandbox()
    a.registerGqlgenServer()
  }

  // Observability components
  a.registerObservability()

  // Developer help components
  a.registerHelpHandler()

  // Run duty HTTP router
  a.runHttpDutyRouter()

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
      log.Errorf("boiler: grpc server run failed: %v", err)
      a.appCloser.CloseAll()
    }
  }()

  // Graceful shutdown for gRPC server
  a.appCloser.Add(func(ctx context.Context) error {
    log.Infof("boiler: grpc server trying graceful shutdown")

    const timeout = 5 * time.Second

    timeoutCtx, cancel := context.WithTimeout(a.appCtx, timeout)
    defer cancel()

    doneCh := make(chan struct{})

    go func() {
      a.grpcServer.GracefulStop()
      doneCh <- struct{}{}
    }()

    for {
      select {
      case <-timeoutCtx.Done():
        log.Infof("boiler: grpc server was not stopped for %s timeout", timeout.String())
        a.grpcServer.Stop()
        log.Infof("boiler: grpc server stopped forced")
        return nil

      case <-doneCh:
        log.Infof("boiler: grpc server stopped gracefully")
        return nil
      }
    }
  })
}

func (a *App) registerGrpcHttpProxyServer(params *GrpcParams) {
  mux := params.GrpcHttpProxyServeMux()
  if mux == nil {
    // gRPC proxy server was not set
    return
  }
  address := fmt.Sprint("localhost", ":", a.grpcHttpProxyPort)

  log.Infof("boiler: grpc http proxy running on port: %d", a.grpcHttpProxyPort)

  go func() {
    if err := http.ListenAndServe(address, mux); err != nil {
      log.Errorf("boiler: grpc http proxy server run failed: %v", err)
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

func (a *App) registerGqlgenSchemaServer(params *GqlgenParams) {
  a.gqlgenServer = handler.NewDefaultServer(params.GqlgenSchema())
  a.gqlgenRouter.Handle("/query", a.gqlgenServer)
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

func (a *App) registerGqlgenAroundMWs() {
  for _, aroundOperations := range a.gqlgenOperationMWs {
    a.gqlgenServer.AroundOperations(aroundOperations)
  }
  for _, aroundFields := range a.gqlgenFieldMWs {
    a.gqlgenServer.AroundFields(aroundFields)
  }
  for _, aroundFields := range a.gqlgenResponseMWs {
    a.gqlgenServer.AroundResponses(aroundFields)
  }
}

// GqlgenRouter USE IN YOUR OWN RISK
func (a *App) GqlgenRouter() chi.Router {
  a.mu.Lock()
  defer a.mu.Unlock()

  return a.gqlgenRouter
}

func (a *App) registerTracer() {
  info := config.ClientConfig(a.appCtx).GetAppInfo()
  shutdowns := tracer.InitTracer(a.appCtx, info.Name, info.Version)

  log.Infof("boiler: tracing registered")
  a.appCloser.Add(shutdowns...)
}

func (a *App) registerMetrics() {
  metrics.WithMetricsHandler(a.dutyHttpRouter)
  mw.InitMetrics()

  log.Infof("boiler: metrics handler registered")
}

func (a *App) registerObservability() {
  // Tracing components
  a.registerTracer()
  // Metrics components
  a.registerMetrics()
}

func (a *App) runHttpDutyRouter() {
  address := fmt.Sprint("localhost", ":", a.dutyHttpPort)

  log.Infof("boiler: http duty server running on port: %d", a.dutyHttpPort)

  go func() {
    if err := http.ListenAndServe(address, a.dutyHttpRouter); err != nil {
      log.Errorf("boiler: http duty server run failed: %v", err)
    }
    a.appCloser.CloseAll()
  }()
}

func (a *App) registerGrpcSwagger(params *GrpcParams) {
  name := config.ClientConfig(a.appCtx).GetAppInfo().Name
  url := fmt.Sprintf("http://localhost:%d/swagger/doc.json", a.grpcHttpProxyPort)

  swagger := httpswagger.Handler(
    httpswagger.InstanceName(name),
    httpswagger.URL(url),
    httpswagger.Doc(params.docOpenAPI),
  )
  mux := params.GrpcHttpProxyServeMux()

  if err := httpgateway.Handle(mux)(http.MethodGet, "/swagger/*", swagger); err != nil {
    log.Fatalf("boiler: swagger handler registering failed: %v", err)
  }
}

func (a *App) registerHelpHandler() {
  marshaledHelp, err := a.marshalHelpInfo()
  if err != nil {
    log.Fatalf("boiler: marshal help info failed: %v", err)
  }
  handleHelp := func(w http.ResponseWriter, r *http.Request) {
    if _, err = w.Write(marshaledHelp); err != nil {
      http.Error(w, "", http.StatusInternalServerError)
    }
  }
  a.dutyHttpRouter.Get("/help", handleHelp)
}

func (a *App) marshalHelpInfo() ([]byte, error) {
  type HelpInfo struct {
    DutyHttpPort      int `json:"duty_http_port"`
    GqlgenPort        int `json:"gqlgen_port"`
    GrpcPort          int `json:"grpc_port"`
    GrpcHttpProxyPort int `json:"grpc_http_proxy_port"`
  }
  const dutyServePort = 8092

  marshaledHelp, err := json.Marshal(&HelpInfo{
    DutyHttpPort:      dutyServePort,
    GqlgenPort:        a.gqlgenPort,
    GrpcPort:          a.grpcPort,
    GrpcHttpProxyPort: a.grpcHttpProxyPort,
  })
  if err != nil {
    return nil, fmt.Errorf("json.Marshal: %w", err)
  }
  return marshaledHelp, nil
}
