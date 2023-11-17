// Package rpc DEPRECATED; DO NOT USE.
package rpc

import (
  "context"
  "encoding/json"
  "fmt"
  "os"
  "path/filepath"
  "strings"

  "github.com/ushakovn/boiler/internal/pkg/filer"
  "github.com/ushakovn/boiler/internal/pkg/stringer"
  "github.com/ushakovn/boiler/internal/pkg/templater"
  "github.com/ushakovn/boiler/templates"
  "gopkg.in/yaml.v3"
)

type Rpc struct {
  rpcDescPath string
  workDirPath string
  rpcDesc     *rootDesc
}

type Config struct {
  RpcDescPath string
}

func NewRpc(config Config) (*Rpc, error) {
  workDirPath, err := filer.WorkDirPath()
  if err != nil {
    return nil, err
  }
  return &Rpc{
    rpcDescPath: config.RpcDescPath,
    workDirPath: workDirPath,
  }, nil
}

func (g *Rpc) Generate(_ context.Context) error {
  if err := g.loadRpcDesc(); err != nil {
    return fmt.Errorf("g.loadRpcDesc: %w", err)
  }
  if err := g.genRpcHandler(); err != nil {
    return fmt.Errorf("g.genRpcHandler: %w", err)
  }
  return nil
}

func (g *Rpc) loadRpcDesc() error {
  buf, err := os.ReadFile(g.rpcDescPath)
  if err != nil {
    return fmt.Errorf("os.ReadFile projectDir: %w", err)
  }
  fileExtension := filer.ExtractFileExtension(g.rpcDescPath)

  desc, err := parseRootDesc(fileExtension, buf)
  if err != nil {
    return fmt.Errorf("parseRootDesc: %w", err)
  }
  g.rpcDesc = desc

  return nil
}

func parseRootDesc(fileExtension string, buf []byte) (*rootDesc, error) {
  var (
    desc *rootDesc
    err  error
  )
  switch fileExtension {
  case "yml", "yaml", "YML", "YAML":
    err = yaml.Unmarshal(buf, &desc)
  case "json", "JSON":
    err = json.Unmarshal(buf, &desc)
  default:
    err = fmt.Errorf("unsupported file extension: %s", fileExtension)
  }
  return desc, err
}

func (g *Rpc) genRpcHandler() error {
  if err := g.rpcDesc.Validate(); err != nil {
    return err
  }
  rpcTemplate := rootDescToRpcTemplates(g.rpcDesc)

  handlerDir, err := g.createHandlerDir()
  if err != nil {
    return fmt.Errorf("g.createHandlerDir: %w", err)
  }
  filePath := filepath.Join(handlerDir, "contracts.go")

  if err = templater.ExecTemplateCopyWithGoFmt(templates.RpcContracts, filePath, rpcTemplate.Contracts, nil); err != nil {
    return fmt.Errorf("executeTemplateCopy: %w", err)
  }
  filePath = filepath.Join(handlerDir, "handler.go")

  if err = templater.ExecTemplateCopyWithGoFmt(templates.RpcHandler, filePath, rpcTemplate, nil); err != nil {
    return fmt.Errorf("executeTemplateCopy: %w", err)
  }

  for _, handle := range rpcTemplate.Handles {
    fileName := stringer.StringToSnakeCase(handle.Name)
    fileName = fmt.Sprint(fileName, ".go")

    filePath = filepath.Join(handlerDir, fileName)

    if err = templater.ExecTemplateCopyWithGoFmt(templates.RpcHandle, filePath, handle, nil); err != nil {
      return fmt.Errorf("executeTemplateCopy: %w", err)
    }
  }
  return nil
}

func (g *Rpc) createHandlerDir() (string, error) {
  projectName := g.projectName()

  defaultDirParts := []string{g.workDirPath, "internal", projectName, "handler"}
  defaultDir := filepath.Join(defaultDirParts...)

  prevDirParts := make([]string, 0, len(defaultDirParts))

  for _, dirPart := range defaultDirParts {
    // Create directories for handler package
    prevDirParts = append(prevDirParts, dirPart)
    path := filepath.Join(prevDirParts...)

    if _, err := os.Stat(path); os.IsNotExist(err) {
      if err = os.Mkdir(path, os.ModePerm); err != nil {
        return "", fmt.Errorf("os.Mkdir: %w", err)
      }
    }
  }
  // Check created directories
  if _, err := os.Stat(defaultDir); os.IsNotExist(err) {
    return "", fmt.Errorf("os.Stat: %s: err: %v", defaultDir, err)
  }

  return defaultDir, nil
}

func (g *Rpc) projectName() string {
  if parts := strings.Split(g.workDirPath, `/`); len(parts) > 0 {
    return parts[len(parts)-1]
  }
  if parts := strings.Split(g.workDirPath, `\`); len(parts) > 0 {
    return parts[len(parts)-1]
  }
  return g.workDirPath
}
