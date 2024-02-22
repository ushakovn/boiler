package kafkaoutbox

import (
  "context"
  "fmt"
  "path/filepath"
  "text/template"
  "time"

  log "github.com/sirupsen/logrus"
  "github.com/ushakovn/boiler/internal/pkg/executor"
  "github.com/ushakovn/boiler/internal/pkg/filer"
  "github.com/ushakovn/boiler/internal/pkg/goose"
  "github.com/ushakovn/boiler/internal/pkg/makefile"
  "github.com/ushakovn/boiler/internal/pkg/stringer"
  "github.com/ushakovn/boiler/internal/pkg/templater"
  "github.com/ushakovn/boiler/templates"
)

type Kafkaoutbox struct {
  workDirPath  string
  goModuleName string
  outboxDesc   *outboxDesc
}

type Config struct{}

var (
  dataPtrStub = any(nil)
  funcMapStub = template.FuncMap(nil)
)

func NewKafkaoutbox(_ Config) (*Kafkaoutbox, error) {
  workDirPath, err := filer.WorkDirPath()
  if err != nil {
    return nil, err
  }
  goModuleName, err := filer.ExtractGoModuleName(workDirPath)
  if err != nil {
    return nil, err
  }
  return &Kafkaoutbox{
    workDirPath:  workDirPath,
    goModuleName: goModuleName,
  }, nil
}

func (g *Kafkaoutbox) Init(_ context.Context) error {
  if err := g.generateConfigTemplate(); err != nil {
    return fmt.Errorf("generateConfigTemplate: %w", err)
  }
  if err := g.generateProtoTemplates(); err != nil {
    return fmt.Errorf("generateProtoTemplates: %w", err)
  }
  return nil
}

func (g *Kafkaoutbox) Generate(ctx context.Context) error {
  if err := g.loadOutboxDesc(); err != nil {
    return fmt.Errorf("loadOutboxDesc: %w", err)
  }
  if err := g.generateMakeTargets(); err != nil {
    return fmt.Errorf("generateMakeTargets: %w", err)
  }
  if err := g.generateOutboxPbGo(ctx); err != nil {
    return fmt.Errorf("generateOutboxPbGo: %w", err)
  }
  if err := g.generateMigrationTemplates(); err != nil {
    return fmt.Errorf("generateMigrationTemplates: %w", err)
  }
  if err := g.generateGoTemplates(); err != nil {
    return fmt.Errorf("generateGoTemplates: %w", err)
  }
  return nil
}

func (g *Kafkaoutbox) loadOutboxDesc() error {
  desc, err := g.buildOutbox()
  if err != nil {
    return fmt.Errorf("buildOutbox: %w", err)
  }
  g.outboxDesc = desc
  return nil
}

func (g *Kafkaoutbox) generateConfigTemplate() error {
  fileDir, err := filer.CreateNestedFolders(g.workDirPath, ".config")
  if err != nil {
    return fmt.Errorf("filer.CreateNestedFolders: %w", err)
  }
  const fileName = "kafkaoutbox_config.yaml"

  filePath := filepath.Join(fileDir, fileName)

  if err = templater.CopyTemplate(templates.KafkaOutboxConfig, filePath); err != nil {
    return fmt.Errorf("templater.CopyTemplate: %w", err)
  }
  return nil
}

func (g *Kafkaoutbox) generateProtoTemplates() error {
  outboxProto := g.buildOutboxProto()

  for _, temp := range protoTemplates {
    fileName := temp.buildFileName(temp.name)

    fileDir, err := filer.CreateNestedFolders(g.workDirPath, "api", outboxProto.ServiceName)
    if err != nil {
      return fmt.Errorf("filer.CreateNestedFolders: %w", err)
    }
    filePath := filepath.Join(fileDir, fileName)

    if err = templater.ExecTemplateCopy(temp.compiled, filePath, outboxProto, funcMapStub); err != nil {
      return fmt.Errorf("templater.ExecTemplateCopyWithGoFmt: %w", err)
    }
  }
  return nil
}

func (g *Kafkaoutbox) generateGoTemplates() error {
  if g.outboxDesc == nil {
    return fmt.Errorf("outbox description is a nil")
  }
  fileDir, err := filer.CreateNestedFolders(g.workDirPath, "internal", "pkg", "kafkaoutbox")
  if err != nil {
    return fmt.Errorf("filer.CreateNestedFolders: %w", err)
  }
  funcMap := template.FuncMap{
    "toSnakeCase":      stringer.StringToSnakeCase,
    "toLowerCamelCase": stringer.StringToLowerCamelCase,
  }
  for _, temp := range goTemplates {
    fileName := temp.buildFileName(temp.name)
    filePath := filepath.Join(fileDir, fileName)

    if err = templater.ExecTemplateCopyWithGoFmt(temp.compiled, filePath, g.outboxDesc, funcMap); err != nil {
      return fmt.Errorf("templater.ExecTemplateCopy: %w", err)
    }
  }
  return nil
}

