package gqlgen

import (
  "context"
  "fmt"
  "path/filepath"

  log "github.com/sirupsen/logrus"
  "github.com/ushakovn/boiler/config"
  "github.com/ushakovn/boiler/internal/boiler/gen"
  "github.com/ushakovn/boiler/internal/pkg/gens/project"
  "github.com/ushakovn/boiler/internal/pkg/utils"
  "github.com/ushakovn/boiler/templates"
)

type gqlgen struct {
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

func NewGqlgen(config Config) (gen.Generator, error) {
  if err := config.Validate(); err != nil {
    return nil, err
  }
  workDirPath, err := utils.WorkDirPath()
  if err != nil {
    return nil, err
  }
  return &gqlgen{
    workDir:            workDirPath,
    gqlgenDescPath:     config.GqlgenDescPath,
    gqlgenDescCompiled: config.GqlgenDescCompiled,
  }, nil
}

func (g *gqlgen) Generate(ctx context.Context) error {
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

func (g *gqlgen) createGqlgenYaml() error {
  filePath := filepath.Join(g.workDir, "gqlgen.yaml")

  if err := utils.CopyTemplate(templates.GqlgenYaml, filePath); err != nil {
    return fmt.Errorf("copyTemplate: %w", err)
  }
  return nil
}
