package goose

import (
  "fmt"
  "regexp"
  "strings"
  "time"

  "github.com/ushakovn/boiler/internal/pkg/filer"
)

const timestampFormat = "20060102150405"

var timestampRegex = regexp.MustCompile(`\d{14}_`)

func BuildFileName(fileName string) (gooseFileName string) {
  return buildFileName(fileName, time.Now())
}

func ExtractFileName(gooseFileName string) (fileName string) {
  gooseFileName = strings.TrimSuffix(gooseFileName, ".sql")
  return timestampRegex.ReplaceAllLiteralString(gooseFileName, "")
}

func IsExistedMigration(filePath string) bool {
  gooseFileName := filer.ExtractFileName(filePath)
  fileName := ExtractFileName(gooseFileName)

  fileDirPath := strings.TrimSuffix(filePath, fmt.Sprint("/", gooseFileName))
  return filer.IsExistedPattern(fileDirPath, fileName)
}

func ChangeTime(gooseFileName string, changeTime func(fileTime time.Time) time.Time) (string, error) {
  fileTimestamp := timestampRegex.FindString(gooseFileName)
  fileTimestamp = strings.TrimSuffix(fileTimestamp, "_")

  fileTime, err := time.Parse(timestampFormat, fileTimestamp)
  if err != nil {
    return "", fmt.Errorf("time parse failed: %w", err)
  }
  fileTime = changeTime(fileTime)

  fileName := ExtractFileName(gooseFileName)
  gooseFileName = buildFileName(fileName, fileTime)

  return gooseFileName, nil
}

func buildFileName(fileName string, fileTime time.Time) (gooseFileName string) {
  fileTimestamp := fileTime.Format(timestampFormat)
  return fmt.Sprintf("%s_%s.sql", fileTimestamp, fileName)
}
