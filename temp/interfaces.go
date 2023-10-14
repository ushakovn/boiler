package temp

import (
  "context"
  "database/sql"
  "errors"
  "fmt"

  "github.com/georgysavva/scany/v2/sqlscan"
)

var ErrZeroRowsRetrieved = errors.New("zero rows retrieved")

type Client interface {
  Pinger
  Execer
  Querier
  Txer
}

type Querier interface {
  QueryContext(context.Context, string, ...any) (*sql.Rows, error)
  QueryRowContext(context.Context, string, ...any) *sql.Row
}

type Execer interface {
  ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type Txer interface {
  Begin() (*sql.Tx, error)
  BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type Pinger interface {
  PingContext(ctx context.Context) error
}

type Builder interface {
  ToSql() (statement string, args []any, err error)
}

func doQueryContext[Model any](ctx context.Context, querier Querier, builder Builder) ([]Model, error) {
  statement, args, err := builder.ToSql()
  if err != nil {
    return nil, fmt.Errorf("builder.ToSql: %w", err)
  }
  var models []Model
  if err = sqlscan.Select(ctx, querier, &models, statement, args); err != nil {
    return nil, fmt.Errorf("sqlscan.Select: %w", err)
  }
  if len(models) == 0 {
    return nil, ErrZeroRowsRetrieved
  }
  return models, nil
}

func doExecContext(ctx context.Context, execer Execer, builder Builder) error {
  statement, args, err := builder.ToSql()
  if err != nil {
    return fmt.Errorf("builder.ToSql: %w", err)
  }
  if _, err = execer.ExecContext(ctx, statement, args); err != nil {
    return fmt.Errorf("execer.ExecContext: %w", err)
  }
  return nil
}