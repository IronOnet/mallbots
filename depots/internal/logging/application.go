package logging

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/irononet/mallbots/depots/internal/application"
	"github.com/irononet/mallbots/depots/internal/application/queries"
	"github.com/irononet/mallbots/depots/internal/application/commands"
	"github.com/irononet/mallbots/depots/internal/domain"
)

type Application struct{
	application.App
	logger zerolog.Logger
}

var _ application.App = (*Application)(nil) 

func LogApplicationAccess(application application.App, logger zerolog.Logger) Application{
	return Application{
		App: application,
		logger: logger,
	}
}

func (a Application) CreateShoppingList(ctx context.Context, cmd commands.CreateShoppingList) (err error){
	a.logger.Info().Msg("--> Depot.CreateShoppingList")
	defer func() { a.logger.Info().Err(err).Msg("<-- Depot.CancelShoppingList")}()
	return a.App.CreateShoppingList(ctx, cmd)
}

func (a Application) CancelShoppingList(ctx context.Context, cmd commands.CancelShoppingList) (err error){
	a.logger.Info().Msg("--> Depot.CancelShoppingList")
	defer func() { a.logger.Info().Err(err).Msg("<-- Depot.CancelShoppingList")}()
	return a.App.CancelShoppingList(ctx, cmd)
}

func (a Application) CompleteShoppingList(ctx context.Context, cmd commands.CompleteShoppingList) (err error){
	a.logger.Info().Msg("--> Depot.CompleteShoppingList")
	defer func() { a.logger.Info().Err(err).Msg("<-- Depot.CompleteShoppingList")}()
	return a.App.CompleteShoppingList(ctx, cmd)
}

func (a Application) GetShoppingList(ctx context.Context, query queries.GetShoppingList) (list *domain.ShoppingList, err error){
	a.logger.Info().Msg("--> Depot.GetShoppingList")
	defer func() { a.logger.Info().Err(err).Msg("<-- Depot.GetShoppingList")}()
	return a.App.GetShoppingList(ctx, query)
}