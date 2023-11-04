package grpc

import (
  "fmt"
  "os"
  "path/filepath"
  "strings"

  "github.com/ushakovn/boiler/internal/pkg/aggr"
  "github.com/ushakovn/boiler/internal/pkg/filer"
)

type grpcServiceDesc struct {
  ServiceName     string
  ServiceCalls    []*grpcServiceCallDesc
  ServicePackages []*goPackageDesc
}

type grpcServiceCallDesc struct {
  ServiceName      string
  CallName         string
  CallInputProto   string
  CallOutputProto  string
  CallStubPackages []*goPackageDesc
}

type goPackageDesc struct {
  CustomName  string
  ImportLine  string
  ImportAlias string
  IsBuiltin   bool
  IsInstall   bool
}

func (g *Grpc) grpcServerInterfaceToDesc(grpcFilePath string, grpcServer *grpcServerInterface) *grpcServiceDesc {
  serviceName := grpcServer.ServerName
  servicePackages := g.buildServicePackages(grpcFilePath)

  serviceCalls := aggr.Map(grpcServer.ServerCalls, func(grpcServerCall *grpcServerInterfaceCall) *grpcServiceCallDesc {
    return g.grpcServerCallToDesc(serviceName, grpcFilePath, grpcServerCall)
  })

  return &grpcServiceDesc{
    ServiceName:     grpcServer.ServerName,
    ServiceCalls:    serviceCalls,
    ServicePackages: servicePackages,
  }
}

func (g *Grpc) buildServicePackages(grpcFilePath string) []*goPackageDesc {
  servicePackages := make([]*goPackageDesc, 0, len(grpcStubPackages)+1)

  servicePackages = append(servicePackages, grpcServicePackages...)
  servicePackages = append(servicePackages, g.buildPbPackage(grpcFilePath))

  return servicePackages
}

func (g *Grpc) grpcServerCallToDesc(serviceName, grpcFilePath string, grpcServerCall *grpcServerInterfaceCall) *grpcServiceCallDesc {
  stubPackages := g.buildCallStubPackages(grpcFilePath)

  return &grpcServiceCallDesc{
    ServiceName:      serviceName,
    CallName:         grpcServerCall.CallName,
    CallInputProto:   grpcServerCall.InputProto,
    CallOutputProto:  grpcServerCall.OutputProto,
    CallStubPackages: stubPackages,
  }
}

func (g *Grpc) buildCallStubPackages(grpcFilePath string) []*goPackageDesc {
  stubPackages := make([]*goPackageDesc, 0, len(grpcStubPackages)+1)

  stubPackages = append(stubPackages, grpcStubPackages...)
  stubPackages = append(stubPackages, g.buildPbPackage(grpcFilePath))

  return stubPackages
}

func (g *Grpc) buildPbPackage(grpcFilePath string) *goPackageDesc {
  // Extract package folder from grpc file path
  grpcFileName := filer.ExtractFileName(grpcFilePath)
  // Trim grpc file name
  grpcFilePath = strings.TrimSuffix(grpcFilePath, fmt.Sprint(string(os.PathSeparator), grpcFileName))
  // Then extract package name
  grpcPackageName := filer.ExtractFileName(grpcFilePath)

  stubPackagePath := filepath.Join(g.goModuleName, "internal", "pb", grpcPackageName)

  return &goPackageDesc{
    CustomName:  "proto/desc",
    ImportLine:  stubPackagePath,
    ImportAlias: "desc",
  }
}

var grpcServicePackages = []*goPackageDesc{
  {
    CustomName:  "boiler/app",
    ImportLine:  "github.com/ushakovn/boiler/pkg/app",
    ImportAlias: "",
    IsInstall:   true,
  },
}

var grpcStubPackages = []*goPackageDesc{
  {
    CustomName:  "go/context",
    ImportLine:  "context",
    ImportAlias: "",
    IsBuiltin:   true,
  },
  {
    CustomName:  "grpc/status",
    ImportLine:  "google.golang.org/grpc/status",
    ImportAlias: "",
    IsInstall:   true,
  },
  {
    CustomName:  "grpc/codes",
    ImportLine:  "google.golang.org/grpc/codes",
    ImportAlias: "",
    IsInstall:   true,
  },
}
