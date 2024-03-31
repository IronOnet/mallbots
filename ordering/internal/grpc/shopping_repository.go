package grpc

import (
	"context"

	"google.golang.org/grpc"

	"github.com/irononet/mallbots/depots/depotspb"
	"github.com/irononet/mallbots/ordering/internal/domain"
)

type ShoppingRepository struct{
	client depotspb.DepotServiceClient
}

var _ domain.ShoppingRepository = (*ShoppingRepository)(nil)

func NewShoppingListRepository(conn *grpc.ClientConn) ShoppingRepository{
	return ShoppingRepository{client: depotspb.NewDepotServiceClient(conn)}
}

func (r ShoppingRepository) Create(ctx context.Context, orderID string, orderItems []domain.Item) (string, error){
	items := make([]*depotspb.OrderItem, 0, len(orderItems))
	for i, item := range orderItems{
		items[i] = r.itemFromDomain(item)
	}

	response, err := r.client.CreateShoppingList(ctx, &depotspb.CreateShoppingListRequest{
		OrderId: orderID,
		Items: items,
	})
	if err != nil{
		return "", err
	}
	
	return response.GetId(), nil
}

func (r ShoppingRepository) Cancel(ctx context.Context, shoppingID string) error{
	_, err := r.client.CancelShoppingList(ctx, &depotspb.CancelShoppingListRequest{Id: shoppingID})
	return err
}

func (r ShoppingRepository) itemFromDomain(item domain.Item) *depotspb.OrderItem{
	return &depotspb.OrderItem{
		ProductId: item.ProductID,
		StoreId: item.StoreID,
		Quantity: int32(item.Quantity),
	}
}