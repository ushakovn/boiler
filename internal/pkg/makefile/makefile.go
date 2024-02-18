package makefile

import (
  "fmt"
  "os"
  "strings"

  "github.com/ushakovn/boiler/internal/pkg/filer"
)

func ContainsTarget(filePath, targetName string) (bool, error) {
  targetName = fmt.Sprint(targetName, ":")
  var ok bool

  err := filer.ScanFileWithBreak(filePath, func(line string) bool {
    line = strings.TrimSpace(line)

    if ok = strings.HasPrefix(line, targetName); ok {
      return false
    }
    return true
  })
  if err != nil {
    if os.IsNotExist(err) {
      return false, nil
    }
    return false, fmt.Errorf("filer.ScanFileWithBreak: %w", err)
  }
  return ok, nil
}
