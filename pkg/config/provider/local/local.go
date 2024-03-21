package local

import (
  "context"

  "github.com/ushakovn/boiler/pkg/config/provider"
  "github.com/ushakovn/boiler/pkg/config/types"
)

type local struct {
  values map[string]types.Value
}

func New(values map[string]types.Value) provider.Values {
  return &local{values: values}
}

func (p *local) Get(_ context.Context, key string) types.Value {
  if p.values == nil {
    return types.NewValue(nil)
  }
  return p.values[key]
}

func (p *local) Watch(_ context.Context, key string, action func(value types.Value)) {
  value := types.NewValue(nil)

  if p.values != nil {
    value = p.values[key]
  }
  action(value)
}
