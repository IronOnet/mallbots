package grpc

import (
	"context"

	"google.golang.org/grpc"

	"github.com/irononet/mallbots/depots/internal/domain"
	"github.com/irononet/mallbots/ordering/orderingpb"
)

type OrderRepository struct{
	client orderingpb.OrderingServiceClient
}

var _ domain.OrderRepository = (*OrderRepository)(nil)

func NewOrderRepository(conn *grpc.ClientConn) OrderRepository{
	return OrderRepository{client: orderingpb.NewOrderingServiceClient(conn)}
}

func (r OrderRepository) Ready(ctx context.Context, orderID string) error{
	_, _err := r.client.ReadyOrder(ctx, &orderingpb.ReadyOrderRequest{Id: orderID})
	return _err
}