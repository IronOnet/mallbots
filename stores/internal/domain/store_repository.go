package domain

import (
	"context"
)

type StoreRepository interface {
	Save(ctx context.Context, store *Store) error
	Load(ctx context.Context, storeID string) (*Store, error)	
}
