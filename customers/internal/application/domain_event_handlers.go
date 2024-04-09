package application

import (
	"context"

	"github.com/irononet/mallbots/internal/ddd"
)

type DomainEventHandler interface{
	OnCustomerRegistered(ctx context.Context, event ddd.Event) error
	OnCustomerAuthorized(ctx context.Context, event ddd.Event) error
	OnCustomerEnabled(ctx context.Context, event ddd.Event) error
	OnCustomerDisabled(ctx context.Context, event ddd.Event) error
}

type ignoreUnimplementedDomainEvents struct{}

var _ DomainEventHandler = (*ignoreUnimplementedDomainEvents)(nil)

func (ignoreUnimplementedDomainEvents) OnCustomerRegistered(ctx context.Context, event ddd.Event) error{
	return nil
}

func (ignoreUnimplementedDomainEvents) OnCustomerAuthorized(ctx context.Context, event ddd.Event) error{
	return nil
}

func (ignoreUnimplementedDomainEvents) OnCustomerEnabled(ctx context.Context, event ddd.Event) error{
	return nil
}

func (ignoreUnimplementedDomainEvents) OnCustomerDisabled(ctx context.Context, event ddd.Event) error{
	return nil
}