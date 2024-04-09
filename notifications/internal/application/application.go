package application

import (
	"context"
)

type (
	OrderCreated struct{
		OrderID string
		CustomerID string
	}

	OrderCanceled struct{
		OrderID string
		CustomerID string
	}

	OrderReady struct{
		OrderID string
		CustomerID string
	}

	App interface{
		NotifyOrderCreated(ctx context.Context, notify OrderCreated) error
		NotifyOrderCanceled(ctx context.Context, notify OrderCanceled) error
		NotifyOrderReady(ctx context.Context, notify OrderReady) error
	}

	Application struct{
		customers CustomerRepository
	}
)

var _ App = (*Application)(nil)

func New(customers CustomerRepository) *Application{
	return &Application{
		customers: customers,
	}
}

func (a Application) NotifyOrderCreated(ctx context.Context, notify OrderCreated) error{
	// Not implemented
	return nil
}

func (a Application) NotifyOrderCanceled(ctx context.Context, notify OrderCanceled) error{
	// Not implemented
	return nil
}

func (a Application) NotifyOrderReady(ctx context.Context, notify OrderReady) error{
	// Not implemented
	return nil
}