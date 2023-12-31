package gqlgen

import (
  "context"
  "fmt"
  "path/filepath"

  "github.com/ushakovn/boiler/config"
  "github.com/ushakovn/boiler/internal/pkg/ast"
  "github.com/ushakovn/boiler/internal/pkg/executor"
  "github.com/ushakovn/boiler/internal/pkg/filer"
  "github.com/ushakovn/boiler/internal/pkg/gens/project"
  "github.com/ushakovn/boiler/internal/pkg/makefile"
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

func (c Config) WithDefault() Config {
  if c.GqlgenDescPath == "" && c.GqlgenDescCompiled == "" {
    c.GqlgenDescCompiled = config.Gqlgen
  }
  return c
}

func NewGqlgen(config Config) (*Gqlgen, error) {
  config = config.WithDefault()

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
  if err = g.createGqlgenConfig(); err != nil {
    return fmt.Errorf("g.createGqlgenConfig: %w", err)
  }
  if err = g.createGqlgenTools(); err != nil {
    return fmt.Errorf("g.createGqlgenTools: %w", err)
  }
  if err = g.createMakeMkTargetIfNotExist(); err != nil {
    return fmt.Errorf("g.createMakeMkTargetIfNotExist: %w", err)
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
    if err = g.regenerateGqlgenService(filePath); err != nil {
      return fmt.Errorf("g.regenerateGqlgenService: %w", err)
    }
  } else {
    if err = g.generateGqlgenService(filePath); err != nil {
      return fmt.Errorf("g.generateGqlgenService: %w", err)
    }
  }
  if err = g.createMakeMkTargetIfNotExist(); err != nil {
    return fmt.Errorf("g.createMakeMkTargetIfNotExist: %w", err)
  }
  return nil
}

func (g *Gqlgen) generateGqlgenSchema(ctx context.Context) error {
  if err := executor.ExecCmdCtx(ctx, "make", "generate-gqlgen"); err != nil {
    return fmt.Errorf("executor.ExecCmdCtx: %w", err)
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

  ok, err := ast.ContainsMethodDecl(filePath, methodName)
  if err != nil {
    return fmt.Errorf("aster.ContainsMethodDecl: %w", err)
  }
  if ok {
    return nil
  }
  if err = filer.AppendStringToFile(filePath, templates.GqlgenRegisterService); err != nil {
    return fmt.Errorf("filer.AppendStringToFile: %w", err)
  }
  return nil
}

func (g *Gqlgen) createMakeMkTargetIfNotExist() error {
  filePath := filepath.Join(g.workDirPath, "make.mk")

  type makeMkTarget struct {
    targetName       string
    compiledTemplate string
  }

  targets := []*makeMkTarget{
    {
      targetName:       templates.GqlgenMakeMkBinDepsName,
      compiledTemplate: templates.GqlgenMakeMkBinDeps,
    },
    {
      targetName:       templates.GqlgenMakeMkGenerateName,
      compiledTemplate: templates.GqlgenMakeMkGenerate,
    },
  }

  for _, target := range targets {
    ok, err := makefile.ContainsTarget(filePath, target.targetName)
    if err != nil {
      return fmt.Errorf("makefile.ContainsTarget: %w", err)
    }
    if ok {
      continue
    }
    if err = g.createMakeMkTarget(target.compiledTemplate); err != nil {
      return fmt.Errorf("g.createMakeMkTarget: %w", err)
    }
  }

  return nil
}

func (g *Gqlgen) createMakeMkTarget(makeMkTemplate string) error {
  makeMkPath := filepath.Join(g.workDirPath, "make.mk")

  if err := filer.AppendStringToFile(makeMkPath, makeMkTemplate); err != nil {
    return fmt.Errorf("filer.AppendStringToFile: %w", err)
  }
  return nil
}

func (g *Gqlgen) createMakefileIfNotExist() error {
  filePath := filepath.Join(g.workDirPath, "Makefile")

  if !filer.IsExistedFile(filePath) {
    if err := templater.ExecTemplateCopy(templates.ProjectMakefile, filePath, nil, nil); err != nil {
      return fmt.Errorf("execTemplateCopy: %w", err)
    }
  }
  return nil
}

func (g *Gqlgen) createGqlgenConfig() error {
  folderPath, err := filer.CreateNestedFolders(g.workDirPath, ".config")
  if err != nil {
    return fmt.Errorf("filer.CreateNestedFolders: %w", err)
  }
  filePath := filepath.Join(folderPath, "gqlgen_config.yaml")

  if err := templater.CopyTemplate(templates.GqlgenConfig, filePath); err != nil {
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
