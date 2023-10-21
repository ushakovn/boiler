package sql

import (
  "errors"
  "fmt"
  "regexp"
  "strings"

  "github.com/ushakovn/boiler/pkg/utils"
)

type state interface {
  next(token string) (state, error)
}

func newTerminateState() state {
  return &terminate{dump: &DumpSQL{
    Tables:    utils.NewStack[*DumpTable](),
    tempStack: utils.NewStack[string](),
  }}
}

func doTransitions(tokens []string) (state, error) {
  var (
    state = newTerminateState()
    err   error
  )
  var collected []string

  for _, token := range tokens {
    token = utils.NormalizeToken(token)

    if state, err = state.next(token); err != nil {
      return nil, fmt.Errorf("state.next: %w", err)
    }

    collected = append(collected, token) //TODO: remove this
  }
  return state, nil
}

type terminate struct {
  dump *DumpSQL
}

func (t *terminate) next(token string) (state, error) {
  switch token {
  case "create":
    return &create{dump: t.dump}, nil
  case "alter":
    return &alter{dump: t.dump}, nil
  default:
    return t, nil
  }
}

type alter struct {
  dump *DumpSQL
}

func (t *alter) next(token string) (state, error) {
  switch token {
  case "table":
    return &table{dump: t.dump}, nil
  default:
    return &terminate{dump: t.dump}, nil
  }
}

type create struct {
  dump *DumpSQL
}

func (t *create) next(token string) (state, error) {
  switch token {
  case "table":
    return &table{dump: t.dump}, nil
  default:
    return &terminate{dump: t.dump}, nil
  }
}

type table struct {
  dump *DumpSQL
}

func (t *table) next(token string) (state, error) {
  switch {
  case token == "only":
    return &only{dump: t.dump}, nil

  case matchTableName(token):
    const (
      withSchema = 2
      schemaPart = 0
      tablePart  = 1
    )
    var (
      schema string
      name   string
    )
    parts := strings.Split(token, ".")

    if len(parts) == withSchema {
      schema = parts[schemaPart]
      name = parts[tablePart]
    }

    t.dump.Tables.Push(&DumpTable{
      RawName: token,
      Name:    name,
      Schema:  schema,
    })

    return &tableName{dump: t.dump}, nil

  default:
    return nil, fmt.Errorf("%w: %s", errUnexpectedToken, token)
  }
}

type only struct {
  dump *DumpSQL
}

func (t *only) next(token string) (state, error) {
  switch {
  case matchTableName(token):

    t.dump.tempStack.Push(token)
    return &tableName{dump: t.dump}, nil

  default:
    return nil, fmt.Errorf("%w: %s", errUnexpectedToken, token)
  }
}

type tableName struct {
  dump *DumpSQL
}

func (t *tableName) next(token string) (state, error) {
  switch token {
  case "(":
    t.dump.tempStack.Pop()
    return &openBracket{dump: t.dump}, nil

  case "add":
    return &add{dump: t.dump}, nil

  case
    "alter":
    t.dump.tempStack.Pop()
    return &terminate{dump: t.dump}, nil

  case
    "owner":
    t.dump.tempStack.Pop()
    t.dump.Tables.Pop()
    return &terminate{dump: t.dump}, nil

  default:
    return nil, fmt.Errorf("%w: %s", errUnexpectedToken, token)
  }
}

type add struct {
  dump *DumpSQL
}

func (t *add) next(token string) (state, error) {
  switch token {
  case "constraint":
    return &constraint{dump: t.dump}, nil
  default:
    return &terminate{dump: t.dump}, nil
  }
}

type constraint struct {
  dump *DumpSQL
}

func (t *constraint) next(token string) (state, error) {
  switch {
  case matchPrimaryKeyConstraint(token):
    return &primaryKeyConstraintName{dump: t.dump}, nil
  default:
    return &terminate{dump: t.dump}, nil
  }
}

type primaryKeyConstraintName struct {
  dump *DumpSQL
}

func (t *primaryKeyConstraintName) next(token string) (state, error) {
  switch token {
  case "primary":
    return &primary{dump: t.dump}, nil
  default:
    return &terminate{dump: t.dump}, nil
  }
}

type primary struct {
  dump *DumpSQL
}

func (t *primary) next(token string) (state, error) {
  switch token {
  case "key":
    return &key{dump: t.dump}, nil
  default:
    return nil, fmt.Errorf("%w: %s", errUnexpectedToken, token)
  }
}

type key struct {
  dump *DumpSQL
}

