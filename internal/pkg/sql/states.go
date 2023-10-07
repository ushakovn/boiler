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
  return &terminate{dump: &DumpSQL{}}
}

func doTransitions(tokens []string) (state, error) {
  var (
    state = newTerminateState()
    err   error
  )
  collected := make([]string, 0, len(tokens))

  for _, token := range tokens {
    token = utils.NormalizeToken(token)
    if state, err = state.next(token); err != nil {
      return nil, fmt.Errorf("state.next: %w", err)
    }

    collected = append(collected, token)
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

  default:
    return t, nil
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

    t.dump.Tables = append(t.dump.Tables, &DumpTable{
      RawName: token,
      Name:    name,
      Schema:  schema,
    })

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
    return &openBracket{dump: t.dump}, nil

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
    table := t.dump.Tables[len(t.dump.Tables)-1]

    table.Columns = append(table.Columns, &DumpColumn{
      Name: token,
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
      table := t.dump.Tables[len(t.dump.Tables)-1]
      column := table.Columns[len(table.Columns)-1]
      column.Typ = token
    }
  }()

  switch {
  case utils.StringOneOfEqual(token,
    "bigint",
    "bit",
    "boolean",
    "bool",
    "bytea",
    "date",
    "integer",
    "int",
    "money",
    "numeric",
    "real",
    "serial",
    "json",
    "jsonb",
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
  table := t.dump.Tables[len(t.dump.Tables)-1]
  column := table.Columns[len(table.Columns)-1]

  switch token {
  case
    "with",
    "without":
    column.TypOpt = token
    return &timeOrTimestampWithOrWithoutOption{dump: t.dump}, nil

  case ",":
    return &openBracket{dump: t.dump}, nil

  case ")":
    return &closeBracket{dump: t.dump}, nil

  case "not":
    column.ColOpt = token
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
      table := t.dump.Tables[len(t.dump.Tables)-1]
      column := table.Columns[len(table.Columns)-1]
      column.TypOpt += " " + token
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
    table := t.dump.Tables[len(t.dump.Tables)-1]
    column := table.Columns[len(table.Columns)-1]
    column.TypOpt += " " + token

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
    table := t.dump.Tables[len(t.dump.Tables)-1]
    column := table.Columns[len(table.Columns)-1]
    column.ColOpt = token

    return &notColumnTypOption{dump: t.dump}, nil

  default:
    return nil, fmt.Errorf("%w: %s", errUnexpectedToken, token)
  }
}

type characterColumnTyp struct {
  dump *DumpSQL
}

func (t *characterColumnTyp) next(token string) (state, error) {
  table := t.dump.Tables[len(t.dump.Tables)-1]
  column := table.Columns[len(table.Columns)-1]

  switch {
  case matchCharacterVaryingOption(token):
    column.TypOpt = token

    return &characterVaryingOption{dump: t.dump}, nil

  case token == ",":
    return &openBracket{dump: t.dump}, nil

  case token == ")":
    return &closeBracket{dump: t.dump}, nil

  case token == "not":
    column.ColOpt = token
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
    table := t.dump.Tables[len(t.dump.Tables)-1]
    column := table.Columns[len(table.Columns)-1]
    column.ColOpt = token

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
    table := t.dump.Tables[len(t.dump.Tables)-1]
    column := table.Columns[len(table.Columns)-1]
    column.ColOpt += " " + token

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

type columnColon struct {
  dump *DumpSQL
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
    table := t.dump.Tables[len(t.dump.Tables)-1]
    column := table.Columns[len(table.Columns)-1]
    column.ColOpt = token

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
)

var errUnexpectedToken = errors.New("unexpected token")
