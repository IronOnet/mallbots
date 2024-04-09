package application

import (
	"context"

	"github.com/irononet/mallbots/internal/ddd"
	"github.com/irononet/mallbots/ordering/internal/domain"
)

type NotificationHandlers[T ddd.AggregateEvent] struct{
	notifications domain.NotificationRepository
	ignoreUnimplementedDomainEvents
}

var _ ddd.EventHandler[ddd.AggregateEvent] = (*NotificationHandlers[ddd.AggregateEvent])(nil)

func NewNotificationHandlers(notifications domain.NotificationRepository) *NotificationHandlers[ddd.AggregateEvent]{
	return &NotificationHandlers[ddd.AggregateEvent]{
		notifications: notifications,
	}
}

func (h NotificationHandlers[T]) HandleEvent(ctx context.Context, event T) error{
	switch event.EventName(){
	case domain.OrderCreatedEvent:
		return h.OnOrderCreated(ctx, event)
	case domain.OrderReadiedEvent:
		return h.OnOrderReadied(ctx, event)
	case domain.OrderCanceledEvent:
		return h.OnOrderCanceled(ctx, event)
	}
	return nil
}

func (h NotificationHandlers[T]) OnOrderCreated(ctx context.Context, event ddd.AggregateEvent) error{
	orderCreated := event.Payload().(*domain.OrderCreated)
	return h.notifications.NotifyOrderCreated(ctx, event.AggregateID(), orderCreated.CustomerID)
}

func (h NotificationHandlers[T]) OnOrderReadied(ctx context.Context, event ddd.AggregateEvent) error{
	orderReadied := event.Payload().(*domain.OrderReadied)
	return h.notifications.NotifyOrderReady(ctx, event.AggregateID(), orderReadied.CustomerID)
}

func (h NotificationHandlers[T]) OnOrderCanceled(ctx context.Context, event ddd.AggregateEvent) error{
	orderCanceled := event.Payload().(*domain.OrderCanceled)
	return h.notifications.NotifyOrderCanceled(ctx, event.AggregateID(), orderCanceled.CustomerID)
}