package handlers


import (
	"github.com/irononet/mallbots/internal/ddd"
	"github.com/irononet/mallbots/ordering/internal/application"
	"github.com/irononet/mallbots/ordering/internal/domain"
)

func RegisterNotificationHandler(notificationHandlers application.DomainEventHandlers, domainSubscriber ddd.EventSubscriber){
	domainSubscriber.Subscribe(domain.OrderCreated{}, notificationHandlers.OnOrderCreated)
	domainSubscriber.Subscribe(domain.OrderReadied{}, notificationHandlers.OnOrderReadied)
	domainSubscriber.Subscribe(domain.OrderCanceled{}, notificationHandlers.OnOrderCanceled)
}