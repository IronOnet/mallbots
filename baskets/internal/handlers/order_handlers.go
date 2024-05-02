package handlers

import (
	"github.com/irononet/mallbots/baskets/internal/domain"
	"github.com/irononet/mallbots/internal/ddd"
)

func RegisterOrderHandlers(orderHandlers ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]) {
	domainSubscriber.Subscribe(orderHandlers, domain.BasketCheckedOutEvent)
}
