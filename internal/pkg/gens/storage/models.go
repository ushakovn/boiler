package storage

import (
  "fmt"
  "regexp"
  "strings"

  "github.com/ushakovn/boiler/internal/pkg/sql"
  "github.com/ushakovn/boiler/pkg/utils"
  "github.com/ushakovn/boiler/templates"
)

type schemaDesc struct {
  SchemaName       string
  Models           []*modelDesc
  BuildersPackages []*goPackageDesc
  ClientPackages   []*goPackageDesc
  ConstsPackages   []*goPackageDesc
  OptionsPackages  []*goPackageDesc
}

type modelDesc struct {
  ModelName         string
  SqlTableName      string
  ModelFields       []*fieldDesc
  ModelPackages     []*goPackageDesc
  InterfacePackages []*goPackageDesc
}

type fieldDesc struct {
  SqlTableFieldName  string
  FieldName          string
  FieldType          string
  FieldZeroType      string
  FieldBuiltinType   string
  FieldTypeSuffix    string
  NotNullField       bool
  FieldBadge         string
  ModelFieldFilters  []*fieldFilterDesc
  ModelsFieldFilters []*fieldFilterDesc
}

type fieldFilterDesc struct {
  FilterName       string
  FilterType       string
  FilterIfStmt     string
  FilterSqOperator string
  FilterTypeSuffix string
}

func (g *storage) loadSchemaDesc() error {
  tables := g.dumpSQL.Tables.Elems()
  models := make([]*modelDesc, 0, len(tables))

  for _, table := range tables {
    columns := table.Columns.Elems()
    fields := make([]*fieldDesc, 0, len(columns))

    for _, column := range columns {
      field, err := tableColumnToFieldDesc(column)
      if err != nil {
        return fmt.Errorf("tableColumnToFieldDesc: err: %w", err)
      }
      fields = append(fields, field)
    }
    modelName := utils.StringToUpperCamelCase(table.Name)

    models = append(models, &modelDesc{
      ModelName:    modelName,
      SqlTableName: table.Name,
      ModelFields:  fields,
    })
  }

  interfacePackages := buildFilePackages(interfaceFileName)
  buildersPackages := buildFilePackages(buildersFileName)
  clientPackages := buildFilePackages(clientFileName)
  constsPackages := buildFilePackages(constsFileName)
  optionsPackages := buildFilePackages(optionsFileName)
  implementationPackages := buildFilePackages(implementationFileName)

  for _, model := range models {
    modelPackages := append([]*goPackageDesc{}, implementationPackages...)
    packagesByCustomNames := map[string]struct{}{}

    for _, field := range model.ModelFields {
      modelPackage, ok := fieldTypToPackageDesc(field.FieldType)
      if !ok {
        continue
      }
      if _, ok = packagesByCustomNames[modelPackage.CustomName]; ok {
        continue
      }
      modelPackages = append(modelPackages, modelPackage)
      packagesByCustomNames[modelPackage.CustomName] = struct{}{}
    }
    model.ModelPackages = modelPackages
    model.InterfacePackages = interfacePackages
  }

  g.schemaDesc = &schemaDesc{
    SchemaName:       g.storageName,
    Models:           models,
    BuildersPackages: buildersPackages,
    ClientPackages:   clientPackages,
    ConstsPackages:   constsPackages,
    OptionsPackages:  optionsPackages,
  }
  return nil
}

