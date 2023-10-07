package utils

import (
  "bufio"
  "fmt"
  "io"
  "regexp"
  "strings"
  "unicode"
)

var (
  regexCamelCase = regexp.MustCompile(`^[a-zA-Z]+[A-Z0-9]*([a-z]+[A-Z0-9]*)*$`)
  regexSnakeCase = regexp.MustCompile(`^[a-z]+([a-z_0-9]+)*[a-z0-9]?$`)
)

func IsCamelCase(s string) bool {
  return regexCamelCase.MatchString(s)
}

func IsSnakeCase(s string) bool {
  return regexSnakeCase.MatchString(s)
}

func SnakeCaseToUpperCamelCase(s string) string {
  out := make([]rune, 0, len([]rune(s)))
  var ok bool

  for index, ch := range s {
    if index == 0 || ok {
      out = append(out, unicode.ToUpper(ch))
      ok = false
      continue
    }
    if ch == '_' {
      ok = true
      continue
    }
    out = append(out, ch)
  }

  return string(out)
}

func CamelCaseToSnakeCase(s string) string {
  out := make([]rune, 0, len([]rune(s)))

  for index, ch := range s {
    if index == 0 || unicode.IsUpper(ch) {
      if index != 0 {
        out = append(out, '_')
      }
      out = append(out, unicode.ToLower(ch))
      continue
    }
    out = append(out, ch)
  }

  return string(out)
}

func CamelCaseToUpperCamelCase(s string) string {
  out := make([]rune, 0, len([]rune(s)))

  for index, ch := range s {
    if index == 0 && unicode.IsLower(ch) {
      ch = unicode.ToUpper(ch)
    }
    out = append(out, ch)
  }

  return string(out)
}

func StringToSnakeCase(s string) string {
  if IsCamelCase(s) {
    return CamelCaseToSnakeCase(s)
  }
  if IsSnakeCase(s) {
    return s
  }
  return ""
}

func StringToUpperCamelCase(s string) string {
  if IsSnakeCase(s) {
    return SnakeCaseToUpperCamelCase(s)
  }
  if IsCamelCase(s) {
    return CamelCaseToUpperCamelCase(s)
  }
  return ""
}

func StringToLowerCase(s string) string {
  out := make([]rune, 0, len([]rune(s)))

  for _, ch := range s {
    if unicode.IsUpper(ch) {
      ch = unicode.ToLower(ch)
    }
    out = append(out, ch)
  }

  return string(out)
}

func IsWrongCase(s string) bool {
  isCamel := IsCamelCase(s)
  isSnake := IsSnakeCase(s)
  return !isCamel && !isSnake
}

func NormalizeToken(s string) string {
  s = strings.TrimSpace(s)
  s = strings.ToLower(s)
  return s
}

func StringOneOfEqual(src string, dst ...string) bool {
  for _, dst := range dst {
    if src == dst {
      return true
    }
  }
  return false
}

func ScanLines(r io.Reader, f func(line string) error) error {
  s := bufio.NewScanner(r)
  s.Split(bufio.ScanLines)
  var err error

  for s.Scan() {
    if err = f(s.Text()); err != nil {
      return fmt.Errorf("f(s.Text()): %w", err)
    }
  }
  if err = s.Err(); err != nil {
    return fmt.Errorf("s.Err: %w", err)
  }
  return nil
}
