package grpc

import (
  "context"
  "fmt"
  "path/filepath"
  "strings"
  "text/template"

  "github.com/ushakovn/boiler/internal/pkg/executor"
  "github.com/ushakovn/boiler/internal/pkg/filer"
  "github.com/ushakovn/boiler/internal/pkg/stringer"
  "github.com/ushakovn/boiler/internal/pkg/templater"
  "github.com/ushakovn/boiler/templates"
)

type Grpc struct {
  workDirPath  string
  goModuleName string
}

type Config struct{}

func NewGrpc(_ Config) (*Grpc, error) {
  workDirPath, err := filer.WorkDirPath()
  if err != nil {
    return nil, err
  }
  goModuleName, err := filer.ExtractGoModuleName(workDirPath)
  if err != nil {
    return nil, err
  }
  return &Grpc{
    workDirPath:  workDirPath,
    goModuleName: goModuleName,
  }, nil
}

func (g *Grpc) Generate(ctx context.Context) error {
  if err := g.createMakeMkTargetIfNotExist(); err != nil {
    return fmt.Errorf("g.createMakeMkTargetIfNotExist: %w", err)
  }
  if err := g.createMakefileIfNotExist(); err != nil {
    return fmt.Errorf("g.createMakefileIfNotExist: %w", err)
  }
  if err := g.generateMakeMkProto(ctx); err != nil {
    return fmt.Errorf("g.generateMakeMkProto: %w", err)
  }
  if err := g.generateGrpcServices(); err != nil {
    return fmt.Errorf("g.generateGrpcServices: %w", err)
  }
  return nil
}

func (g *Grpc) Init(_ context.Context) error {
  protoDirPath, err := g.createProtoDirectory()
  if err != nil {
    return fmt.Errorf("g.createProtoDirectory: %w", err)
  }
  if err = g.createServiceProtoFile(protoDirPath); err != nil {
    return fmt.Errorf("g.createServiceProtoFile: %w", err)
  }
  if err = g.createMakeMkTargetIfNotExist(); err != nil {
    return fmt.Errorf("g.createMakeMkTargetIfNotExist: %w", err)
  }
  if err = g.createMakefileIfNotExist(); err != nil {
    return fmt.Errorf("g.createMakefileIfNotExist: %w", err)
  }
  return nil
}

func (g *Grpc) generateGrpcServices() error {
  // Collect grpc files
  grpcFilesPath, err := g.collectGrpcFilesPath()
  if err != nil {
    return fmt.Errorf("g.collectGrpcFilesPath: %w", err)
  }
  // For each found grpc file
  for _, grpcFilePath := range grpcFilesPath {
    // Create grpc service and call stubs
    if err = g.generateServiceWithCallStubs(grpcFilePath); err != nil {
      return fmt.Errorf("g.generateServiceWithCallStubs: %w", err)
    }
  }
  return nil
}

func (g *Grpc) generateServiceWithCallStubs(grpcFilePath string) error {
  serverInterface, err := scanGrpcServerInterface(grpcFilePath)
  if err != nil {
    return fmt.Errorf("scanGrpcServerInterface: %w", err)
  }
  serviceDesc := g.grpcServerInterfaceToDesc(grpcFilePath, serverInterface)
  serviceName := stringer.StringToSnakeCase(serviceDesc.ServiceName)

  serviceFolderPath, err := filer.CreateNestedFolders(g.workDirPath, "internal", "app", serviceName)
  if err != nil {
    return fmt.Errorf("filer.CreateNestedFolders: %w", err)
  }
  serviceFilePath := filepath.Join(serviceFolderPath, "service.go")

  templateFuncMap := template.FuncMap{
    "toSnakeCase": stringer.StringToSnakeCase,
  }

  if filer.IsExistedFile(serviceFilePath) {
    // If service file exist - analyze it with ast
    if err = regenerateGrpcService(serviceFilePath); err != nil {
      return fmt.Errorf("regenerateGrpcService: %w", err)
    }
  } else {
    // If service file not exist - generated it
    if err = templater.ExecTemplateCopyWithGoFmt(templates.GrpcService, serviceFilePath, serviceDesc, templateFuncMap); err != nil {
      return fmt.Errorf("executeTemplateCopy: %w", err)
    }
  }

  // For each grpc service calls
  for _, serviceCallDesc := range serviceDesc.ServiceCalls {
    callName := stringer.StringToSnakeCase(serviceCallDesc.CallName)

    callStubFileName := fmt.Sprint(callName, ".go")
    callStubFilePath := filepath.Join(serviceFolderPath, callStubFileName)

    if filer.IsExistedFile(callStubFilePath) {
      // If call stub file exist - analyze it with ast
      if err = regenerateGrpcServiceStub(callStubFilePath, serviceCallDesc); err != nil {
        return fmt.Errorf("regenerateGrpcServiceStub: %w", err)
      }
    } else {
      // If call stub file not exist - generated it
      if err = templater.ExecTemplateCopyWithGoFmt(templates.GrpcStub, callStubFilePath, serviceCallDesc, templateFuncMap); err != nil {
        return fmt.Errorf("executeTemplateCopy: %w", err)
      }
    }
  }
  return nil
}

