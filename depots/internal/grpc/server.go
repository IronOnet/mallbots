package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/irononet/mallbots/depots/depotspb"
	"github.com/irononet/mallbots/depots/internal/application"
	"github.com/irononet/mallbots/depots/internal/application/commands"
)

type server struct{
	app application.App
	depotspb.UnimplementedDepotServiceServer
}

var _ depotspb.DepotServiceServer = (*server)(nil)

func Register(_ context.Context, app application.App, registrar grpc.ServiceRegistrar) error{
	depotspb.RegisterDepotServiceServer(registrar, server{app: app})
	return nil
}

func (s server) CreateShoppingList(ctx context.Context, request *depotspb.CreateShoppingListRequest) (*depotspb.CreateShoppingListResponse, error){
	id := uuid.New().String()

	items := make([]commands.OrderItem, 0, len(request.GetItems()))
	for _, item := range request.GetItems(){
		items = append(items, s.itemToDomain(item))
	}

	err := s.app.CreateShoppingList(ctx, commands.CreateShoppingList{
		ID: id,
		OrderID: request.GetOrderId(),
		Items: items,
	})

	return &depotspb.CreateShoppingListResponse{Id: id}, err
}

func (s server) CancelShoppingList(ctx context.Context, request *depotspb.CancelShoppingListRequest) (*depotspb.CancelShoppingListResponse, error){
	err := s.app.CancelShoppingList(ctx, commands.CancelShoppingList{
		ID: request.GetId(),
	})

	return &depotspb.CancelShoppingListResponse{}, err
}

func (s server) AssignShoppingList(ctx context.Context, request *depotspb.AssignShoppingListRequest) (*depotspb.AssignShoppingListResponse, error){
	err := s.app.AssignShoppingList(ctx, commands.AssignShoppingList{
		ID: request.GetId(),
		BotID: request.GetBotId(),
	})
	return &depotspb.AssignShoppingListResponse{}, err
}

func (s server) CompleteShoppingList(ctx context.Context, request *depotspb.CompleteShoppingListRequest) (*depotspb.CompleteShoppingListResponse, error){
	err := s.app.CompleteShoppingList(ctx, commands.CompleteShoppingList{ID: request.GetId()})
	return &depotspb.CompleteShoppingListResponse{}, err
}

func (s server) itemToDomain(item *depotspb.OrderItem) commands.OrderItem{
	return commands.OrderItem{
		StoreID: item.GetStoreId(),
		ProductID: item.GetProductId(),
		Quantity: int(item.GetQuantity()),
	}
}