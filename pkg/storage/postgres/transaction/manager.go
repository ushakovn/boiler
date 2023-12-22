package transaction

import (
  "context"

  "github.com/jackc/pgx/v5"
)

type contextKey struct{}

type Manager struct{}

func NewManager() *Manager {
  return &Manager{}
}

func (m *Manager) FromContext(txCtx context.Context) pgx.Tx {
  tx, _ := txCtx.Value(contextKey{}).(pgx.Tx)
  return tx
}

func (m *Manager) ToContext(ctx context.Context, tx pgx.Tx) context.Context {
  return context.WithValue(ctx, contextKey{}, tx)
}
