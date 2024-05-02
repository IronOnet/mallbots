package grpc

import (
	"context"

	"google.golang.org/grpc"

	"github.com/irononet/mallbots/customers/customerspb"
	"github.com/irononet/mallbots/search/internal/application"
	"github.com/irononet/mallbots/search/internal/models"
)

type CustomerRepository struct{
	client customerspb.CustomersServiceClient
}

var _ application.CustomerRepository = (*CustomerRepository)(nil)

func NewCustomerRepository(conn *grpc.ClientConn) CustomerRepository{
	return CustomerRepository{
		client: customerspb.NewCustomersServiceClient(conn),
	}
}

func (r CustomerRepository) Find(ctx context.Context, customerID string) (*models.Customer, error){
	resp, err := r.client.GetCustomer(ctx, &customerspb.GetCustomerRequest{Id: customerID})
	if err != nil{
		return nil, err
	}

	return &models.Customer{
		ID: resp.GetCustomer().GetId(),
		Name: resp.GetCustomer().GetName(),
	}, nil
}

