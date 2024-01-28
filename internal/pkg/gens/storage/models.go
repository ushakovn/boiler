package storage

import (
  "fmt"
  "path"
  "regexp"
  "strings"

  "github.com/samber/lo"
  "github.com/ushakovn/boiler/internal/pkg/pgdump"
  "github.com/ushakovn/boiler/internal/pkg/stringer"
  "github.com/ushakovn/boiler/templates"
)

type schemaDesc struct {
  Models           []*modelDesc
  ModelsPackages   []*goPackageDesc
  BuildersPackages []*goPackageDesc
  OptionsPackages  []*goPackageDesc
}

type modelDesc struct {
  ModelName            string
  SqlTableName         string
  ModelFields          []*fieldDesc
  ModelOptionsPackages []*goPackageDesc
  ModelMethodsPackages []*goPackageDesc
}

type fieldDesc struct {
  // Sql column field attributes
  SqlTableFieldName string
  NotNullField      bool
  WithDefaultField  bool
  FieldBadge        string

  // Core field attributes
  FieldName       string
  FieldType       string
  FieldTypeSuffix string
  FieldIfStmt     string

  // Builtin field attributes
  FieldBuiltinType string

  // Zero package field attributes
  FieldZeroType       string
  FieldZeroTypeSuffix string
  FieldZeroTypeIfStmt string

  // Filter field attributes
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

func (g *Storage) loadSchemaDesc() error {
  tables := g.dumpSQL.Tables.Elems()
  models := make([]*modelDesc, 0, len(tables))

  for _, table := range tables {
    columns := table.Columns.Elems()
    fields := make([]*fieldDesc, 0, len(columns))

    for _, column := range columns {
      field, err := g.tableColumnToFieldDesc(column)
      if err != nil {
        return fmt.Errorf("tableColumnToFieldDesc: err: %w", err)
      }
      fields = append(fields, field)
    }
    modelName := buildModelName(table.Name)

    models = append(models, &modelDesc{
      ModelName:    modelName,
      SqlTableName: table.Name,
      ModelFields:  fields,
    })
  }

  buildersPackages := buildFilePackages(buildersFileName)
  optionsPackages := buildFilePackages(optionsFileName)
  modelsPackages := buildFilePackages(modelsFileName)

  modelMethodsPackages := mergeGoPackages(
    buildFilePackages(modelMethodsFileName),
    buildCrossFilePackages(g.goModuleName, modelMethodsFileName),
  )

  modelUniquePackages := map[string]struct{}{}

  for _, model := range models {
    modelOptionsPackages := mergeGoPackages(
      buildFilePackages(modelOptionsFileName),
      buildCrossFilePackages(g.goModuleName, modelOptionsFileName),
    )
    modelOptionsUnique := map[string]struct{}{}

    for _, field := range model.ModelFields {
      if fieldPackages, ok := buildFieldPackages(field.FieldType); ok {
        for _, fieldPackage := range fieldPackages {
          if _, ok = modelUniquePackages[fieldPackage.CustomName]; ok {
            continue
          }
          modelsPackages = append(modelsPackages, fieldPackage)
          modelUniquePackages[fieldPackage.CustomName] = struct{}{}
        }
      }
      for _, filter := range field.ModelsFieldFilters {
        if fieldPackages, ok := buildFilterFieldPackages(filter.FilterType); ok {
          for _, fieldPackage := range fieldPackages {
            if _, ok = modelOptionsUnique[fieldPackage.CustomName]; ok {
              continue
            }
            modelOptionsPackages = append(modelOptionsPackages, fieldPackage)
            modelOptionsUnique[fieldPackage.CustomName] = struct{}{}
          }
        }
      }
    }

    model.ModelOptionsPackages = modelOptionsPackages
    model.ModelMethodsPackages = modelMethodsPackages
  }

  g.schemaDesc = &schemaDesc{
    Models:           models,
    ModelsPackages:   modelsPackages,
    BuildersPackages: buildersPackages,
    OptionsPackages:  optionsPackages,
  }

  return nil
}

func (g *Storage) tableColumnToFieldDesc(column *pgdump.DumpColumn) (*fieldDesc, error) {
  sqlTableFieldName := column.Name
  fieldName := stringer.StringToUpperCamelCase(column.Name)

  fieldZeroTyp, ok := g.columnNullableTypToFieldTyp(column.Typ)
  if !ok {
    return nil, fmt.Errorf("field zero type not found for column: %s type: %s", column.Name, column.Typ)
  }

  fieldBuiltinTyp, ok := g.columnNotNullToFieldTypMapping(column.Typ)
  if !ok {
    return nil, fmt.Errorf("field builtin type not found for column: %s type: %s", column.Name, column.Typ)
  }

  fieldTyp := lo.Ternary(column.IsNotNull, fieldBuiltinTyp, fieldZeroTyp)
  fieldBadge := lo.Ternary(column.IsPrimaryKey, fieldBadgePk, "")

  fieldIfStmt := buildFieldIfStmt(fieldName, fieldTyp)
  fieldTypSuffix := buildFieldTypeSuffix(fieldTyp)

  fieldZeroTypIfStmt := buildFieldIfStmt(fieldName, fieldZeroTyp)
  fieldZeroTypSuffix := buildFieldTypeSuffix(fieldZeroTyp)

  modelFilters, modelsFilters := buildFieldFilters(fieldName, fieldTyp, fieldZeroTyp, fieldBuiltinTyp)

  return &fieldDesc{
    SqlTableFieldName: sqlTableFieldName,
    NotNullField:      column.IsNotNull,
    WithDefaultField:  column.WithDefault,
    FieldBadge:        fieldBadge,

    FieldName:       fieldName,
    FieldType:       fieldTyp,
    FieldIfStmt:     fieldIfStmt,
    FieldTypeSuffix: fieldTypSuffix,

    FieldBuiltinType: fieldBuiltinTyp,

    FieldZeroType:       fieldZeroTyp,
    FieldZeroTypeIfStmt: fieldZeroTypIfStmt,
    FieldZeroTypeSuffix: fieldZeroTypSuffix,

    ModelFieldFilters:  modelFilters,
    ModelsFieldFilters: modelsFilters,
  }, nil
}

func buildModelName(tableName string) string {
  tableName = stringer.StringToUpperCamelCase(tableName)
  modelName := stringer.NormalizeName(tableName)
  return modelName
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
    // Filter type suffix at the same that field type suffix
    filterTypSuffix := buildFieldTypeSuffix(filterTyp)
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
  // Filter type suffix at the same that field type suffix
  filterTypSuffix := buildFieldTypeSuffix(filterTyp)
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
    // Filter type suffix at the same that field type suffix
    filterTypSuffix := buildFieldTypeSuffix(filterTyp)
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

func (g *Storage) columnNotNullToFieldTypMapping(columnTyp string) (string, bool) {
  fieldTyp, ok1 := map[string]string{
    // Scalar column types

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

    "money":   "float64",
    "real":    "float32",
    "float":   "float32",
    "double":  "float64",
    "decimal": "float64",
    "numeric": "float64",

    "varchar":   "string",
    "varying":   "string",
    "character": "string",

    "uuid": "string",
    "text": "string",

    "date":      "time.Time",
    "time":      "time.Time",
    "timestamp": "time.Time",

    // Slice of bytes column types

    "bytea": "[]byte",
    "json":  "[]byte",
    "jsonb": "[]byte",

    // Slice of integers column types

    "integer[]":  "[]int",
    "smallint[]": "[]int",
    "bigint[]":   "[]int",
    "int[]":      "[]int",

    // Slice of floats column types

    "real[]":    "[]float64",
    "float[]":   "[]float64",
    "double[]":  "[]float64",
    "decimal[]": "[]float64",
    "numeric[]": "[]float64",

    // Slice of strings column types

    "text[]": "[]string",
  }[columnTyp]

  // Try override type from pg type config
  if customTyp, ok2 := g.config.PgTypeConfig[columnTyp]; ok2 {
    return customTyp.GoType, true
  }

  return fieldTyp, ok1
}

func (g *Storage) columnNullableTypToFieldTyp(columnTyp string) (string, bool) {
  var fieldTyp string
  ok1 := true

  switch columnTyp {
  // Scalar column types

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
    "money",
    "real",
    "float",
    "double",
    "decimal",
    "numeric":
    fieldTyp = "zero.Float"

  case
    "varchar",
    "varying",
    "character",
    "uuid",
    "text":
    fieldTyp = "zero.String"

  case
    "date",
    "time",
    "timestamp":
    fieldTyp = "zero.Time"

  case
    "bool",
    "boolean":
    fieldTyp = "zero.Bool"

  // Slice of bytes column types

  case
    "bytea",
    "json",
    "jsonb":
    fieldTyp = "[]byte"

  // Slice of integers column types

  case
    "integer[]",
    "smallint[]",
    "bigint[]",
    "int[]":
    fieldTyp = "[]int"

  // Slice of floats column types

  case
    "real[]",
    "float[]",
    "double[]",
    "decimal[]",
    "numeric[]":
    fieldTyp = "[]float64"

  // Slice of strings column types

  case "text[]":
    fieldTyp = "[]string"

  default:
    ok1 = false
  }

  // Try override type from pg type config
  if customTyp, ok2 := g.config.PgTypeConfig[columnTyp]; ok2 {
    return customTyp.GoZeroType, ok2
  }

  return fieldTyp, ok1
}

func buildFieldTypeSuffix(fieldTyp string) string {
  var fieldTypSuffix string

  if matchZeroTyp(fieldTyp) {
    // zero.Int{} and zero.Float{}
    // Contains zero.Int{}.Int64 and zero.Float{}.Float64 fields
    if matchZeroNumericTyp(fieldTyp) {
      fieldTyp = fmt.Sprint(fieldTyp, 64)
    }
    // Trim zero package prefix
    fieldTypSuffix = strings.TrimPrefix(fieldTyp, zeroTypPackagePrefix)
  }
  return fieldTypSuffix
}

func buildFieldIfStmt(fieldName, fieldTyp string) string {
  var fieldIfStmt string

  if matchSliceTyp(fieldTyp) {
    fieldIfStmt = fmt.Sprintf(templates.StorageInputIfStmtWithLen, fieldName)
  }
  if matchZeroTyp(fieldTyp) {
    fieldIfStmt = fmt.Sprintf(templates.StorageInputIfStmtWithPtr, fieldName)
  }
  return fieldIfStmt
}

func matchSliceTyp(fieldTyp string) bool {
  return strings.HasPrefix(fieldTyp, sliceTypPrefix)
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
  return fieldTyp == "zero.Int" || fieldTyp == "zero.Float"
}

func matchTimeTyp(fieldTyp string) bool {
  return fieldTyp == "time.Time" || fieldTyp == "zero.Time"
}

func matchZeroTyp(fieldTyp string) bool {
  return regexZeroPackageTyp.MatchString(fieldTyp)
}

var (
  regexNumericTyp     = regexp.MustCompile(`^(int(16|32|64)?)|float(32|64)$`)
  regexZeroPackageTyp = regexp.MustCompile(`^zero\.[A-Z][a-z]+$`)
)

func buildNumericFilterName(fieldName, filterOperator string) string {
  return fmt.Sprint(fieldName, filterOperator)
}

func buildNumericFilterType(fieldZeroTyp, fieldBuiltinTyp, filterOperator string) string {
  if filterOperator == numericFilterOperatorIn || filterOperator == numericFilterOperatorNotIn {
    return fmt.Sprint(sliceTypPrefix, fieldBuiltinTyp)
  }
  return fieldZeroTyp
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
    filterIfStmt = fmt.Sprintf(templates.StorageFilterIfStmtWithPtr, filterName)
  case
    numericFilterOperatorIn,
    numericFilterOperatorNotIn:
    filterIfStmt = fmt.Sprintf(templates.StorageFilterIfStmtWithLen, filterName)
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
    numericFilterOperatorIn:    "Eq",
    numericFilterOperatorNotIn: "NotEq",
  }[filterOperator]
}

func buildBoolFilterName(fieldName string) string {
  return fmt.Sprint(boolFilterOperatorWhere, fieldName)
}

func buildBoolFilterType(fieldZeroTyp string) string {
  return fieldZeroTyp
}

func buildBoolFilterIfStmt(filterName string) string {
  return fmt.Sprintf(templates.StorageFilterIfStmtWithPtr, filterName)
}

func buildBoolFilterSqOperator() string {
  return "Eq"
}

func buildStringFilterName(fieldName, filterOperator string) string {
  return fmt.Sprint(fieldName, filterOperator)
}

func buildStringFilterType(fieldZeroTyp, fieldBuiltinTyp, filterOperator string) string {
  if filterOperator == stringFilterOperatorIn || filterOperator == stringFilterOperatorNotIn {
    return fmt.Sprint(sliceTypPrefix, fieldBuiltinTyp)
  }
  return fieldZeroTyp
}

func buildStringFilterIfStmt(filterName, filterOperator string) string {
  var filterIfStmt string

  switch filterOperator {
  case
    stringFilterOperatorEq:
    filterIfStmt = fmt.Sprintf(templates.StorageFilterIfStmtWithPtr, filterName)
  case
    stringFilterOperatorIn,
    stringFilterOperatorNotIn:
    filterIfStmt = fmt.Sprintf(templates.StorageFilterIfStmtWithLen, filterName)
  }
  return filterIfStmt
}

func buildStringFilterSqOperator(filterOperator string) string {
  return map[string]string{
    stringFilterOperatorEq:    "Eq",
    stringFilterOperatorIn:    "Eq",
    stringFilterOperatorNotIn: "NotEq",
  }[filterOperator]
}

const (
  fieldBadgePk         = "pk"
  sliceTypPrefix       = "[]"
  zeroTypPackagePrefix = "zero"
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
    desc := importPackagesByNames[packageName]
    packagesDesc = append(packagesDesc, desc)
  }
  return packagesDesc
}

