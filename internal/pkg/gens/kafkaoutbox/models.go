package kafkaoutbox

import (
  "fmt"
  "path/filepath"

  "github.com/ushakovn/boiler/internal/pkg/filext"
  "github.com/ushakovn/boiler/internal/pkg/goose"
  "github.com/ushakovn/boiler/internal/pkg/stringer"
  "github.com/ushakovn/boiler/templates"
)

const (
  outboxFileName  = "outbox"
  configFileName  = "config"
  modelsFileName  = "models"
  storageFileName = "storage"

  extensionUUIDOsspFileName = "extension_uuid_ossp"
  outboxTableFileName       = "kafka_outbox_table"
  outboxFuncFileName        = "kafka_outbox_func"
  outboxTriggerFileName     = "kafka_outbox_trigger"

  kafkaoutboxFileName = "kafkaoutbox"
  protoOptionFileName = "option"
)

const (
  contextPackageName = "context"
  fmtPackageName     = "fmt"
  timePackageName    = "time"
  osPackageName      = "os"
  yamlPackageName    = "yaml"
  errorsPackageName  = "errors"

  logrusPackageName   = "logrus"
  squirrelPackageName = "squirrel"

  reflectPackageName      = "reflect"
  protoreflectPackageName = "protoreflect"
  protojsonPackageName    = "protojson"

  ozzoValidationPackageName = "ozzo-validation"
  saramaIBMPackageName      = "IBM-samara"

  pgExecutorPackageName = "pg-executor"
)

type outboxDesc struct {
  OutboxTables          []*outboxTableDesc
  OutboxPackages        []*goPackageDesc
  OutboxConfigPackages  []*goPackageDesc
  OutboxModelsPackages  []*goPackageDesc
  OutboxStoragePackages []*goPackageDesc
}

type outboxTableDesc struct {
  SourceTableName   string
  OutboxProtoTyp    string
  OutboxTopicName   string
  OutboxTableName   string
  OutboxFuncName    string
  OutboxTriggerName string
}

type protoFileDesc struct {
  ServiceName string
  GoPackage   string
}

type goPackageDesc struct {
  CustomName  string
  ImportLine  string
  ImportAlias string
  IsBuiltin   bool
  IsInstall   bool
}

type templateDesc struct {
  name          string
  compiled      string
  buildCheck    func(filePath string) bool
  buildFileName func(fileName string) string
}

type makeMkTargetDesc struct {
  targetName       string
  compiledTemplate string
}

var protoTemplates = []*templateDesc{
  {
    name:          kafkaoutboxFileName,
    compiled:      templates.KafkaOutboxProto,
    buildCheck:    buildCheckStub,
    buildFileName: filext.Proto,
  },
  {
    name:          protoOptionFileName,
    compiled:      templates.KafkaOutboxProtoOption,
    buildCheck:    buildCheckStub,
    buildFileName: filext.Proto,
  },
}

var goTemplates = []*templateDesc{
  {
    name:          outboxFileName,
    compiled:      templates.KafkaOutbox,
    buildCheck:    buildCheckStub,
    buildFileName: filext.Go,
  },
  {
    name:          configFileName,
    compiled:      templates.KafkaOutboxConfig,
    buildCheck:    buildCheckStub,
    buildFileName: filext.Go,
  },
  {
    name:          modelsFileName,
    compiled:      templates.KafkaOutboxModels,
    buildCheck:    buildCheckStub,
    buildFileName: filext.Go,
  },
  {
    name:          storageFileName,
    compiled:      templates.KafkaOutboxStorage,
    buildCheck:    buildCheckStub,
    buildFileName: filext.Go,
  },
}

var migrationTableTemplates = []*templateDesc{
  {
    name:       outboxTableFileName,
    compiled:   templates.KafkaOutboxMigrationOutboxTable,
    buildCheck: buildCheckMigration,
    buildFileName: func(tableName string) string {
      fileName := fmt.Sprintf("%s_%s", outboxTableFileName, tableName)
      return goose.BuildFileName(fileName)
    },
  },
  {
    name:       outboxFuncFileName,
    compiled:   templates.KafkaOutboxMigrationOutboxFunc,
    buildCheck: buildCheckMigration,
    buildFileName: func(tableName string) string {
      fileName := fmt.Sprintf("%s_%s", outboxFuncFileName, tableName)
      return goose.BuildFileName(fileName)
    },
  },
  {
    name:       outboxTriggerFileName,
    compiled:   templates.KafkaOutboxMigrationOutboxTrigger,
    buildCheck: buildCheckMigration,
    buildFileName: func(tableName string) string {
      fileName := fmt.Sprintf("%s_%s", outboxTriggerFileName, tableName)
      return goose.BuildFileName(fileName)
    },
  },
}

