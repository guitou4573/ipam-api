package backends

import(
  "database/sql"
  "fmt"
  "log"
  "time"

  "github.com/mattn/go-sqlite3" // required by sqlx to connect to the db
  "github.com/jmoiron/sqlx"
)

type SQLite struct {
	db *sqlx.DB
	maxLifetime time.Duration
	maxOpenConns int
	maxIdleConns int
}

func (s *SQLite) Init() {
	s.maxLifetime = time.Hour * 2
	s.maxOpenConns = 50
	s.maxIdleConns = 25
  var err error
  // sqlite3.EnableExtensionOnOff = 1

  sql.Register("sqlite3_ext",
		&sqlite3.SQLiteDriver{
			Extensions: []string{
				"static/inet",
			},
		},
  )
	s.db, err = sqlx.Connect("sqlite3_ext", ":memory:")
	if err != nil {
		log.Println(fmt.Errorf("[connect]: %s", err))
		return
	}

	err = s.db.Ping()
	if err != nil {
		log.Println("[connect] Ping error")
		return
	}

	s.db.SetMaxOpenConns(s.maxOpenConns)
	s.db.SetMaxIdleConns(s.maxIdleConns)
	s.db.SetConnMaxLifetime(s.maxLifetime)
}

func (s *SQLite) GetConnection() (*sqlx.DB){
	return s.db
}

func (s *SQLite) Close() (error) {
	if s.db == nil {
		return fmt.Errorf("[DB Close] DB not initialised")
	}
	err := s.db.Close()
	return err
}
