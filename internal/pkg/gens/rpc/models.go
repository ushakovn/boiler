package rpc

import (
	"fmt"
	"regexp"

	"github.com/ushakovn/boiler/pkg/utils"
)

type rootDesc struct {
	Handlers []*handlerDesc `json:"handlers"`
}

type handlerDesc struct {
	Name     string        `json:"name"`
	Request  *requestDesc  `json:"request"`
	Response *responseDesc `json:"response"`
}

type requestDesc struct {
	Fields []*fieldDesc `json:"fields"`
}

type responseDesc struct {
	Fields []*fieldDesc `json:"fields"`
}

type fieldDesc struct {
	Name          string             `json:"name"`
	ScalarType    string             `json:"scalar_type"`
	CompositeType *compositeTypeDesc `json:"composite_type"`
}

type compositeTypeDesc struct {
	Desc string   `json:"desc"`
	Def  *defDesc `json:"def"`
}

type defDesc struct {
	Name   string       `json:"name"`
	Fields []*fieldDesc `json:"fields"`
}

func (d *rootDesc) Validate() error {
	if len(d.Handlers) == 0 {
		return fmt.Errorf("rpc handler not specfied")
	}
	for _, handler := range d.Handlers {
		if err := handler.Validate(); err != nil {
			return fmt.Errorf("handler error: %v", err)
		}
	}
	return nil
}

func (d *handlerDesc) Validate() error {
	if d.Name == "" {
		return fmt.Errorf("name not specified")
	}
	if d.Request == nil {
		return fmt.Errorf("request not specified")
	}
	if err := d.Request.Validate(); err != nil {
		return fmt.Errorf("request error: %v", err)
	}
	if err := d.Response.Validate(); err != nil {
		return fmt.Errorf("response error: %v", err)
	}
	return nil
}

func (d *responseDesc) Validate() error {
	for _, field := range d.Fields {
		if err := field.Validate(); err != nil {
			return fmt.Errorf("field: %s: error: %v", field.Name, err)
		}
	}
	return nil
}

func (d *requestDesc) Validate() error {
	for _, field := range d.Fields {
		if err := field.Validate(); err != nil {
			return fmt.Errorf("field: %s: error: %v", field.Name, err)
		}
	}
	return nil
}

func (d *fieldDesc) Validate() error {
	if d.ScalarType == "" && d.CompositeType == nil {
		return fmt.Errorf("scalar type or composite type not specfied")
	}
	if d.ScalarType != "" && d.CompositeType != nil {
		return fmt.Errorf("must has one scalar or composite type only")
	}
	if d.ScalarType != "" {
		if _, ok := scalarTypesMap[d.ScalarType]; !ok {
			return fmt.Errorf("has unsupported scalar type name: %s", d.ScalarType)
		}
	}
	if d.CompositeType != nil {
		if err := d.CompositeType.Validate(); err != nil {
			return fmt.Errorf("composite type: %v", err)
		}
	}
	return nil
}

func (d *compositeTypeDesc) Validate() error {
	var (
		err error
		ok  bool
	)
	for _, rule := range compTypeValidations {
		if ok = rule.match(d.Desc); !ok {
			continue
		}
		if err = rule.validate(d.Def); err != nil {
			return fmt.Errorf("composite type definition error: %s", err)
		}
	}
	if !ok {
		return fmt.Errorf("has unsupported composite type description: %s", d.Desc)
	}
	return nil
}

type compTypeValidation struct {
	validate func(desc *defDesc) (err error)
	match    func(typ string) (ok bool)
}

var compTypeValidations = []*compTypeValidation{
	{
		validate: func(desc *defDesc) error {
			return desc.Validate()
		},
		match: compTypeStructMatch,
	},
	{
		validate: func(desc *defDesc) error {
			if desc != nil {
				return fmt.Errorf("definition must be null")
			}
			return nil
		},
		match: compTypeSlicesMatch,
	},
}

func (d *defDesc) Validate() error {
	if d.Name == "" {
		return fmt.Errorf("name not specified")
	}
	if utils.IsCamelCase(d.Name) || utils.IsSnakeCase(d.Name) {
		return fmt.Errorf("invalid name")
	}
	for _, field := range d.Fields {
		if err := field.Validate(); err != nil {
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
	"byte":    {},
	"string":  {},
}

var (
	compTypeStructMatch = regexp.MustCompile(`struct\{\}`).MatchString
	compTypeSlicesMatch = regexp.MustCompile(`^\[\]((bool|string|byte|int)|(float|int)(32|64))$`).MatchString
)

type rpcTemplates struct {
	Handles   []*rpcHandle
	Contracts *rpcContracts
}

type rpcHandle struct {
	Name string
}

type rpcContracts struct {
	Requests       []*rpcRequest
	Responses      []*rpcResponse
	CompositeTypes []*rpcCompositeType
}

type rpcRequest struct {
	Fields []*rpcStructField
}

type rpcResponse struct {
	Fields []*rpcStructField
}

type rpcCompositeType struct {
	Name   string
	Fields []*rpcStructField
}

type rpcStructField struct {
	Name string
	Type string
	Tag  string
}

func rootDescToRpcTemplates(desc *rootDesc) (*rpcTemplates, error) {
	var (
		handles []*rpcHandle

		// TODO: complete this function
		_ []*rpcRequest
		_ []*rpcResponse
		_ []*rpcCompositeType
	)
	for _, handlerDesc := range desc.Handlers {
		// collect handles names
		handles = append(handles, &rpcHandle{
			Name: handlerDesc.Name,
		})
		// init request
		_ = &rpcRequest{}

		// collect request fields
		for _, fieldDesc := range handlerDesc.Request.Fields {
			reqField := &rpcStructField{}
			// set name and json tag for field
			reqField.Name = fieldDesc.Name
			reqField.Tag = utils.ToStructTag(fieldDesc.Name)

			if typ := fieldDesc.ScalarType; typ != "" {
				reqField.Type = typ
			}
			if typ := fieldDesc.CompositeType; typ != nil {
				_ = &rpcCompositeType{}

				// if composite type match slice
				if compTypeSlicesMatch(typ.Desc) {
					reqField.Type = typ.Desc
				}
				// if composite type match struct
				if compTypeStructMatch(typ.Desc) {
					reqField.Type = typ.Def.Name

				}
			}
		}
	}
	// TODO: complete this function
	return nil, nil
}

func fieldDescToRpc() {
	// TODO: complete this function
}

func compTypeDescToRpcCompType(desc *compositeTypeDesc) []*rpcCompositeType {
	// collect rpc composite types
	var _ []*rpcCompositeType

	rpcCompTyp := &rpcCompositeType{}
	rpcCompTyp.Name = desc.Def.Name

	for _, fieldDesc := range desc.Def.Fields {
		rpcStructField := &rpcStructField{}

		if typ := fieldDesc.ScalarType; typ != "" {
			rpcStructField.Type = typ
		}
		if typ := fieldDesc.CompositeType; typ != nil {
			// if composite type match slice
			if compTypeSlicesMatch(typ.Desc) {
				rpcStructField.Type = typ.Desc
			}
			// if composite type match struct
			if compTypeStructMatch(typ.Desc) {
				rpcStructField.Type = typ.Def.Name
				rpcStructField.Tag = utils.ToStructTag(typ.Def.Name)

				// TODO: complete this function
				for _, fieldDesc = range typ.Def.Fields {
					return compTypeDescToRpcCompType(fieldDesc.CompositeType)
				}
			}
		}
	}
	return nil // TODO: complete this function
}
