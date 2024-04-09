package application

import (
	"context"

	"github.com/irononet/mallbots/stores/internal/application/commands"
	"github.com/irononet/mallbots/stores/internal/application/queries"
	"github.com/irononet/mallbots/stores/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}

	Commands interface {
		CreateStore(ctx context.Context, cmd commands.CreateStore) error
		EnableParticipation(ctx context.Context, cmd commands.EnableParticipation) error
		DisableParticipation(ctx context.Context, cmd commands.DisableParticipation) error
		AddProduct(ctx context.Context, cmd commands.AddProduct) error
		RebrandProduct(ctx context.Context, cmd commands.RebrandProduct) error
		RebrandStore(ctx context.Context, cmd commands.RebrandStore) error
		IncreaseProductPrice(ctx context.Context, cmd commands.IncreaseProductPrice) error
		DecreaseProductPrice(ctx context.Context, cmd commands.DecreaseProductPrice) error
		RemoveProduct(ctx context.Context, cmd commands.RemoveProduct) error
	}

	Queries interface {
		GetStore(ctx context.Context, query queries.GetStore) (*domain.MallStore, error)
		GetStores(ctx context.Context, query queries.GetStores) ([]*domain.MallStore, error)
		GetParticipatingStores(ctx context.Context, query queries.GetParticipatingStores) ([]*domain.MallStore, error)
		GetCatalog(ctx context.Context, query queries.GetCatalog) ([]*domain.CatalogProduct, error)
		GetProduct(ctx context.Context, query queries.GetProduct) (*domain.CatalogProduct, error)
	}

	Application struct {
		appCommands
		appQueries
	}

	appCommands struct {
		commands.CreateStoreHandler
		commands.EnableParticipationHandler
		commands.DisableParticipationHandler
		commands.RebrandStoreHandler
		commands.RebrandProductHandler
		commands.IncreaseProductPriceHandler
		commands.DecreaseProductPriceHandler
		commands.AddProductHandler
		commands.RemoveProductHandler
	}

	appQueries struct {
		queries.GetStoreHandler
		queries.GetStoresHandler
		queries.GetParticipatingStoresHandler
		queries.GetCatalogHandler
		queries.GetProductHandler
	}
)

var _ App = (*Application)(nil)

func New(stores domain.StoreRepository, products domain.ProductRepository,
	catalog domain.CatalogRepository, mall domain.MallRepository,
) *Application {
	return &Application{
		appCommands: appCommands{
			CreateStoreHandler:          commands.NewCreateStoreHandler(stores),
			EnableParticipationHandler:  commands.NewEnableParticipationHanler(stores),
			DisableParticipationHandler: commands.NewDisableParticipationHandler(stores),
			RebrandStoreHandler: commands.NewRebrandStoreHandler(stores),
			AddProductHandler:           commands.NewAddProductHandler(products),
			RebrandProductHandler: commands.NewRebrandProductHandler(products),
			IncreaseProductPriceHandler: commands.NewIncreaseProductPriceHandler(products),
			DecreaseProductPriceHandler: commands.NewDecreaseProductPriceHandler(products),
			RemoveProductHandler:        commands.NewRemoveProductHandler(products),
		},

		appQueries: appQueries{
			GetStoreHandler:               queries.NewGetStoreHandler(mall),
			GetStoresHandler:              queries.NewGetStoresHandler(mall),
			GetParticipatingStoresHandler: queries.NewGetParticipatingStoresHandler(mall),
			GetCatalogHandler:             queries.NewGetCatalogHandler(catalog),
			GetProductHandler:             queries.NewGetProductHandler(catalog),
		},
	}
}
