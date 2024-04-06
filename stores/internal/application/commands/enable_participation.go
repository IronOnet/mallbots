package commands

import (
	"context"

	"github.com/irononet/mallbots/stores/internal/domain"
)

type EnableParticipation struct{
	ID string
}

type EnableParticipationHandler struct{
	stores domain.StoreRepository
}

func NewEnableParticipationHanler(stores domain.StoreRepository) EnableParticipationHandler{
	return EnableParticipationHandler{
		stores: stores,
	}
}

func (h EnableParticipationHandler) EnableParticipation(ctx context.Context, cmd EnableParticipation) error{
	store, err := h.stores.Load(ctx, cmd.ID)
	if err != nil{
		return err
	}

	if err = store.EnableParticipation(); err != nil{
		return err
	}

	if err = h.stores.Save(ctx, store); err != nil{
		return err
	}

	return nil
}