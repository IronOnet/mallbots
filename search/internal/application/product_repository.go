package application

import (
	"context"

	"github.com/irononet/mallbots/search/internal/models"
)

type ProductRepository interface{
	Find(ctx context.Context, productID string) (*models.Product, error)
}