func (g *Grpc) collectGrpcFilesPath() ([]string, error) {
  parentPath := filepath.Join(g.workDirPath, "internal", "pb")

  genProtoFilesPath, err := filer.CollectNestedFilesPath("go", parentPath)
  if err != nil {
    return nil, fmt.Errorf("filer.CollectNestedFilesPath: %w", err)
  }
  const (
    genPbSuffix   = ".pb.go"
    genGrpcSuffix = "_grpc"
  )
  var grpcFilesPath []string

  for _, genProtoFilePath := range genProtoFilesPath {
    if !strings.HasSuffix(genProtoFilePath, genPbSuffix) {
      continue
    }
    protoFileName := strings.TrimSuffix(genProtoFilePath, genPbSuffix)
    protoFileName = filer.ExtractFileName(protoFileName)

    if !strings.HasSuffix(protoFileName, genGrpcSuffix) {
      continue
    }
    grpcFilesPath = append(grpcFilesPath, genProtoFilePath)
  }

  return grpcFilesPath, nil
}

func (g *Grpc) generateMakeMkProto(ctx context.Context) error {
  if err := executor.ExecCommandContext(ctx, "make", "generate-protoc"); err != nil {
    return fmt.Errorf("executor.ExecCommandContext: %w", err)
  }
  return nil
}

func (g *Grpc) createMakeMkTargetIfNotExist() error {
  const fileName = "make.mk"
  filePath := filepath.Join(g.workDirPath, fileName)

  if !filer.IsExistedFile(filePath) {
    if err := g.createMakeMkTarget(); err != nil {
      return fmt.Errorf("g.createMakeMkTarget: %w", err)
    }
  }
  return nil
}

func (g *Grpc) createMakeMkTarget() error {
  const fileName = "make.mk"
  goPackageTrim := g.goModuleName

  templateData := map[string]any{
    "goPackageTrim": goPackageTrim,
  }
  executedBuf, err := templater.ExecTemplate(templates.GrpcMakeMk, templateData, nil)
  if err != nil {
    return fmt.Errorf("executeTemplate")
  }
  executedTarget := string(executedBuf)
  makeMkPath := filepath.Join(g.workDirPath, fileName)

  if err = filer.AppendStringToFile(makeMkPath, executedTarget); err != nil {
    return fmt.Errorf("filer.AppendStringToFile: %w", err)
  }
  return nil
}

func (g *Grpc) createMakefileIfNotExist() error {
  const fileName = "Makefile"
  filePath := filepath.Join(g.workDirPath, fileName)

  if !filer.IsExistedFile(filePath) {
    if err := templater.ExecTemplateCopy(templates.ProjectMakefile, filePath, nil, nil); err != nil {
      return fmt.Errorf("execTemplateCopy: %w", err)
    }
  }
  return nil
}

func (g *Grpc) createProtoDirectory() (string, error) {
  protoDirPath, err := filer.CreateNestedFolders(g.workDirPath, "api", g.workDirFolder())
  if err != nil {
    return "", fmt.Errorf("filer.CreateNestedFolders: %w", err)
  }
  return protoDirPath, nil
}

func (g *Grpc) createServiceProtoFile(protoDirPath string) error {
  serviceName := g.workDirFolder()

  goPackage := filepath.Join(g.goModuleName, "internal", "pb", serviceName)
  goPackageWithSuffix := fmt.Sprint(goPackage, ";", serviceName)

  templateData := map[string]any{
    "serviceName": serviceName,
    "goPackage":   goPackageWithSuffix,
  }
  protoFileName := fmt.Sprint(serviceName, ".proto")

  protoFilePath := filepath.Join(protoDirPath, protoFileName)

  if err := templater.ExecTemplateCopy(templates.GrpcProto, protoFilePath, templateData, nil); err != nil {
    return fmt.Errorf("executeTemplateCopy: %w", err)
  }
  return nil
}

func (g *Grpc) workDirFolder() string {
  return filer.ExtractFileName(g.workDirPath)
}
