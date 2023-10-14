package temp

import (
  "context"
  "fmt"

  "gopkg.in/guregu/null.v4/zero"
)

// Other storages for entities

type ProductStorage interface {
  Products(ctx context.Context, input ProductsInput) ([]*Product, error)
  Product(ctx context.Context, input ProductInput) (*Product, error)
  PutProduct(ctx context.Context, input PutProductInput) (*Product, error)
  UpdateProduct(ctx context.Context, input UpdateProductInput) (*Product, error)
  DeleteProduct(ctx context.Context, input DeleteProductInput) error
}

type ProductInput struct {
  Filters *ProductFilter
}

type ProductFilter struct {
  // Primary key without pointer
  IdEq         zero.String
  NameEq       zero.String
  WhereInStock zero.Bool
  // Other fields
}

type PutProductInput struct {
  // Without primary key it is serial or uuid
  // If sequential exists
  Id          string
  Name        string
  Description zero.String
  InStock     bool
  Cost        float64
  Warehouse   zero.String
  CreatedAt   zero.Time
}

type DeleteProductInput struct {
  // Only primary key
  Id string
}

type UpdateProductInput struct {
  // Without primary key
  Name        string
  Description zero.String
  InStock     bool
  Cost        float64
  Warehouse   zero.String
  CreatedAt   zero.Time
}

type ProductsInput struct {
  Filters    *ProductsFilters
  Sort       *ProductsSort
  Pagination *Pagination
}

type tableName string

const ProductTableName tableName = "product"

type productField string

const (
  IdProductField          productField = "id"
  NameProductField        productField = "name"
  DescriptionProductField productField = "description"
  InStockProductField     productField = "in_stock"
  CreatedAtProductField   productField = "created_at"
  // Other fields
)

type sortOrder string

const (
  SortOrderAsc  sortOrder = "ASC"
  SortOrderDesc sortOrder = "DESC"
  SortOrderRand sortOrder = "RAND()"
)

type Pagination struct {
  Page    uint64
  PerPage uint64
}

func (p *Pagination) orDefault() *Pagination {
  if p != nil {
    return p
  }
  return &Pagination{
    Page:    0,
    PerPage: 100,
  }
}

func (p *Pagination) validate() error {
  if p == nil {
    return nil
  }
  if p.Page < 0 {
    return fmt.Errorf("pagination.Page=%d must be non-negative", p.Page)
  }
  if p.PerPage < 0 {
    return fmt.Errorf("pagination.PerPage=%d must be positive", p.PerPage)
  }
  return nil
}

func (p *Pagination) toOffsetLimit() (offset uint64, limit uint64) {
  offset = p.Page * p.PerPage
  limit = p.PerPage
  return
}

type ProductsSort struct {
  Field productField
  Order sortOrder
}

func (p *ProductsSort) productSort() string {
  return fmt.Sprintf("%s %s", p.Field, p.Order)
}

type ProductsFilters struct {
  IdsIn            []string
  IdsNotIn         []string
  NameIn           []string
  NameNotIn        []string
  DescriptionIn    []string
  DescriptionNotIn []string
  WhereInStock     zero.Bool
  CostIn           []float64
  CostNotIn        []float64
  CostGt           zero.Float
  CostGte          zero.Float
  CostLt           zero.Float
  CostLte          zero.Float
  CostNotEq        []float64
  WarehouseIn      []string
  WarehouseNotIn   []string
  CreatedAtIn      []zero.Time
  CreatedAtNotIn   []zero.Time
  CreatedAtGt      zero.Time
  CreatedAtGte     zero.Time
  CreatedAtLt      zero.Time
  CreatedAtLte     zero.Time
  // Other fields
}
