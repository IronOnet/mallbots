package application

import (
	"context"

	"github.com/irononet/mallbots/search/internal/models"
)

type CustomerRepository interface{
	Find(ctx context.Context, customerID string) (*models.Customer, error)
}