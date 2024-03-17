package application

import (
	"context"

	"github.com/irononet/mallbots/notifications/internal/models"
)

type CustomerRepository interface{
	Find(ctx context.Context, customerID string) (*models.Customer, error)
}