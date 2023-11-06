package gqlgen

import (
  "context"
  "fmt"
  "path/filepath"

  log "github.com/sirupsen/logrus"
  "github.com/ushakovn/boiler/config"
  "github.com/ushakovn/boiler/internal/pkg/executor"
  "github.com/ushakovn/boiler/internal/pkg/filer"
  "github.com/ushakovn/boiler/internal/pkg/gens/project"
  "github.com/ushakovn/boiler/internal/pkg/templater"
  "github.com/ushakovn/boiler/templates"
)

type Gqlgen struct {
  workDirPath        string
  gqlgenDescPath     string
  gqlgenDescCompiled string
}

type Config struct {
  GqlgenDescPath     string
  GqlgenDescCompiled string
}

func (c *Config) Validate() error {
  if c.GqlgenDescPath == "" && c.GqlgenDescCompiled == "" {
    log.Infof("boiler: using default gqlgen description")
    c.GqlgenDescCompiled = config.Gqlgen
  }
  return nil
}

func NewGqlgen(config Config) (*Gqlgen, error) {
  if err := config.Validate(); err != nil {
    return nil, err
  }
  workDirPath, err := filer.WorkDirPath()
  if err != nil {
    return nil, err
  }
  return &Gqlgen{
    workDirPath:        workDirPath,
    gqlgenDescPath:     config.GqlgenDescPath,
    gqlgenDescCompiled: config.GqlgenDescCompiled,
  }, nil
}

func (g *Gqlgen) Init(ctx context.Context) error {
  // Use project generator for create gqlgen dirs
  p, err := project.NewProject(project.Config{
    ProjectDescPath:     g.gqlgenDescPath,
    ProjectDescCompiled: g.gqlgenDescCompiled,
  })
  if err != nil {
    return fmt.Errorf("project.NewProject: %w", err)
  }
  // Generate gqlgen project dirs
  if err = p.Init(ctx); err != nil {
    return fmt.Errorf("p.Init: %w", err)
  }
  if err = g.createGqlgenYaml(); err != nil {
    return fmt.Errorf("g.createGqlgenYaml: %w", err)
  }
  if err = g.createGqlgenTools(); err != nil {
    return fmt.Errorf("g.createGqlgenTools: %w", err)
  }
  if err = g.createMakeMkTarget(); err != nil {
    return fmt.Errorf("g.createMakeMkTarget: %w", err)
  }
  if err = g.createMakefileIfNotExist(); err != nil {
    return fmt.Errorf("g.createMakefileIfNotExist: %w", err)
  }
  return nil
}

func (g *Gqlgen) Generate(ctx context.Context) error {
  if err := executor.ExecCommandContext(ctx, "make", "generate-gqlgen"); err != nil {
    return fmt.Errorf("executor.ExecCommandContext: %w", err)
  }
  return nil
}

func (g *Gqlgen) createMakeMkTargetIfNotExist() error {
  filePath := filepath.Join(g.workDirPath, "make.mk")

  if !filer.IsExistedFile(filePath) {
    if err := g.createMakeMkTarget(); err != nil {
      return fmt.Errorf("g.createMakeMkTarget: %w", err)
    }
  }
  return nil
}

func (g *Gqlgen) createMakeMkTarget() error {
  makeMkPath := filepath.Join(g.workDirPath, "make.mk")

  if err := filer.AppendStringToFile(makeMkPath, templates.GqlgenMakeMk); err != nil {
    return fmt.Errorf("filer.AppendStringToFile: %w", err)
  }
  return nil
}

func (g *Gqlgen) createMakefileIfNotExist() error {
  filePath := filepath.Join(g.workDirPath, "Makefile")

  if !filer.IsExistedFile(filePath) {
    if err := templater.ExecTemplateCopy(templates.Makefile, filePath, nil, nil); err != nil {
      return fmt.Errorf("execTemplateCopy: %w", err)
    }
  }
  return nil
}

func (g *Gqlgen) createGqlgenYaml() error {
  filePath := filepath.Join(g.workDirPath, "gqlgen.yaml")

  if err := templater.CopyTemplate(templates.GqlgenYaml, filePath); err != nil {
    return fmt.Errorf("copyTemplate: %w", err)
  }
  return nil
}

func (g *Gqlgen) createGqlgenTools() error {
  filePath := filepath.Join(g.workDirPath, "tools.go")

  if err := templater.CopyTemplate(templates.GqlgenTools, filePath); err != nil {
    return fmt.Errorf("copyTemplate: %w", err)
  }
  return nil
}
