package stringer

import (
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

func StringToLowerCamelCase(s string) string {
  if IsSnakeCase(s) {
    return SnakeCaseToLowerCamelCase(s)
  }
  if IsCamelCase(s) {
    return CamelCaseToLowerCamelCase(s)
  }
  return ""
}

func CamelCaseToLowerCamelCase(s string) string {
  out := make([]rune, 0, len([]rune(s)))

  for index, ch := range s {
    if index == 0 && unicode.IsUpper(ch) {
      ch = unicode.ToLower(ch)
    }
    out = append(out, ch)
  }
  return string(out)
}

func SnakeCaseToLowerCamelCase(s string) string {
  out := make([]rune, 0, len([]rune(s)))
  var ok bool

  for _, ch := range s {
    if ch == '_' {
      ok = true
      continue
    }
    if ok {
      ch = unicode.ToUpper(ch)
      ok = false
    }
    out = append(out, ch)
  }
  return string(out)
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

func StringToCapitalizeCase(s string) string {
  if IsCamelCase(s) {
    return CamelCaseToCapitalizeCase(s)
  }
  if IsSnakeCase(s) {
    return SnakeCaseToCapitalizeCase(s)
  }
  return ""
}

func CamelCaseToCapitalizeCase(s string) string {
  count := len([]rune(s))
  out := make([]rune, 0, len([]rune(s)))

  for index, ch := range s {
    if index != 0 && index != count-1 && unicode.IsUpper(ch) {
      out = append(out, '_')
    }
    out = append(out, unicode.ToUpper(ch))
  }
  return string(out)
}

func SnakeCaseToCapitalizeCase(s string) string {
  out := make([]rune, 0, len([]rune(s)))

  for _, ch := range s {
    out = append(out, unicode.ToUpper(ch))
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
