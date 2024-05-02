package handlers

import (
	"context"

	"github.com/irononet/mallbots/internal/am"
	"github.com/irononet/mallbots/internal/ddd"
	"github.com/irononet/mallbots/stores/storespb"
)

func RegisterProductHandlers(productHandlers ddd.EventHandler[ddd.Event], stream am.EventSubscriber) error{
	evtMsgHandler := am.MessageHandlerFunc[am.EventMessage](func(ctx context.Context, evtMsg am.EventMessage) error{
		return productHandlers.HandleEvent(ctx, evtMsg)
	})

	return stream.Subscribe(storespb.ProductAggregateChannel, evtMsgHandler, am.MessageFilter{
		storespb.ProductAddedEvent,
		storespb.ProductRebrandedEvent,
		storespb.ProductRemovedEvent,
	}, am.GroupName("search-products"))
}