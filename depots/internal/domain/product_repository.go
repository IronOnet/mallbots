package domain

import (
	"context"
)

type ProductRepository interface {
	Find(ctx context.Context, productId string) (*Product, error)
}
