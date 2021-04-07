package models

import(
	"time"
)

type Token struct {
	Token string `json:"token" db:"token"`
	Created time.Time `json:"created_at" db:"created_at"`
}

func (m *Token) GetToken() string {
	return m.Token
}
