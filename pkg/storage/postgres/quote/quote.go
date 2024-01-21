package quote

import (
  "encoding/hex"
  "fmt"
  "strconv"
  "strings"
  "time"
)

func String(arg any) string {
  var quote string

  switch t := arg.(type) {
  case nil:
    quote = "null"
  case int64:
    quote = strconv.FormatInt(t, 10)
  case float64:
    quote = strconv.FormatFloat(t, 'f', -1, 64)
  case bool:
    quote = strconv.FormatBool(t)
  case time.Time:
    quote = t.Format("'2006-01-02 15:04:05.999999999Z07:00:00'")
  case time.Duration:
    quote = fmt.Sprintf("'%d milliseconds'", t.Milliseconds())
  case []byte:
    quote = fmt.Sprintf(`'\x%s'::bytea`, hex.EncodeToString(t))
  case string:
    quote = fmt.Sprintf(`'%s'`, strings.Replace(t, "'", "''", -1))
  }
  return quote
}
