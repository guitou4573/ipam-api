package stores

// VPCStoreOption is the configuration options type for token store
type VPCStoreOption func(s *VPCStore)

// WithVPCStoreTableName returns option that sets token store table name
func WithVPCStoreTableName(tableName string) VPCStoreOption {
	return func(s *VPCStore) {
		s.tableName = tableName
	}
}
