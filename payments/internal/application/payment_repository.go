package application

import (
	"context"

	"github.com/irononet/mallbots/payments/internal/models"
)

type PaymentRepository interface {
	Find(ctx context.Context, paymentID string) (*models.Payment, error)
	Save(ctx context.Context, payment *models.Payment) error
}