func tableColumnToFieldDesc(column *sql.DumpColumn) (*fieldDesc, error) {
  sqlTableFieldName := column.Name
  fieldName := utils.StringToUpperCamelCase(column.Name)

  fieldZeroTyp, ok := columnNullableTypToFieldTyp(column.Typ)
  if !ok {
    return nil, fmt.Errorf("field zero type not found for column: %s type: %s", column.Name, column.Typ)
  }

  fieldBuiltinTyp, ok := columnNotNullToFieldTypMapping(column.Typ)
  if !ok {
    return nil, fmt.Errorf("field builtin type not found for column: %s type: %s", column.Name, column.Typ)
  }

  var (
    fieldTyp     string
    notNullField bool
  )
  if column.IsNotNull {
    fieldTyp = fieldBuiltinTyp
    notNullField = true
  } else {
    fieldTyp = fieldZeroTyp
    notNullField = false
  }

  var fieldBadge string

  if column.IsPrimaryKey {
    fieldBadge = fieldBadgePk
  }

  fieldTypSuffix := fieldTypeSuffix(fieldTyp)
  modelFilters, modelsFilters := buildFieldFilters(fieldName, fieldTyp, fieldZeroTyp, fieldBuiltinTyp)

  return &fieldDesc{
    SqlTableFieldName:  sqlTableFieldName,
    FieldName:          fieldName,
    FieldType:          fieldTyp,
    FieldZeroType:      fieldZeroTyp,
    FieldBuiltinType:   fieldBuiltinTyp,
    FieldTypeSuffix:    fieldTypSuffix,
    NotNullField:       notNullField,
    FieldBadge:         fieldBadge,
    ModelFieldFilters:  modelFilters,
    ModelsFieldFilters: modelsFilters,
  }, nil
}

func buildFieldFilters(fieldName, fieldTyp, fieldZeroTyp, fieldBuiltinTyp string) (modelFilters []*fieldFilterDesc, modelsFilters []*fieldFilterDesc) {
  if matchNumericTyp(fieldTyp) || matchZeroNumericTyp(fieldTyp) || matchTimeTyp(fieldTyp) {
    modelFilters = buildNumericFilters(modelNumericFilterOperators, fieldName, fieldZeroTyp, fieldBuiltinTyp)
    modelsFilters = buildNumericFilters(modelsNumericFilterOperators, fieldName, fieldZeroTyp, fieldBuiltinTyp)
  }
  if matchStringTyp(fieldTyp) {
    modelFilters = buildStringFilters(modelStringFilterOperators, fieldName, fieldZeroTyp, fieldBuiltinTyp)
    modelsFilters = buildStringFilters(modelsStringFilterOperators, fieldName, fieldZeroTyp, fieldBuiltinTyp)
  }
  if matchBoolTyp(fieldTyp) {
    modelFilters = buildBoolFilters(fieldName, fieldZeroTyp)
    modelsFilters = buildBoolFilters(fieldName, fieldZeroTyp)
  }
  return
}

func buildStringFilters(stringFilterOperators []string, fieldName, fieldZeroTyp, fieldBuiltinTyp string) []*fieldFilterDesc {
  stringFilters := make([]*fieldFilterDesc, 0, len(stringFilterOperators))

  for _, filterOperator := range stringFilterOperators {
    filterName := buildStringFilterName(fieldName, filterOperator)
    filterTyp := buildStringFilterType(fieldZeroTyp, fieldBuiltinTyp, filterOperator)
    filterTypSuffix := buildStringFilterTypeSuffix(filterTyp)
    filterIfStmt := buildStringFilterIfStmt(filterName, filterOperator)
    filterSqOperator := buildStringFilterSqOperator(filterOperator)

    stringFilters = append(stringFilters, &fieldFilterDesc{
      FilterName:       filterName,
      FilterType:       filterTyp,
      FilterIfStmt:     filterIfStmt,
      FilterSqOperator: filterSqOperator,
      FilterTypeSuffix: filterTypSuffix,
    })
  }
  return stringFilters
}

func buildBoolFilters(fieldName, fieldZeroTyp string) []*fieldFilterDesc {
  filterName := buildBoolFilterName(fieldName)
  filterTyp := buildBoolFilterType(fieldZeroTyp)
  filterTypSuffix := buildBoolFilterTypeSuffix(fieldZeroTyp)
  filterIfStmt := buildBoolFilterIfStmt(filterName)
  filterSqOperator := buildBoolFilterSqOperator()

  return []*fieldFilterDesc{
    {
      FilterName:       filterName,
      FilterType:       filterTyp,
      FilterIfStmt:     filterIfStmt,
      FilterSqOperator: filterSqOperator,
      FilterTypeSuffix: filterTypSuffix,
    },
  }
}

