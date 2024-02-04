package storage

import (
  "fmt"

  validation "github.com/go-ozzo/ozzo-validation"
  "github.com/ushakovn/boiler/internal/pkg/filer"
)

const (
  pgColFilterOpGt      = "Gt"
  pgColFilterOpGtOrEq  = "GtOrEq"
  pgColFilterOpLt      = "Lt"
  pgColFilterOpLtOrEq  = "LtOrEq"
  pgColFilterOpIn      = "In"
  pgColFilterOpNotIn   = "NotIn"
  pgColFilterOpLike    = "Like"
  pgColFilterOpNotLike = "NotLike"
)

func (c *Config) Validate() error {
  if c.PgConfig != nil {
    if err := c.PgConfig.Validate(); err != nil {
      return fmt.Errorf("pg_config: %w", err)
    }
  }
  if c.PgDumpPath != "" {
    if !filer.IsExistedFile(c.PgDumpPath) {
      return fmt.Errorf("pg_dump_path: file not found")
    }
  }
  if c.PgTableConfig != nil {
    if err := c.PgTableConfig.Validate(); err != nil {
      return fmt.Errorf("pg_table_config: %w", err)
    }
  }
  if len(c.PgTypeConfig) > 0 {
    if err := validation.Validate(c.PgTypeConfig, validation.Each(validation.Required)); err != nil {
      return fmt.Errorf("pg_type_config: %w", err)
    }
  }
  return nil
}

func (c *PgConfig) Validate() error {
  rule := validation.Required

  return validation.ValidateStruct(c,
    validation.Field(&c.Host, rule),
    validation.Field(&c.Port, rule),
    validation.Field(&c.User, rule),
    validation.Field(&c.DBName, rule),
  )
}

func (c *PgTypeConfig) Validate() error {
  rule := validation.Required

  return validation.ValidateStruct(c,
    validation.Field(&c.GoType, rule),
    validation.Field(&c.GoZeroType, rule),
  )
}

func (c *PgTableConfig) Validate() error {
  return validation.ValidateStruct(c, validation.Field(&c.PgColumnFilter))
}

func (c *PgColFilter) Validate() error {
  ruleString := validation.Each(validation.In(pgColStringOperators...))
  ruleNumeric := validation.Each(validation.In(pgColNumericOperators...))

  if err := validation.ValidateStruct(c,
    validation.Field(&c.String, ruleString),
    validation.Field(&c.Numeric, ruleNumeric),
  ); err != nil {
    return err
  }
  ruleFilters := validation.Each(validation.In(pgColOperators...))

  for tableName, tableFilters := range c.Overrides {
    if tableName == "" {
      return fmt.Errorf("overrides: table name cannot be blank")
    }
    for columnName, columnFilters := range tableFilters {
      if columnName == "" {
        return fmt.Errorf("overrides: %s: column name cannot be blank", tableName)
      }
      if err := validation.Validate(columnFilters, ruleFilters); err != nil {
        return fmt.Errorf("overrides: %s.%s: filters: %w", tableName, columnName, err)
      }
    }
  }
  return nil
}

var pgColOperators = []any{
  pgColFilterOpIn,
  pgColFilterOpNotIn,
  pgColFilterOpLike,
  pgColFilterOpNotLike,
  pgColFilterOpGt,
  pgColFilterOpGtOrEq,
  pgColFilterOpLt,
  pgColFilterOpLtOrEq,
  pgColFilterOpIn,
  pgColFilterOpNotIn,
}

var pgColStringOperators = []any{
  pgColFilterOpIn,
  pgColFilterOpNotIn,
  pgColFilterOpLike,
  pgColFilterOpNotLike,
}

var pgColNumericOperators = []any{
  pgColFilterOpGt,
  pgColFilterOpGtOrEq,
  pgColFilterOpLt,
  pgColFilterOpLtOrEq,
  pgColFilterOpIn,
  pgColFilterOpNotIn,
}
