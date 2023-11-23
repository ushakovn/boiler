package config

import (
  "context"
  "fmt"
  "path/filepath"
  "text/template"

  "github.com/ushakovn/boiler/internal/pkg/filer"
  "github.com/ushakovn/boiler/internal/pkg/stringer"
  "github.com/ushakovn/boiler/internal/pkg/templater"
  "github.com/ushakovn/boiler/templates"
)

type GenConfig struct {
  workDirPath string
}

// Config
// used only for type assertion in gens factory
type Config struct{}

func NewGenConfig(_ Config) (*GenConfig, error) {
  workDirPath, err := filer.WorkDirPath()
  if err != nil {
    return nil, err
  }
  return &GenConfig{
    workDirPath: workDirPath,
  }, nil
}

func (g *GenConfig) Init(_ context.Context) error {
  configFolder, err := filer.CreateNestedFolders(g.workDirPath, ".boiler")
  if err != nil {
    return fmt.Errorf("filer.CreateNestedFolders: %w", err)
  }
  configPath := filepath.Join(configFolder, "config.yaml")

  if err = templater.ExecTemplateCopy(templates.GenConfigEmpty, configPath, nil, nil); err != nil {
    return fmt.Errorf("execTemplateCopy: %w", err)
  }
  return nil
}

func (g *GenConfig) Generate(ctx context.Context) error {
  // TODO

  configDesc, err := g.loadGenConfigDesc()
  if err != nil {
    return fmt.Errorf("g.loadGenConfigDesc: %w", err)
  }

  configFolder, err := filer.CreateNestedFolders(g.workDirPath, "internal", "config")
  if err != nil {
    return fmt.Errorf("filer.CreateNestedFolders: %w", err)
  }

  templatesFuncMap := template.FuncMap{
    "toLowerCamelCase": stringer.StringToLowerCamelCase,
    "toUpperCamelCase": stringer.StringToUpperCamelCase,
    "toSnakeCase":      stringer.StringToSnakeCase,
  }

  for _, cnf := range configTemplates {
    filePath := filepath.Join(configFolder, cnf.fileName)

    if err = templater.ExecTemplateCopyWithGoFmt(cnf.compiledTemplate, filePath, configDesc, templatesFuncMap); err != nil {
      return fmt.Errorf("execTemplateCopy: %w", err)
    }
  }

  return nil
}

type configTemplate struct {
  fileName         string
  compiledTemplate string
}

var configTemplates = []*configTemplate{
  {
    fileName:         "groups.go",
    compiledTemplate: "", // TODO
  },
  {
    fileName:         "config.go",
    compiledTemplate: "", // TODO
  },
}
