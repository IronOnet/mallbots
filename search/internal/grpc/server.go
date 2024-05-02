package grpc

import (
	"context"

	"google.golang.org/grpc"

	"github.com/irononet/mallbots/search/internal/application"
	"github.com/irononet/mallbots/search/searchpb"
)

type server struct{
	app application.Application
	searchpb.UnimplementedSearchServiceServer
}

func RegisterServer(_ context.Context, app application.Application, registrar grpc.ServiceRegistrar) error{
	searchpb.RegisterSearchServiceServer(registrar, server{app: app})
	return nil
}

func (s server) SearchOrders(ctx context.Context, request *searchpb.SearchOrdersRequest) (*searchpb.SearchOrdersResponse, error){
	panic("not implemented")
}

func (s server) GetOrder(ctx context.Context, request *searchpb.GetOrderRequest) (*searchpb.GetOrderResponse, error){
	panic("not implemented")
}