func (t *key) next(token string) (state, error) {
  switch {
  case matchPrimaryKeyName(token):
    columnName := strings.Trim(token, "()")

    rawTableName, ok := t.dump.tempStack.Pop()
    if !ok {
      return nil, fmt.Errorf("table name not found: primary key: %s", token)
    }
    ok = false

    for tableIdx, table := range t.dump.Tables.Elems() {
      if ok {
        break
      }
      if rawTableName != table.RawName {
        continue
      }
      for columnIdx, column := range table.Columns.Elems() {
        if columnName != column.Name {
          continue
        }
        t.dump.Tables.ElemWith(tableIdx, func(table *DumpTable) {
          table.Columns.ElemWith(columnIdx, func(column *DumpColumn) {
            column.IsPrimaryKey = true
          })
        })
        ok = true
        break
      }
    }
    if !ok {
      return nil, fmt.Errorf("invalid primary key: table name: %s column: %s", rawTableName, columnName)
    }
    return &primaryKeyName{dump: t.dump}, nil

  default:
    return nil, fmt.Errorf("%w: %s", errUnexpectedToken, token)
  }
}

type primaryKeyName struct {
  dump *DumpSQL
}

func (t *primaryKeyName) next(token string) (state, error) {
  switch token {
  case ";":
    return &terminate{dump: t.dump}, nil
  default:
    return nil, fmt.Errorf("%w: %s", errUnexpectedToken, token)
  }
}

type openBracket struct {
  dump *DumpSQL
}

func (t *openBracket) next(token string) (state, error) {
  switch {
  case matchColumnName(token):
    t.dump.Tables.PeekWith(func(table *DumpTable) {
      table.Columns.Push(&DumpColumn{
        Name: token,
      })
    })
    return &columnName{dump: t.dump}, nil

  case token == ")":
    return &closeBracket{dump: t.dump}, nil

  default:
    return nil, fmt.Errorf("%w: %s", errUnexpectedToken, token)
  }
}

type columnName struct {
  dump *DumpSQL
}

func (t *columnName) next(token string) (state, error) {
  var err error

  defer func() {
    if err == nil {
      t.dump.Tables.PeekWith(func(table *DumpTable) {
        table.Columns.PeekWith(func(column *DumpColumn) {
          column.Typ = token
        })
      })
    }
  }()

  switch {
  case utils.StringOneOfEqual(token,
    "integer",

    "smallint",
    "int",
    "bigint",

    "smallserial",
    "serial",
    "bigserial",

    "bit",
    "bool",
    "boolean",

    "money",
    "real",
    "float",
    "double",
    "decimal",
    "numeric",

    "bytea",
    "json",
    "jsonb",

    "text",
    "uuid",
  ) ||
    matchNVarcharColumnTyp(token) ||
    matchCharacterBracketsColumnTyp(token):

    return &scalarColumnTyp{dump: t.dump}, nil

  case token == "character":
    return &characterColumnTyp{dump: t.dump}, nil

  case utils.StringOneOfEqual(token,
    "timestamp",
    "time",
  ):
    return &timeOrTimestampColumnTyp{dump: t.dump}, nil

  default:
    err = fmt.Errorf("%w: %s", errUnexpectedToken, token)
    return nil, err
  }
}

type timeOrTimestampColumnTyp struct {
  dump *DumpSQL
}

func (t *timeOrTimestampColumnTyp) next(token string) (state, error) {
  switch token {
  case
    "with",
    "without":
    t.dump.Tables.PeekWith(func(table *DumpTable) {
      table.Columns.PeekWith(func(column *DumpColumn) {
        column.TypOptions = token
      })
    })
    return &timeOrTimestampWithOrWithoutOption{dump: t.dump}, nil

  case ",":
    return &openBracket{dump: t.dump}, nil

  case ")":
    return &closeBracket{dump: t.dump}, nil

  case "not":
    return &notColumnTypOption{dump: t.dump}, nil

  default:
    return nil, fmt.Errorf("%w: %s", errUnexpectedToken, token)
  }
}

type timeOrTimestampWithOrWithoutOption struct {
  dump *DumpSQL
}

func (t *timeOrTimestampWithOrWithoutOption) next(token string) (state, error) {
  var err error

  defer func() {
    if err == nil {
      t.dump.Tables.PeekWith(func(table *DumpTable) {
        table.Columns.PeekWith(func(column *DumpColumn) {
          column.TypOptions += " " + token
        })
      })
    }
  }()

  switch token {
  case "time":
    return &timeOrTimestampTimeOption{dump: t.dump}, nil

  case "timezone":
    return &timeOrTimestampTimezoneOption{dump: t.dump}, nil

  default:
    return nil, fmt.Errorf("%w: %s", errUnexpectedToken, token)
  }
}

type timeOrTimestampTimeOption struct {
  dump *DumpSQL
}