func buildNumericFilters(numericFilterOperators []string, fieldName, fieldZeroTyp, fieldBuiltinTyp string) []*fieldFilterDesc {
  numericFilters := make([]*fieldFilterDesc, 0, len(numericFilterOperators))

  for _, filterOperator := range numericFilterOperators {
    filterName := buildNumericFilterName(fieldName, filterOperator)
    filterTyp := buildNumericFilterType(fieldZeroTyp, fieldBuiltinTyp, filterOperator)
    filterTypSuffix := buildNumericFilterTypeSuffix(filterTyp)
    filterIfStmt := buildNumericFilterIfStmt(filterName, filterOperator)
    filterSqOperator := buildNumericFilterSqOperator(filterOperator)

    numericFilters = append(numericFilters, &fieldFilterDesc{
      FilterName:       filterName,
      FilterType:       filterTyp,
      FilterIfStmt:     filterIfStmt,
      FilterSqOperator: filterSqOperator,
      FilterTypeSuffix: filterTypSuffix,
    })
  }
  return numericFilters
}

func columnNotNullToFieldTypMapping(columnTyp string) (string, bool) {
  fieldTyp, ok := map[string]string{
    "integer": "int",

    "smallint": "int16",
    "int":      "int",
    "bigint":   "int64",

    "smallserial": "int16",
    "serial":      "int",
    "bigserial":   "int64",

    "bit":     "bool",
    "bool":    "bool",
    "boolean": "bool",

    "real":    "float32",
    "money":   "float64",
    "numeric": "float64",

    "bytea": "[]byte",
    "json":  "[]byte",
    "jsonb": "[]byte",

    "varchar":   "string",
    "varying":   "string",
    "character": "string",

    "text": "string",

    "time":      "time.Time",
    "timestamp": "time.Time",
  }[columnTyp]

  return fieldTyp, ok
}

func columnNullableTypToFieldTyp(columnTyp string) (string, bool) {
  var fieldTyp string
  ok := true

  switch columnTyp {
  case
    "integer",
    "smallint",
    "int",
    "bigint",
    "smallserial",
    "serial",
    "bigserial":
    fieldTyp = "zero.Int"

  case
    "real",
    "money",
    "numeric":
    fieldTyp = "zero.Float"

  case
    "bytea",
    "json",
    "jsonb":
    fieldTyp = "[]byte"

  case
    "varchar",
    "varying",
    "character",
    "text":
    fieldTyp = "zero.String"

  case
    "time",
    "timestamp":
    fieldTyp = "zero.Time"

  default:
    ok = false
  }
  return fieldTyp, ok
}

func fieldTypeSuffix(fieldTyp string) string {
  var fieldTypSuffix string

  if matchZeroPackageTyp(fieldTyp) {
    fieldTypSuffix = strings.TrimPrefix(fieldTyp, zeroTypPackagePrefix)
  }
  return fieldTypSuffix
}

func matchBoolTyp(fieldTyp string) bool {
  return fieldTyp == "bool"
}

func matchStringTyp(fieldTyp string) bool {
  return fieldTyp == "string" || fieldTyp == "zero.String"
}

func matchNumericTyp(fieldTyp string) bool {
  return regexNumericTyp.MatchString(fieldTyp)
}

func matchZeroNumericTyp(fieldTyp string) bool {
  return regexZeroNumericTyp.MatchString(fieldTyp)
}

func matchTimeTyp(fieldTyp string) bool {
  return fieldTyp == "time.Time" || fieldTyp == "zero.Time"
}

func matchZeroPackageTyp(fieldTyp string) bool {
  return regexZeroPackageTyp.MatchString(fieldTyp)
}

func matchTimePackageTyp(fieldTyp string) bool {
  return fieldTyp == "time.Time"
}

var (
  regexNumericTyp     = regexp.MustCompile(`^(int(16|32|64)?)|float(32|64)$`)
  regexZeroPackageTyp = regexp.MustCompile(`^zero\.[A-Z][a-z]+$`)
  regexZeroNumericTyp = regexp.MustCompile(`^zero\.(Int|Float)$`)
)

func buildNumericFilterName(fieldName, filterOperator string) string {
  return fmt.Sprint(fieldName, filterOperator)
}

