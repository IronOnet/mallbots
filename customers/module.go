package customers

import (
	"context"

	"github.com/irononet/mallbots/customers/internal/application"
	"github.com/irononet/mallbots/customers/internal/grpc"
	"github.com/irononet/mallbots/customers/internal/logging"
	"github.com/irononet/mallbots/customers/internal/postgres"
	"github.com/irononet/mallbots/customers/internal/rest"
	"github.com/irononet/mallbots/internal/ddd"
	"github.com/irononet/mallbots/internal/monolith"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Monolith) error{
	// Setup Driven adapters
	domainDispatcher := ddd.NewEventDispatcher[ddd.AggregateEvent]()
	customers := postgres.NewCustomerRepository("customers.customers", mono.DB())

	// Setup application
	app := logging.LogApplicationAccess(
		application.New(customers, domainDispatcher),
		mono.Logger(),
	)

	if err := grpc.RegisterServer(app, mono.RPC()); err != nil{
		return err
	}

	if err := rest.RegisterGateway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil{
		return err
	}

	if err := rest.RegisterSwagger(mono.Mux()); err != nil{
		return err
	}

	return nil
}