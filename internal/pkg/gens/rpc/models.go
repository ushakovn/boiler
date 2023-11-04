package rpc

import (
  "fmt"
  "regexp"
  "strings"

  "github.com/ushakovn/boiler/internal/pkg/aggr"
  "github.com/ushakovn/boiler/internal/pkg/stringer"
)

type rootDesc struct {
  Handles  []*handleDesc   `json:"handles" yaml:"handles"`
  TypeDefs []*typeDefsDesc `json:"type_defs" yaml:"type_defs"`
}

type handleDesc struct {
  Name     string        `json:"name" yaml:"name"`
  Route    string        `json:"route" yaml:"route"`
  Request  *contractDesc `json:"request" yaml:"request"`
  Response *contractDesc `json:"response" yaml:"response"`
}

type contractDesc struct {
  Fields []*fieldDesc `json:"fields" yaml:"fields"`
}

type fieldDesc struct {
  Name string `json:"name" yaml:"name"`
  Type string `json:"type" yaml:"type"`
}

type typeDefsDesc struct {
  Name   string       `json:"name" yaml:"name"`
  Fields []*fieldDesc `json:"fields" yaml:"fields"`
}

func (d *rootDesc) Validate() error {
  defTypesMap := map[string]struct{}{}

  for _, typeDef := range d.TypeDefs {
    if typeDef.Name == "" {
      return fmt.Errorf("type definition name not specified")
    }
    if stringer.IsWrongCase(typeDef.Name) {
      return fmt.Errorf("invalid type definition name: %s", typeDef.Name)
    }
    if _, ok := defTypesMap[typeDef.Name]; ok {
      return fmt.Errorf("duplicated type definition for type: %s", typeDef.Name)
    }
    defTypesMap[typeDef.Name] = struct{}{}
  }

  for _, typeDef := range d.TypeDefs {
    if err := typeDef.Validate(defTypesMap); err != nil {
      return fmt.Errorf("type definition: %w", err)
    }
  }

  if len(d.Handles) == 0 {
    return fmt.Errorf("rpc handler not specfied")
  }
  handlerNamesMap := map[string]struct{}{}
  handlerRouterMap := map[string]struct{}{}

  for _, handler := range d.Handles {
    if stringer.IsWrongCase(handler.Name) {
      return fmt.Errorf("invalid handler name: %s", handler.Name)
    }
    if _, ok := handlerNamesMap[handler.Name]; ok {
      return fmt.Errorf("duplicated handler name: %s", handler.Name)
    }
    handlerNamesMap[handler.Name] = struct{}{}

    if stringer.IsWrongCase(handler.Route) {
      return fmt.Errorf("invalid handler route: %s", handler.Route)
    }
    if _, ok := handlerRouterMap[stringer.StringToLowerCase(handler.Route)]; ok {
      return fmt.Errorf("duplicated handler route: %s", handler.Route)
    }
    handlerRouterMap[stringer.StringToLowerCase(handler.Route)] = struct{}{}

    if err := handler.Validate(defTypesMap); err != nil {
      return fmt.Errorf("handler error: %v", err)
    }
  }
  return nil
}

func (d *handleDesc) Validate(types map[string]struct{}) error {
  if d.Name == "" {
    return fmt.Errorf("name not specified")
  }
  if stringer.IsWrongCase(d.Name) {
    return fmt.Errorf("invalid name: %s", d.Name)
  }
  if d.Request == nil {
    return fmt.Errorf("request not specified")
  }
  if err := d.Request.Validate(types); err != nil {
    return fmt.Errorf("request error: %v", err)
  }
  if err := d.Response.Validate(types); err != nil {
    return fmt.Errorf("response error: %v", err)
  }
  return nil
}

func (d *contractDesc) Validate(types map[string]struct{}) error {
  fieldsNamesMap := map[string]struct{}{}

  for _, field := range d.Fields {
    if field.Name == "" {
      return fmt.Errorf("field name not specified")
    }
    if stringer.IsWrongCase(field.Name) {
      return fmt.Errorf("invalid field name: %s", field.Name)
    }
    if _, ok := fieldsNamesMap[field.Name]; ok {
      return fmt.Errorf("duplicated field: %s", field.Name)
    }
    fieldsNamesMap[field.Name] = struct{}{}

    if err := field.Validate(types); err != nil {
      return fmt.Errorf("field: %s: error: %v", field.Name, err)
    }
  }
  return nil
}

func (d *fieldDesc) Validate(types map[string]struct{}) error {
  if d.Name == "" {
    return fmt.Errorf("field name not specified")
  }
  if d.Type == "" {
    return fmt.Errorf("field type not specified")
  }
  if _, ok := scalarTypesMap[d.Type]; ok {
    return nil
  }
  if typeScalarSliceRegex.MatchString(d.Type) {
    return nil
  }
  if typeDefSliceRegex.MatchString(d.Type) {
    if _, ok := types[strings.TrimPrefix(d.Type, slicePrefix)]; !ok {
      return fmt.Errorf("unexpected field type: %s", d.Type)
    }
  }
  return nil
}

