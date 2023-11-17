package validator

import (
  "errors"
  "strings"
)

type Validator interface {
  Validate() error
}

func Validate(validators ...Validator) error {
  mergedErrStrParts := make([]string, 0, len(validators))

  for _, validator := range validators {
    if err := validator.Validate(); err != nil {
      mergedErrStrParts = append(mergedErrStrParts, err.Error())
    }
  }
  if len(mergedErrStrParts) != 0 {
    mergedErrStr := strings.Join(mergedErrStrParts, "\n")
    return errors.New(mergedErrStr)
  }
  return nil
}
