package main

import (
  "context"
  "flag"
  "fmt"

  log "github.com/sirupsen/logrus"
  "github.com/ushakovn/boiler/internal/pkg/gens/factory"
  "github.com/ushakovn/boiler/internal/pkg/gens/project"
  "github.com/ushakovn/boiler/internal/pkg/gens/rpc"
)

func main() {
  logBoilerMark()
  ctx := context.Background()

  typ, config, err := parseFlags()
  if err != nil {
    log.Fatalf("boiler validation error: %v", err)
  }

  gens, err := factory.NewGenerator(config, typ)
  if err != nil {
    log.Fatalf("boiler initialization error: %v", err)
  }

  log.Infof("boiler info: generation started")

  if err = gens.Generate(ctx); err != nil {
    log.Fatalf("boiler generation error: %v", err)
  }

  log.Infof("boiler info: generation finished")
}

func parseFlags() (factory.Typ, factory.CommonConfig, error) {
  genType := flag.String("type", "", "generator type")

  projDescPath := flag.String("project", "", "project description json/yaml path")
  rpcDescPath := flag.String("rpc", "", "rpc description json/yaml path")

  flag.Parse()
  if *genType == "" {
    return "", factory.CommonConfig{}, fmt.Errorf("generator type not specified")
  }

  return factory.Typ(*genType), factory.CommonConfig{
    Project: project.Config{
      ProjectDescPath: *projDescPath,
    },
    Rpc: rpc.Config{
      RpcDescPath: *rpcDescPath,
    },
  }, nil
}

func logBoilerMark() {
  log.Infof(`

 _           _ _           
| |         (_) |          
| |__   ___  _| | ___ _ __ 
| '_ \ / _ \| | |/ _ \ '__|
| |_) | (_) | | |  __/ |   
|_.__/ \___/|_|_|\___|_|


`)
}
