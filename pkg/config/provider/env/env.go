package env

import (
  "context"
  "os"

  "github.com/ushakovn/boiler/internal/pkg/stringer"
  "github.com/ushakovn/boiler/pkg/config/provider"
  "github.com/ushakovn/boiler/pkg/config/types"
)

type env struct{}

func New() provider.Values {
  return new(env)
}

func (p *env) Get(_ context.Context, key string) types.Value {
  if value, ok := os.LookupEnv(key); ok {
    return types.NewValue(value)
  }
  key = stringer.StringToCapitalizeCase(key)

  if value, ok := os.LookupEnv(key); ok {
    return types.NewValue(value)
  }
  return types.NewNilValue()
}

func (p *env) Watch(ctx context.Context, key string, action func(value types.Value)) {
  value := p.Get(ctx, key)
  action(value)
}
