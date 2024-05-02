package handlers

import (
	"context"

	"github.com/irononet/mallbots/customers/customerspb"
	"github.com/irononet/mallbots/internal/am"
	"github.com/irononet/mallbots/internal/ddd"
)

func RegisterCustomerHandlers(customerHandlers ddd.EventHandler[ddd.Event], stream am.EventSubscriber) error{
	evtMsgHandler := am.MessageHandlerFunc[am.EventMessage](func(ctx context.Context, evtMsg am.EventMessage) error {
		return customerHandlers.HandleEvent(ctx, evtMsg)
	})

	return stream.Subscribe(customerspb.CustomerAggregateChannel, evtMsgHandler, am.MessageFilter{
		customerspb.CustomerRegisteredEvent,
	}, am.GroupName("search-customer"))
}