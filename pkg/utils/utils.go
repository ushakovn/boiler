package utils

import (
  "fmt"
  "os"
  "strings"
)

func Map[T any, S any](t []T, f func(T) S) []S {
  s := make([]S, 0, len(t))

  for _, t := range t {
    s = append(s, f(t))
  }
  return s
}

func Env(key string) (string, error) {
  if val := os.Getenv(key); val != "" {
    return val, nil
  }
  return "", fmt.Errorf("%s value not found", key)
}

func ExtractFileExtension(fullName string) string {
  parts := strings.Split(fullName, ".")
  return parts[len(parts)-1]
}

func MapLookup[K comparable, V any, M map[K]V](m M, k K) bool {
  _, ok := m[k]
  return ok
}
