package types

import (
  "fmt"

  "github.com/spf13/cast"
)

func CastValue(typ, rawValue string) (Value, error) {
  var (
    value any
    err   error
  )
  switch typ {
  case "int":
    value, err = cast.ToIntE(rawValue)
  case "int32":
    value, err = cast.ToInt32E(rawValue)
  case "int64":
    value, err = cast.ToInt64E(rawValue)
  case "float32":
    value, err = cast.ToFloat32E(rawValue)
  case "float64":
    value, err = cast.ToFloat64E(rawValue)
  case "uint32":
    value, err = cast.ToUint32E(rawValue)
  case "uint64":
    value, err = cast.ToUint64E(rawValue)
  case "string":
    value, err = cast.ToStringE(rawValue)
  case "bool":
    value, err = cast.ToBoolE(rawValue)
  case "time":
    value, err = cast.ToTimeE(rawValue)
  case "duration":
    value, err = cast.ToDurationE(rawValue)
  default:
    err = fmt.Errorf("unexpected config value type: %s", typ)
  }
  return NewValue(value), err
}
