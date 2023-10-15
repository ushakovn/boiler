package storage

import (
  "bytes"
  "context"
  "fmt"
  "go/format"
  "os"
  "path"
  "path/filepath"
  "text/template"

  "github.com/ushakovn/boiler/internal/boiler/gen"
  "github.com/ushakovn/boiler/internal/pkg/sql"
  "github.com/ushakovn/boiler/pkg/utils"
  "github.com/ushakovn/boiler/templates"
)

type storage struct {
  storageName string
  dumpSQL     *sql.DumpSQL
  schemaDesc  *schemaDesc
  workDirPath string
}

type Config struct {
  StorageName  string
  PgConfigPath string
  PgDumpPath   string
}

func (c *Config) Validate() error {
  if c.StorageName == "" {
    return fmt.Errorf("storage name must be specified")
  }
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
    storageName: config.StorageName,
    dumpSQL:     dumpSQL,
    workDirPath: workDirPath,
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
    "toLowerCamelCase": utils.StringToLowerCase,
    "toUpperCamelCase": utils.StringToUpperCamelCase,
    "withDot":          utils.StringWithDotPrefix,
  }

  type storageCommonTemplate struct {
    templateName     string
    compiledTemplate string
    fileNameBuilder  func() string
  }

  commonTemplates := []*storageCommonTemplate{
    {
      templateName:     "Options",
      compiledTemplate: templates.Options,
      fileNameBuilder:  buildOptionsFileName,
    },
    {
      templateName:     "Consts",
      compiledTemplate: templates.Consts,
      fileNameBuilder:  buildConstsFileName,
    },
    {
      templateName:     "Models",
      compiledTemplate: templates.Models,
      fileNameBuilder:  buildModelsFileName,
    },
  }

  for _, commonTemplate := range commonTemplates {
    filePath := path.Join(storagePath, commonTemplate.fileNameBuilder())

    if err := executeTemplateCopy(commonTemplate.compiledTemplate, filePath, g.schemaDesc, templatesFuncMap); err != nil {
      return fmt.Errorf("executeTemplateCopy templates.%s: %w", commonTemplate.templateName, err)
    }
  }

  for _, model := range g.schemaDesc.Models {
    interfaceFileName := buildStorageInterfaceFileName(model.ModelName)
    interfaceFilePath := path.Join(storagePath, interfaceFileName)
    if err := executeTemplateCopy(templates.Interface, interfaceFilePath, model, templatesFuncMap); err != nil {
      return fmt.Errorf("executeTemplateCopy templates.Interface: %w", err)
    }

    implementationFileName := buildStorageImplementationFileName(model.ModelName)
    implementationFolderName := buildStorageImplementationFolderName(model.ModelName)
    implementationFolderPath, err := createStorageFolders(storagePath, implementationFolderName)
    if err != nil {
      return fmt.Errorf("createStorageFolders: %w", err)
    }

    implementationFilePath := path.Join(implementationFolderPath, implementationFileName)
    if err := executeTemplateCopy(templates.Implementation, implementationFilePath, model, templatesFuncMap); err != nil {
      return fmt.Errorf("executeTemplateCopy templates.Interface: %w", err)
    }
  }

  return nil
}

func buildConstsFileName() string {
  return "consts.go"
}

func buildOptionsFileName() string {
  return "options.go"
}

func buildModelsFileName() string {
  return "models.go"
}

func buildStorageImplementationFileName(modelName string) string {
  modelName = utils.StringToLowerCase(modelName)
  return fmt.Sprint(modelName, "_implementation.go")
}

func buildStorageImplementationFolderName(modelName string) string {
  return utils.StringToLowerCase(modelName)
}

func buildStorageInterfaceFileName(modelName string) string {
  modelName = utils.StringToLowerCase(modelName)
  return fmt.Sprint(modelName, "_interface.go")
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
