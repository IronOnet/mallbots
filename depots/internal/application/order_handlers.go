package application

import (
	"context"

	"github.com/irononet/mallbots/depots/internal/domain"
	"github.com/irononet/mallbots/internal/ddd"
)

type OrderHandler struct{
	orders domain.OrderRepository
	ignoreUnimplementedDomainEvents
}

func NewOrderHandlers(orders domain.OrderRepository) OrderHandler{
	return OrderHandler{
		orders: orders,
	}
}

func (h OrderHandler) OnShoppingListCompleted(ctx context.Context, event ddd.Event) error{
	completed := event.(*domain.ShoppingListCompleted)
	return h.orders.Ready(ctx, completed.ShoppingList.OrderId)
}