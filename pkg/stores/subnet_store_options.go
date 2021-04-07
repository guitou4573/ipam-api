package stores

// SubnetStoreOption is the configuration options type for client store
type SubnetStoreOption func(s *SubnetStore)

// WithSubnetStoreTableName returns option that sets client store table name
func WithSubnetStoreTableName(tableName string) SubnetStoreOption {
	return func(s *SubnetStore) {
		s.tableName = tableName
	}
}
