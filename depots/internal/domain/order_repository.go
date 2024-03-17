package domain

import (
	"context"
)

type OrderRepository interface{
	Ready(ctx context.Context, orderId string) error
}