package client

import (
  "context"
  "fmt"

  "github.com/jackc/pgx/v5"
  log "github.com/sirupsen/logrus"
)

type Txer interface {
  WithTransaction(ctx context.Context, txFn func(txCtx context.Context) error) error
  WithTransactionOpt(ctx context.Context, txOptions pgx.TxOptions, txFn func(txCtx context.Context) error) error
}

func (c *client) WithTransaction(ctx context.Context, txFn func(txCtx context.Context) error) error {
  return c.WithTransactionOpt(ctx, pgx.TxOptions{}, txFn)
}

func (c *client) WithTransactionOpt(ctx context.Context, txOptions pgx.TxOptions, txFn func(txCtx context.Context) error) error {
  defer func() {
    if rec := recover(); rec != nil {
      log.Errorf("client.WithTransaction: panic recovered: %v", rec)
    }
  }()

  tx, err := c.BeginTx(ctx, txOptions)
  if err != nil {
    return fmt.Errorf("c.BeginTx: %w", err)
  }
  txCtx := c.ToContext(ctx, tx)

  if err = txFn(txCtx); err != nil {
    if txErr := tx.Rollback(ctx); txErr != nil {
      log.Errorf("client.WithTransaction: tx.Rollback: %v", err)
    }
    return fmt.Errorf("txFn: %w", err)
  }

  if txErr := tx.Commit(ctx); txErr != nil {
    return fmt.Errorf("tx.Commit: %w", err)
  }
  return nil
}
