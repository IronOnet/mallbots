package stores

import (
	"context"

	"github.com/irononet/mallbots/internal/ddd"
	"github.com/irononet/mallbots/internal/monolith"
	"github.com/irononet/mallbots/stores/internal/application"
	"github.com/irononet/mallbots/stores/internal/grpc"
	"github.com/irononet/mallbots/stores/internal/logging"
	"github.com/irononet/mallbots/stores/internal/postgres"
	"github.com/irononet/mallbots/stores/internal/rest"
)

type Module struct{}

func (Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	domainDispatcher := ddd.NewEventDispatcher()
	stores := postgres.NewStoreRepository("stores.stores", mono.DB())
	participatingStores := postgres.NewParticipatingStoreRepository("stores.stores", mono.DB())
	products := postgres.NewProductRepository("stores.products", mono.DB())

	// setup app
	var app application.App
	app = application.New(stores, participatingStores, products, domainDispatcher)
	app = logging.LogApplicationAccess(app, mono.Logger())

	// setup driver adapters
	if err := grpc.RegisterServer(ctx, app, mono.RPC()); err != nil {
		return err
	}
	if err := rest.RegisterGateway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil {
		return err
	}
	if err := rest.RegisterSwagger(mono.Mux()); err != nil {
		return err
	}

	return nil
}