func buildNumericFilterType(fieldZeroTyp, fieldBuiltinTyp, filterOperator string) string {
  if filterOperator == numericFilterOperatorIn || filterOperator == numericFilterOperatorNotIn {
    return fmt.Sprint(filterSlicePrefix, fieldBuiltinTyp)
  }
  return fieldZeroTyp
}

func buildNumericFilterTypeSuffix(filterTyp string) string {
  var filterTypSuffix string

  if matchZeroPackageTyp(filterTyp) {
    filterTypSuffix = strings.TrimPrefix(filterTyp, zeroTypPackagePrefix)
  }
  return filterTypSuffix
}

func buildNumericFilterIfStmt(filterName, filterOperator string) string {
  var filterIfStmt string

  switch filterOperator {
  case
    numericFilterOperatorGt,
    numericFilterOperatorGte,
    numericFilterOperatorLt,
    numericFilterOperatorLte,
    numericFilterOperatorEq:
    filterIfStmt = fmt.Sprintf(templates.FilterIfStmtWithPtr, filterName)
  case
    numericFilterOperatorIn,
    numericFilterOperatorNotIn:
    filterIfStmt = fmt.Sprintf(templates.FilterIfStmtWithLen, filterName)
  }
  return filterIfStmt
}

func buildNumericFilterSqOperator(filterOperator string) string {
  return map[string]string{
    numericFilterOperatorGt:    "Gt",
    numericFilterOperatorGte:   "GtOrEq",
    numericFilterOperatorLt:    "Lt",
    numericFilterOperatorLte:   "LtOrEq",
    numericFilterOperatorEq:    "Eq",
    numericFilterOperatorIn:    "In",
    numericFilterOperatorNotIn: "NotEq",
  }[filterOperator]
}

func buildBoolFilterName(fieldName string) string {
  return fmt.Sprint(boolFilterOperatorWhere, fieldName)
}

func buildBoolFilterType(fieldZeroTyp string) string {
  return fieldZeroTyp
}

func buildBoolFilterTypeSuffix(fieldZeroTyp string) string {
  return strings.TrimPrefix(fieldZeroTyp, zeroTypPackagePrefix)
}

func buildBoolFilterIfStmt(filterName string) string {
  return fmt.Sprintf(templates.FilterIfStmtWithPtr, filterName)
}

func buildBoolFilterSqOperator() string {
  return "Eq"
}

func buildStringFilterName(fieldName, filterOperator string) string {
  return fmt.Sprint(fieldName, filterOperator)
}

func buildStringFilterType(fieldZeroTyp, fieldBuiltinTyp, filterOperator string) string {
  if filterOperator == stringFilterOperatorIn || filterOperator == stringFilterOperatorNotIn {
    return fmt.Sprint(filterSlicePrefix, fieldBuiltinTyp)
  }
  return fieldZeroTyp
}

func buildStringFilterTypeSuffix(filterTyp string) string {
  var filterTypSuffix string

  if matchZeroPackageTyp(filterTyp) {
    filterTypSuffix = strings.TrimPrefix(filterTyp, zeroTypPackagePrefix)
  }
  return filterTypSuffix
}

func buildStringFilterIfStmt(filterName, filterOperator string) string {
  var filterIfStmt string

  switch filterOperator {
  case
    stringFilterOperatorEq:
    filterIfStmt = fmt.Sprintf(templates.FilterIfStmtWithPtr, filterName)
  case
    stringFilterOperatorIn,
    stringFilterOperatorNotIn:
    filterIfStmt = fmt.Sprintf(templates.FilterIfStmtWithLen, filterName)
  }
  return filterIfStmt
}

func buildStringFilterSqOperator(filterOperator string) string {
  return map[string]string{
    stringFilterOperatorEq:    "Eq",
    stringFilterOperatorIn:    "In",
    stringFilterOperatorNotIn: "NotEq",
  }[filterOperator]
}

const (
  fieldBadgePk         = "pk"
  filterSlicePrefix    = "[]"
  zeroTypPackagePrefix = "zero."
)

const (
  boolFilterOperatorWhere = "Where"
)

const (
  stringFilterOperatorIn    = "In"
  stringFilterOperatorNotIn = "NotIn"
  stringFilterOperatorEq    = "Eq"
)

