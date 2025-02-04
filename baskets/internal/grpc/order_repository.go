package grpc

import (
	"context"

	"github.com/stackus/errors"
	"google.golang.org/grpc"

	"github.com/irononet/mallbots/baskets/internal/domain"
	"github.com/irononet/mallbots/ordering/orderingpb"
)

type OrderRepository struct {
	client orderingpb.OrderingServiceClient
}

var _ domain.OrderRepository = (*OrderRepository)(nil)

func NewOrderRepository(conn *grpc.ClientConn) OrderRepository {
	return OrderRepository{client: orderingpb.NewOrderingServiceClient(conn)}
}

func (r OrderRepository) Save(ctx context.Context, paymentID, customerID string, basketItems map[string]domain.Item) (string, error) {
	items := make([]*orderingpb.Item, 0, len(basketItems))
	for _, item := range basketItems {
		items = append(items, &orderingpb.Item{
			StoreId:     item.StoreID,
			ProductId:   item.ProductID,
			StoreName:   item.ProductName,
			ProductName: item.ProductName,
			Price:       item.ProductPrice,
			Quantity:    int32(item.Quantity),
		})
	}

	resp, err := r.client.CreateOrder(ctx, &orderingpb.CreateOrderRequest{
		Items:      items,
		CustomerId: customerID,
		PaymentId:  paymentID,
	})

	if err != nil {
		return "", errors.Wrap(err, "saving order")
	}

	return resp.GetId(), nil
}
