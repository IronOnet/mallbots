package handlers

import (
	"github.com/irononet/mallbots/depots/internal/application"
	"github.com/irononet/mallbots/depots/internal/domain"
	"github.com/irononet/mallbots/internal/ddd"
)

func RegisterOrderHandler(orderHandlers application.DomainEventHandlers, domainSubscriber ddd.EventSubscriber){
	domainSubscriber.Subscribe(domain.ShoppingListCompleted{}, orderHandlers.OnShoppingListCompleted)
}