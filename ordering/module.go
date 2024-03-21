package ordering

import (
	"context"

	"github.com/irononet/mallbots/internal/ddd"
	"github.com/irononet/mallbots/internal/monolith"
	"github.com/irononet/mallbots/ordering/internal/application"
	"github.com/irononet/mallbots/ordering/internal/grpc"
	"github.com/irononet/mallbots/ordering/internal/handlers"
	"github.com/irononet/mallbots/ordering/internal/logging"
	"github.com/irononet/mallbots/ordering/internal/postgres"
	"github.com/irononet/mallbots/ordering/internal/rest"
)

type Module struct{}

func (Module) Startup(ctx context.Context, mono monolith.Monolith) error{
	// Setup Driven adapters
	domainDispatcher := ddd.NewEventDispatcher()
	orders := postgres.NewOrderRepository("ordering.orders", mono.DB())
	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil{
		return err
	}
	customers := grpc.NewCustomerRepository(conn)
	payments := grpc.NewPaymentRepository(conn) 
	invoices := grpc.NewInvoiceRepository(conn)
	shopping := grpc.NewShoppingListRepository(conn)
	notifications := grpc.NewNotificationRepository(conn)

	// setup application 
	var app application.App
	app = application.New(orders, customers, payments, shopping, domainDispatcher)
	app = logging.LogApplicationAccess(app, mono.Logger())
	// setup application handlers
	notificationHandlers := logging.LogDomainEventHandlerAccess(
		application.NewNotificationHandlers(notifications),
		mono.Logger(),
	)
	invoiceHandlers := logging.LogDomainEventHandlerAccess(
		application.NewInvoiceHandlers(invoices),
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