package domain

import (
	"context"
)

type StoreRepository interface {
	Save(ctx context.Context, store *Store) error
	Find(ctx context.Context, storeID string) (*Store, error)
	Update(ctx context.Context, store *Store) error
	FindAll(ctx context.Context) ([]*Store, error)
}
