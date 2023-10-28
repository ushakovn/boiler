package utils

import (
  "bufio"
  "bytes"
  "fmt"
  "go/format"
  "os"
  "path/filepath"
  "strings"
  "text/template"

  log "github.com/sirupsen/logrus"
)

func IsExistedDirectory(dirPath string) bool {
  if info, err := os.Stat(dirPath); err == nil && info.IsDir() {
    return true
  }
  return false
}

func ExecuteTemplateCopy(templateCompiled, filePath string, structPtr any, funcMap template.FuncMap) error {
  t := template.New("")
  if len(funcMap) > 0 {
    t = t.Funcs(funcMap)
  }
  t, err := t.Parse(templateCompiled)
  if err != nil {
    return fmt.Errorf("template.New().Parse: %w", err)
  }
  var (
    buffer bytes.Buffer
    buf    []byte
  )
  if err = t.Execute(&buffer, structPtr); err != nil {
    return fmt.Errorf("t.Execute: %w", err)
  }
  if buf, err = format.Source(buffer.Bytes()); err != nil {
    return fmt.Errorf("format.Source: %w", err)
  }
  if err = os.WriteFile(filePath, buf, os.ModePerm); err != nil {
    return fmt.Errorf("os.WriteFile: %w", err)
  }
  return nil
}

func CopyTemplate(templateCompiled string, filePath string) error {
  if err := os.WriteFile(filePath, []byte(templateCompiled), os.ModePerm); err != nil {
    return fmt.Errorf("os.WriteFile: %w", err)
  }
  return nil
}

func WorkDirPath() (string, error) {
  const env = "PWD"
  path, err := Env(env)
  if err != nil {
    return "", err
  }
  return path, nil
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

func Map[T any, S any](t []T, f func(T) S) []S {
  s := make([]S, 0, len(t))

  for _, t := range t {
    s = append(s, f(t))
  }
  return s
}

func MapLookup[K comparable, V any, M map[K]V](m M, k K) bool {
  _, ok := m[k]
  return ok
}
