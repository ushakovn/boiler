package storage

import (
  "fmt"
  "path/filepath"

  "github.com/ushakovn/boiler/internal/pkg/templater"
  "github.com/ushakovn/boiler/templates"
)

type customSchemaDesc struct {
  CustomModels         []*customModelDesc
  CustomModelsPackages []*goPackageDesc
}

type customModelDesc struct {
  ModelName              string
  StructDescription      string
  ModelPackages          []*goPackageDesc
  InterfacePackages      []*goPackageDesc
  ImplementationPackages []*goPackageDesc
}

func (g *Storage) generateCustomStorages() error {
  customSchema, err := g.buildCustomSchemaDesc()
  if err != nil {
    return fmt.Errorf("g.buildCustomSchemaDesc: %w", err)
  }
  storagePath := filepath.Join(g.workDirPath, "internal", "pkg", "storage")

  filePath := filepath.Join(storagePath, "models", "custom_models.go")

  if err = templater.ExecTemplateCopyWithGoFmt(templates.StorageCustomModels, filePath, customSchema, nil); err != nil {
    return fmt.Errorf("executeTemplateCopy templates.StorageCustomModels: %w", err)
  }

  for _, model := range customSchema.CustomModels {

    modelTemplates, ok := storageTemplatesByCustomModelNames[model.ModelName]
    if !ok {
      return fmt.Errorf("storage templates not found for custom model: %s", model.ModelName)
    }

    for _, modelTemplate := range modelTemplates {

      filePath, err = createStorageFolders(storagePath, modelTemplate.filePathParts...)
      if err != nil {
        return fmt.Errorf("createStorageFolders: %w", err)
      }
      filePath = filepath.Join(filePath, modelTemplate.fileNameBuild(model.ModelName))

      execFunc := templater.ExecTemplateCopyFunc(modelTemplate.notGoTemplate)

      if err = execFunc(modelTemplate.compiledTemplate, filePath, model, nil); err != nil {
        return fmt.Errorf("executeTemplateCopy templates.%s: %w", modelTemplate.templateName, err)
      }
    }
  }
  return nil
}

func (g *Storage) buildCustomSchemaDesc() (*customSchemaDesc, error) {
  modelNames := []string{
    rocketLock,
  }
  customModelsDesc := make([]*customModelDesc, 0, len(modelNames))

  for _, modelName := range modelNames {
    customModel, err := g.buildCustomModelDesc(modelName)
    if err != nil {
      return nil, fmt.Errorf("g.buildCustomModelDesc: %w", err)
    }
    customModelsDesc = append(customModelsDesc, customModel)
  }

  // TODO: collect distinct packages
  commonPackages := customModelsDesc[0].ModelPackages

  return &customSchemaDesc{
    CustomModels:         customModelsDesc,
    CustomModelsPackages: commonPackages,
  }, nil
}

func (g *Storage) buildCustomModelDesc(modelName string) (*customModelDesc, error) {
  packagesNames, ok := packagesNamesByCustomModelsNames[modelName]
  if !ok {
    return nil, fmt.Errorf("imports not found for custom model: %s", modelName)
  }
  structDesc, ok := structDescByCustomModelNames[modelName]
  if !ok {
    return nil, fmt.Errorf("struct description not found for custom model: %s", modelName)
  }

  interfacePackages := mergeGoPackages(
    buildPackagesForNames(packagesNames.Interface),
    buildCrossFilePackages(g.goModuleName, interfaceFileName),
  )
  implementationPackages := mergeGoPackages(
    buildPackagesForNames(packagesNames.Implementation),
    buildCrossFilePackages(g.goModuleName, implementationFileName),
  )
  modelPackages := buildPackagesForNames(packagesNames.Model)

  return &customModelDesc{
    ModelName:              modelName,
    StructDescription:      structDesc,
    ModelPackages:          modelPackages,
    InterfacePackages:      interfacePackages,
    ImplementationPackages: implementationPackages,
  }, nil
}

type customModelPackagesNames struct {
  Model          []string
  Interface      []string
  Implementation []string
}

const (
  rocketLock = "boiler_rocket_lock"
)

var rocketLockPackagesNames = &customModelPackagesNames{
  Model: []string{
    timePackageName,
  },
  Interface: []string{
    errorsPackageName,
    timePackageName,
    contextPackageName,
  },
  Implementation: []string{
    contextPackageName,
    errorsPackageName,
    logrusPackageName,
    fmtPackageName,
  },
}

var structDescByCustomModelNames = map[string]string{
  rocketLock: templates.StorageRocketLockModel,
}

var packagesNamesByCustomModelsNames = map[string]*customModelPackagesNames{
  rocketLock: rocketLockPackagesNames,
}

var storageTemplatesByCustomModelNames = map[string][]*storageTemplate{
  rocketLock: rocketLockStorageTemplates,
}

var rocketLockStorageTemplates = []*storageTemplate{
  {
    templateName:     "Interface",
    compiledTemplate: templates.StorageRocketLockInterface,
    fileNameBuild:    buildInterfaceFileName,
  },
  {
    templateName:     "Implementation",
    compiledTemplate: templates.StorageRocketLockImplementation,
    fileNameBuild:    buildImplementationFileName,
  },
  {
    templateName:     "Migration",
    compiledTemplate: templates.StorageRocketLockMigration,

    filePathParts: []string{"../../../migrations"},

    fileNameBuild: func(modelName string) string {
      return "boiler_rocket_lock.sql"
    },

    notGoTemplate: true,
  },
}
