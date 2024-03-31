package handlers


import (
	"github.com/irononet/mallbots/internal/ddd"
	"github.com/irononet/mallbots/ordering/internal/domain"
)

func RegisterNotificationHandler(notificationHandlers ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]){
	domainSubscriber.Subscribe(domain.OrderCreatedEvent, notificationHandlers)
	domainSubscriber.Subscribe(domain.OrderReadiedEvent, notificationHandlers)
	domainSubscriber.Subscribe(domain.OrderCanceledEvent, notificationHandlers)
}