package models

import(
	"time"
)

type Env struct {
  Id int64 `json:"id" db:"id"`
	Created time.Time `json:"created_at" db:"created_at"`
	Description string `json:"description" db:"description"`
}

func (m *Env) GetId() int64 {
	return m.Id
}

func (m *Env) GetDescription() string {
  return m.Description
}
