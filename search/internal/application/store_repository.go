package application

import (
	"context"

	"github.com/irononet/mallbots/search/internal/models"
)

type StoreRepository interface{
	Find(ctx context.Context, storeID string) (*models.Store, error)
}