func (t *timeOrTimestampTimeOption) next(token string) (state, error) {
  switch token {
  case "zone":
    t.dump.Tables.PeekWith(func(table *DumpTable) {
      table.Columns.PeekWith(func(column *DumpColumn) {
        column.TypOptions += " " + token
      })
    })

    return &timeOrTimestampTimezoneOption{dump: t.dump}, nil

  default:
    return nil, fmt.Errorf("%w: %s", errUnexpectedToken, token)
  }
}

type timeOrTimestampTimezoneOption struct {
  dump *DumpSQL
}

func (t *timeOrTimestampTimezoneOption) next(token string) (state, error) {
  switch token {
  case ",":
    return &openBracket{dump: t.dump}, nil

  case ")":
    return &closeBracket{dump: t.dump}, nil

  case "not":
    return &notColumnTypOption{dump: t.dump}, nil

  default:
    return nil, fmt.Errorf("%w: %s", errUnexpectedToken, token)
  }
}

type characterColumnTyp struct {
  dump *DumpSQL
}

func (t *characterColumnTyp) next(token string) (state, error) {
  switch {
  case matchCharacterVaryingOption(token):
    t.dump.Tables.PeekWith(func(table *DumpTable) {
      table.Columns.PeekWith(func(column *DumpColumn) {
        column.TypOptions = token
      })
    })

    return &characterVaryingOption{dump: t.dump}, nil

  case token == ",":
    return &openBracket{dump: t.dump}, nil

  case token == ")":
    return &closeBracket{dump: t.dump}, nil

  case token == "not":
    return &notColumnTypOption{dump: t.dump}, nil

  default:
    return nil, fmt.Errorf("%w: %s", errUnexpectedToken, token)
  }
}

type characterVaryingOption struct {
  dump *DumpSQL
}

func (t *characterVaryingOption) next(token string) (state, error) {
  switch token {
  case ",":
    return &openBracket{dump: t.dump}, nil

  case ")":
    return &closeBracket{dump: t.dump}, nil

  case "not":
    return &notColumnTypOption{dump: t.dump}, nil

  default:
    return nil, fmt.Errorf("%w: %s", errUnexpectedToken, token)
  }
}

type notColumnTypOption struct {
  dump *DumpSQL
}

func (t *notColumnTypOption) next(token string) (state, error) {
  switch token {
  case "null":
    t.dump.Tables.PeekWith(func(table *DumpTable) {
      table.Columns.PeekWith(func(column *DumpColumn) {
        column.IsNotNull = true
      })
    })
    return &nullColumnTypOption{dump: t.dump}, nil

  default:
    return nil, fmt.Errorf("%w: %s", errUnexpectedToken, token)
  }
}

type nullColumnTypOption struct {
  dump *DumpSQL
}

func (t *nullColumnTypOption) next(token string) (state, error) {
  switch token {
  case ",":
    return &openBracket{dump: t.dump}, nil

  case ")":
    return &closeBracket{dump: t.dump}, nil

  default:
    return nil, fmt.Errorf("%w: %s", errUnexpectedToken, token)
  }
}

type scalarColumnTyp struct {
  dump *DumpSQL
}

func (t *scalarColumnTyp) next(token string) (state, error) {
  switch token {
  case ",":
    return &openBracket{dump: t.dump}, nil

  case ")":
    return &closeBracket{dump: t.dump}, nil

  case "not":
    return &notColumnTypOption{dump: t.dump}, nil

  default:
    return nil, fmt.Errorf("%w: %s", errUnexpectedToken, token)
  }
}

type closeBracket struct {
  dump *DumpSQL
}

func (t *closeBracket) next(token string) (state, error) {
  switch token {
  case ";":
    return &terminate{dump: t.dump}, nil

  default:
    return nil, fmt.Errorf("%w: %s", errUnexpectedToken, token)
  }
}

var (
  matchTableName                  = regexp.MustCompile(`^(\w+\.)?\w+$`).MatchString
  matchColumnName                 = regexp.MustCompile(`^\w+$`).MatchString
  matchCharacterBracketsColumnTyp = regexp.MustCompile(`^char(acter)?(\(\d+\))$`).MatchString
  matchNVarcharColumnTyp          = regexp.MustCompile(`^(n)+varchar(\(\d+\))?$`).MatchString
  matchCharacterVaryingOption     = regexp.MustCompile(`^varying(\(\d+\)?)*$`).MatchString
  matchPrimaryKeyConstraint       = regexp.MustCompile(`^\w+_pkey$`).MatchString
  matchPrimaryKeyName             = regexp.MustCompile(`^(\w+)|\(\w+\)$`).MatchString
)

var errUnexpectedToken = errors.New("unexpected token")
