package domain

import (
	"github.com/stackus/errors"

	"github.com/irononet/mallbots/internal/ddd"
)

var (
	ErrOrderHasNoItems = errors.Wrap(errors.ErrBadRequest, "the order has no items")
	ErrOrderCannotBeCancelled = errors.Wrap(errors.ErrBadRequest, "the order cannot be cancelled")
	ErrCustomerIDCannotBeBlank = errors.Wrap(errors.ErrBadRequest, "the customer id cannot be blank")
	ErrPaymentIDCannotBeBlank = errors.Wrap(errors.ErrBadRequest, "the payment id cannot be blank")
	ErrOrderIsAlreadyCanceled = errors.Wrap(errors.ErrBadRequest, "the order is already canceled")
	ErrOrderIsInReadyState = errors.Wrap(errors.ErrBadRequest, "the order is already in ready state")
	ErrOrderUknown = errors.Wrap(errors.ErrBadRequest, "the status of the order is unknown")
	ErrOrderAlreadyCompleted = errors.Wrap(errors.ErrBadRequest, "the order is already completed")
)

type Order struct{
	ddd.AggregateBase
	CustomerID string
	PaymentID string
	InvoiceID string
	ShoppingID string
	Items []*Item
	Status OrderStatus
}

func CreateOrder(id, customerID, paymentID string, items []*Item) (*Order, error){
	if len(items) == 0{
		return nil, ErrOrderHasNoItems
	}

	if customerID == ""{
		return nil, ErrCustomerIDCannotBeBlank
	}

	if paymentID == ""{
		return nil, ErrPaymentIDCannotBeBlank
	}

	order := &Order{
		AggregateBase: ddd.AggregateBase{
			ID: id,
		},
		CustomerID: customerID,
		PaymentID: paymentID,
		Items: items,
		Status: OrderIsPending,
	}

	order.AddEvent(&OrderCreated{
		Order: order,
	})

	return order, nil
}

func (o *Order) Cancel() error{
	if o.Status != OrderIsPending{
		return ErrOrderCannotBeCancelled
	}

	o.Status = OrderIsCancelled

	o.AddEvent(&OrderCanceled{
		Order: o,
	})
	return nil
}

func (o *Order) Ready() error{
	switch o.Status{
	case OrderIsCancelled:
		return ErrOrderIsAlreadyCanceled
	case OrderIsPending, OrderIsInProcess:
		return ErrOrderIsInReadyState
	case OrderIsCompleted:
		return ErrOrderAlreadyCompleted
	default:
		o.Status = OrderIsReady
	}

	o.AddEvent(&OrderReadied{
		Order: o,
	})

	return nil
}

func (o *Order) Complete(invoiceID string) error{
	// Validate invoice exists

	// Validate status
	o.InvoiceID = invoiceID
	o.Status = OrderIsCompleted

	o.AddEvent(&OrderCompleted{
		Order: o,
	})

	return nil
}

func (o Order) GetTotal() float64{
	var total float64

	for _, item := range o.Items{
		total += item.Price * float64(item.Quantity)
	}

	return total
}