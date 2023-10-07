package sql

import (
  "bytes"
  "fmt"
  "os"
  "os/exec"
  "unicode"

  "github.com/ushakovn/boiler/pkg/utils"
)

type DumpSQL struct {
  Tables []*DumpTable
}

type DumpTable struct {
  RawName string
  Name    string
  Schema  string
  Columns []*DumpColumn
}

type DumpColumn struct {
  Name   string
  Typ    string
  TypOpt string
  ColOpt string
}

// DumpSchemaSQL returns SQL dump including table definitions from CREATE TABLE statements
func DumpSchemaSQL(conn *PgConn) (*DumpSQL, error) {
  pgDumpBuf, err := execPgDump(conn)
  if err != nil {
    return nil, fmt.Errorf("execPgDump: %w", err)
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

type PgConn struct {
  DB   string
  Host string
  Port string
  User string
  Pass string
}

func (c *PgConn) pgDump() (name string, args []string, err error) {
  if err := os.Setenv("PGPASSWORD", c.Pass); err != nil {
    return "", nil, fmt.Errorf("os.Setenv: %w", err)
  }
  args = []string{
    "--no-owner",
    "-h", c.Host,
    "-p", c.Port,
    "-U", c.User,
    c.DB,
  }
  name = "pg_dump"

  return name, args, nil
}

func execPgDump(conn *PgConn) ([]byte, error) {
  name, args, err := conn.pgDump()
  if err != nil {
    return nil, fmt.Errorf("conn.pgDump: %w", err)
  }
  cmd := exec.Command(name, args...)
  if cmd.Err != nil {
    return nil, fmt.Errorf("exec.CommandContext %w", cmd.Err)
  }
  buf, err := cmd.Output()
  if err != nil {
    return nil, fmt.Errorf("cmd.Output: %w", err)
  }
  return buf, nil
}