func (g *Kafkaoutbox) generateMigrationTemplates() error {
  if g.outboxDesc == nil || g.outboxDesc.OutboxTables == nil {
    return fmt.Errorf("outbox tables description is a nil")
  }

  fileDir, err := filer.CreateNestedFolders(g.workDirPath, "migrations")
  if err != nil {
    return fmt.Errorf("filer.CreateNestedFolders: %w", err)
  }

  var lastTime time.Time

  fixFileName := func(fileName string) (string, error) {
    return goose.ChangeTime(fileName, func(fileTime time.Time) time.Time {
      if lastTime.Equal(time.Time{}) {
        lastTime = fileTime
      }
      lastTime = lastTime.Add(time.Second)

      return lastTime
    })
  }

  for _, temp := range migrationCommonTemplates {
    fileName := temp.buildFileName(temp.name)
    if err != nil {
      return fmt.Errorf("goose.ChangeTime: %w", err)
    }
    if fileName, err = fixFileName(fileName); err != nil {
      return fmt.Errorf("fixFileName: %w", err)
    }
    filePath := filepath.Join(fileDir, fileName)

    if ok := temp.buildCheck(filePath); !ok {
      continue
    }
    if err = templater.ExecTemplateCopy(temp.compiled, filePath, dataPtrStub, funcMapStub); err != nil {
      return fmt.Errorf("templater.ExecTemplateCopy: %w", err)
    }
  }

  for _, outboxTable := range g.outboxDesc.OutboxTables {
    tableName := stringer.StringToSnakeCase(outboxTable.SourceTableName)

    for _, temp := range migrationTableTemplates {
      fileName := temp.buildFileName(tableName)

      if fileName, err = fixFileName(fileName); err != nil {
        return fmt.Errorf("fixFileName: %w", err)
      }
      filePath := filepath.Join(fileDir, fileName)

      if ok := temp.buildCheck(filePath); !ok {
        continue
      }
      if err = templater.ExecTemplateCopy(temp.compiled, filePath, outboxTable, funcMapStub); err != nil {
        return fmt.Errorf("templater.ExecTemplateCopy: %w", err)
      }
    }
  }

  return nil
}

func (g *Kafkaoutbox) generateOutboxPbGo(ctx context.Context) error {
  for _, target := range makeMkTargets {
    out, err := executor.ExecCmdCtxWithOut(ctx, "make", target.targetName)
    if err != nil {
      log.Errorf(`boiler: try add "include make.mk" to your Makefile if "make rule not found"`)
      return fmt.Errorf("executor.ExecCmdCtx: target: %s: %w", target.targetName, err)
    }
    log.Debugf("boiler: executed make target: %s\n%s\n", target.targetName, string(out))
  }
  return nil
}

func (g *Kafkaoutbox) generateMakeTargets() error {
  if err := g.createMakeMkTargetsIfNotExist(); err != nil {
    return fmt.Errorf("createMakeMkTargetsIfNotExist: %w", err)
  }
  if err := g.createMakefileIfNotExist(); err != nil {
    return fmt.Errorf("createMakefileIfNotExist: %w", err)
  }
  return nil
}

func (g *Kafkaoutbox) createMakeMkTargetsIfNotExist() error {
  const fileName = "make.mk"
  filePath := filepath.Join(g.workDirPath, fileName)

  for _, target := range makeMkTargets {
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

func (g *Kafkaoutbox) createMakeMkTarget(makeMkTemplate string) error {
  const fileName = "make.mk"
  goPackageTrim := g.goModuleName

  templateData := map[string]any{
    "goPackageTrim": goPackageTrim,
  }
  executedBuf, err := templater.ExecTemplate(makeMkTemplate, templateData, funcMapStub)
  if err != nil {
    return fmt.Errorf("executeTemplate")
  }
  executedTarget := string(executedBuf)
  makeMkPath := filepath.Join(g.workDirPath, fileName)

  if err = filer.AppendStringToFile(makeMkPath, executedTarget); err != nil {
    return fmt.Errorf("filer.AppendStringToFile: %w", err)
  }
  return nil
}

func (g *Kafkaoutbox) createMakefileIfNotExist() error {
  const fileName = "Makefile"
  filePath := filepath.Join(g.workDirPath, fileName)

  if !filer.IsExistedFile(filePath) {
    if err := templater.ExecTemplateCopy(templates.ProjectMakefile, filePath, dataPtrStub, funcMapStub); err != nil {
      return fmt.Errorf("execTemplateCopy: %w", err)
    }
  }
  return nil
}

func (g *Kafkaoutbox) workDirFolder() string {
  return filer.ExtractFileName(g.workDirPath)
}
