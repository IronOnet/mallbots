package depots

import (
	"context"
	"github.com/irononet/mallbots/depots/internal/application"
	"github.com/irononet/mallbots/depots/internal/grpc"
	"github.com/irononet/mallbots/depots/internal/handlers"
	"github.com/irononet/mallbots/depots/internal/logging"
	"github.com/irononet/mallbots/depots/internal/postgres"
	"github.com/irononet/mallbots/depots/internal/rest"
	"github.com/irononet/mallbots/internal/ddd"
	"github.com/irononet/mallbots/internal/monolith"
)

type Module struct{}

func (Module) Startup(ctx context.Context, mono monolith.Monolith) error{
	// Setup driven adapters
	domainDispatcher := ddd.NewEventDispatcher()
	shoppingLists := postgres.NewShoppingListRepository("depot.shopping_lists", mono.DB())
	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil{
		return err
	}
	stores := grpc.NewStoreRepository(conn)
	products := grpc.NewProductRepository(conn)
	orders := grpc.NewOrderRepository(conn)

	// Setup application 
	app := logging.LogApplicationAccess(application.New(shoppingLists, stores, products, domainDispatcher),
	mono.Logger())

	orderHandlers := logging.LogDomainEventHandlerAccess(
		application.NewOrderHandlers(orders),
		mono.Logger(),
	)

	// setup Driver adapters
	if err := grpc.Register(ctx, app, mono.RPC()); err != nil{
		return err
	}
	if err := rest.RegisterGatway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil{
		return err
	}
	if err := rest.RegisterSwagger(mono.Mux()); err != nil{
		return err
	}

	handlers.RegisterOrderHandler(orderHandlers, domainDispatcher)

	return nil
}