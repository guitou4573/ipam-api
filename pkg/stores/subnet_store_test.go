package stores

import (
  "testing"

	"github.com/guitou4573/ipam/pkg/backends"
	"github.com/guitou4573/ipam/cmd/migration"
  "github.com/guitou4573/ipam/pkg/models"
)

// t.Log
// t.Error
// t.Fail
// t.Errorf

func TestSubnetNewSubnetStore(t *testing.T) {
  db := &backends.SQLite{}
  db.Init()
  defer db.Close()

  subnetStore, err := NewSubnetStore(db)
  if subnetStore == nil {
    t.Errorf("subnet store nil: %s", err)
  }
}

func TestSubnetInitTable(t *testing.T) {
  db := &backends.SQLite{}
  db.Init()
  defer db.Close()
  subnetStore, _ := NewSubnetStore(db)
  res := subnetStore.initTable()
  if res != nil {
    t.Error()
  }
}

func TestSubnetGetByVPCid(t *testing.T) {
  db := &backends.SQLite{}
  db.Init()
  defer db.Close()
  mig := migration.NewMigration(db)
  mig.InitMigration()
  mig.MigrateFile("../../migrations/1.up.sqlite.sql")
  db.GetConnection().Exec(`INSERT INTO subnet (subnetname, created_at) VALUES ("subnet@domain",CURRENT_TIMESTAMP)`)
  subnetStore, _ := NewSubnetStore(db)
  subnet, err := subnetStore.GetByVPCId("vpc-123456")
  if err != nil || subnet == nil || subnet.IdVPC != "vpc-123456" {
      t.Error()
  }
}

func TestSubnetCreate(t *testing.T) {
  db := &backends.SQLite{}
  db.Init()
  defer db.Close()
  mig := migration.NewMigration(db)
  mig.InitMigration()
  mig.MigrateFile("../../migrations/1.up.sqlite.sql")
  subnetStore, _ := NewSubnetStore(db)
  subnetModel :=  models.Subnet{
  	IdVPC: "vpc-123456",
  	NetAddr: "10.0.0.0",
  	Mask: 27,
  }
  id, err := subnetStore.Create(&subnetModel)
  if err != nil || id == "" {
    t.Error(err)
  }
  var item models.Subnet
	err = db.GetConnection().QueryRowx(`SELECT * FROM subnet WHERE subnetname = "subnet@domain"`).StructScan(&item)
  if err != nil || item.NetAddr == ""{
    t.Error(err)
  }
}
