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
  "github.com/ushakovn/boiler/internal/pkg/filer"
  "github.com/ushakovn/boiler/templates"
)

type Project struct {
  projectDescPath     string
  projectDescCompiled string
  workDirPath         string
  projectDesc         *projectDesc
  execFunctions       map[string]execFunc
}

type Config struct {
  ProjectDescPath     string
  ProjectDescCompiled string
}

func (c *Config) Validate() error {
  if c.ProjectDescPath == "" && c.ProjectDescCompiled == "" {
    log.Infof("boiler: using default project directories")
    c.ProjectDescCompiled = config.Project
  }
  return nil
}

func NewProject(config Config) (*Project, error) {
  if err := config.Validate(); err != nil {
    return nil, err
  }
  workDirPath, err := filer.WorkDirPath()
  if err != nil {
    return nil, err
  }
  proj := &Project{
    projectDescPath:     config.ProjectDescPath,
    projectDescCompiled: config.ProjectDescCompiled,
    workDirPath:         workDirPath,
  }
  proj.setExecFunctions()

  return proj, nil
}

func (g *Project) Generate(ctx context.Context) error {
  if err := g.loadProjectDesc(); err != nil {
    return fmt.Errorf("g.loadProjectDesc: %w", err)
  }
  for _, file := range g.projectDesc.Root.Files {
    file.Path = filepath.Join(g.workDirPath, file.Name)

    if err := g.genFile(file); err != nil {
      return fmt.Errorf("g.genFile: %w", err)
    }
  }
  for _, dir := range g.projectDesc.Root.Dirs {
    if err := g.genDirectory(ctx, dir, g.workDirPath); err != nil {
      return fmt.Errorf("g.genDirectory %w", err)
    }
  }
  return nil
}

func (g *Project) loadProjectDesc() error {
  var (
    buf []byte
    err error
  )
  if g.projectDescCompiled != "" {
    buf = []byte(g.projectDescCompiled)
  }
  if g.projectDescPath != "" {
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

func (g *Project) genFile(file *fileDesc) error {
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

func (g *Project) genDirectory(ctx context.Context, dir *directoryDesc, parentPath string) error {
  path := filepath.Join(parentPath, dir.Name.Execute(g.execFunctions))

  if err := os.Mkdir(path, os.ModePerm); err != nil && !filer.IsExistedDirectory(path) {
    return fmt.Errorf("os.Mkdir dir: %w", err)
  }
  for _, file := range dir.Files {
    file.Path = filepath.Join(parentPath, dir.Name.Execute(g.execFunctions), file.Name)

    if err := g.genFile(file); err != nil {
      return fmt.Errorf("g.genFile file: %w", err)
    }
  }
  for _, nested := range dir.Dirs {
    if err := g.genDirectory(ctx, nested, path); err != nil {
      return fmt.Errorf("g.genDirectory nested: %w", err)
    }
  }
  return nil
}

func (g *Project) workDirFolder() string {
  return filer.ExtractFileName(g.workDirPath)
}
