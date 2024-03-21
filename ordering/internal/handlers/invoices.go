package handlers

import (
	"github.com/irononet/mallbots/internal/ddd"
	"github.com/irononet/mallbots/ordering/internal/application"
	"github.com/irononet/mallbots/ordering/internal/domain"
)

func RegisterInvoiceHandlers(invoiceHandler application.DomainEventHandlers, domainSubscriber ddd.EventSubscriber){
	domainSubscriber.Subscribe(domain.OrderReadied{}, invoiceHandler.OnOrderReadied)
}