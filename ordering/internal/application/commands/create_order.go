package commands

import (
	"context"

	"github.com/stackus/errors"
	
	
	"github.com/irononet/mallbots/ordering/internal/domain"
)

type CreateOrder struct{
	ID string
	CustomerID string
	PaymentID string
	Items []domain.Item
}

type CreateOrderHandler struct{
	orders domain.OrderRepository
	customers domain.CustomerRepository
	payments domain.PaymentRepository
	shopping domain.ShoppingRepository
}

func NewCreateOrderHandler(orders domain.OrderRepository, customers domain.CustomerRepository, payments domain.PaymentRepository, shopping domain.ShoppingRepository) CreateOrderHandler{
	return CreateOrderHandler{
		orders: orders,
		customers: customers,
		payments: payments,
		shopping: shopping,
	}
}

func (h CreateOrderHandler) CreateOrder(ctx context.Context, cmd CreateOrder) error{
	order, err := h.orders.Load(ctx, cmd.ID)
	if err != nil{
		return err
	}

	// AuthorizeCustomer
	if err = h.customers.Authorize(ctx, order.CustomerID); err != nil{
		return errors.Wrap(err, "order customer authorization")
	}

	if err = h.payments.Confirm(ctx, order.PaymentID); err != nil{
		return errors.Wrap(err, "order payment confirmation")
	}

	if order.ShoppingID, err = h.shopping.Create(ctx, cmd.ID, cmd.Items); err != nil{
		return errors.Wrap(err, "order shopping scheduling")
	}

	// orderCreation
	if err = h.orders.Save(ctx, order); err != nil{
		return errors.Wrap(err, "order creation")
	}


	return nil
}