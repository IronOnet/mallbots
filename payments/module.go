package payments

import (
	"context"

	"github.com/irononet/mallbots/internal/monolith"
	"github.com/irononet/mallbots/payments/internal/application"
	"github.com/irononet/mallbots/payments/internal/grpc"
	"github.com/irononet/mallbots/payments/internal/logging"
	"github.com/irononet/mallbots/payments/internal/postgres"
	"github.com/irononet/mallbots/payments/internal/rest"
)

type Module struct{}

func (Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	//Setup driven adapters
	invoices := postgres.NewInvoiceRepository("payments.invoices", mono.DB())
	payments := postgres.NewPaymentRepository("payments.payments", mono.DB())
	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}
	orders := grpc.NewOrderRepository(conn)

	// Setup application
	var app application.App
	app = application.New(invoices, payments, orders)
	app = logging.LogApplicationAccess(app, mono.Logger())

	// setup Driver adapters
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
