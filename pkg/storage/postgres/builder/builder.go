package builder

import (
  "fmt"

  sq "github.com/Masterminds/squirrel"
  "github.com/ushakovn/boiler/pkg/storage/postgres/quote"
)

func NewSelectBuilder() sq.SelectBuilder {
  return sq.SelectBuilder{}.PlaceholderFormat(sq.Dollar)
}

func NewInsertBuilder() sq.InsertBuilder {
  return sq.InsertBuilder{}.PlaceholderFormat(sq.Dollar)
}

func NewUpdateBuilder() sq.UpdateBuilder {
  return sq.UpdateBuilder{}.PlaceholderFormat(sq.Dollar)
}

func NewDeleteBuilder() sq.DeleteBuilder {
  return sq.DeleteBuilder{}.PlaceholderFormat(sq.Dollar)
}

func NewBuildedExpr(sql string, args ...any) sq.Sqlizer {
  quoted := make([]any, 0, len(args))

  for _, arg := range args {
    quoted = append(quoted, quote.String(arg))
  }
  sql = fmt.Sprintf(sql, quoted...)

  return sq.Expr(sql)
}
