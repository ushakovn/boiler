package sql

import (
  "bytes"
  "fmt"
  "unicode"

  "github.com/ushakovn/boiler/pkg/utils"
)

type DumpSQL struct {
  Tables    *utils.Stack[*DumpTable]
  tempStack *utils.Stack[string]
}

type DumpTable struct {
  RawName string
  Name    string
  Schema  string
  Columns utils.Stack[*DumpColumn]
}

type DumpColumn struct {
  Name         string
  Typ          string
  TypOptions   string
  IsNotNull    bool
  IsPrimaryKey bool
}

// DumpSchemaSQL returns SQL dump including table definitions from CREATE TABLE statements
func DumpSchemaSQL(option PgDumpOption) (*DumpSQL, error) {
  pgDumpBuf, err := option.Call()
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
  if state, ok := state.(*terminate); ok {
    return state.dump, nil
  }
  return nil, fmt.Errorf("not a terminate state: %T", state)
}

func scanSchemaSQLTokens(pgDump []byte) ([]string, error) {
  var tokens []string

  // Scan tokens from file
  if err := utils.ScanLines(bytes.NewReader(pgDump), func(line string) error {
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
    return nil, fmt.Errorf("utils.ScanLines: %w", err)
  }
  return tokens, nil
}

var sqlStickyTokens = map[rune]struct{}{
  ',': {},
  ';': {},
}
