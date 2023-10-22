package project

import (
  "context"
  "encoding/json"
  "fmt"
  "os"
  "path/filepath"
  "strings"

  log "github.com/sirupsen/logrus"
  "github.com/ushakovn/boiler/config"
  "github.com/ushakovn/boiler/internal/boiler/gen"
  "github.com/ushakovn/boiler/pkg/utils"
  "github.com/ushakovn/boiler/templates"
)

type project struct {
  projectDescPath  string
  workDirPath      string
  projectDesc      *projectDesc
  execFunctions    map[string]execFunc
  withCompiledDesc bool
}

type Config struct {
  ProjectDescPath  string
  withCompiledDesc bool
}

func (c *Config) Validate() error {
  if c.ProjectDescPath == "" {
    log.Warnf("boiler: using default project directories")
    c.withCompiledDesc = true
  }
  return nil
}

func NewProject(config Config) (gen.Generator, error) {
  if err := config.Validate(); err != nil {
    return nil, err
  }
  workDirPath, err := utils.Env("PWD")
  if err != nil {
    return nil, err
  }
  proj := &project{
    withCompiledDesc: config.withCompiledDesc,
    projectDescPath:  config.ProjectDescPath,
    workDirPath:      workDirPath,
  }
  proj.setExecFunctions()

  return proj, nil
}

func (g *project) Generate(ctx context.Context) error {
  if err := g.loadProjectDesc(); err != nil {
    return fmt.Errorf("g.loadProjectDesc: %w", err)
  }
  for _, file := range g.projectDesc.Root.Files {
    file.Path = g.buildPath(file.Name)

    if err := g.genFile(file); err != nil {
      return fmt.Errorf("g.genFile: %w", err)
    }
  }
  for _, dir := range g.projectDesc.Root.Dirs {
    if err := g.genDirectory(ctx, dir, ""); err != nil {
      return fmt.Errorf("g.genDirectory %w", err)
    }
  }
  return nil
}

func (g *project) loadProjectDesc() error {
  var (
    buf []byte
    err error
  )
  if g.withCompiledDesc {
    buf = []byte(config.Project)
  } else {
    if buf, err = os.ReadFile(g.projectDescPath); err != nil {
      return fmt.Errorf("os.ReadFile projectDir: %w", err)
    }
  }
  proj := &projectDesc{}

  if err = json.Unmarshal(buf, proj); err != nil {
    return fmt.Errorf("json.Unmarshal: %w", err)
  }
  g.projectDesc = proj

  return nil
}

func (g *project) genFile(file *fileDesc) error {
  var path string

  if path = file.Path; path == "" {
    return fmt.Errorf("file.Path not specified")
  }
  if extent := file.Extension; extent != "" {
    extent = strings.TrimPrefix(file.Extension, ".")
    path = fmt.Sprintf("%s.%s", file.Path, extent)
  }
  if _, err := os.Create(path); err != nil {
    return fmt.Errorf("os.CreateFile: %w", err)
  }
  if template := file.Template; template != nil {
    buf := loadFileTemplate(template)

    if err := os.WriteFile(path, buf, os.ModePerm); err != nil {
      return fmt.Errorf("os.WriteFile: %w", err)
    }
  }
  return nil
}

func loadFileTemplate(desc *templateDesc) []byte {
  var compiled string

  switch desc.Name {
  case templates.NameMain:
    compiled = templates.Main
  case templates.NameGomod:
    compiled = templates.Gomod
  }
  return []byte(compiled)
}

func (g *project) genDirectory(ctx context.Context, dir *directoryDesc, parentPath string) error {
  path := g.buildPath(parentPath, dir.Name.Execute(g.execFunctions))

  if err := os.Mkdir(path, os.ModePerm); err != nil {
    return fmt.Errorf("os.Mkdir dir: %w", err)
  }
  for _, file := range dir.Files {
    file.Path = g.buildPath(parentPath, dir.Name.Execute(g.execFunctions), file.Name)

    if err := g.genFile(file); err != nil {
      return fmt.Errorf("g.genFile file: %w", err)
    }
  }
  for _, nested := range dir.Dirs {
    if err := g.genDirectory(ctx, nested, dir.Name.Execute(g.execFunctions)); err != nil {
      return fmt.Errorf("g.genDirectory nested: %w", err)
    }
  }
  return nil
}

func (g *project) buildPath(parts ...string) string {
  pd := []string{g.workDirPath}
  pd = append(pd, parts...)
  p := filepath.Join(pd...)
  return p
}

func (g *project) workDirFolder() string {
  if parts := strings.Split(g.workDirPath, `/`); len(parts) > 0 {
    return parts[len(parts)-1]
  }
  if parts := strings.Split(g.workDirPath, `\`); len(parts) > 0 {
    return parts[len(parts)-1]
  }
  return g.workDirPath
}
