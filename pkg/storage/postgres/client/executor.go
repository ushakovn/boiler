package client

import (
  "context"
  "fmt"

  "github.com/georgysavva/scany/v2/pgxscan"
  "github.com/jackc/pgx/v5"
  "github.com/jackc/pgx/v5/pgconn"
  "github.com/ushakovn/boiler/pkg/storage/postgres/errors"
)

func (c *client) Executor(ctx context.Context) Executor {
  if tx := c.FromContext(ctx); tx != nil {
    return tx
  }
  return c.Pool
}

type Executor interface {
  Execer
  Querier
}

type Querier interface {
  Query(context.Context, string, ...any) (pgx.Rows, error)
  QueryRow(context.Context, string, ...any) pgx.Row
}

type Execer interface {
  Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
}

type Builder interface {
  ToSql() (statement string, args []any, err error)
}

func DoQueryContext[Model any](ctx context.Context, querier Querier, builder Builder) ([]Model, error) {
  statement, args, err := builder.ToSql()
  if err != nil {
    return nil, fmt.Errorf("builder.ToSql: %w", err)
  }
  var models []Model

  if err = pgxscan.Select(ctx, querier, &models, statement, args...); err != nil {
    return nil, fmt.Errorf("sqlscan.Select: %w", err)
  }
  if len(models) == 0 {
    return nil, errors.ErrModelNotFound
  }
  return models, nil
}

func DoExecContext(ctx context.Context, execer Execer, builder Builder) error {
  statement, args, err := builder.ToSql()
  if err != nil {
    return fmt.Errorf("builder.ToSql: %w", err)
  }
  if _, err = execer.Exec(ctx, statement, args...); err != nil {
    return fmt.Errorf("execer.ExecContext: %w", err)
  }
  return nil
}
