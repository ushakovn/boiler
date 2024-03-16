package stringer

import (
  "regexp"
  "strings"
  "unicode"

  "github.com/iancoleman/strcase"
)

var (
  regexCamelCase      = regexp.MustCompile(`^[a-zA-Z]+[A-Z0-9]*([a-z]+[A-Z0-9]*)*$`)
  regexSnakeCase      = regexp.MustCompile(`^[a-z]+([a-z_0-9]+)*[a-z0-9]?$`)
  regexCapitalizeCase = regexp.MustCompile(`^[A-Z]+(_?[A-Z]+)*$`)
)

func IsCamelCase(s string) bool {
  return regexCamelCase.MatchString(s)
}

func IsSnakeCase(s string) bool {
  return regexSnakeCase.MatchString(s)
}

func IsCapitalizeCase(s string) bool {
  return regexCapitalizeCase.MatchString(s)
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
  runes := []rune(s)
  count := len(runes)

  out := make([]rune, 0, count)
  var prevUpper bool

  for index := 0; index < count; index++ {
    ch := runes[index]

    if index == 0 || unicode.IsUpper(ch) {
      if index != 0 && (!prevUpper || (index+1 < count && unicode.IsLower(runes[index+1]))) {
        prevUpper = true

        out = append(out, '_')
      }
      out = append(out, unicode.ToLower(ch))

      continue
    }
    prevUpper = false

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
  // Fallback to strcase package
  return strcase.ToSnake(s)
}

func StringToUpperCamelCase(s string) string {
  if IsSnakeCase(s) {
    return SnakeCaseToUpperCamelCase(s)
  }
  if IsCamelCase(s) {
    return CamelCaseToUpperCamelCase(s)
  }
  // Fallback to strcase package
  return strcase.ToCamel(s)
}

func StringToLowerCamelCase(s string) string {
  if IsSnakeCase(s) {
    return SnakeCaseToLowerCamelCase(s)
  }
  if IsCamelCase(s) {
    return CamelCaseToLowerCamelCase(s)
  }
  // Fallback to strcase package
  return strcase.ToLowerCamel(s)
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
  if IsCapitalizeCase(s) {
    return s
  }
  if IsCamelCase(s) {
    return CamelCaseToCapitalizeCase(s)
  }
  if IsSnakeCase(s) {
    return SnakeCaseToCapitalizeCase(s)
  }
  // Fallback to strcase package
  return strcase.ToScreamingSnake(s)
}

func CamelCaseToCapitalizeCase(s string) string {
  s = CamelCaseToSnakeCase(s)
  s = SnakeCaseToCapitalizeCase(s)
  return s
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

func NormalizeName(s string) string {
  s = TrimPluralForm(s)
  return s
}

func TrimPluralForm(src string) string {
  return strings.TrimSuffix(src, "s")
}

func UnquoteString(s string) string {
  return strings.Trim(s, `"`)
}
