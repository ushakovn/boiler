package provider

import (
  "context"

  "github.com/ushakovn/boiler/pkg/config/types"
)

type Values interface {
  Get(ctx context.Context, key string) types.Value
  Watch(ctx context.Context, key string, action func(value types.Value))
}
