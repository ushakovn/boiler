package pgdump

import (
  "bytes"
  "fmt"
  "regexp"
  "strings"
  "unicode"

  log "github.com/sirupsen/logrus"
  "github.com/ushakovn/boiler/internal/pkg/filer"
  "github.com/ushakovn/boiler/internal/pkg/stack"
)

type DumpSQL struct {
  Tables    *stack.Stack[*DumpTable]
  tempStack *stack.Stack[string]
  option    DumpOption
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
  WithDefault  bool
}

// Do return SQL dump including table definitions from CREATE TABLE statements
func (o DumpOption) Do() (*DumpSQL, error) {
  if err := o.Validate(); err != nil {
    return nil, fmt.Errorf("pg dump option invalid: %w", err)
  }
  // Sanitize pg dump before scanning tokens
  o.pgDump = sanitizePgDump(o.pgDump)

  tokens, err := scanSchemaSQLTokens(o.pgDump)
  if err != nil {
    return nil, fmt.Errorf("scanSchemaSQLTokens: %w", err)
  }
  state, err := doTransitions(tokens, o)
  if err != nil {
    return nil, fmt.Errorf("doTransitions: %w", err)
  }
  var dump *DumpSQL

  // Confirm a terminate state
  if state, ok := state.(*terminate); ok {
    dump = state.dump
  }
  if dump == nil {
    return nil, fmt.Errorf("sql dump: not a terminate state: %T", state)
  }
  // Sanitize dumped sql table definitions
  dump = sanitizeDumpSQL(dump)

  return dump, nil
}

func sanitizeDumpSQL(dump *DumpSQL) *DumpSQL {
  sanitized := &DumpSQL{
    Tables: stack.NewStack[*DumpTable](),
  }
  for _, table := range dump.Tables.Elems() {
    var ok bool

    if _, ok = systemTablesNames[table.Name]; ok {
      continue
    }
    for _, column := range table.Columns.Elems() {
      if ok = column.IsPrimaryKey; ok {
        break
      }
    }
    if !ok {
      log.Warnf("pg_dump: table '%s' skipped: does not contain a primary key", table.Name)
      // Skip tables without primary key
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
  "rocket_locks":     {},
}

func sanitizePgDump(pgDump []byte) []byte {
  pgDumpStr := strings.ToLower(string(pgDump))

  pgDumpStr = regexSqlTimestamp.ReplaceAllLiteralString(pgDumpStr, "")

  matchCns := regexSqlConstraint.FindAllString(pgDumpStr, -1)

  for _, matchCn := range matchCns {
    if regexSqlPkConstraint.MatchString(matchCn) {
      continue
    }
    pgDumpStr = strings.Replace(pgDumpStr, matchCn, "", 1)
  }

  return []byte(pgDumpStr)
}

var (
  regexSqlTimestamp    = regexp.MustCompile(`timezone\(.+\)|now\(.*\)`)
  regexSqlConstraint   = regexp.MustCompile(`constraint\s.*`)
  regexSqlPkConstraint = regexp.MustCompile(`constraint\s.*_pkey\s.*`)
)
