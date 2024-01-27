package storage

import (
  "context"
  "fmt"
  "os"
  "path/filepath"
  "text/template"

  log "github.com/sirupsen/logrus"
  "github.com/ushakovn/boiler/internal/pkg/filer"
  "github.com/ushakovn/boiler/internal/pkg/pgdump"
  "github.com/ushakovn/boiler/internal/pkg/stringer"
  "github.com/ushakovn/boiler/internal/pkg/templater"
  "github.com/ushakovn/boiler/templates"
)

type Storage struct {
  config  *Config
  dumpSQL *pgdump.DumpSQL

  schemaDesc   *schemaDesc
  workDirPath  string
  goModuleName string
}

func NewStorage(configPath ConfigPath) (*Storage, error) {
  ctx := context.Background()

  workDirPath, err := filer.WorkDirPath()
  if err != nil {
    return nil, err
  }
  goModuleName, err := filer.ExtractGoModuleName(workDirPath)
  if err != nil {
    return nil, err
  }

  var (
    config  *Config
    dumpSQL *pgdump.DumpSQL
  )

  if configPath.String() != "" {
    // If config path specified

    config, err = configPath.Parse()
    if err != nil {
      log.Fatalf("storage config parsing failed: %v", err)
    }
    if err = config.Validate(); err != nil {
      log.Fatalf("storage config validation failed: %v", err)
    }
    dumpOption := buildDumpOption(ctx, config)

    // Dump schema SQL
    dumpSQL, err = dumpOption.Do()
    if err != nil {
      return nil, err
    }
  }

  return &Storage{
    config:  config,
    dumpSQL: dumpSQL,

    workDirPath:  workDirPath,
    goModuleName: goModuleName,
  }, nil
}

func buildDumpOption(ctx context.Context, config *Config) pgdump.DumpOption {
  option := pgdump.NewDumpOption()

  if config.PgConfig != nil {
    option = option.WithPgConfig(ctx, pgdump.PgConfig(*config.PgConfig))
  }
  if config.PgDumpPath != "" {
    option = option.WithPgDumpPath(config.PgDumpPath)
  }

  if len(config.PgTypeConfig) > 0 {
    pgTypes := make([]string, 0, len(config.PgTypeConfig))

    for pgType := range config.PgTypeConfig {
      pgTypes = append(pgTypes, pgType)
    }
    option = option.WithCustomTypes(pgTypes)
  }

  return option
}

func (g *Storage) Init(_ context.Context) error {
  if err := g.createPgConfig(); err != nil {
    return fmt.Errorf("g.createPgConfig: %w", err)
  }
  return nil
}

func (g *Storage) Generate(_ context.Context) error {
  if err := g.generateStorages(); err != nil {
    return fmt.Errorf("g.generateStorages: %w", err)
  }
  if err := g.generateCustomStorages(); err != nil {
    return fmt.Errorf("g.generateCustomStorages: %w", err)
  }
  return nil
}

func (g *Storage) generateStorages() error {
  if err := g.loadSchemaDesc(); err != nil {
    return fmt.Errorf("loadSchemaDesc: %w", err)
  }

  storagePath, err := createStorageFolders(g.workDirPath, "internal", "pkg", "storage")
  if err != nil {
    return fmt.Errorf("g.createStorageFolders: %w", err)
  }

  templatesFuncMap := template.FuncMap{
    "toLowerCamelCase": stringer.StringToLowerCamelCase,
    "toUpperCamelCase": stringer.StringToUpperCamelCase,
    "toSnakeCase":      stringer.StringToSnakeCase,
  }

  for _, commonTemplate := range storageCommonTemplates {
    filePath, err := createStorageFolders(storagePath, commonTemplate.filePathParts...)
    if err != nil {
      return fmt.Errorf("createStorageFolders: %w", err)
    }
    filePath = filepath.Join(filePath, commonTemplate.fileNameBuild(""))

    if err = templater.ExecTemplateCopyWithGoFmt(commonTemplate.compiledTemplate, filePath, g.schemaDesc, templatesFuncMap); err != nil {
      return fmt.Errorf("executeTemplateCopy templates.%s: %w", commonTemplate.templateName, err)
    }
  }

  for _, model := range g.schemaDesc.Models {
    for _, modelTemplate := range storageModelTemplates {
      filePath, err := createStorageFolders(storagePath, modelTemplate.filePathParts...)
      if err != nil {
        return fmt.Errorf("createStorageFolders: %w", err)
      }
      filePath = filepath.Join(filePath, modelTemplate.fileNameBuild(model.ModelName))

      if err = templater.ExecTemplateCopyWithGoFmt(modelTemplate.compiledTemplate, filePath, model, templatesFuncMap); err != nil {
        return fmt.Errorf("executeTemplateCopy templates.%s: %w", modelTemplate.templateName, err)
      }
    }
  }

  return nil
}

func (g *Storage) createPgConfig() error {
  folderPath, err := filer.CreateNestedFolders(g.workDirPath, ".config")
  if err != nil {
    return fmt.Errorf("filer.CreateNestedFolders: %w", err)
  }
  filePath := filepath.Join(folderPath, "storage_config.yaml")

  if err = templater.ExecTemplateCopy(templates.StorageConfig, filePath, nil, nil); err != nil {
    return fmt.Errorf("execTemplateCopy: %w", err)
  }
  return nil
}

func createStorageFolders(sourcePath string, destNestedFolders ...string) (string, error) {
  defaultDirParts := append([]string{sourcePath}, destNestedFolders...)
  defaultDir := filepath.Join(defaultDirParts...)

  prevDirParts := make([]string, 0, len(defaultDirParts))

  for _, dirPart := range defaultDirParts {
    // Create directories for storage package
    prevDirParts = append(prevDirParts, dirPart)
    curPath := filepath.Join(prevDirParts...)

    if _, err := os.Stat(curPath); os.IsNotExist(err) {
      if err = os.Mkdir(curPath, os.ModePerm); err != nil {
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

type storageTemplate struct {
  templateName     string
  compiledTemplate string

  filePathParts []string
  fileNameBuild func(modelName string) string

  preBuildCheck func(filePath string) bool
  notGoTemplate bool
}

var storageModelTemplates = []*storageTemplate{
  {
    templateName:     "model_options",
    compiledTemplate: templates.StorageModelOptions,
    fileNameBuild:    buildModelOptionsFileName,
  },
  {
    templateName:     "model_methods",
    compiledTemplate: templates.StorageModelMethods,
    fileNameBuild:    buildModelMethodsFileName,
  },
}

var storageCommonTemplates = []*storageTemplate{
  {
    templateName:     "options",
    compiledTemplate: templates.StorageOptions,
    fileNameBuild: func(modelName string) string {
      return "options.go"
    },
  },
  {
    templateName:     "consts",
    compiledTemplate: templates.StorageConsts,
    fileNameBuild: func(modelName string) string {
      return "consts.go"
    },
  },
  {
    templateName:     "models",
    compiledTemplate: templates.StorageModels,
    filePathParts:    []string{"models"},
    fileNameBuild: func(modelName string) string {
      return "models.go"
    },
  },
}

func buildModelOptionsFileName(modelName string) string {
  modelName = stringer.StringToSnakeCase(modelName)
  return fmt.Sprint(modelName, ".options.go")
}

func buildModelMethodsFileName(modelName string) string {
  modelName = stringer.StringToSnakeCase(modelName)
  return fmt.Sprint(modelName, ".methods.go")
}
