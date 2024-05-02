package handlers

import (
	"context"

	"github.com/irononet/mallbots/internal/am"
	"github.com/irononet/mallbots/internal/ddd"
	"github.com/irononet/mallbots/ordering/orderingpb"
)

func RegisterOrderHandlers(orderHandlers ddd.EventHandler[ddd.Event], stream am.EventSubscriber) error{
	evtMsgHandler := am.MessageHandlerFunc[am.EventMessage](func(ctx context.Context, evtMsg am.EventMessage) error{
		return orderHandlers.HandleEvent(ctx, evtMsg)
	})

	return stream.Subscribe(orderingpb.OrderAggregateChannel, evtMsgHandler, am.MessageFilter{
		orderingpb.OrderCreatedEvent,
		orderingpb.OrderReadiedEvent,
		orderingpb.OrderCanceledEvent,
		orderingpb.OrderCompletedEvent,
	}, am.GroupName("notification-orders"))
}