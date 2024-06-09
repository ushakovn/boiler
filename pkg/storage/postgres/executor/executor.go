package executor

import (
  "context"
  "fmt"

  "github.com/georgysavva/scany/v2/pgxscan"
  "github.com/jackc/pgx/v5"
  "github.com/jackc/pgx/v5/pgconn"
  "github.com/jackc/pgx/v5/pgxpool"
  log "github.com/sirupsen/logrus"
  "github.com/ushakovn/boiler/pkg/ctxdetach"
  "github.com/ushakovn/boiler/pkg/retries"
  "github.com/ushakovn/boiler/pkg/stack"
  "github.com/ushakovn/boiler/pkg/storage/postgres/errors"
)

var txContextKey struct{}

type Executor struct {
  *pgxpool.Pool
  pgx.Tx
}

type Querier interface {
  Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
}

type Execer interface {
  Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
}

type Builder interface {
  ToSql() (query string, args []any, err error)
}

type TxStack interface {
  Push(tx pgx.Tx)
  Pop() pgx.Tx
  Peek() pgx.Tx
}

func New(ctx context.Context, dsn string) (*Executor, error) {
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

  return &Executor{
    Pool: pool,
  }, nil
}

func (e *Executor) BeginTx(ctx context.Context) (context.Context, error) {
  txStack := txStackFromContext(ctx)

  tx, err := e.Pool.Begin(ctx)
  if err != nil {
    return nil, err
  }
  txStack.Push(tx)

  return txStackToContext(ctx, txStack), nil
}

func (e *Executor) BeginTxWithOptions(ctx context.Context, options pgx.TxOptions) (context.Context, error) {
  txStack := txStackFromContext(ctx)

  tx, err := e.Pool.BeginTx(ctx, options)
  if err != nil {
    return nil, err
  }
  txStack.Push(tx)

  return txStackToContext(ctx, txStack), nil
}

func (e *Executor) CommitTx(ctx context.Context) error {
  tx := popTxFromContext(ctx)
  if tx == nil {
    return nil
  }
  return tx.Commit(ctx)
}

func (e *Executor) RollbackTx(ctx context.Context) {
  tx := popTxFromContext(ctx)
  if tx == nil {
    return
  }
  ctx = ctxdetach.Do(ctx)

  if err := tx.Rollback(ctx); err != nil {
    log.Error(err)
  }
}

func (e *Executor) Querier(ctx context.Context) Querier {
  if tx := peekTxFromContext(ctx); tx != nil {
    return tx
  }
  return e.Pool
}

func (e *Executor) Execer(ctx context.Context) Execer {
  if tx := peekTxFromContext(ctx); tx != nil {
    return tx
  }
  return e.Pool
}

func Select[T any](ctx context.Context, querier Querier, builder Builder) ([]T, error) {
  query, args, err := builder.ToSql()
  if err != nil {
    return nil, fmt.Errorf("builder.ToSql: %w", err)
  }
  var models []T

  if err = pgxscan.Select(ctx, querier, &models, query, args...); err != nil {
    return nil, fmt.Errorf("sqlscan.Select: %w", err)
  }
  return models, nil
}

func Get[T any](ctx context.Context, querier Querier, builder Builder) (T, error) {
  models, err := Select[T](ctx, querier, builder)
  if err != nil {
    return *new(T), err
  }
  if len(models) == 0 {
    return *new(T), errors.ErrModelNotFound
  }
  return models[0], nil
}

func Exec(ctx context.Context, execer Execer, builder Builder) error {
  query, args, err := builder.ToSql()
  if err != nil {
    return fmt.Errorf("builder.ToSql: %w", err)
  }
  if _, err = execer.Exec(ctx, query, args...); err != nil {
    return fmt.Errorf("execer.ExecContext: %w", err)
  }
  return nil
}

func txStackFromContext(ctx context.Context) TxStack {
  if txStack, ok := ctx.Value(txContextKey).(TxStack); ok {
    return txStack
  }
  return stack.New[pgx.Tx]()
}

func txStackToContext(ctx context.Context, txStack TxStack) context.Context {
  return context.WithValue(ctx,
    txContextKey,
    txStack,
  )
}

func peekTxFromContext(ctx context.Context) pgx.Tx {
  txStack := txStackFromContext(ctx)
  return txStack.Peek()
}

func popTxFromContext(ctx context.Context) pgx.Tx {
  txStack := txStackFromContext(ctx)
  return txStack.Pop()
}
