package commands

import (
	"context"

	"github.com/irononet/mallbots/stores/internal/domain"
)

type DisableParticipation struct {
	ID string
}

type DisableParticipationHandler struct{
	stores domain.StoreRepository
}

func NewDisableParticipationHandler(stores domain.StoreRepository,) DisableParticipationHandler{
	return DisableParticipationHandler{
		stores: stores, 
	}
}

func (h DisableParticipationHandler) DisableParticipation(ctx context.Context, cmd DisableParticipation) error{
	store, err := h.stores.Load(ctx, cmd.ID)
	if err != nil{
		return err
	}

	if err = store.DisableParticipation(); err != nil{
		return err
	}

	if err = h.stores.Save(ctx, store); err != nil{
		return err
	}

	return nil
}