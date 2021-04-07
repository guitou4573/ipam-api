package models

import(
  "net"
	"time"
)

type Subnet struct {
  NetAddr string `json:"netaddr" db:"netaddr"`
  Mask int64 `json:"mask" db:"mask"`
  IdVPC string `json:"idvpc" db:"idvpc"`
	Created time.Time `json:"created_at" db:"created_at"`
	Description string `json:"description" db:"description"`
}

func (m *Subnet) GetNetAddr() string {
  return m.NetAddr
}

func (m *Subnet) GetNetMask() int64 {
	return m.Mask
}

func (m *Subnet) GetIdVPC() string {
  return m.IdVPC
}

func (m *Subnet) GetDescription() string {
  return m.Description
}

func (m *Subnet) GetFirstLastAddress() []string {
  addr := net.ParseIP(m.NetAddr)
    if addr != nil {
      return []string{}
    }
    return []string{}
}
