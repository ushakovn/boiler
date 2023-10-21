package storage

import (
  "bytes"
  "context"
  "fmt"
  "go/format"
  "os"
  "path/filepath"
  "text/template"

  "github.com/ushakovn/boiler/internal/boiler/gen"
  "github.com/ushakovn/boiler/internal/pkg/sql"
  "github.com/ushakovn/boiler/pkg/utils"
  "github.com/ushakovn/boiler/templates"
)

type storage struct {
  dumpSQL      *sql.DumpSQL
  schemaDesc   *schemaDesc
  workDirPath  string
  goModuleName string
}

type Config struct {
  PgConfigPath string
  PgDumpPath   string
}

func (c *Config) Validate() error {
  if (c.PgDumpPath == "" && c.PgConfigPath == "") || (c.PgDumpPath != "" && c.PgConfigPath != "") {
    return fmt.Errorf("pg dump path OR pg config path must be specified")
  }
  return nil
}

func NewStorage(config Config) (gen.Generator, error) {
  if err := config.Validate(); err != nil {
    return nil, err
  }
  workDirPath, err := utils.Env("PWD")
  if err != nil {
    return nil, err
  }
  goModuleName, err := utils.ExtractGoModuleName(workDirPath)
  if err != nil {
    return nil, err
  }
  var (
    filePath string
    option   sql.PgDumpOption
  )
  if filePath = config.PgConfigPath; filePath != "" {
    option = sql.NewPgDumpOption(sql.WithPgConfigFile, filePath)
  }
  if filePath = config.PgDumpPath; filePath != "" {
    option = sql.NewPgDumpOption(sql.WithPgDumpFile, filePath)
  }
  dumpSQL, err := sql.DumpSchemaSQL(option)
  if err != nil {
    return nil, fmt.Errorf("sql.DumpSchemaSQL: %w", err)
  }

  return &storage{
    dumpSQL:      dumpSQL,
    workDirPath:  workDirPath,
    goModuleName: goModuleName,
  }, nil
}

func (g *storage) Generate(ctx context.Context) error {
  if err := g.loadSchemaDesc(); err != nil {
    return fmt.Errorf("loadSchemaDesc: %w", err)
  }

  storagePath, err := createStorageFolders(g.workDirPath, "internal", "pkg", "storage")
  if err != nil {
    return fmt.Errorf("g.createStorageFolders: %w", err)
  }

  templatesFuncMap := template.FuncMap{
    "toLowerCamelCase": utils.StringToLowerCamelCase,
    "toUpperCamelCase": utils.StringToUpperCamelCase,
    "toSnakeCase":      utils.StringToSnakeCase,
  }

  for _, commonTemplate := range storageCommonTemplates {
    filePath, err := createStorageFolders(storagePath, commonTemplate.filePathParts...)
    if err != nil {
      return fmt.Errorf("createStorageFolders: %w", err)
    }
    filePath = filepath.Join(filePath, commonTemplate.fileNameBuild(""))

    if err := executeTemplateCopy(commonTemplate.compiledTemplate, filePath, g.schemaDesc, templatesFuncMap); err != nil {
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

      if err := executeTemplateCopy(modelTemplate.compiledTemplate, filePath, model, templatesFuncMap); err != nil {
        return fmt.Errorf("executeTemplateCopy templates.%s: %w", modelTemplate.templateName, err)
      }
    }
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

func executeTemplateCopy(templateCompiled, filePath string, structPtr any, funcMap template.FuncMap) error {
  t := template.New("")
  if len(funcMap) > 0 {
    t = t.Funcs(funcMap)
  }
  t, err := t.Parse(templateCompiled)
  if err != nil {
    return fmt.Errorf("template.New().Parse: %w", err)
  }
  var (
    buffer bytes.Buffer
    buf    []byte
  )
  if err = t.Execute(&buffer, structPtr); err != nil {
    return fmt.Errorf("t.Execute: %w", err)
  }
  if buf, err = format.Source(buffer.Bytes()); err != nil {
    return fmt.Errorf("format.Source: %w", err)
  }
  if err = os.WriteFile(filePath, buf, os.ModePerm); err != nil {
    return fmt.Errorf("os.WriteFile: %w", err)
  }
  return nil
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
    compiledTemplate: templates.Interface,
    fileNameBuild: func(modelName string) string {
      modelName = utils.StringToSnakeCase(modelName)
      return fmt.Sprint(modelName, ".ifc.go")
    },
  },
  {
    templateName:     "Implementation",
    compiledTemplate: templates.Implementation,
    fileNameBuild: func(modelName string) string {
      modelName = utils.StringToSnakeCase(modelName)
      return fmt.Sprint(modelName, ".imp.go")
    },
  },
}

var storageCommonTemplates = []*storageTemplate{
  {
    templateName:     "Builders",
    compiledTemplate: templates.Builders,
    filePathParts:    []string{"client"},
    fileNameBuild: func(modelName string) string {
      return "builders.storage.go"
    },
  },
  {
    templateName:     "Client",
    compiledTemplate: templates.Client,
    filePathParts:    []string{"client"},
    fileNameBuild: func(modelName string) string {
      return "client.storage.go"
    },
  },
  {
    templateName:     "Options",
    compiledTemplate: templates.Options,
    fileNameBuild: func(modelName string) string {
      return "options.storage.go"
    },
  },
  {
    templateName:     "Consts",
    compiledTemplate: templates.Consts,
    fileNameBuild: func(modelName string) string {
      return "storage.storage.go"
    },
  },
  {
    templateName:     "Models",
    compiledTemplate: templates.Models,
    filePathParts:    []string{"models"},
    fileNameBuild: func(modelName string) string {
      return "models.storage.go"
    },
  },
}
