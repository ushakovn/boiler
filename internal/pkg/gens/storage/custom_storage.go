package storage

import (
  "fmt"
  "path/filepath"
  "strings"

  "github.com/ushakovn/boiler/internal/pkg/filer"
  "github.com/ushakovn/boiler/internal/pkg/goose"
  "github.com/ushakovn/boiler/internal/pkg/templater"
  "github.com/ushakovn/boiler/templates"
)

type customSchemaDesc struct {
  CustomModels []*customModelDesc
}

type customModelDesc struct {
  ModelName            string
  StructDescription    string
  ModelPackages        []*goPackageDesc
  ModelOptionsPackages []*goPackageDesc
  ModelMethodsPackages []*goPackageDesc
}

func (g *Storage) generateCustomStorages() error {
  customSchema, err := g.buildCustomSchemaDesc()
  if err != nil {
    return fmt.Errorf("g.buildCustomSchemaDesc: %w", err)
  }
  storagePath := filepath.Join(g.workDirPath, "internal", "pkg", "storage")

  for _, model := range customSchema.CustomModels {
    modelTemplates, ok := storageTemplatesByCustomModelNames[model.ModelName]
    if !ok {
      return fmt.Errorf("storage templates not found for custom model: %s", model.ModelName)
    }
    for _, modelTemplate := range modelTemplates {
      filePath, err := createStorageFolders(storagePath, modelTemplate.filePathParts...)
      if err != nil {
        return fmt.Errorf("createStorageFolders: %w", err)
      }
      filePath = filepath.Join(filePath, modelTemplate.fileNameBuild(model.ModelName))

      if modelTemplate.preBuildCheck != nil && !modelTemplate.preBuildCheck(filePath) {
        continue
      }
      execFunc := templater.ExecTemplateCopyFunc(modelTemplate.notGoTemplate)

      if err = execFunc(modelTemplate.compiledTemplate, filePath, model, nil); err != nil {
        return fmt.Errorf("executeTemplateCopy templates.%s: %w", modelTemplate.templateName, err)
      }
    }
  }
  return nil
}

func (g *Storage) buildCustomSchemaDesc() (*customSchemaDesc, error) {
  models := make([]*customModelDesc, 0, len(customModelsNames))

  packagesNames := map[string]struct{}{}
  var modelsPackages []*goPackageDesc

  for _, modelName := range customModelsNames {
    model, err := g.buildCustomModelDesc(modelName)
    if err != nil {
      return nil, fmt.Errorf("g.buildCustomModelDesc: %w", err)
    }
    models = append(models, model)

    for _, goPackage := range model.ModelPackages {
      if _, ok := packagesNames[goPackage.CustomName]; ok {
        continue
      }
      modelsPackages = append(modelsPackages, goPackage)
      packagesNames[goPackage.CustomName] = struct{}{}
    }
  }

  return &customSchemaDesc{
    CustomModels: models,
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

  modelOptionsPackages := mergeGoPackages(
    buildPackagesForNames(packagesNames.ModelOptions),
    buildCrossFilePackages(g.goModuleName, modelOptionsFileName),
  )
  modelMethodsPackages := mergeGoPackages(
    buildPackagesForNames(packagesNames.ModelMethods),
    buildCrossFilePackages(g.goModuleName, modelMethodsFileName),
  )
  modelPackages := buildPackagesForNames(packagesNames.Model)

  return &customModelDesc{
    ModelName:            modelName,
    StructDescription:    structDesc,
    ModelPackages:        modelPackages,
    ModelOptionsPackages: modelOptionsPackages,
    ModelMethodsPackages: modelMethodsPackages,
  }, nil
}

type customModelPackagesNames struct {
  Model        []string
  ModelOptions []string
  ModelMethods []string
}

var customModelsNames = []string{
  rocketLock,
}

const (
  rocketLock = "rocket_lock"
)

var rocketLockPackagesNames = &customModelPackagesNames{
  Model: []string{
    timePackageName,
  },
  ModelOptions: []string{
    errorsPackageName,
    timePackageName,
  },
  ModelMethods: []string{
    contextPackageName,
    errorsPackageName,
    logrusPackageName,
    fmtPackageName,
    pgClientPackageName,
    pgErrorsPackageName,
    pgBuilderPackageName,
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
    templateName:     "model_options",
    compiledTemplate: templates.StorageRocketLockModelOptions,
    fileNameBuild:    buildModelOptionsFileName,
  },
  {
    templateName:     "model_methods",
    compiledTemplate: templates.StorageRocketLockModelMethods,
    fileNameBuild:    buildModelMethodsFileName,
  },
  {
    templateName:     "model",
    compiledTemplate: templates.StorageCustomModel,
    filePathParts:    []string{"models"},
    fileNameBuild:    buildModelFileName,
  },
  {
    templateName:     "migration",
    compiledTemplate: templates.StorageRocketLockMigration,

    filePathParts: []string{"../../../migrations"},

    fileNameBuild: func(modelName string) string {
      gooseFileName := goose.BuildFileName(modelName)
      return gooseFileName
    },

    preBuildCheck: func(filePath string) bool {
      gooseFileName := filer.ExtractFileName(filePath)
      fileName := goose.ExtractFileName(gooseFileName)

      fileDirPath := strings.TrimSuffix(filePath, fmt.Sprint("/", gooseFileName))
      return !filer.IsExistedPattern(fileDirPath, fileName)
    },

    notGoTemplate: true,
  },
}
