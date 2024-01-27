package validator

import (
  "errors"
  "fmt"
  "reflect"
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

func ValidateStructWithTags[T any](structPtr T, tagKeys ...string) error {
  if len(tagKeys) == 0 {
    return fmt.Errorf("tag keys not specified")
  }
  if err := validateStructPtr(structPtr); err != nil {
    return fmt.Errorf("validateStructPtr: %w", err)
  }
  refVal, refType := structPtrReflection(structPtr)

  for fieldIdx := 0; fieldIdx < refVal.NumField(); fieldIdx++ {
    for _, tagKey := range tagKeys {
      if tagVal := refType.Field(fieldIdx).Tag.Get(tagKey); tagVal != "" {
        var fieldVal reflect.Value

        if fieldVal = refVal.Field(fieldIdx); !fieldVal.IsValid() || fieldVal.IsZero() {
          return fmt.Errorf("not specified value for field \"%s\" with tags `%s:\"%s\"`",
            refType.Field(fieldIdx).Name, tagKeys, tagVal)
        }
        if fieldIface := fieldVal.Interface(); validateStructPtr(fieldIface) == nil {
          return ValidateStructWithTags(fieldIface, tagKeys...)
        }
      }
    }
  }
  return nil
}

func structPtrReflection(structPtr any) (reflect.Value, reflect.Type) {
  refVal := reflect.ValueOf(structPtr).Elem()
  refType := reflect.TypeOf(structPtr).Elem()
  return refVal, refType
}

func validateStructPtr(structPtr any) error {
  if refStructPtr := reflect.ValueOf(structPtr); refStructPtr.Kind() == reflect.Ptr {
    if refStruct := refStructPtr.Elem(); refStruct.Kind() == reflect.Struct {
      return nil
    }
    return fmt.Errorf("function argument not a struct pointer")
  }
  return fmt.Errorf("function argument not a pointer")
}

func OneOfNotZero(values ...any) bool {
  var ok bool

  for _, value := range values {
    if !IsZeroValue(value) {
      if ok {
        return false
      }
      ok = true
    }
  }
  return ok
}

func IsZeroValue(value any) bool {
  v := reflect.ValueOf(value)
  return !v.IsValid() || reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}