var migrationCommonTemplates = []*templateDesc{
  {
    name:          extensionUUIDOsspFileName,
    compiled:      templates.KafkaOutboxMigrationUUIDOssp,
    buildCheck:    buildCheckMigration,
    buildFileName: goose.BuildFileName,
  },
}

var makeMkTargets = []*makeMkTargetDesc{
  {
    targetName:       templates.GrpcMakeMkBinDepsName,
    compiledTemplate: templates.GrpcMakeMkBinDeps,
  },
  {
    targetName:       templates.GrpcMakeMkGenerateName,
    compiledTemplate: templates.GrpcMakeMkGenerate,
  },
}

var packagesByFiles = map[string][]string{
  outboxFileName: {
    contextPackageName,
    fmtPackageName,
    timePackageName,
    logrusPackageName,
    reflectPackageName,
    protoreflectPackageName,
    protojsonPackageName,
    ozzoValidationPackageName,
    saramaIBMPackageName,
  },
  configFileName: {
    fmtPackageName,
    osPackageName,
    timePackageName,
    ozzoValidationPackageName,
    yamlPackageName,
  },
  modelsFileName: {
    timePackageName,
  },
  storageFileName: {
    contextPackageName,
    fmtPackageName,
    timePackageName,
    squirrelPackageName,
    pgExecutorPackageName,
  },
}

var packagesByNames = map[string]*goPackageDesc{
  reflectPackageName: {
    CustomName: "go/reflect",
    ImportLine: "reflect",
    IsBuiltin:  true,
  },
  protoreflectPackageName: {
    CustomName: "protobuf/protoreflect",
    ImportLine: "google.golang.org/protobuf/reflect/protoreflect",
    IsInstall:  true,
  },
  protojsonPackageName: {
    CustomName: "protobuf/protojson",
    ImportLine: "google.golang.org/protobuf/encoding/protojson",
    IsInstall:  true,
  },
  saramaIBMPackageName: {
    CustomName: "IBM/sarama",
    ImportLine: "github.com/IBM/sarama",
    IsInstall:  true,
  },
  ozzoValidationPackageName: {
    CustomName: "ozzo/validation",
    ImportLine: "github.com/go-ozzo/ozzo-validation",
    IsInstall:  true,
  },
  contextPackageName: {
    CustomName: "go/context",
    ImportLine: "context",
    IsBuiltin:  true,
  },
  fmtPackageName: {
    CustomName: "go/fmt",
    ImportLine: "fmt",
    IsBuiltin:  true,
  },
  timePackageName: {
    CustomName: "go/time",
    ImportLine: "time",
    IsBuiltin:  true,
  },
  osPackageName: {
    CustomName: "go/os",
    ImportLine: "os",
    IsBuiltin:  true,
  },
  yamlPackageName: {
    CustomName: "gopkg/yaml",
    ImportLine: "gop" +
      "kg.in/yaml.v3",
    IsInstall: true,
  },
  squirrelPackageName: {
    CustomName:  "masterminds/squirrel",
    ImportLine:  "github.com/Masterminds/squirrel",
    ImportAlias: "sq",
    IsInstall:   true,
  },
  errorsPackageName: {
    CustomName: "go/errors",
    ImportLine: "errors",
    IsBuiltin:  true,
  },
  logrusPackageName: {
    CustomName:  "sirupsen/logrus",
    ImportLine:  "github.com/sirupsen/logrus",
    ImportAlias: "log",
    IsInstall:   true,
  },
  pgExecutorPackageName: {
    CustomName:  "boiler/pg-executor",
    ImportLine:  "github.com/ushakovn/boiler/pkg/storage/postgres/executor",
    ImportAlias: "pg",
    IsInstall:   true,
  },
}

