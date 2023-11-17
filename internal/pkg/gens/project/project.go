package project

import (
  "context"
  "encoding/json"
  "fmt"
  "os"
  "path/filepath"
  "strings"

  "github.com/ushakovn/boiler/config"
  "github.com/ushakovn/boiler/internal/pkg/filer"
  "github.com/ushakovn/boiler/internal/pkg/templater"
  "github.com/ushakovn/boiler/templates"
)

type Project struct {
  projectDescPath     string
  projectDescCompiled string
  goModVersion        string
  workDirPath         string
  projectDesc         *projectDesc
}

type Config struct {
  GoModVersion        string
  ProjectDescPath     string
  ProjectDescCompiled string
}

func (c Config) WithDefault() Config {
  const goModVersionDefault = "1.19"

  if c.GoModVersion == "" {
    c.GoModVersion = goModVersionDefault
  }
  if c.ProjectDescPath == "" && c.ProjectDescCompiled == "" {
    c.ProjectDescCompiled = config.Project
  }
  return c
}

func NewProject(config Config) (*Project, error) {
  config = config.WithDefault()

  workDirPath, err := filer.WorkDirPath()
  if err != nil {
    return nil, err
  }
  proj := &Project{
    projectDescPath:     config.ProjectDescPath,
    projectDescCompiled: config.ProjectDescCompiled,
    goModVersion:        config.GoModVersion,
    workDirPath:         workDirPath,
  }
  return proj, nil
}

func (g *Project) Init(ctx context.Context) error {
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
    globalTemplate := loadGlobalCompiledTemplate(template)
    globalData := g.loadGlobalTemplatesData()

    buf, err := templater.ExecTemplate(globalTemplate, globalData, nil)
    if err != nil {
      return fmt.Errorf("execTemplate: %w", err)
    }
    if err = os.WriteFile(path, buf, os.ModePerm); err != nil {
      return fmt.Errorf("os.WriteFile: %w", err)
    }
  }
  return nil
}

func loadGlobalCompiledTemplate(desc *templateDesc) string {
  var compiled string

  // Global file template name
  switch desc.Name {

  // Project templates

  case templates.NameMain:
    compiled = templates.ProjectMain

  case templates.NameGomod:
    compiled = templates.ProjectGomod

  case templates.NameMakefile:
    compiled = templates.ProjectMakefile

  // Gqlgen Graphql Schema templates

  case templates.NameGqlgenSchema:
    compiled = templates.GqlgenSchema

  case templates.NameGqlgenMutation:
    compiled = templates.GqlgenMutation

  case templates.NameGqlgenQuery:
    compiled = templates.GqlgenQuery

  case templates.NameGqlgenTypes:
    compiled = templates.GqlgenTypes

  case templates.NameGqlgenEnums:
    compiled = templates.GqlgenEnums

  case templates.NameGqlgenScalars:
    compiled = templates.GqlgenScalars
  }

  return compiled
}

func (g *Project) loadGlobalTemplatesData() map[string]any {
  templatesData := map[string]any{
    "goModVersion": g.goModVersion,
  }
  return templatesData
}

func (g *Project) genDirectory(ctx context.Context, dir *directoryDesc, parentPath string) error {
  path := filepath.Join(parentPath, dir.Name.Value)

  if err := os.Mkdir(path, os.ModePerm); err != nil && !filer.IsExistedDirectory(path) {
    return fmt.Errorf("os.Mkdir dir: %w", err)
  }
  for _, file := range dir.Files {
    file.Path = filepath.Join(parentPath, dir.Name.Value, file.Name)

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