func buildPackagesForNames(packagesNames []string) []*goPackageDesc {
  packagesDesc := make([]*goPackageDesc, 0, len(packagesNames))

  for _, packageName := range packagesNames {
    desc := importPackagesByNames[packageName]
    packagesDesc = append(packagesDesc, desc)
  }
  return packagesDesc
}

func buildFieldPackages(fieldTyp string) ([]*goPackageDesc, bool) {
  var (
    fieldPackages []*goPackageDesc
  )
  switch {
  case matchZeroTyp(fieldTyp):
    fieldPackages = append(fieldPackages, importPackagesByNames[zeroPackageName])
  case matchTimeTyp(fieldTyp):
    fieldPackages = append(fieldPackages, importPackagesByNames[timePackageName])
  }
  return fieldPackages, len(fieldPackages) != 0
}

func buildFilterFieldPackages(filterTyp string) ([]*goPackageDesc, bool) {
  var (
    fieldPackages []*goPackageDesc
  )
  if matchZeroTyp(filterTyp) {
    fieldPackages = append(fieldPackages, importPackagesByNames[zeroPackageName])
  }
  if matchTimeTyp(filterTyp) {
    fieldPackages = append(fieldPackages, importPackagesByNames[timePackageName])
  }
  return fieldPackages, len(fieldPackages) != 0
}

