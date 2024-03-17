package notifications

import (
	"context"

	"github.com/irononet/mallbots/internal/monolith"
	"github.com/irononet/mallbots/notifications/internal/application"
	"github.com/irononet/mallbots/notifications/internal/grpc"
	"github.com/irononet/mallbots/notifications/internal/logging"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Monolith) error{
	// Setup Driven adapters
	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil{
		return err
	}
	customers := grpc.NewCustomerRepository(conn)

	// setup app
	var app application.App
	app = application.New(customers)
	app = logging.LogApplicationAccess(app, mono.Logger())

	// Setup Driver adapters
	if err := grpc.RegisterServer(ctx, app, mono.RPC()); err != nil{
		return err
	}

	return nil
}