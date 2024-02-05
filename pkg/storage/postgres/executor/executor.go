package executor

import (
  "context"
  "fmt"

  "github.com/georgysavva/scany/v2/pgxscan"
  "github.com/jackc/pgx/v5"
  "github.com/jackc/pgx/v5/pgconn"
  "github.com/jackc/pgx/v5/pgxpool"
  "github.com/ushakovn/boiler/pkg/retries"
  "github.com/ushakovn/boiler/pkg/storage/postgres/errors"
)

type Executor interface {
  Querier
  Execer
  Txer
}

type Querier interface {
  Query(context.Context, string, ...any) (pgx.Rows, error)
  QueryRow(context.Context, string, ...any) pgx.Row
}

type Execer interface {
  Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
}

type Txer interface {
  Begin(ctx context.Context) (pgx.Tx, error)
}

type Builder interface {
  ToSql() (statement string, args []any, err error)
}

type executor struct {
  *pgxpool.Pool
}

func NewExecutor(ctx context.Context, dsn string) (Executor, error) {
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

  return &executor{
    Pool: pool,
  }, nil
}

func SelectCtx[T any](ctx context.Context, querier Querier, builder Builder) ([]T, error) {
  statement, args, err := builder.ToSql()
  if err != nil {
    return nil, fmt.Errorf("builder.ToSql: %w", err)
  }
  var models []T

  if err = pgxscan.Select(ctx, querier, &models, statement, args...); err != nil {
    return nil, fmt.Errorf("sqlscan.Select: %w", err)
  }
  return models, nil
}

func GetCtx[T any](ctx context.Context, querier Querier, builder Builder) (T, error) {
  models, err := SelectCtx[T](ctx, querier, builder)
  if err != nil {
    return *new(T), err
  }
  if len(models) == 0 {
    return *new(T), errors.ErrModelNotFound
  }
  return models[0], nil
}

func ExecCtx(ctx context.Context, execer Execer, builder Builder) error {
  statement, args, err := builder.ToSql()
  if err != nil {
    return fmt.Errorf("builder.ToSql: %w", err)
  }
  if _, err = execer.Exec(ctx, statement, args...); err != nil {
    return fmt.Errorf("execer.ExecContext: %w", err)
  }
  return nil
}
