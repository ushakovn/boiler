package kafkaoutbox

import (
  "fmt"
  "os"

  "github.com/samber/lo"
  "github.com/ushakovn/boiler/internal/pkg/pgdump"
  "github.com/ushakovn/boiler/internal/pkg/stringer"
  "github.com/yoheimuta/go-protoparser"
  "github.com/yoheimuta/go-protoparser/parser"
)

type parsedProtoMessage struct {
  tableName     string
  messageName   string
  messageFields []*parsedProtoField
}

type parsedProtoField struct {
  fieldName  string
  fieldType  string
  isOptional bool
  isRepeated bool
}

type parsedProto struct {
  messages []*parsedProtoMessage
}

func parseProto(filePath, optionName string) (*parsedProto, error) {
  file, err := os.Open(filePath)
  if err != nil {
    return nil, fmt.Errorf("file open failed: %w", err)
  }
  parsed, err := protoparser.Parse(file)
  if err != nil {
    return nil, fmt.Errorf("protoparser.Parse: %w", err)
  }
  var (
    messages []*parsedProtoMessage
    token    parser.Visitee
  )
  for _, token = range parsed.ProtoBody {
    message, ok := token.(*parser.Message)
    if !ok {
      continue
    }
    enums := make(map[string]struct{})

    for _, token = range message.MessageBody {
      enum, ok := token.(*parser.Enum)
      if !ok {
        continue
      }
      enums[enum.EnumName] = struct{}{}
    }

    m := &parsedProtoMessage{
      messageName: message.MessageName,
    }
    for _, token = range message.MessageBody {
      switch t := token.(type) {

      case *parser.Option:
        if t.OptionName != optionName {
          continue
        }
        m.tableName = stringer.UnquoteString(t.Constant)

      case *parser.Field:
        fieldType := unwrapProtoFieldTyp(t.Type, enums)

        m.messageFields = append(m.messageFields, &parsedProtoField{
          fieldName:  t.FieldName,
          fieldType:  fieldType,
          isOptional: t.IsOptional,
          isRepeated: t.IsRepeated,
        })

      default:
        continue
      }
    }
    messages = append(messages, m)
  }

  return &parsedProto{
    messages: messages,
  }, nil
}

func validateProto(parsed *parsedProto, dump *pgdump.DumpSQL) error {
  tables := lo.Associate(dump.Tables.Elems(),
    func(table *pgdump.DumpTable) (string, *pgdump.DumpTable) {
      return table.Name, table
    })

  for _, message := range parsed.messages {
    table, ok := tables[message.tableName]
    if !ok {
      return fmt.Errorf("message %s: table: %s: table not found is schema",
        message.messageName, message.tableName)
    }

    columns := lo.Associate(table.Columns.Elems(),
      func(column *pgdump.DumpColumn) (string, *pgdump.DumpColumn) {
        return column.Name, column
      })

    for _, field := range message.messageFields {
      column, ok := columns[field.fieldName]
      if !ok {
        return fmt.Errorf("message: %s: field: %s: column not found in table: %s",
          message.messageName, field.fieldName, message.tableName)
      }

      if column.IsNotNull && field.isOptional {
        return fmt.Errorf("message: %s: field: %s: field must be not optional",
          message.messageName, field.fieldName)
      }

      if !column.IsNotNull && !field.isOptional {
        return fmt.Errorf("message: %s: field: %s: field must be optional",
          message.messageName, field.fieldName)
      }

      fieldTypes, isRepeated := pgTypeToProtoTypes(column.Typ)

      if len(fieldTypes) == 0 {
        return fmt.Errorf("message: %s: field: %s: column type not supported: %s",
          message.messageName, field.fieldName, column.Typ)
      }

      if isRepeated && !field.isRepeated {
        return fmt.Errorf("message: %s: field: %s: field must be repeated",
          message.messageName, field.fieldName)
      }

      if !lo.Contains(fieldTypes, field.fieldType) {
        return fmt.Errorf("message: %s: field: %s: invalid field type: %s",
          message.messageName, field.fieldName, field.fieldType)
      }
    }
  }
  return nil
}

func pgTypeToProtoTypes(pgType string) (protoTypes []string, isRepeated bool) {
  switch pgType {
  case
    "smallserial",
    "smallint",
    "integer",
    "int",
    "serial",
    "bigint",
    "bigserial":
    return []string{"int32", "int64"}, false

  case
    "bit",
    "bool",
    "boolean":
    return []string{"bool"}, false

  case
    "money",
    "real",
    "float",
    "double",
    "decimal",
    "numeric":
    return []string{"float", "double"}, false

  case
    "varchar",
    "varying",
    "character",
    "uuid",
    "text":
    return []string{"string"}, false

  case
    "date",
    "time",
    "timestamp":
    return []string{"google.protobuf.Timestamp"}, false

  case
    "bytea",
    "json",
    "jsonb":
    return []string{"bytes"}, false

  case
    "smallint[]",
    "integer[]",
    "bigint[]",
    "int[]":
    return []string{"int32", "int64"}, true

  case
    "real[]",
    "float[]",
    "double[]",
    "decimal[]",
    "numeric[]":
    return []string{"float", "double"}, true

  case
    "text[]":
    return []string{"string"}, true

  // Extension types section

  case "citext":
    return []string{"string"}, false

  default:
    return nil, false
  }
}

func unwrapProtoFieldTyp(protoType string, protoEnums map[string]struct{}) string {
  if _, ok := protoEnums[protoType]; ok {
    return "int32"
  }
  return protoType
}
