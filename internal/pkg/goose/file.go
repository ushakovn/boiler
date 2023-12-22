package goose

import (
  "fmt"
  "regexp"
  "time"
)

const timestampFormat = "20060102150405"

var timestampRegex = regexp.MustCompile(`\d{14}_`)

func BuildFileName(fileName string) (gooseFileName string) {
  now := time.Now().Format(timestampFormat)
  return fmt.Sprintf("%s_%s.sql", now, fileName)
}

func ExtractFileName(gooseFileName string) (fileName string) {
  return timestampRegex.ReplaceAllLiteralString(gooseFileName, "")
}
