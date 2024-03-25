package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/stackus/errors"

	"github.com/irononet/mallbots/baskets/internal/domain"
	"github.com/irononet/mallbots/internal/ddd"
)

type BasketRepository struct{
	tableName string
	db *sql.DB
}

var _ domain.BasketRepository = (*BasketRepository)(nil)

func NewBasketRepository(tableName string, db *sql.DB) BasketRepository{
	return BasketRepository{tableName: tableName, db: db}
}

func (r BasketRepository) Find(ctx context.Context, basketID string) (*domain.Basket, error){
	const query = "SELECT customer_id, payment_id, items, status FROM %s WHERE id = $1 LIMIT 1"

	basket := &domain.Basket{
		AggregateBase: ddd.AggregateBase{
			ID: basketID,
		},
	}

	var items []byte
	var status string

	err := r.db.QueryRowContext(ctx, r.table(query), basketID).Scan(&basket.CustomerID, &basket.PaymentID, &items, &status)
	if err != nil{
		return nil, errors.ErrInternalServerError.Err(err)
	}

	basket.Status, err = r.statusToDomain(status)
	if err != nil{
		return nil, errors.ErrInternalServerError.Err(err)
	}

	err = json.Unmarshal(items, &basket.Items)
	if err != nil{
		return nil, errors.ErrInternalServerError.Err(err)
	}

	return basket, nil
}

func (r BasketRepository) Save(ctx context.Context, basket *domain.Basket) error{
	const query = "INSERT INTO %s (id, customer_id, payment_id, items, status) VALUES ($1, $2, $3, $4, $5)"

	items, err := json.Marshal(basket.Items)
	if err != nil{
		return errors.ErrInternalServerError.Err(err)
	}

	_, err = r.db.ExecContext(ctx, r.table(query), basket.ID, basket.CustomerID, basket.PaymentID, items, basket.Status.String())

	return errors.ErrInternalServerError.Err(err)
}

func (r BasketRepository) Update(ctx context.Context, basket *domain.Basket) error{
	const query = "UPDATE %s SET customer_id = $2, payment_id = $3, items = $4, status = $5 WHERE id = $1"

	items, err := json.Marshal(basket.Items)
	if err != nil{
		return errors.ErrInternalServerError.Err(err)
	}

	_, err = r.db.ExecContext(ctx, r.table(query), basket.ID, basket.CustomerID, basket.PaymentID, items, basket.Status.String())

	return errors.ErrInternalServerError.Err(err)
}

func (r BasketRepository) DeleteBasket(ctx context.Context, basketID string) error{
	const query = "DELETE FROM %s WHERE id = $1 LIMIT 1"

	_, err := r.db.ExecContext(ctx, r.table(query), basketID)

	return errors.ErrInternalServerError.Err(err)
}

func (r BasketRepository) table(query string) string{
	return fmt.Sprintf(query, r.tableName)
}

func (r BasketRepository) statusToDomain(status string) (domain.BasketStatus, error){
	switch status{
	case domain.BasketIsOpen.String():
		return domain.BasketIsOpen, nil 
	case domain.BasketIsCanceled.String():
		return domain.BasketIsCanceled, nil
	case domain.BasketIsCheckout.String():
		return domain.BasketIsCheckout, nil
	default:
		return domain.BasketUnknown, fmt.Errorf("unkown basket status: %s", status)
	}

}