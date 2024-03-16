package filer

import (
  "bufio"
  "bytes"
  "fmt"
  "io"
  "os"
  "path/filepath"
  "strings"

  log "github.com/sirupsen/logrus"
  "github.com/ushakovn/boiler/internal/pkg/env"
)

func IsExistedPattern(fileDirPath, fileNamePattern string) bool {
  entries, err := os.ReadDir(fileDirPath)
  if err != nil {
    return false
  }
  for _, entry := range entries {
    entryName := entry.Name()

    if strings.Contains(entryName, fileNamePattern) {
      return true
    }
  }
  return false
}

func IsExistedDirectory(dirPath string) bool {
  if info, err := os.Stat(dirPath); err == nil && info.IsDir() {
    return true
  }
  return false
}

func IsExistedFile(filePath string) bool {
  if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
    return true
  }
  return false
}

func WorkDirPath() (string, error) {
  path, err := env.Env("PWD")
  if err != nil {
    return "", err
  }
  return path, nil
}

func ExtractFileExtension(fullName string) string {
  parts := strings.Split(fullName, ".")
  ext := strings.TrimSpace(parts[len(parts)-1])
  ext = strings.ToLower(ext)
  return ext
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
    if closeErr := file.Close(); closeErr != nil {
      log.Infof("ExtractModuleName: file.Close: %v", closeErr)
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

func CreateNestedFolders(sourcePath string, destNestedFolders ...string) (string, error) {
  defaultDirParts := append([]string{sourcePath}, destNestedFolders...)
  defaultDir := filepath.Join(defaultDirParts...)

  prevDirParts := make([]string, 0, len(defaultDirParts))

  for _, dirPart := range defaultDirParts {
    // Create directories
    prevDirParts = append(prevDirParts, dirPart)
    curPath := filepath.Join(prevDirParts...)

    if _, err := os.Stat(curPath); os.IsNotExist(err) {
      if err = os.Mkdir(curPath, os.ModePerm); err != nil {
        return "", fmt.Errorf("os.Mkdir: %w", err)
      }
    }
  }
  // Check created directories
  if _, err := os.Stat(defaultDir); os.IsNotExist(err) {
    return "", fmt.Errorf("os.Stat: %s: err: %v", defaultDir, err)
  }

  return defaultDir, nil
}

func CollectNestedFilesPath(extension, parentPath string) ([]string, error) {
  filesPath := make([]string, 0)

  if err := collectNestedFilesPath(extension, parentPath, &filesPath); err != nil {
    return nil, fmt.Errorf("collectNestedFilesPath: %w", err)
  }
  return filesPath, nil
}

func collectNestedFilesPath(extension, parentPath string, filesPath *[]string) error {
  entries, err := os.ReadDir(parentPath)
  if err != nil {
    return fmt.Errorf("os.ReadDir: %w", err)
  }
  for _, entry := range entries {
    if !entry.IsDir() && ExtractFileExtension(entry.Name()) == extension {
      path := filepath.Join(parentPath, entry.Name())

      *filesPath = append(*filesPath, path)
      continue
    }
  }
  for _, entry := range entries {
    if entry.IsDir() {
      childPath := filepath.Join(parentPath, entry.Name())

      if err = collectNestedFilesPath(extension, childPath, filesPath); err != nil {
        return fmt.Errorf("g.collectNestedFilesPath: %w", err)
      }
    }
  }
  return nil
}

func ExtractFileName(filePath string) string {
  if parts := strings.Split(filePath, `/`); len(parts) > 0 {
    return parts[len(parts)-1]
  }
  if parts := strings.Split(filePath, `\`); len(parts) > 0 {
    return parts[len(parts)-1]
  }
  return filePath
}

func AppendStringToFile(filePath, rawString string) error {
  // If the file does not exist create it or append to the file
  file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
  if err != nil {
    return fmt.Errorf("os.OpenFile: %w", err)
  }
  // Close file
  defer func() {
    if closeErr := file.Close(); closeErr != nil {
      log.Errorf("file.Close: %v", err)
    }
  }()
  // Append blank lines around target raw string
  targetStrings := []string{
    "\n", rawString, "\n",
  }
  for _, target := range targetStrings {
    // WriteWithBreak any target string
    if _, err = file.WriteString(target); err != nil {
      return fmt.Errorf("file.WriteString: %w", err)
    }
  }
  return nil
}

func ScanFileWithBreak(filePath string, f func(line string) bool) error {
  buf, err := os.ReadFile(filePath)
  if err != nil {
    return err
  }
  return ScanLinesWithBreak(bytes.NewReader(buf), f)
}

func ScanLinesWithBreak(r io.Reader, f func(line string) bool) error {
  s := bufio.NewScanner(r)
  s.Split(bufio.ScanLines)

  for s.Scan() {
    if !f(s.Text()) {
      break
    }
  }
  if err := s.Err(); err != nil {
    return fmt.Errorf("s.Err: %w", err)
  }
  return nil
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

func EditFile(filePath string, f func(line string) string) error {
  file, err := os.Open(filePath)
  if err != nil {
    return err
  }
  s := bufio.NewScanner(file)
  s.Split(bufio.ScanLines)

  w := bufio.NewWriter(file)

  for s.Scan() {
    text := f(s.Text())

    if _, err = w.Write([]byte(text)); err != nil {
      return fmt.Errorf("writer.Write: %w", err)
    }
  }
  if err = s.Err(); err != nil {
    return fmt.Errorf("scanner.Err: %w", err)
  }
  return nil
}