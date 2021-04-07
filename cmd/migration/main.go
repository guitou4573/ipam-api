package migration

import (
	"errors"
	"io/ioutil"
	"log"
	"strings"

	"github.com/guitou4573/ipam/pkg/backends"

)

type Migration struct {
	db backends.Backend
}

func NewMigration(db backends.Backend) (*Migration){
	return &Migration{
		db: db,
	}
}

// Applies queries in one transaction
func (m *Migration) Migrate(stmts []string) error{
	tx, err := m.db.GetConnection().Begin()
	for _, stmt := range stmts {
		trimStmt := strings.TrimSpace(stmt)
		if trimStmt != "" {
      // log.Println(trimStmt)
			if _, err := tx.Exec(trimStmt); err != nil {
				return err
			}
		}
	}
	err = tx.Commit()
	return err
}

// Creates an entry in the migration table and then applies queries in one transaction
func (m *Migration) MigrateFile(file string) error{
	mig, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	migstring := string(mig[:])
	stmts := strings.Split(migstring ,";")
	_, err = m.db.GetConnection().Exec("INSERT INTO `migration` (name) VALUES (?)", file)
	if err != nil {
		return err
	}
	err = m.Migrate(stmts)
	return err
}

// Creates migration table
func (m *Migration) InitMigration() error {
	switch m.db.GetConnection().DriverName() {
	case "mysql":
		return m.InitMigrationSQL()
	case "sqlite3","sqlite3_ext":
		return m.InitMigrationSQLLite()
	default:
		return errors.New("unknown db driver")
	}
}


func (m *Migration) InitMigrationSQL() error {
	var migTable = `CREATE TABLE IF NOT EXISTS `+"`migration`"+` (
  `+"`idmigration`"+` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `+"`name`"+` VARCHAR(255) NOT NULL UNIQUE,
  `+"`created_at`"+` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY (`+"`idmigration`"+`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4;`
	return m.Migrate([]string{migTable})
}

func (m *Migration) InitMigrationSQLLite() error {
	var migTable = `CREATE TABLE IF NOT EXISTS `+"`migration`"+` (
  `+"`idmigration`"+` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  `+"`name`"+` VARCHAR(255) NOT NULL UNIQUE,
  `+"`created_at`"+` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP);`
	return m.Migrate([]string{migTable})
}

// Makes sure the database is ready
func (m *Migration) Execute() {
	if m.db == nil {
	  log.Fatal("No backend provided")
	}
	// Create migration table
	m.InitMigration()

	// Create app tables
  var err error
  if (m.db.GetConnection().DriverName() == "mysql") {
    err = m.MigrateFile("./static/migrations/1.up.sql")
  } else {
    err = m.MigrateFile("./static/migrations/1.up.sqlite.sql")
  }

	if err != nil {
		log.Print(err)
	}

}
