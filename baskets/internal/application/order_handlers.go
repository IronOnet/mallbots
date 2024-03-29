package application

import (
	"context"

	"github.com/irononet/mallbots/baskets/internal/domain"
	"github.com/irononet/mallbots/internal/ddd"
)

type OrderHandlers[T ddd.AggregateEvent] struct {
	orders domain.OrderRepository
}

var _ ddd.EventHandler[ddd.AggregateEvent] = (*OrderHandlers[ddd.AggregateEvent])(nil)

func NewOrderHandlers(orders domain.OrderRepository) OrderHandlers[ddd.AggregateEvent] {
	return OrderHandlers[ddd.AggregateEvent]{
		orders: orders,
	}
}

func (h OrderHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.BasketCheckedOutEvent:
		return h.OnBasketCheckedOut(ctx, event)
	}
	return nil
}

func (h OrderHandlers[T]) OnBasketCheckedOut(ctx context.Context, event ddd.AggregateEvent) error {
	checkedOut := event.Payload().(*domain.BasketCheckedOut)
	_, err := h.orders.Save(ctx, checkedOut.PaymentID, checkedOut.CustomerID, checkedOut.Items)
	return err
}
