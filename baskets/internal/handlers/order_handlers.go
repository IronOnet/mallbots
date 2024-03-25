package handlers

import (
	"github.com/irononet/mallbots/baskets/internal/application"
	"github.com/irononet/mallbots/baskets/internal/domain"
	"github.com/irononet/mallbots/internal/ddd"
)

func RegisterOrderHandlers(orderHandlers application.DomainEventHandlers, domainSubscriber ddd.EventSubscriber) {
	domainSubscriber.Subscribe(domain.BasketCheckedOut{}, orderHandlers.OnBasketCheckedOut)
}
