package gqlgen

import (
  "context"
  "fmt"
  "path/filepath"

  log "github.com/sirupsen/logrus"
  "github.com/ushakovn/boiler/config"
  "github.com/ushakovn/boiler/internal/pkg/aster"
  "github.com/ushakovn/boiler/internal/pkg/executor"
  "github.com/ushakovn/boiler/internal/pkg/filer"
  "github.com/ushakovn/boiler/internal/pkg/gens/project"
  "github.com/ushakovn/boiler/internal/pkg/templater"
  "github.com/ushakovn/boiler/templates"
)

type Gqlgen struct {
  workDirPath        string
  goModuleName       string
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
  goModuleName, err := filer.ExtractGoModuleName(workDirPath)
  if err != nil {
    return nil, err
  }
  return &Gqlgen{
    workDirPath:        workDirPath,
    goModuleName:       goModuleName,
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
  if err := g.generateGqlgenSchema(ctx); err != nil {
    return fmt.Errorf("g.generateGqlgenSchema: %w", err)
  }
  folderPath, err := filer.CreateNestedFolders(g.workDirPath, "internal", "app", "graph")
  if err != nil {
    return fmt.Errorf("filer.CreateNestedFolders: %w", err)
  }
  filePath := filepath.Join(folderPath, "service.go")

  if filer.IsExistedFile(filePath) {
    if err := g.regenerateGqlgenService(filePath); err != nil {
      return fmt.Errorf("g.regenerateGqlgenService: %w", err)
    }
  } else {
    if err := g.generateGqlgenService(filePath); err != nil {
      return fmt.Errorf("g.generateGqlgenService: %w", err)
    }
  }
  return nil
}

func (g *Gqlgen) generateGqlgenSchema(ctx context.Context) error {
  if err := executor.ExecCommandContext(ctx, "make", "generate-gqlgen"); err != nil {
    return fmt.Errorf("executor.ExecCommandContext: %w", err)
  }
  return nil
}

func (g *Gqlgen) generateGqlgenService(filePath string) error {
  templateData := g.buildGqlgenServiceDesc()

  if err := templater.ExecTemplateCopyWithGoFmt(templates.GqlgenService, filePath, templateData, nil); err != nil {
    return fmt.Errorf("execTemplateCopy: %w", err)
  }
  return nil
}

func (g *Gqlgen) regenerateGqlgenService(filePath string) error {
  const methodName = "RegisterService"

  methodFound, err := aster.FindMethodDeclaration(filePath, methodName)
  if err != nil {
    return fmt.Errorf("aster.FindMethodDeclaration: %w", err)
  }
  if methodFound {
    return nil
  }
  if err = filer.AppendStringToFile(filePath, templates.GqlgenRegisterService); err != nil {
    return fmt.Errorf("filer.AppendStringToFile: %w", err)
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
