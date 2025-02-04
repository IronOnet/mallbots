package grpc

import (
	"context"

	"google.golang.org/grpc"

	"github.com/irononet/mallbots/search/internal/application"
	"github.com/irononet/mallbots/search/internal/models"
	"github.com/irononet/mallbots/stores/storespb"
)

type StoreRepository struct{
	client storespb.StoresServiceClient
}

var _ application.StoreRepository = (*StoreRepository)(nil)

func NewStoreRepository(conn *grpc.ClientConn) StoreRepository{
	return StoreRepository{
		client: storespb.NewStoresServiceClient(conn),
	}
}

func (r StoreRepository) Find(ctx context.Context, storeID string) (*models.Store, error){
	resp, err := r.client.GetStore(ctx, &storespb.GetStoreRequest{Id: storeID })
	if err != nil{
		return nil, err
	}

	return r.storeToDomain(resp.Store), nil
}

func (r StoreRepository) storeToDomain(store *storespb.Store) *models.Store{
	return &models.Store{
		ID: store.GetId(),
		Name: store.GetName(),
	}
}