func buildCrossFilePackages(goModuleName, fileName string) []*goPackageDesc {
  crossFileNames := map[string][]string{
    modelMethodsFileName: {modelsFileName},
  }[fileName]

  var (
    nestedFileParts   = []string{"internal", "pkg", "storage"}
    crossFilePackages []*goPackageDesc
  )

  for _, crossFileName := range crossFileNames {
    var packagePathParts []string

    packagePathParts = append(packagePathParts, goModuleName)
    packagePathParts = append(packagePathParts, nestedFileParts...)
    packagePathParts = append(packagePathParts, crossFileName)

    packageImportPath := path.Join(packagePathParts...)
    packageImportAlias := crossFileName

    crossFilePackages = append(crossFilePackages, &goPackageDesc{
      CustomName:  "boiler/cross-package",
      ImportLine:  packageImportPath,
      ImportAlias: packageImportAlias,
      IsBuiltin:   false,
      IsInstall:   false,
    })
  }

  return crossFilePackages
}

func mergeGoPackages(goPackages ...[]*goPackageDesc) []*goPackageDesc {
  var count int

  for _, p := range goPackages {
    count += len(p)
  }
  merged := make([]*goPackageDesc, 0, count)

  for _, p := range goPackages {
    merged = append(merged, p...)
  }
  return merged
}

