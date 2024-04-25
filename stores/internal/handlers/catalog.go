package handlers

import (
	"github.com/irononet/mallbots/internal/ddd"
	"github.com/irononet/mallbots/stores/internal/domain"
)

func RegisterCatalogHandlers(catalogHandlers ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]){
	domainSubscriber.Subscribe(catalogHandlers,
	domain.ProductAddedEvent,
	domain.ProductRebrandedEvent,
	domain.ProductPriceIncreasedEvent,
	domain.ProductPriceDecreasedEvent,
	domain.ProductRemovedEvent,
	)
}