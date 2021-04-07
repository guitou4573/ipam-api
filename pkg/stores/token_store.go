package stores

import (
	"fmt"

  "github.com/guitou4573/ipam/pkg/backends"
	"github.com/guitou4573/ipam/pkg/models"
)

// TokenStore db structure for subnets
type TokenStore struct {
	db backends.Backend
	tableName string
	initTableDisabled bool
}

type TokenInfo interface {
	GetToken() string
}

// NewTokenStore creates SQL store instance
func NewTokenStore(db backends.Backend, options ...TokenStoreOption) (*TokenStore, error) {
	if db.GetConnection() == nil {
		return nil, fmt.Errorf("[TokenStore] init failed, db missing")
	}

	store := &TokenStore{
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

func (s *TokenStore) initTable() error {
	return nil
}

// GetByTokenName retrieves and returns client information by TokenName
func (s *TokenStore) Get(token string) (*models.Token, error) {
	if token == "" {
		return nil, nil
	}

	var item models.Token
	err := s.db.GetConnection().QueryRowx(fmt.Sprintf("SELECT * FROM %s WHERE token = ? LIMIT 1", s.tableName), token).StructScan(&item)
	switch {
	case err == backends.SQLDbErrorNumRows:
			return nil, nil
		case err != nil:
			return nil, err
	}

	return &item, nil
}

// Create creates and stores the new subnet information
func (s *TokenStore) Create(info TokenInfo) (string, error) {
	res, errsql := s.db.GetConnection().Exec(fmt.Sprintf("INSERT INTO %s (token, created_at) VALUES (?,CURRENT_TIMESTAMP)", s.tableName),
		info.GetToken(),
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
