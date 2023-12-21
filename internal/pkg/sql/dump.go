package sql

import (
  "bytes"
  "context"
  "fmt"
  "strings"
  "unicode"

  "github.com/ushakovn/boiler/internal/pkg/filer"
  "github.com/ushakovn/boiler/internal/pkg/stack"
)

type DumpSQL struct {
  Tables    *stack.Stack[*DumpTable]
  tempStack *stack.Stack[string]
}

type DumpTable struct {
  RawName string
  Name    string
  Schema  string
  Columns stack.Stack[*DumpColumn]
}

type DumpColumn struct {
  Name         string
  Typ          string
  TypOptions   string
  IsNotNull    bool
  IsPrimaryKey bool
}

// DumpSchemaSQL returns SQL dump including table definitions from CREATE TABLE statements
func DumpSchemaSQL(ctx context.Context, option PgDumpOption) (*DumpSQL, error) {
  pgDumpBuf, err := option.Call(ctx)
  if err != nil {
    return nil, fmt.Errorf("option.Call: %w", err)
  }
  tokens, err := scanSchemaSQLTokens(pgDumpBuf)
  if err != nil {
    return nil, fmt.Errorf("scanSchemaSQLTokens: %w", err)
  }
  state, err := doTransitions(tokens)
  if err != nil {
    return nil, fmt.Errorf("doTransitions: %w", err)
  }
  var dump *DumpSQL

  if state, ok := state.(*terminate); ok {
    dump = state.dump
  }
  if dump == nil {
    return nil, fmt.Errorf("sql dump: not a terminate state: %T", state)
  }
  dump = sanitizeDumpSQL(dump)

  return dump, nil
}

func sanitizeDumpSQL(dump *DumpSQL) *DumpSQL {
  sanitized := &DumpSQL{
    Tables: stack.NewStack[*DumpTable](),
  }
  for _, table := range dump.Tables.Elems() {
    if _, ok := systemTablesNames[table.Name]; ok || hasBoilerSystemPrefix(table.Name) {
      continue
    }
    sanitized.Tables.Push(table)
  }
  return sanitized
}

func scanSchemaSQLTokens(pgDump []byte) ([]string, error) {
  var tokens []string

  // Scan tokens from file
  if err := filer.ScanLines(bytes.NewReader(pgDump), func(line string) error {
    var token []rune

    // Scan current token
    for _, ch := range line {
      if _, ok := sqlStickyTokens[ch]; ok || unicode.IsSpace(ch) {
        if len(token) != 0 {
          // Push previous token
          tokens = append(tokens, string(token))
        }
        if ok {
          // Push current char as separate token
          tokens = append(tokens, string(ch))
        }
        token = token[:0]

        continue
      }
      // Convert all runes to lower case
      if unicode.IsUpper(ch) {
        ch = unicode.ToLower(ch)
      }
      token = append(token, ch)
    }
    if len(token) != 0 {
      tokens = append(tokens, string(token))
    }
    return nil

  }); err != nil {
    return nil, fmt.Errorf("filer.ScanLines: %w", err)
  }
  return tokens, nil
}

var sqlStickyTokens = map[rune]struct{}{
  ',': {},
  ';': {},
}

var systemTablesNames = map[string]struct{}{
  "goose_db_version": {},
  "db_version":       {},
}

func hasBoilerSystemPrefix(tableName string) bool {
  return strings.HasPrefix(tableName, "__boiler")
}
