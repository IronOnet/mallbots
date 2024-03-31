package domain

import (
	"github.com/stackus/errors"

	"github.com/irononet/mallbots/internal/ddd"
	"github.com/irononet/mallbots/internal/es"
)

const OrderAggregate = "ordering.Order"

var (
	ErrOrderAlreadyCreated = errors.Wrap(errors.ErrBadRequest, "the order cannot be created")
	ErrOrderHasNoItems = errors.Wrap(errors.ErrBadRequest, "the order has no items")
	ErrOrderCannotBeCancelled = errors.Wrap(errors.ErrBadRequest, "the order cannot be cancelled")
	ErrCustomerIDCannotBeBlank = errors.Wrap(errors.ErrBadRequest, "the customer id cannot be blank")
	ErrPaymentIDCannotBeBlank = errors.Wrap(errors.ErrBadRequest, "the payment id cannot be blank")
	ErrOrderIsAlreadyCanceled = errors.Wrap(errors.ErrBadRequest, "the order is already canceled")
	ErrOrderIsInReadyState = errors.Wrap(errors.ErrBadRequest, "the order is already in ready state")
	ErrOrderUknown = errors.Wrap(errors.ErrBadRequest, "the status of the order is unknown")
	ErrOrderAlreadyCompleted = errors.Wrap(errors.ErrBadRequest, "the order is already completed")
)

var _ interface{
	es.EventApplier
	es.Snapshotter
} = (*Order)(nil)

func NewOrder(id string) *Order{
	return &Order{
		Aggregate: es.NewAggregate(id, OrderAggregate),
	}
}

type Order struct{
	es.Aggregate
	CustomerID string
	PaymentID string
	InvoiceID string
	ShoppingID string
	Items []Item
	Status OrderStatus
}

func (Order) Key() string { return OrderAggregate }

func(o *Order) CreateOrder(id, customerID, paymentID, shoppingID string, items []Item) error{
	if o.Status != OrderUnknown{
		return ErrOrderAlreadyCreated
	}

	if len(items) == 0{
		return ErrOrderHasNoItems
	}

	if customerID == ""{
		return ErrCustomerIDCannotBeBlank
	}

	if paymentID == ""{
		return ErrPaymentIDCannotBeBlank
	}

	o.AddEvent(OrderCreatedEvent, &OrderCreated{
		CustomerID: customerID,
		PaymentID: paymentID,
		ShoppingID: shoppingID,
		Items: items,
	})
	
	return nil
}

func (o *Order) Cancel() error{
	if o.Status != OrderIsPending{
		return ErrOrderCannotBeCancelled
	}

	o.Status = OrderIsCancelled

	o.AddEvent(OrderCanceledEvent, &OrderCanceled{
		CustomerID: o.CustomerID,
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

	o.AddEvent(OrderReadiedEvent, &OrderReadied{
		CustomerID: o.CustomerID,
		PaymentID: o.PaymentID,
		Total: o.GetTotal(),
	})

	return nil
}

func (o *Order) Complete(invoiceID string) error{
	// Validate invoice exists

	// Validate status
	o.InvoiceID = invoiceID
	o.Status = OrderIsCompleted

	o.AddEvent(OrderCompletedEvent, &OrderCompleted{
		InvoiceID: invoiceID,
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

func (o *Order) ApplyEvent(event ddd.Event) error{
	switch payload := event.Payload().(type){
	case *OrderCreated:
		o.CustomerID = payload.CustomerID
		o.PaymentID = payload.PaymentID
		o.ShoppingID = payload.ShoppingID
		o.Items = payload.Items
		o.Status = OrderIsPending

	case *OrderCanceled:
		o.Status = OrderIsCancelled

	case *OrderReadied:
		o.Status = OrderIsReady
	
	case *OrderCompleted:
		o.InvoiceID = payload.InvoiceID
		o.Status = OrderIsCompleted

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", o, event.EventName(), payload)
	}

	return nil
}

func (o *Order) ApplySnapshot(snapshot es.Snapshot) error{
	switch ss := snapshot.(type){
	case *OrderV1:
		o.CustomerID = ss.CustomerID
		o.PaymentID = ss.PaymentID
		o.InvoiceID = ss.InvoiceID
		o.ShoppingID = ss.ShoppingID
		o.Items = ss.Items
		o.Status = ss.Status

	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", o, snapshot)
	}

	return nil
}

func (o *Order) ToSnapshot() es.Snapshot{
	return &OrderV1{
		CustomerID: o.CustomerID,
		PaymentID: o.PaymentID,
		InvoiceID: o.InvoiceID,
		ShoppingID: o.ShoppingID,
		Items: o.Items,
		Status: o.Status,
	}
}