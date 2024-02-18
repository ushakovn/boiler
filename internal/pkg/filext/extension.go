package filext

import "fmt"

const (
  goExt    = "go"
  protoExt = "proto"

  pattern = "%s.%s"
)

func build(fileName, fileExt string) string {
  return fmt.Sprintf(pattern, fileName, fileExt)
}

func Go(fileName string) string {
  return build(fileName, goExt)
}

func Proto(fileName string) string {
  return build(fileName, protoExt)
}
