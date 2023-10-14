package storage

import (
  "fmt"
  "regexp"
  "strings"

  "github.com/ushakovn/boiler/internal/pkg/sql"
)

type schemaDesc struct {
  Models []*modelDesc
}

type modelDesc struct {
  ModelName   string
  ModelFields []*fieldDesc
}

type fieldDesc struct {
  SqlTableFieldName  string
  FieldName          string
  FieldType          string
  FiledTypeSuffix    string
  NotNullField       bool
  FieldBadge         string
  ModelFieldFilters  []*fieldFilterDesc
  ModelsFieldFilters []*fieldFilterDesc
}

type fieldFilterDesc struct {
  FilterName       string
  FilterType       string
  FilterTypeSuffix string
}

func (g *storage) loadSchemaDesc() error {
  for range g.dumpSQL.Tables.Elems() {
    return nil
  }
  return nil // TODO
}

//func tableColumnToFieldDesc(fieldName, fieldTyp string) *fieldDesc {
//  var (
//    modelFilters  []*fieldFilterDesc
//    modelsFilters []*fieldFilterDesc
//  )
//  if matchNumericTyp(fieldTyp) {
//    modelFilters = buildNumericFilters(modelNumericFilterOperators, fieldName, fieldTyp)
//    modelsFilters = buildNumericFilters(modelNumericFilterOperators, fieldName, fieldTyp)
//  }
//  if matchBoolTyp(fieldTyp) {
//    modelFilters = buildBoolFilters(fieldName)
//    modelsFilters = buildBoolFilters(fieldName)
//  }
//  return nil // TODO
//}

func buildStringFilter(fieldName string) *fieldFilterDesc {
  return &fieldFilterDesc{
    FilterName:       buildStringFilterName(fieldName),
    FilterType:       buildStringFilterType(),
    FilterTypeSuffix: buildStringFilterTypeSuffix(),
  }
}

func buildBoolFilters(fieldName string) []*fieldFilterDesc {
  return []*fieldFilterDesc{
    {
      FilterName:       buildBoolFilterName(fieldName),
      FilterType:       buildBoolFilterType(fieldName),
      FilterTypeSuffix: buildBoolFilterTypeSuffix(),
    },
  }
}

func buildNumericFilters(numericFilterOperators []string, fieldName, fieldTyp string) []*fieldFilterDesc {
  numericFilters := make([]*fieldFilterDesc, 0, len(numericFilterOperators))

  for _, filterOperator := range numericFilterOperators {
    numericFilter := &fieldFilterDesc{
      FilterName:       buildNumericFilterName(fieldName, filterOperator),
      FilterType:       buildNumericFilterType(fieldTyp, filterOperator),
      FilterTypeSuffix: buildNumericFilterTypeSuffix(fieldTyp),
    }
    numericFilters = append(numericFilters, numericFilter)
  }
  return numericFilters
}

func buildModelFieldFilterForColumn(fieldName, fieldTyp string) []*fieldFilterDesc {
  return nil // TODO
}

func buildModelsFieldFilterForColumn(column *sql.DumpColumn) []*fieldFilterDesc {
  return nil // TODO
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

func columnNullableToFieldTypMapping(columnTyp string) (string, bool) {
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

type goPackage struct {
  name     string
  external bool
}

func fieldTypToPackage(fieldTyp string) (*goPackage, bool) {
  var pack *goPackage
  ok := true

  switch {
  case matchZeroPackageTyp(fieldTyp):
    pack = &goPackage{
      name:     "gopkg.in/guregu/null.v4/zero",
      external: true,
    }
  case fieldTyp == "time.Time":
    pack = &goPackage{
      name:     "time.Time",
      external: false,
    }
  default:
    ok = false
  }
  return pack, ok
}

func matchBoolTyp(fieldTyp string) bool {
  return fieldTyp == "bool"
}

var (
  matchZeroPackageTyp = regexp.MustCompile(`^zero\.[A-Z][a-z]+$`).MatchString
  matchNumericTyp     = regexp.MustCompile(`^(int(16|32|64)?)|float(32|64)$`).MatchString
  matchZeroNumericTyp = regexp.MustCompile(`^zero\.(Int|Float)$`).MatchString
  matchSliceTyp       = regexp.MustCompile(`^\[\]\w+`).MatchString
)

func buildNumericFilterName(fieldName, filterOperator string) string {
  return fmt.Sprint(fieldName, filterOperator)
}

func buildNumericFilterType(fieldTyp, filterOperator string) string {
  if filterOperator == numericFilterOperatorIn {
    return fmt.Sprint(filterSlicePrefix, fieldTyp)
  }
  return fieldTyp
}

func buildNumericFilterTypeSuffix(fieldTyp string) string {
  var filterTypSuffix string

  if matchZeroNumericTyp(fieldTyp) {
    filterTypSuffix = strings.TrimPrefix(fieldTyp, zeroTypPackagePrefix)
  }
  return filterTypSuffix
}

func buildBoolFilterName(fieldName string) string {
  const boolFilterOperator = "Where"
  return fmt.Sprint(boolFilterOperator, fieldName)
}

func buildBoolFilterType(fieldTyp string) string {
  return fieldTyp
}

func buildBoolFilterTypeSuffix() string {
  return ""
}

func buildStringFilterName(fieldName string) string {
  const stringFilterOperator = "In"
  return fmt.Sprint(fieldName, stringFilterOperator)
}

func buildStringFilterType() string {
  const stringFilterTyp = "string"
  return fmt.Sprint(filterSlicePrefix, stringFilterTyp)
}

func buildStringFilterTypeSuffix() string {
  return ""
}

const (
  filterSlicePrefix    = "[]"
  zeroTypPackagePrefix = "zero."
)

const (
  numericFilterOperatorGt  = "Gt"
  numericFilterOperatorGte = "Gte"
  numericFilterOperatorLt  = "Lt"
  numericFilterOperatorLte = "Lte"
  numericFilterOperatorIn  = "In"
  numericFilterOperatorEq  = "Eq"
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
}
