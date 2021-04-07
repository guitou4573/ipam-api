package stores

// TokenStoreOption is the configuration options type for token store
type TokenStoreOption func(s *TokenStore)

// WithTokenStoreTableName returns option that sets token store table name
func WithTokenStoreTableName(tableName string) TokenStoreOption {
	return func(s *TokenStore) {
		s.tableName = tableName
	}
}
