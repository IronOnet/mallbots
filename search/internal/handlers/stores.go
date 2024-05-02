package handlers

import (
	"context"

	"github.com/irononet/mallbots/internal/am"
	"github.com/irononet/mallbots/internal/ddd"
	"github.com/irononet/mallbots/stores/storespb"
)

func RegisterStoreHandlers(storeHandlers ddd.EventHandler[ddd.Event], stream am.EventSubscriber) error{
	evtMsgHandler := am.MessageHandlerFunc[am.EventMessage](func(ctx context.Context, evtMessage am.EventMessage) error{
		return storeHandlers.HandleEvent(ctx, evtMessage)
	})

	return stream.Subscribe(storespb.StoreAggregateChannel, evtMsgHandler, am.MessageFilter{
		storespb.StoreCreatedEvent,
		storespb.StoreRebrandedEvent,
	}, am.GroupName("search-stores"))
}