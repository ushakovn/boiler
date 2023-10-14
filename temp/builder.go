package temp

import sq "github.com/Masterminds/squirrel"

func newSelectBuilder() sq.SelectBuilder {
  return sq.SelectBuilder{}.PlaceholderFormat(sq.Dollar)
}

func newInsertBuilder() sq.InsertBuilder {
  return sq.InsertBuilder{}.PlaceholderFormat(sq.Dollar)
}

func newUpdateBuilder() sq.UpdateBuilder {
  return sq.UpdateBuilder{}.PlaceholderFormat(sq.Dollar)
}

func newDeleteBuilder() sq.DeleteBuilder {
  return sq.DeleteBuilder{}.PlaceholderFormat(sq.Dollar)
}
