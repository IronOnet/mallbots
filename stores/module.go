package stores

import (
	"context"

	"github.com/irononet/mallbots/internal/ddd"
	"github.com/irononet/mallbots/internal/es"
	"github.com/irononet/mallbots/internal/monolith"
	pg "github.com/irononet/mallbots/internal/postgres"
	"github.com/irononet/mallbots/internal/registry"
	"github.com/irononet/mallbots/internal/registry/serdes"
	"github.com/irononet/mallbots/stores/internal/application"
	"github.com/irononet/mallbots/stores/internal/domain"
	"github.com/irononet/mallbots/stores/internal/grpc"
	"github.com/irononet/mallbots/stores/internal/handlers"
	"github.com/irononet/mallbots/stores/internal/logging"
	"github.com/irononet/mallbots/stores/internal/postgres"
	"github.com/irononet/mallbots/stores/internal/rest"
)

type Module struct{}

func (Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	reg := registry.New()
	err := registrations(reg)
	if err != nil{
		return err
	}
	domainDispatcher := ddd.NewEventDispatcher[ddd.AggregateEvent]()
	aggregateStore := es.AggregateStoreWithMiddleware(
		pg.NewEventStore("stores.events", mono.DB(), reg),
		es.NewEventPublisher(domainDispatcher),
		pg.NewSnapshotStore("stores.snapshots", mono.DB(), reg),
	)


	stores := es.NewAggregaterRepository[*domain.Store](domain.StoreAggregate, reg, aggregateStore)
	products := es.NewAggregaterRepository[*domain.Product](domain.ProductAggregate,reg, aggregateStore)
	catalog := postgres.NewCatalogRepository("stores.products", mono.DB())
	mall := postgres.NewMallRepository("stores.stores", mono.DB())

	// setup app
	app := logging.LogApplicationAccess(
		application.New(stores, products, catalog, mall),
		mono.Logger(),
	)

	catalogHandlers := logging.LogEventHandlerAccess[ddd.AggregateEvent](
		application.NewCatalogHandlers(catalog),
		"catalog", mono.Logger(),
	)

	mallHandlers := logging.LogEventHandlerAccess[ddd.AggregateEvent](
		application.NewMallHandlers(mall),
		"Mall", mono.Logger(),
	)

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
	handlers.RegisterCatalogHandlers(catalogHandlers, domainDispatcher)
	handlers.RegisterCatalogHandlers(mallHandlers, domainDispatcher)

	return nil
}

func registrations(reg registry.Registry) (err error){
	serde := serdes.NewJsonSerde(reg)

	// Store
	if err = serde.Register(domain.StoreCreated{}); err != nil{
		return
	}

	if err = serde.RegisterKey(domain.StoreParticipationEnabledEvent, domain.StoreParticipationToggled{}); err != nil{
		return 
	}
	if err = serde.RegisterKey(domain.StoreParticipationDisabledEvent, domain.StoreParticipationToggled{}); err != nil{
		return
	}
	if err = serde.Register(domain.StoreRebranded{}); err != nil{
		return
	}
	// store snapshots
	if err = serde.RegisterKey(domain.StoreV1{}.SnapshotName(), domain.StoreV1{}); err != nil{
		return 
	}

	// Product
	if err = serde.Register(domain.Product{}, func(v any) error{
		store := v.(*domain.Product)
		store.Aggregate = es.NewAggregate("", domain.ProductAggregate)
		return nil
	}); err != nil{
		return
	}

	// Product events
	if err = serde.Register(domain.ProductAdded{}); err != nil{
		return 
	}

	if err = serde.Register(domain.ProductRebranded{}); err != nil{
		return
	}
	if err = serde.RegisterKey(domain.ProductPriceIncreasedEvent, domain.ProductPriceChanged{}); err != nil{
		return
	}
	if err = serde.RegisterKey(domain.ProductPriceDecreasedEvent, domain.ProductPriceChanged{}); err != nil{
		return 
	}
	if err = serde.Register(domain.ProductRemoved{}); err != nil{
		return
	}

	// product snapshots
	if err = serde.RegisterKey(domain.ProductV1{}.SnapshotName(), domain.ProductV1{}); err != nil{
		return
	}

	return 
}