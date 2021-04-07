package backends

import(
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql" // required by sqlx to connect to the db
	"github.com/guitou4573/ipam/pkg/utils"
	"github.com/jmoiron/sqlx"
)

type SQLDb struct {
	db *sqlx.DB
	maxLifetime time.Duration
	maxOpenConns int
	maxIdleConns int
}

var (
		SQLDbErrorNumRows = sql.ErrNoRows
)

func (s *SQLDb) Init() {
	s.maxLifetime = time.Hour * 2
	s.maxOpenConns = 50
	s.maxIdleConns = 25
	tries, _ := strconv.Atoi(utils.EnvGet("MYSQL_RETRIES", "8"))
	sqlparams := "?parseTime=true"
	// sqlparams += "&tls=skip-verify"
	constring := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s%s",
		utils.EnvGet("MYSQL_USER", ""),
		utils.EnvGet("MYSQL_PASSWORD", ""),
		utils.EnvGet("MYSQL_HOST", "db"),
		string(utils.EnvGet("MYSQL_PORT", "3306")),
		utils.EnvGet("MYSQL_DB", "wire"),
		sqlparams)

	var err error
	s.db, err = sqlx.Connect("mysql", constring)
	for err != nil && tries > 0 {
		log.Println("[connect]", tries, (4 * time.Second), err)
		tries = tries-1
		time.Sleep(4 * time.Second)
		s.db, err = sqlx.Connect("mysql", constring)
	}
	if err != nil {
		log.Println(fmt.Errorf("[connect]: %s", err))
		return
	}
	log.Println("[connect] DB connection established")
	err = s.db.Ping()
	if err != nil {
		log.Println("[connect] Ping error")
		return
	}

	s.db.SetMaxOpenConns(s.maxOpenConns)
	s.db.SetMaxIdleConns(s.maxIdleConns)
	s.db.SetConnMaxLifetime(s.maxLifetime)
}

func (s *SQLDb) GetConnection() (*sqlx.DB){
	return s.db
}

func (s *SQLDb) Close() (error) {
	if s.db == nil {
		return fmt.Errorf("[DB Close] DB not initialised")
	}
	err := s.db.Close()
	return err
}
