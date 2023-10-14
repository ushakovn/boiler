package temp

import (
  "context"
  "fmt"

  sq "github.com/Masterminds/squirrel"
)

type productStorage struct {
  client Client
}

func NewProductStorage(client Client) ProductStorage {
  return &productStorage{
    client: client,
  }
}

func (s *productStorage) Products(ctx context.Context, input ProductsInput) ([]*Product, error) {
  query := newSelectBuilder().
    Columns(
      string(IdProductField),
      string(NameProductField),
    ).
    From(string(ProductTableName))

  if err := input.Pagination.validate(); err != nil {
    return nil, fmt.Errorf("pagination.Validate: %w", err)
  }
  offset, limit := input.Pagination.orDefault().toOffsetLimit()
  query = query.Offset(offset).Limit(limit)

  if input.Sort != nil {
    query = query.OrderBy(input.Sort.productSort())
  }
  filters := input.Filters

  if len(filters.IdsIn) > 0 {
    query = query.Where(sq.Eq{string(IdProductField): filters.IdsIn})
  }
  if len(filters.IdsNotIn) > 0 {
    query = query.Where(sq.NotEq{string(IdProductField): filters.IdsIn})
  }
  if filters.WhereInStock.Ptr() != nil {
    query = query.Where(sq.Eq{string(InStockProductField): filters.WhereInStock.Bool})
  }
  if filters.CreatedAtGte.Ptr() != nil {
    query = query.Where(sq.GtOrEq{string(CreatedAtProductField): filters.CreatedAtGte.Time})
  }
  if filters.CreatedAtGt.Ptr() != nil {
    query = query.Where(sq.Gt{string(CreatedAtProductField): filters.CreatedAtGt.Time})
  }
  return doQueryContext[*Product](ctx, s.client, query)
}

func (s *productStorage) Product(ctx context.Context, input ProductInput) (*Product, error) {
  query := newSelectBuilder().
    Columns(
      string(IdProductField),
      string(NameProductField),
    ).
    From(string(ProductTableName))

  filters := input.Filters

  if filters.IdEq.Ptr() != nil {
    query = query.Where(sq.Eq{string(IdProductField): filters.IdEq.String})
  }
  if filters.NameEq.Ptr() != nil {
    query = query.Where(sq.Eq{string(IdProductField): filters.NameEq.String})
  }
  if filters.WhereInStock.Ptr() != nil {
    query = query.Where(sq.Eq{string(InStockProductField): filters.WhereInStock.Bool})
  }

  models, err := doQueryContext[*Product](ctx, s.client, query)
  if err != nil {
    return nil, err
  }
  return models[0], nil
}

func (s *productStorage) PutProduct(ctx context.Context, input PutProductInput) (*Product, error) {
  query := newInsertBuilder().
    Into(string(ProductTableName))

  fields := map[string]any{
    string(IdProductField):      input.Id,
    string(NameProductField):    input.Name,
    string(InStockProductField): input.InStock,
    // Other fields
  }
  if input.Description.Ptr() != nil {
    fields[string(DescriptionProductField)] = input.Description.String
  }
  if input.CreatedAt.Ptr() != nil {
    fields[string(CreatedAtProductField)] = input.CreatedAt.Time
  }
  query = query.SetMap(fields)

  if err := doExecContext(ctx, s.client, query); err != nil {
    return nil, err
  }
  model := &Product{
    Id:          input.Id,
    Name:        input.Name,
    Description: input.Description,
    InStock:     input.InStock,
    Cost:        input.Cost,
    Warehouse:   input.Warehouse,
    CreatedAt:   input.CreatedAt,
    // Other fields
  }
  return model, nil
}

func (s *productStorage) UpdateProduct(ctx context.Context, input UpdateProductInput) (*Product, error) {
  query := newUpdateBuilder().
    Table(string(ProductTableName))

  fields := map[string]any{
    string(NameProductField):    input.Name,
    string(InStockProductField): input.InStock,
    // Other fields
  }
  if input.Description.Ptr() != nil {
    fields[string(DescriptionProductField)] = input.Description.String
  }
  if input.CreatedAt.Ptr() != nil {
    fields[string(CreatedAtProductField)] = input.CreatedAt.Time
  }
  query = query.SetMap(fields)

  if err := doExecContext(ctx, s.client, query); err != nil {
    return nil, err
  }
  model := &Product{
    Name:        input.Name,
    Description: input.Description,
    InStock:     input.InStock,
    Cost:        input.Cost,
    Warehouse:   input.Warehouse,
    CreatedAt:   input.CreatedAt,
    // Other fields
  }
  return model, nil
}

func (s *productStorage) DeleteProduct(ctx context.Context, input DeleteProductInput) error {
  query := newDeleteBuilder().
    From(string(ProductTableName)).
    Where(sq.Eq{string(IdProductField): input.Id})

  return doExecContext(ctx, s.client, query)
}
