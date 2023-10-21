package utils

import (
  "bufio"
  "fmt"
  "os"
  "path/filepath"
  "strings"

  log "github.com/sirupsen/logrus"
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

func ExtractGoModuleName(workDirPath string) (string, error) {
  const (
    moduleFile   = "go.mod"
    modulePrefix = "module"
  )
  filePath := filepath.Join(workDirPath, moduleFile)

  file, err := os.Open(filePath)
  if err != nil {
    return "", fmt.Errorf("os.Open: %w", err)
  }

  defer func() {
    if err := file.Close(); err != nil {
      log.Infof("ExtractModuleName: file.Close: %v", err)
    }
  }()
  sc := bufio.NewScanner(file)

  for sc.Scan() {
    if moduleLine := sc.Text(); strings.HasPrefix(moduleLine, modulePrefix) {
      moduleLine = strings.TrimPrefix(moduleLine, modulePrefix)
      moduleLine = strings.TrimSpace(moduleLine)
      return moduleLine, nil
    }
  }

  return "", fmt.Errorf("module name not found in file: %s", filePath)
}
