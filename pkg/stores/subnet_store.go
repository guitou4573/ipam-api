package stores

import (
	"fmt"

  "github.com/guitou4573/ipam/pkg/backends"
	"github.com/guitou4573/ipam/pkg/models"
)

// SubnetStore db structure for subnets
type SubnetStore struct {
	db backends.Backend
	tableName string
	initTableDisabled bool
}

type SubnetInfo interface {
  GetIdVPC() string
	GetNetAddr() int64
  GetNetMask() int64
}

// NewSubnetStore creates SQL store instance
func NewSubnetStore(db backends.Backend, options ...SubnetStoreOption) (*SubnetStore, error) {
	if db.GetConnection() == nil {
		return nil, fmt.Errorf("[SubnetStore] init failed, db missing")
	}

	store := &SubnetStore{
		db:           db,
		tableName:    "subnet",
	}
	for _, o := range options {
		o(store)
	}

	var err error
	if !store.initTableDisabled {
		err = store.initTable()
	}

	return store, err
}

func (s *SubnetStore) initTable() error {
	return nil
}

// GetBySubnetName retrieves and returns client information by SubnetName
func (s *SubnetStore) GetByVPCId(vpcId string) (*[]models.Subnet, error) {
	if vpcId == "" {
		return nil, nil
	}

  var subnets []models.Subnet
  err := s.db.GetConnection().Select(&subnets, fmt.Sprintf("SELECT * FROM %s WHERE idvpc = ?", s.tableName), vpcId)
	switch {
	case err == backends.SQLDbErrorNumRows:
			return nil, nil
		case err != nil:
			return nil, err
	}

	return &subnets, nil
}

// Create creates and stores the new subnet information
func (s *SubnetStore) Create(info SubnetInfo) (string, error) {
	res, errsql := s.db.GetConnection().Exec(fmt.Sprintf("INSERT INTO %s (idvpc, netaddr, mask, created_at) VALUES (?,CURRENT_TIMESTAMP)", s.tableName),
		info.GetIdVPC(),
    info.GetNetAddr(),
    info.GetNetMask(),
	)
	if errsql != nil {
		return "", errsql
	}

	id, errid := res.LastInsertId()
	if errid != nil {
		return "", errid
	} else {
		return string(id), nil
	}

}
