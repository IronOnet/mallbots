package search

import (
	"context"

	"github.com/irononet/mallbots/search/internal/handlers"
	"github.com/irononet/mallbots/search/internal/logging"
	"github.com/irononet/mallbots/customers/customerspb"
	"github.com/irononet/mallbots/internal/am"
	"github.com/irononet/mallbots/internal/ddd"
	"github.com/irononet/mallbots/internal/jetstream"
	"github.com/irononet/mallbots/internal/monolith"
	"github.com/irononet/mallbots/internal/registry"
	"github.com/irononet/mallbots/ordering/orderingpb"
	"github.com/irononet/mallbots/search/internal/application"
	"github.com/irononet/mallbots/search/internal/grpc"
	"github.com/irononet/mallbots/search/internal/postgres"
	"github.com/irononet/mallbots/search/internal/rest"
	"github.com/irononet/mallbots/stores/storespb"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Monolith) (err error){
	reg := registry.New()
	if err = orderingpb.Registrations(reg); err != nil{
		return err
	}
	if err = customerspb.Registrations(reg); err != nil{
		return err
	}
	if err = storespb.Registrations(reg); err != nil{
		return err
	}
	eventStream := am.NewEventStream(reg, jetstream.NewStream(mono.Config().Nats.Stream, mono.JS()))
	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil{
		return err
	}
	customers := postgres.NewCustomerCacheRepository("search.customers_cache", mono.DB(), grpc.NewCustomerRepository(conn))
	stores := postgres.NewStoreCacheRepository("search.stores_cache", mono.DB(), grpc.NewStoreRepository(conn))
	products := postgres.NewProductCacheRepository("search.products_cache", mono.DB(), grpc.NewProductRepository(conn))
	orders := postgres.NewOrderRepository("search.orders", mono.DB())

	app := logging.LogApplicationAccess(
		application.New(orders),
		mono.Logger(),
	)

	orderHandlers := logging.LogEventHandlerAccess[ddd.Event](
		application.NewOrderHandlers(orders, customers, stores, products),
		"Order", mono.Logger(),
	)

	customerHandlers := logging.LogEventHandlerAccess[ddd.Event](
		application.NewCustomerHandler(customers),
		"Customer", mono.Logger(),
	)

	storeHandlers := logging.LogEventHandlerAccess[ddd.Event](
		application.NewStoreHandlers(stores),
		"Store", mono.Logger(),
	)

	productHandlers := logging.LogEventHandlerAccess[ddd.Event](
		application.NewProductHandlers(products),
		"Product", mono.Logger(),
	)

	// setup driver adapters
	if err = grpc.RegisterServer(ctx, app, mono.RPC()); err != nil{
		return err
	}
	if err = rest.RegisterGateway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil{
		return err
	}
	if err = rest.RegisterSwagger(mono.Mux()); err != nil{
		return err
	}
	if err = handlers.RegisterOrderHandlers(orderHandlers, eventStream); err != nil{
		return err
	}
	if err = handlers.RegisterCustomerHandlers(customerHandlers, eventStream); err != nil{
		return err
	}
	if err = handlers.RegisterStoreHandlers(storeHandlers, eventStream); err != nil{
		return err
	}
	if err = handlers.RegisterProductHandlers(productHandlers, eventStream); err != nil{
		return err
	}

	return nil
}