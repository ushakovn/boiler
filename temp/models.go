package temp

import (
  "gopkg.in/guregu/null.v4/zero"
)

type Product struct {
  Id          string
  Name        string
  Description zero.String
  InStock     bool
  Cost        float64
  Warehouse   zero.String
  CreatedAt   zero.Time
  // Other fields
}
