package types

import (
  "time"

  "github.com/spf13/cast"
)

type Value interface {
  IsNil() bool
  Int() int
  Int64() int64
  Int32() int32
  Uint32() uint32
  Uint64() uint64
  Float32() float32
  Float64() float64
  String() string
  Bool() bool
  Time() time.Time
  Duration() time.Duration
}

type configValue struct {
  value any
}

func NewValue(value any) Value {
  return &configValue{value: value}
}

func NewNilValue() Value {
  return configValue{}
}

func (c configValue) IsNil() bool {
  return c.value == nil
}

func (c configValue) Int() int {
  return cast.ToInt(c.value)
}

func (c configValue) Int64() int64 {
  return cast.ToInt64(c.value)
}

func (c configValue) Int32() int32 {
  return cast.ToInt32(c.value)
}

func (c configValue) Uint32() uint32 {
  return cast.ToUint32(c.value)
}

func (c configValue) Uint64() uint64 {
  return cast.ToUint64(c.value)
}

func (c configValue) Float32() float32 {
  return cast.ToFloat32(c.value)
}

func (c configValue) Float64() float64 {
  return cast.ToFloat64(c.value)
}

func (c configValue) String() string {
  return cast.ToString(c.value)
}

func (c configValue) Bool() bool {
  return cast.ToBool(c.value)
}

func (c configValue) Time() time.Time {
  return cast.ToTime(c.value)
}

func (c configValue) Duration() time.Duration {
  return cast.ToDuration(c.value)
}
