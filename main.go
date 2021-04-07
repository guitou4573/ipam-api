package main

import(
  "github.com/guitou4573/ipam/cmd/apiservice"
  "github.com/guitou4573/ipam/pkg/backends"
	"github.com/guitou4573/ipam/pkg/utils"
	"github.com/guitou4573/ipam/cmd/migration"
	"github.com/guitou4573/ipam/pkg/stores"
)

var (
	db backends.Backend
	subnetStore stores.SubnetStore
)

func main() {
  if(utils.EnvGet("DB_TYPE", "mysql") == "mysql") {
    db = &backends.SQLDb{}
  } else {
    db = &backends.SQLite{}
  }

	db.Init()
	defer db.Close()
  mig := migration.NewMigration(db)
  mig.Execute()

  subnetStore, _ := stores.NewSubnetStore(db)

  apiservice.Execute(apiservice.ApiServiceOptions{
    SubnetStore: subnetStore,
  })
}
