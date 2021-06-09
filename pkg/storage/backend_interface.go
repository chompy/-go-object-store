package storage

import "gitlab.com/contextualcode/storage-backend/pkg/types"

// Backend defines a backend connector.
type Backend interface {
	Put(*types.Object) error
	Delete(*types.Object) error
	Get(string) (*types.Object, error)
	Query(types.Query) ([]*types.Object, error)
}
