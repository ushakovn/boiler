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
    return types.NewNilValue()
  }
  if value, ok := p.values[key]; ok {
    return value
  }
  return types.NewNilValue()
}

func (p *local) Watch(ctx context.Context, key string, action func(value types.Value)) {
  value := p.Get(ctx, key)
  action(value)
}
