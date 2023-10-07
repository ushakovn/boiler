package rpc

import (
  "bytes"
  "context"
  "encoding/json"
  "fmt"
  "go/format"
  "os"
  "path/filepath"
  "strings"
  "text/template"

  "github.com/ushakovn/boiler/internal/boiler/gen"
  "github.com/ushakovn/boiler/pkg/utils"
  "github.com/ushakovn/boiler/templates"
  "gopkg.in/yaml.v3"
)

type rpc struct {
  rpcDescPath string
  workDirPath string
  rpcDesc     *rootDesc
}

type Config struct {
  RpcDescPath string
}

func (c *Config) Validate() error {
  if c.RpcDescPath == "" {
    return fmt.Errorf("rpc description path not specfied")
  }
  return nil
}

func NewRpc(config Config) (gen.Generator, error) {
  if err := config.Validate(); err != nil {
    return nil, err
  }
  workDirPath, err := utils.Env("PWD")
  if err != nil {
    return nil, err
  }
  return &rpc{
    rpcDescPath: config.RpcDescPath,
    workDirPath: workDirPath,
  }, nil
}

func (g *rpc) Generate(context.Context) error {
  if err := g.loadRpcDesc(); err != nil {
    return fmt.Errorf("g.loadRpcDesc: %w", err)
  }
  if err := g.genRpcHandler(); err != nil {
    return fmt.Errorf("g.genRpcHandler: %w", err)
  }
  return nil
}

func (g *rpc) loadRpcDesc() error {
  buf, err := os.ReadFile(g.rpcDescPath)
  if err != nil {
    return fmt.Errorf("os.ReadFile projectDir: %w", err)
  }
  format := utils.FileFormat(g.rpcDescPath)

  desc, err := parseRootDesc(format, buf)
  if err != nil {
    return fmt.Errorf("parseRootDesc: %w", err)
  }
  g.rpcDesc = desc

  return nil
}

func parseRootDesc(format string, buf []byte) (*rootDesc, error) {
  var (
    desc *rootDesc
    err  error
  )
  switch format {
  case "yml", "yaml", "YML", "YAML":
    err = yaml.Unmarshal(buf, &desc)
  case "json", "JSON":
    err = json.Unmarshal(buf, &desc)
  default:
    err = fmt.Errorf("unsupported format: %s", format)
  }
  return desc, err
}

func (g *rpc) genRpcHandler() error {
  if err := g.rpcDesc.Validate(); err != nil {
    return err
  }
  rpcTemplate := rootDescToRpcTemplates(g.rpcDesc)

  handlerDir, err := g.createHandlerDir()
  if err != nil {
    return fmt.Errorf("g.createHandlerDir: %w", err)
  }
  filePath := filepath.Join(handlerDir, "contracts.go")

  if err = executeTemplateCopy(templates.Contracts, filePath, rpcTemplate.Contracts); err != nil {
    return fmt.Errorf("executeTemplateCopy: %w", err)
  }
  filePath = filepath.Join(handlerDir, "handler.go")

  if err = executeTemplateCopy(templates.Handler, filePath, rpcTemplate); err != nil {
    return fmt.Errorf("executeTemplateCopy: %w", err)
  }

  for _, handle := range rpcTemplate.Handles {
    fileName := utils.StringToSnakeCase(handle.Name)
    fileName = fmt.Sprint(fileName, ".go")

    filePath = filepath.Join(handlerDir, fileName)

    if err = executeTemplateCopy(templates.Handle, filePath, handle); err != nil {
      return fmt.Errorf("executeTemplateCopy: %w", err)
    }
  }
  return nil
}

func executeTemplateCopy(templateCompiled, filePath string, structPtr any) error {
  t, err := template.New("").Parse(templateCompiled)
  if err != nil {
    return fmt.Errorf("template.New().Parse: %w", err)
  }
  var (
    buffer bytes.Buffer
    buf    []byte
  )
  if err = t.Execute(&buffer, structPtr); err != nil {
    return fmt.Errorf("t.Execute: %w", err)
  }
  if buf, err = format.Source(buffer.Bytes()); err != nil {
    return fmt.Errorf("format.Source: %w", err)
  }
  if err = os.WriteFile(filePath, buf, os.ModePerm); err != nil {
    return fmt.Errorf("os.WriteFile: %w", err)
  }
  return nil
}

func (g *rpc) createHandlerDir() (string, error) {
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

func (g *rpc) projectName() string {
  if parts := strings.Split(g.workDirPath, `/`); len(parts) > 0 {
    return parts[len(parts)-1]
  }
  if parts := strings.Split(g.workDirPath, `\`); len(parts) > 0 {
    return parts[len(parts)-1]
  }
  return g.workDirPath
}
