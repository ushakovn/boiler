package ctxdetach

import (
  "context"
  "time"
)

type detached struct {
  ctx context.Context
}

func (d detached) Deadline() (deadline time.Time, ok bool) {
  return time.Time{}, false
}

func (d detached) Done() <-chan struct{} {
  return nil
}

func (d detached) Err() error {
  return nil
}

func (d detached) Value(key any) any {
  return d.ctx.Value(key)
}

// Do отключает контекст от родительского
// (таймаут или отмена контекста перестают действовать на возвращаемый контекст)
func Do(ctx context.Context) context.Context {
  return &detached{ctx: ctx}
}
