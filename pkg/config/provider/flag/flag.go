package flag

import (
  "context"
  "os"
  "strings"
  "sync"

  "github.com/ushakovn/boiler/internal/pkg/stringer"
  "github.com/ushakovn/boiler/pkg/config/provider"
  "github.com/ushakovn/boiler/pkg/config/types"
)

type flag struct {
  mu     sync.Mutex
  values map[string]string
}

func New() provider.Values {
  return &flag{
    values: collectValues(),
  }
}

func (p *flag) Get(_ context.Context, key string) types.Value {
  if value, ok := p.values[key]; ok {
    return types.NewValue(value)
  }
  key = stringer.StringToKebabCase(key)

  if value, ok := p.values[key]; ok {
    return types.NewValue(value)
  }
  return types.NewNilValue()
}

func (p *flag) Watch(ctx context.Context, key string, action func(value types.Value)) {
  value := p.Get(ctx, key)
  action(value)
}

func collectValues() map[string]string {
  values := make(map[string]string)

  for _, arg := range os.Args {
    if !strings.HasPrefix(arg, "--") {
      continue
    }
    parts := strings.SplitN(arg, "=", 2)

    if len(parts) != 2 {
      continue
    }
    key := strings.Trim(parts[0], " --")
    value := strings.TrimSpace(parts[1])

    values[key] = value
  }

  return values
}