func (g *Kafkaoutbox) buildOutbox() (*outboxDesc, error) {
  outboxTables, err := g.buildOutboxTables()
  if err != nil {
    return nil, fmt.Errorf("buildOutboxTables: %w", err)
  }
  outboxPackages := buildTemplatePackages(outboxFileName)
  outboxPackages = append(outboxPackages, g.buildPbGoPackage())

  outbox := &outboxDesc{
    OutboxTables:          outboxTables,
    OutboxPackages:        outboxPackages,
    OutboxConfigPackages:  buildTemplatePackages(configFileName),
    OutboxModelsPackages:  buildTemplatePackages(modelsFileName),
    OutboxStoragePackages: buildTemplatePackages(storageFileName),
  }

  return outbox, nil
}

func (g *Kafkaoutbox) buildOutboxProto() *protoFileDesc {
  serviceName := g.buildProtoServiceName()

  goPackage := filepath.Join(g.goModuleName, "internal", "pb", serviceName)
  goPackageWithSuffix := fmt.Sprint(goPackage, ";", serviceName)

  return &protoFileDesc{
    ServiceName: serviceName,
    GoPackage:   goPackageWithSuffix,
  }
}

func (g *Kafkaoutbox) buildOutboxTables() ([]*outboxTableDesc, error) {
  filePath := g.buildProtoServicePath()
  optionName := g.buildProtoOptionName()

  parsed, err := parseProto(filePath, optionName)
  if err != nil {
    return nil, fmt.Errorf("proto parsing failed: %w", err)
  }

  if g.validateProto && g.storageGen != nil {
    if err = validateProto(parsed, g.storageGen.DumpSQL()); err != nil {
      return nil, fmt.Errorf("proto validation failed: %w", err)
    }
  }
  outboxTables := make([]*outboxTableDesc, 0, len(parsed.messages))

  for _, message := range parsed.messages {
    sourceTableName := message.tableName
    snakeTableName := stringer.StringToSnakeCase(sourceTableName)

    outboxTable := &outboxTableDesc{
      SourceTableName:   sourceTableName,
      OutboxProtoTyp:    buildOutboxProtoTyp(message.messageName),
      OutboxTopicName:   buildOutboxTopicName(snakeTableName),
      OutboxTableName:   buildOutboxTableName(snakeTableName),
      OutboxFuncName:    buildOutboxFuncName(snakeTableName),
      OutboxTriggerName: buildOutboxTriggerName(snakeTableName),
    }
    outboxTables = append(outboxTables, outboxTable)
  }
  return outboxTables, nil
}

func (g *Kafkaoutbox) buildProtoServicePath() string {
  serviceName := g.buildProtoServiceName()

  fileDir := filepath.Join(g.workDirPath, "api", serviceName)
  filePath := filext.Proto(filepath.Join(fileDir, kafkaoutboxFileName))

  return filePath
}

func (g *Kafkaoutbox) buildProtoServiceName() string {
  workDirFolder := stringer.StringToSnakeCase(g.workDirFolder())
  return fmt.Sprintf("%s_%s", kafkaoutboxFileName, workDirFolder)
}

func (g *Kafkaoutbox) buildProtoOptionName() string {
  serviceName := g.buildProtoServiceName()
  return fmt.Sprintf("(%s.table_name)", serviceName)
}

func (g *Kafkaoutbox) buildPbGoPackage() *goPackageDesc {
  path := filepath.Join(g.goModuleName, "internal", "pb", g.buildProtoServiceName())

  return &goPackageDesc{
    CustomName:  "pb/go",
    ImportLine:  path,
    ImportAlias: "desc",
  }
}

func buildOutboxProtoTyp(messageName string) string {
  return fmt.Sprintf("desc.%s{}", messageName)
}

func buildOutboxTableName(tableName string) string {
  return fmt.Sprintf("%s_outbox", tableName)
}

func buildOutboxTriggerName(tableName string) string {
  return fmt.Sprintf("%s_outbox_trigger", tableName)
}

func buildOutboxTopicName(tableName string) string {
  return fmt.Sprintf("kafka_outbox_%s", tableName)
}

func buildOutboxFuncName(tableName string) string {
  return fmt.Sprintf("%s_outbox_func", tableName)
}

func buildTemplatePackages(fileName string) []*goPackageDesc {
  var tmplPackages []*goPackageDesc

  for _, packageName := range packagesByFiles[fileName] {
    tmplPackage, ok := packagesByNames[packageName]
    if !ok {
      continue
    }
    tmplPackages = append(tmplPackages, tmplPackage)
  }
  return tmplPackages
}

func buildCheckStub(string) bool { return true }

func buildCheckMigration(filePath string) bool {
  return !goose.IsExistedMigration(filePath)
}
