package models

import "time"

type ActionType string

const (
  CreateActionType ActionType = "create"
  UpdateActionType ActionType = "update"
  DeleteActionType ActionType = "delete"
)

type Record struct {
  ID         string        `db:"id"`
  ActionType ActionType    `db:"action_type"`
  JSONOut    []byte        `db:"db_out"`
  CreatedAt  time.Duration `db:"created_at"`
}

func (t ActionType) String() string {
  return string(t)
}
