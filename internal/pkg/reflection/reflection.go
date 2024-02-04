package reflection

import (
  "fmt"
  "reflect"
)

func InitStruct(structPtr any) error {
  if err := validateStructPtr(structPtr); err != nil {
    return fmt.Errorf("validateStructPtr: %w", err)
  }
  refVal, refTyp := structPtrReflection(structPtr)
  initStruct(refVal, refTyp)
  return nil
}

func initStruct(refV reflect.Value, refT reflect.Type) {
  for i := 0; i < refV.NumField(); i++ {
    fieldVal := refV.Field(i)
    fieldTyp := refT.Field(i)

    switch fieldTyp.Type.Kind() {
    case reflect.Map:
      fieldVal.Set(reflect.MakeMap(fieldTyp.Type))
    case reflect.Slice:
      fieldVal.Set(reflect.MakeSlice(fieldTyp.Type, 0, 0))
    case reflect.Chan:
      fieldVal.Set(reflect.MakeChan(fieldTyp.Type, 0))
    case reflect.Struct:
      initStruct(fieldVal, fieldTyp.Type)
    case reflect.Ptr:
      fv := reflect.New(fieldTyp.Type.Elem())
      initStruct(fv.Elem(), fieldTyp.Type.Elem())
      fieldVal.Set(fv)
    }
  }
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
