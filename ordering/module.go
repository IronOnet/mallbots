package ordering

import (
	"context"

	"github.com/irononet/mallbots/internal/ddd"
	"github.com/irononet/mallbots/internal/es"
	"github.com/irononet/mallbots/internal/monolith"
	"github.com/irononet/mallbots/internal/registry"
	"github.com/irononet/mallbots/internal/registry/serdes"
	"github.com/irononet/mallbots/ordering/internal/application"
	"github.com/irononet/mallbots/ordering/internal/domain"
	"github.com/irononet/mallbots/ordering/internal/grpc"
	"github.com/irononet/mallbots/ordering/internal/handlers"
	"github.com/irononet/mallbots/ordering/internal/logging"
	"github.com/irononet/mallbots/ordering/internal/rest"
	pg "github.com/irononet/mallbots/internal/postgres"
)

type Module struct{}

func (Module) Startup(ctx context.Context, mono monolith.Monolith) error{
	// Setup Driven adapters
	reg := registry.New()
	err := registrations(reg)
	if err != nil{
		return err
	}
	domainDispatcher := ddd.NewEventDispatcher[ddd.AggregateEvent]()
	aggregateStore := es.AggregateStoreWithMiddleware(
		pg.NewEventStore("ordering.events", mono.DB(), reg),
		es.NewEventPublisher(domainDispatcher),
		pg.NewSnapshotStore("ordering.snapshots", mono.DB(), reg),
	)
	orders := es.NewAggregaterRepository[*domain.Order](domain.OrderAggregate, reg, aggregateStore)
	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())

	customers := grpc.NewCustomerRepository(conn)
	payments := grpc.NewPaymentRepository(conn) 
	invoices := grpc.NewInvoiceRepository(conn)
	shopping := grpc.NewShoppingListRepository(conn)
	notifications := grpc.NewNotificationRepository(conn)

	// setup application 
	var app application.App
	app = application.New(orders, customers, payments, shopping)
	app = logging.LogApplicationAccess(app, mono.Logger())
	// setup application handlers
	notificationHandlers := logging.LogEventHandlerAcces(
		application.NewNotificationHandlers(notifications),
		"Notification",
		mono.Logger(),
	)
	invoiceHandlers := logging.LogEventHandlerAcces(
		application.NewInvoiceHandlers(invoices),
		"Invoice",
		mono.Logger(),
	)

	// setup Driver adapters
	if err := grpc.RegisterServer(app, mono.RPC()); err != nil{
		return err
	}
	if err := rest.RegisterGateway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil{
		return err
	}
	if err := rest.RegisterSwagger(mono.Mux()); err != nil{
		return err
	}
	handlers.RegisterNotificationHandler(notificationHandlers, domainDispatcher)
	handlers.RegisterInvoiceHandlers(invoiceHandlers, domainDispatcher)

	return nil
}

func registrations(reg registry.Registry) error{
	serde := serdes.NewJsonSerde(reg)

	// Order
	if err := serde.Register(domain.Order{}, func(v any) error{
		order := v.(*domain.Order)
		order.Aggregate = es.NewAggregate("", domain.OrderAggregate)
		return nil
	}); err != nil{
		return err
	}
	// Order events
	if err := serde.Register(domain.OrderCreated{}); err != nil{
		return err
	}
	if err := serde.Register(domain.OrderCanceled{}); err != nil{
		return err
	}
	if err := serde.Register(domain.OrderReadied{}); err != nil{
		return err
	}
	if err := serde.Register(domain.OrderCompleted{}); err != nil{
		return err
	}
	// order snapshots
	if err := serde.RegisterKey(domain.OrderV1{}.SnapshotName(), domain.OrderV1{}); err != nil{
		return err
	}
	return nil
}