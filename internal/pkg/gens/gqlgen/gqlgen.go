package gqlgen

import (
  "context"
  "fmt"
  "path/filepath"

  log "github.com/sirupsen/logrus"
  "github.com/ushakovn/boiler/config"
  "github.com/ushakovn/boiler/internal/pkg/filer"
  "github.com/ushakovn/boiler/internal/pkg/gens/project"
  "github.com/ushakovn/boiler/internal/pkg/templater"
  "github.com/ushakovn/boiler/templates"
)

type Gqlgen struct {
  workDir            string
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
    workDir:            workDirPath,
    gqlgenDescPath:     config.GqlgenDescPath,
    gqlgenDescCompiled: config.GqlgenDescCompiled,
  }, nil
}

func (g *Gqlgen) Generate(ctx context.Context) error {
  // Use project generator for create gqlgen dirs
  p, err := project.NewProject(project.Config{
    ProjectDescPath:     g.gqlgenDescPath,
    ProjectDescCompiled: g.gqlgenDescCompiled,
  })
  if err != nil {
    return err
  }
  // Generate gqlgen project dirs
  if err = p.Generate(ctx); err != nil {
    return err
  }
  // Create yaml config for project
  if err = g.createGqlgenYaml(); err != nil {
    return err
  }
  return nil
}

func (g *Gqlgen) createGqlgenYaml() error {
  filePath := filepath.Join(g.workDir, "gqlgen.yaml")

  if err := templater.CopyTemplate(templates.GqlgenYaml, filePath); err != nil {
    return fmt.Errorf("copyTemplate: %w", err)
  }
  return nil
}
