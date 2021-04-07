package models

import(
	"time"
)

type VPC struct {
  Id string `json:"id" db:"id"`
	IdSection string `json:"idsection" db:"idsection"`
	Created time.Time `json:"created_at" db:"created_at"`
	Description string `json:"description" db:"description"`
}

func (m *VPC) GetId() string {
	return m.Id
}

func (m *VPC) GetIdSection() string {
  return m.IdSection
}

func (m *VPC) GetDescription() string {
  return m.Description
}
