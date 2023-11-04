package gen

import "context"

type Initor interface {
  Init(ctx context.Context) error
}
