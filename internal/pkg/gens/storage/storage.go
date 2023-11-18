package storage

import (
  "context"
  "fmt"
  "os"
  "path/filepath"
  "text/template"

  "github.com/ushakovn/boiler/internal/pkg/filer"
  "github.com/ushakovn/boiler/internal/pkg/sql"
  "github.com/ushakovn/boiler/internal/pkg/stringer"
  "github.com/ushakovn/boiler/internal/pkg/templater"
  "github.com/ushakovn/boiler/templates"
)

type Storage struct {
  dumpSQL      *sql.DumpSQL
  schemaDesc   *schemaDesc
  workDirPath  string
  goModuleName string
}

type Config struct {
  PgConfigPath string
  PgDumpPath   string
}

func NewStorage(config Config) (*Storage, error) {
  workDirPath, err := filer.WorkDirPath()
  if err != nil {
    return nil, err
  }
  goModuleName, err := filer.ExtractGoModuleName(workDirPath)
  if err != nil {
    return nil, err
  }

  var (
    dumpSQL *sql.DumpSQL
    option  sql.PgDumpOption
  )
  if filePath := config.PgConfigPath; filePath != "" {
    option = sql.NewPgDumpOption(sql.WithPgConfigFile, filePath)
  }
  if filePath := config.PgDumpPath; filePath != "" {
    option = sql.NewPgDumpOption(sql.WithPgDumpFile, filePath)
  }
  if option != nil {
    // If option was set
    dumpSQL, err = sql.DumpSchemaSQL(context.Background(), option)
    if err != nil {
      return nil, fmt.Errorf("sql.DumpSchemaSQL: %w", err)
    }
  }

  return &Storage{
    dumpSQL:      dumpSQL,
    workDirPath:  workDirPath,
    goModuleName: goModuleName,
  }, nil
}

func (g *Storage) Init(_ context.Context) error {
  if err := g.createPgConfig(); err != nil {
    return fmt.Errorf("g.createPgConfig: %w", err)
  }
  return nil
}

func (g *Storage) Generate(_ context.Context) error {
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
  filePath := filepath.Join(g.workDirPath, "pg_config.yaml")

  if err := templater.ExecTemplateCopy(templates.StorageConfig, filePath, nil, nil); err != nil {
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
  filePathParts    []string
  fileNameBuild    func(modelName string) string
}

var storageModelTemplates = []*storageTemplate{
  {
    templateName:     "Interface",
    compiledTemplate: templates.StorageInterface,
    fileNameBuild: func(modelName string) string {
      modelName = stringer.StringToSnakeCase(modelName)
      return fmt.Sprint(modelName, ".interface.go")
    },
  },
  {
    templateName:     "Implementation",
    compiledTemplate: templates.StorageImplementation,
    fileNameBuild: func(modelName string) string {
      modelName = stringer.StringToSnakeCase(modelName)
      return fmt.Sprint(modelName, ".implementation.go")
    },
  },
}

var storageCommonTemplates = []*storageTemplate{
  {
    templateName:     "Builders",
    compiledTemplate: templates.StorageBuilders,
    filePathParts:    []string{"client"},
    fileNameBuild: func(modelName string) string {
      return "builders.go"
    },
  },
  {
    templateName:     "Client",
    compiledTemplate: templates.StorageClient,
    filePathParts:    []string{"client"},
    fileNameBuild: func(modelName string) string {
      return "client.go"
    },
  },
  {
    templateName:     "Options",
    compiledTemplate: templates.StorageOptions,
    fileNameBuild: func(modelName string) string {
      return "options.go"
    },
  },
  {
    templateName:     "Consts",
    compiledTemplate: templates.StorageConsts,
    fileNameBuild: func(modelName string) string {
      return "consts.go"
    },
  },
  {
    templateName:     "Models",
    compiledTemplate: templates.StorageModels,
    filePathParts:    []string{"models"},
    fileNameBuild: func(modelName string) string {
      return "models.go"
    },
  },
}
