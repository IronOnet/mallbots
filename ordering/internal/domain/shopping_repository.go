package domain

import (
	"context"
)

type ShoppingRepository interface{
	Create(ctx context.Context, orderId string, items []Item) (string, error)
	Cancel(ctx context.Context, shoppingID string) error
}