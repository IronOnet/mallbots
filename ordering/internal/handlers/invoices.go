package handlers

import (
	"github.com/irononet/mallbots/internal/ddd"
	"github.com/irononet/mallbots/ordering/internal/domain"
)

func RegisterInvoiceHandlers(invoiceHandlers ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]){
	domainSubscriber.Subscribe(domain.OrderReadiedEvent, invoiceHandlers)
}