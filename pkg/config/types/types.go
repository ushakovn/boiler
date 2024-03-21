package types

var valueTypes = map[string]struct{}{
  "int":      {},
  "int32":    {},
  "int64":    {},
  "float32":  {},
  "float64":  {},
  "uint32":   {},
  "uint64":   {},
  "string":   {},
  "bool":     {},
  "time":     {},
  "duration": {},
}

func IsValid(typ string) bool {
  _, ok := valueTypes[typ]
  return ok
}
