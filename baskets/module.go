package baskets

import (
	"context"

	"github.com/irononet/mallbots/baskets/internal/application"
	"github.com/irononet/mallbots/baskets/internal/domain"
	"github.com/irononet/mallbots/baskets/internal/grpc"
	"github.com/irononet/mallbots/baskets/internal/handlers"
	"github.com/irononet/mallbots/baskets/internal/logging"
	"github.com/irononet/mallbots/baskets/internal/rest"
	"github.com/irononet/mallbots/internal/ddd"
	"github.com/irononet/mallbots/internal/es"
	"github.com/irononet/mallbots/internal/monolith"
	pg "github.com/irononet/mallbots/internal/postgres"
	"github.com/irononet/mallbots/internal/registry"
	"github.com/irononet/mallbots/internal/registry/serdes"
)

type Module struct{}

func (Module) Startup(ctx context.Context, mono monolith.Monolith) (err error) {
	reg := registry.New()
	err = registrations(reg)
	if err != nil {
		return err
	}

	domainDispatcher := ddd.NewEventDispatcher[ddd.AggregateEvent]()
	aggregateStore := es.AggregateStoreWithMiddleware(
		pg.NewEventStore("baskets.events", mono.DB(), reg),
		es.NewEventPublisher(domainDispatcher),
		pg.NewSnapshotStore("baskets.snapshots", mono.DB(), reg),
	)

	baskets := es.NewAggregaterRepository[*domain.Basket](domain.BasketAggregate, reg, aggregateStore)
	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}
	stores := grpc.NewStoreRepository(conn)
	products := grpc.NewProductRepository(conn)
	orders := grpc.NewOrderRepository(conn)

	// setup application
	app := logging.LogApplicationAccess(
		application.New(baskets, stores, products, orders),
		mono.Logger(),
	)

	orderHandlers := logging.LogEventHandlerAccess[ddd.AggregateEvent](
		application.NewOrderHandlers(orders),
		"Order", mono.Logger(),
	)

	// Setup driver adapters
	if err := grpc.RegisterServer(app, mono.RPC()); err != nil {
		return err
	}

	if err := rest.RegisterGateway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil {
		return err
	}

	if err := rest.RegisterSwagger(mono.Mux()); err != nil {
		return err
	}

	handlers.RegisterOrderHandlers(orderHandlers, domainDispatcher)

	return
}

func registrations(reg registry.Registry) error {
	serde := serdes.NewJsonSerde(reg)

	// basket
	if err := serde.Register(domain.Basket{}, func(v interface{}) error {
		basket := v.(*domain.Basket)
		basket.Items = make(map[string]domain.Item)
		return nil
	}); err != nil {
		return err
	}

	// Baskets events
	if err := serde.Register(domain.BasketStarted{}); err != nil {
		return err
	}
	if err := serde.Register(domain.BasketCanceled{}); err != nil {
		return err
	}

	if err := serde.Register(domain.BasketCheckedOut{}); err != nil {
		return err
	}

	if err := serde.Register(domain.BasketItemAdded{}); err != nil {
		return err
	}

	if err := serde.Register(domain.BasketItemRemoved{}); err != nil {
		return err
	}

	if err := serde.RegisterKey(domain.BasketV1{}.SnapshotName(), domain.BasketV1{}); err != nil {
		return err
	}
	return nil
}
