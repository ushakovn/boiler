package makefile

import (
  "bytes"
  "fmt"
  "os"
  "strings"

  "github.com/ushakovn/boiler/internal/pkg/filer"
)

func ContainsTarget(filePath, targetName string) (bool, error) {
  buf, err := os.ReadFile(filePath)
  if err != nil {
    if os.IsNotExist(err) {
      return false, nil
    }
    return false, fmt.Errorf("os.ReadFile: %w", err)
  }
  targetName = fmt.Sprint(targetName, ":")
  var foundTarget bool

  if err = filer.ScanLinesWithBreak(bytes.NewReader(buf), func(line string) bool {
    line = strings.TrimSpace(line)

    if foundTarget = strings.HasPrefix(line, targetName); foundTarget {
      return false
    }
    return true

  }); err != nil {
    return false, fmt.Errorf("filer.ScanLinesWithBreak: %w", err)
  }
  return foundTarget, nil
}
