package handlers


import (
	"github.com/irononet/mallbots/internal/ddd"
	"github.com/irononet/mallbots/ordering/internal/domain"
)

func RegisterNotificationHandler(notificationHandlers ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]){
	domainSubscriber.Subscribe(notificationHandlers,
	domain.OrderCreatedEvent,
	domain.OrderReadiedEvent,
	domain.OrderCanceledEvent,
	)
}