package stores

import (
	"fmt"

  "github.com/guitou4573/ipam/pkg/backends"
	"github.com/guitou4573/ipam/pkg/models"
)

// VPCStore db structure for vpcs
type VPCStore struct {
	db backends.Backend
	tableName string
	initTableDisabled bool
}

type VPCInfo interface {
  GetId() string
	GetIdSection() string
  GetDescription() int64
}

// NewVPCStore creates SQL store instance
func NewVPCStore(db backends.Backend, options ...VPCStoreOption) (*VPCStore, error) {
	if db.GetConnection() == nil {
		return nil, fmt.Errorf("[VPCStore] init failed, db missing")
	}

	store := &VPCStore{
		db:           db,
		tableName:    "vpc",
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

func (s *VPCStore) initTable() error {
	return nil
}

// GetBySection retrieves and returns a list of VPCs for a given section
func (s *VPCStore) GetBySection(section string, offset int, limit int) (*[]models.VPC, error) {
	if section == "" {
		return nil, nil
	}

	var vpcs []models.VPC
  err := s.db.GetConnection().Select(&vpcs, fmt.Sprintf("SELECT * FROM %s WHERE idsection = ? LIMIT ?,?", s.tableName),
    section,
    offset,
    limit,
  )
	switch {
	case err == backends.SQLDbErrorNumRows:
			return nil, nil
		case err != nil:
			return nil, err
	}

	return &vpcs, nil
}

// Create creates and stores the new vpc information
func (s *VPCStore) Create(info VPCInfo) (string, error) {
	res, errsql := s.db.GetConnection().Exec(fmt.Sprintf(
    "INSERT INTO %s (id, idsection, description, created_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP)",
    s.tableName,
    ),
		info.GetId(),
    info.GetIdSection(),
    info.GetDescription(),
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