var modelStringFilterOperators = []string{
  stringFilterOperatorEq,
}

var modelsStringFilterOperators = []string{
  stringFilterOperatorIn,
  stringFilterOperatorNotIn,
}

const (
  numericFilterOperatorGt    = "Gt"
  numericFilterOperatorGte   = "Gte"
  numericFilterOperatorLt    = "Lt"
  numericFilterOperatorLte   = "Lte"
  numericFilterOperatorIn    = "In"
  numericFilterOperatorNotIn = "NotIn"
  numericFilterOperatorEq    = "Eq"
)

var modelNumericFilterOperators = []string{
  numericFilterOperatorEq,
}

var modelsNumericFilterOperators = []string{
  numericFilterOperatorGt,
  numericFilterOperatorGte,
  numericFilterOperatorLt,
  numericFilterOperatorLte,
  numericFilterOperatorIn,
  numericFilterOperatorNotIn,
}

type goPackageDesc struct {
  CustomName  string
  ImportLine  string
  ImportAlias string
  IsBuiltin   bool
  IsInstall   bool
}

func buildFilePackages(fileName string) []*goPackageDesc {
  packagesNames := importPackagesByFiles[fileName]
  packagesDesc := make([]*goPackageDesc, 0, len(packagesNames))

  for _, packageName := range packagesNames {
    interfacePackage := importPackagesByNames[packageName]
    packagesDesc = append(packagesDesc, interfacePackage)
  }
  return packagesDesc
}

func fieldTypToPackageDesc(fieldTyp string) (*goPackageDesc, bool) {
  var p *goPackageDesc
  ok := true

  switch {
  case matchZeroPackageTyp(fieldTyp):
    p = importPackagesByNames[zeroPackageName]
  case matchTimePackageTyp(fieldTyp):
    p = importPackagesByNames[timePackageName]
  default:
    ok = false
  }
  return p, ok
}

const (
  constsFileName         = "consts"
  clientFileName         = "client"
  buildersFileName       = "builders"
  optionsFileName        = "options"
  modelsFileName         = "models"
  interfaceFileName      = "interface"
  implementationFileName = "implementation"
)

var importPackagesByFiles = map[string][]string{
  clientFileName: {
    contextPackageName,
    databaseSqlPackageName,
    errorsPackageName,
    fmtPackageName,
    sqlscanPackageName,
  },
  buildersFileName: {
    squirrelPackageName,
  },
  optionsFileName: {
    fmtPackageName,
  },
  interfaceFileName: {
    contextPackageName,
    errorsPackageName,
    fmtPackageName,
  },
  implementationFileName: {
    contextPackageName,
    fmtPackageName,
    squirrelPackageName,
  },
  modelsFileName: {},
  constsFileName: {},
}

const (
  contextPackageName     = "context"
  fmtPackageName         = "fmt"
  zeroPackageName        = "zero"
  timePackageName        = "time"
  squirrelPackageName    = "squirrel"
  sqlscanPackageName     = "sqlscan"
  databaseSqlPackageName = "sql"
  errorsPackageName      = "errors"
)

var importPackagesByNames = map[string]*goPackageDesc{
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
  zeroPackageName: {
    CustomName: "guregu/zero",
    ImportLine: "gopkg.in/guregu/null.v4/zero",
    IsInstall:  true,
  },
  timePackageName: {
    CustomName: "go/time",
    ImportLine: "time",
    IsBuiltin:  true,
  },
  squirrelPackageName: {
    CustomName:  "masterminds/squirrel",
    ImportLine:  "github.com/Masterminds/squirrel",
    ImportAlias: "sq",
    IsInstall:   true,
  },
  sqlscanPackageName: {
    CustomName:  "scany/sqlscan",
    ImportLine:  "github.com/georgysavva/scany/v2/sqlscan",
    ImportAlias: "sc",
    IsInstall:   true,
  },
  databaseSqlPackageName: {
    CustomName: "database/sql",
    ImportLine: "database/sql",
    IsBuiltin:  true,
  },
  errorsPackageName: {
    CustomName: "go/errors",
    ImportLine: "errors",
    IsBuiltin:  true,
  },
}
