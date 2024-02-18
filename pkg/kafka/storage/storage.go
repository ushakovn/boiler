package storage

import (
  "context"
  "fmt"
  "time"

  sq "github.com/Masterminds/squirrel"
  "github.com/ushakovn/boiler/pkg/kafka/models"
  pg "github.com/ushakovn/boiler/pkg/storage/postgres/executor"
)

type Outbox struct {
  executor  pg.Executor
  lockTTL   time.Duration
  batchSize uint32
}

func NewStorage(executor pg.Executor, lockTTL time.Duration, batchSize uint32) *Outbox {
  return &Outbox{
    executor:  executor,
    lockTTL:   lockTTL,
    batchSize: batchSize,
  }
}

func (o *Outbox) BatchRecords(ctx context.Context, tableName string) ([]*models.Record, error) {
  query := `update %s
    set locked_until = now() + interval '%d' millisecond
    where id in (select id
        from %[1]s
        where locked_until is null or locked_until < now() 
        limit %d)
    returning *`

  query = fmt.Sprintf(query, tableName, o.lockTTL.Milliseconds(), o.batchSize)

  return pg.SelectCtx[*models.Record](ctx, o.executor, sq.Expr(query))
}

func (o *Outbox) DeleteRecord(ctx context.Context, tableName, recordID string) error {
  query := fmt.Sprintf(`delete from %s where id = '%s'`, tableName, recordID)

  return pg.ExecCtx(ctx, o.executor, sq.Expr(query))
}
