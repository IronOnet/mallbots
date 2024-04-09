package handlers

import (
	"github.com/irononet/mallbots/depots/internal/domain"
	"github.com/irononet/mallbots/internal/ddd"
)

func RegisterOrderHandler(orderHandlers ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]) {
	domainSubscriber.Subscribe(domain.ShoppingListCompletedEvent, orderHandlers)
}
