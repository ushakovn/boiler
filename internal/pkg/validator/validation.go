package validator

import (
  "errors"
  "strings"
)

type ValidateFunc func() error

func Validate(fs ...ValidateFunc) error {
  errStrParts := make([]string, 0, len(fs))

  for _, f := range fs {
    if err := f(); err != nil {
      errStrParts = append(errStrParts, err.Error())
    }
  }
  if len(errStrParts) != 0 {
    errStr := strings.Join(errStrParts, "\n")
    return errors.New(errStr)
  }
  return nil
}
