package client

import (
  "context"
  "fmt"

  "github.com/jackc/pgx/v5"
  "github.com/jackc/pgx/v5/pgxpool"
  _ "github.com/jackc/pgx/v5/stdlib"
  "github.com/ushakovn/boiler/pkg/retries"

  _ "github.com/georgysavva/scany/v2/pgxscan"
)

type Client interface {
  Txer
  Executor
  Executor(ctx context.Context) Executor
}

type client struct {
  *pgxpool.Pool
  TxManager
}

type TxManager interface {
  FromContext(txCtx context.Context) pgx.Tx
  ToContext(ctx context.Context, tx pgx.Tx) context.Context
}

func NewClient(ctx context.Context, dsn string, manager TxManager) (Client, error) {
  pool, err := pgxpool.New(ctx, dsn)
  if err != nil {
    return nil, fmt.Errorf("pgxpool.New: %w", err)
  }

  err = retries.DoWithRetries(ctx,
    retries.Options{},

    func(ctx context.Context) error {
      if err = pool.Ping(ctx); err != nil {
        return fmt.Errorf("pool.Ping: %w", err)
      }
      return nil
    })

  if err != nil {
    return nil, fmt.Errorf("retries.DoWithRetries: %w", err)
  }

  return &client{
    Pool:      pool,
    TxManager: manager,
  }, nil
}
