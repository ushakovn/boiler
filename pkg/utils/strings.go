package utils

import (
	"regexp"
	"unicode"
)

var (
	regexCamelCase = regexp.MustCompile(`^[a-zA-Z]+([A-Z][a-z]+)+$`)
	regexSnakeCase = regexp.MustCompile(`^[a-z]+(\_[a-z]+)?$`)
)

func ToStructTag(s string) string {
	if IsCamelCase(s) {
		return CamelCaseToSnakeCase(s)
	}
	if IsSnakeCase(s) {
		return s
	}
	return ""
}

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
