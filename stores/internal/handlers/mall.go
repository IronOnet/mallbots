package handlers

import (
	"github.com/irononet/mallbots/internal/ddd"
	"github.com/irononet/mallbots/stores/internal/domain"
)

func RegisterMallHandlers(mallHandlers ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]){
	domainSubscriber.Subscribe(mallHandlers,
	domain.StoreCreatedEvent,
	domain.StoreParticipationEnabledEvent,
	domain.StoreParticipationDisabledEvent,
	domain.StoreRebrandedEvent,
	)
}