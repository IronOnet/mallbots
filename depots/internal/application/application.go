package application

import (
	"context"

	"github.com/irononet/mallbots/depots/internal/application/queries"
	"github.com/irononet/mallbots/depots/internal/application/commands"
	"github.com/irononet/mallbots/depots/internal/domain"
	"github.com/irononet/mallbots/internal/ddd"
)

type (
	App interface{
		Commands
		Queries
	}

	Commands interface{
		CreateShoppingList(ctx context.Context, cmd commands.CreateShoppingList) error
		CancelShoppingList(ctx context.Context, cmd commands.CancelShoppingList) error
		AssignShoppingList(ctx context.Context, cmd commands.AssignShoppingList) error
		CompleteShoppingList(ctx context.Context, cmd commands.CompleteShoppingList) error
	}

	Queries interface{
		GetShoppingList(ctx context.Context, query queries.GetShoppingList) (*domain.ShoppingList, error)
	}

	Application struct{
		appCommands
		appQueries
	}

	appCommands struct{
		commands.CreateShoppingListHandler
		commands.CancelShoppingListHandler
		commands.AssignShoppingListHandler
		commands.CompleteShoppingListHandler
	}

	appQueries struct{
		queries.GetShoppingListHandler
	}
)

var _ App = (*Application)(nil)

func New(shoppingLists domain.ShoppingListRepository, stores domain.StoreRepository, products domain.ProductRepository, domainPublisher ddd.EventPublisher[ddd.AggregateEvent]) *Application{
	return &Application{
		appCommands: appCommands{
			CreateShoppingListHandler: commands.NewCreateShoppingListHandler(shoppingLists, stores, products, domainPublisher),
			CancelShoppingListHandler: commands.NewCancelShoppingListHandler(shoppingLists, domainPublisher),
			AssignShoppingListHandler: commands.NewAssignShoppingListHandler(shoppingLists, domainPublisher),
			CompleteShoppingListHandler: commands.NewCompleteShoppingListHandler(shoppingLists, domainPublisher),
		},

		appQueries: appQueries{
			GetShoppingListHandler: queries.NewGetShoppingListHandler(shoppingLists),
		},
	}
}