func (d *typeDefsDesc) Validate(types map[string]struct{}) error {
  if d.Name == "" {
    return fmt.Errorf("type name not specified")
  }
  fieldsNamesMap := map[string]struct{}{}

  for _, field := range d.Fields {
    if _, ok := fieldsNamesMap[field.Name]; ok {
      return fmt.Errorf("duplicated field: %s", field.Name)
    }
    fieldsNamesMap[field.Name] = struct{}{}

    if err := field.Validate(types); err != nil {
      return fmt.Errorf("field: %s: error: %v", field.Name, err)
    }
  }
  return nil
}

var scalarTypesMap = map[string]struct{}{
  "int":     {},
  "int32":   {},
  "int64":   {},
  "float32": {},
  "float64": {},
  "bool":    {},
  "byte":    {},
  "string":  {},
}

var (
  typeDefSliceRegex    = regexp.MustCompile(`^\[\]\w+$`)
  typeScalarSliceRegex = regexp.MustCompile(`^\[\]((bool|string|byte|int)|(float|int)(32|64))$`)
)

type rpcTemplates struct {
  Handles   []*rpcHandle
  Contracts *rpcContracts
}

type rpcHandle struct {
  Name  string
  Route string
}

type rpcContracts struct {
  Requests  []*rpcContract
  Responses []*rpcContract
  TypeDefs  []*rpcTypeDef
}

type rpcContract struct {
  Name   string
  Fields []*rpcStructField
}

type rpcTypeDef struct {
  Name   string
  Fields []*rpcStructField
}

type rpcStructField struct {
  Name string
  Type string
  Tag  string
}

func rootDescToRpcTemplates(desc *rootDesc) *rpcTemplates {
  handlesCount := len(desc.Handles)

  handles := make([]*rpcHandle, 0, handlesCount)
  reqs := make([]*rpcContract, 0, handlesCount)
  resp := make([]*rpcContract, 0, handlesCount)

  for _, handleDesc := range desc.Handles {
    handles = append(handles, handleDescToRpcHandle(handleDesc))
    reqs = append(reqs, contractDescToRpcContract(handleDesc.Name, handleDesc.Request))
    resp = append(resp, contractDescToRpcContract(handleDesc.Name, handleDesc.Response))
  }
  typeDefs := aggr.Map(desc.TypeDefs, typeDefDescToRpcTypeDef)

  return &rpcTemplates{
    Handles: handles,
    Contracts: &rpcContracts{
      Requests:  reqs,
      Responses: resp,
      TypeDefs:  typeDefs,
    },
  }
}

func contractDescToRpcContract(name string, desc *contractDesc) *rpcContract {
  return &rpcContract{
    Name:   name,
    Fields: aggr.Map(desc.Fields, fieldDescToRpcStructField),
  }
}

func handleDescToRpcHandle(desc *handleDesc) *rpcHandle {
  return &rpcHandle{
    Name:  stringer.StringToUpperCamelCase(desc.Name),
    Route: desc.Route,
  }
}

func typeDefDescToRpcTypeDef(desc *typeDefsDesc) *rpcTypeDef {
  return &rpcTypeDef{
    Name:   stringer.StringToUpperCamelCase(desc.Name),
    Fields: aggr.Map(desc.Fields, fieldDescToRpcStructField),
  }
}

func fieldDescToRpcStructField(desc *fieldDesc) *rpcStructField {
  structField := &rpcStructField{}

  structField.Name = stringer.StringToUpperCamelCase(desc.Name)
  structField.Tag = stringer.StringToSnakeCase(desc.Name)

  switch {
  case aggr.MapLookup(scalarTypesMap, desc.Type) || typeScalarSliceRegex.MatchString(desc.Type):
    // Not convert
    structField.Type = desc.Type

  case !typeScalarSliceRegex.MatchString(desc.Type) && typeDefSliceRegex.MatchString(desc.Type):
    sliceType := strings.TrimPrefix(desc.Type, slicePrefix)
    // Convert to upper case
    sliceType = stringer.StringToUpperCamelCase(sliceType)
    // Convert to slice of pointers
    sliceType = strings.Join([]string{slicePrefix, ptrPrefix, sliceType}, dummySep)
    structField.Type = sliceType

  default:
    // Convert to upper case
    structType := stringer.StringToUpperCamelCase(desc.Type)
    // Convert to pointer
    structField.Type = strings.Join([]string{ptrPrefix, structType}, dummySep)
  }
  return structField
}

const (
  dummySep    = ""
  ptrPrefix   = "*"
  slicePrefix = "[]"
)