const (
  constsFileName       = "consts"
  buildersFileName     = "builders"
  optionsFileName      = "options"
  modelsFileName       = "models"
  modelOptionsFileName = "model_options"
  modelMethodsFileName = "model_methods"
)

var importPackagesByFiles = map[string][]string{
  buildersFileName: {
    squirrelPackageName,
  },
  optionsFileName: {
    fmtPackageName,
    contextPackageName,
    pgPgxPackageName,
    pgClientPackageName,
  },
  modelOptionsFileName: {
    fmtPackageName,
  },
  modelMethodsFileName: {
    contextPackageName,
    fmtPackageName,
    squirrelPackageName,
    pgClientPackageName,
    pgBuilderPackageName,
  },
  modelsFileName: {},
  constsFileName: {},
}

const (
  contextPackageName = "context"
  fmtPackageName     = "fmt"
  zeroPackageName    = "zero"
  timePackageName    = "time"
  errorsPackageName  = "errors"

  logrusPackageName      = "logrus"
  databaseSqlPackageName = "sql"
  squirrelPackageName    = "squirrel"

  pgPgxPackageName     = "pg-pgx"
  pgClientPackageName  = "pg-client"
  pgBuilderPackageName = "pg-builder"
  pgErrorsPackageName  = "pg-errors"
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
  logrusPackageName: {
    CustomName:  "sirupsen/logrus",
    ImportLine:  "github.com/sirupsen/logrus",
    ImportAlias: "log",
    IsInstall:   true,
  },
  pgClientPackageName: {
    CustomName:  "boiler/pg-client",
    ImportLine:  "github.com/ushakovn/boiler/pkg/storage/postgres/client",
    ImportAlias: "pg",
    IsInstall:   true,
  },
  pgBuilderPackageName: {
    CustomName:  "boiler/pg-builder",
    ImportLine:  "github.com/ushakovn/boiler/pkg/storage/postgres/builder",
    ImportAlias: "br",
    IsInstall:   true,
  },
  pgErrorsPackageName: {
    CustomName:  "boiler/pg-errors",
    ImportLine:  "github.com/ushakovn/boiler/pkg/storage/postgres/errors",
    ImportAlias: "pe",
    IsInstall:   true,
  },
  pgPgxPackageName: {
    CustomName:  "jackc/pgx",
    ImportLine:  "github.com/jackc/pgx/v5",
    ImportAlias: "pgx",
    IsInstall:   true,
  },
}
