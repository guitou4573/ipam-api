package backends

import(
  "github.com/jmoiron/sqlx"
)

type Backend interface {
  Init()
  Close() (error)
  GetConnection() *sqlx.DB
}
