package logging

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/irononet/mallbots/baskets/internal/application"
	"github.com/irononet/mallbots/internal/ddd"
)

type EventHandlers[T ddd.Event] struct {
	ddd.EventHandler[T]
	label  string
	logger zerolog.Logger
}

var _ ddd.EventHandler[ddd.Event] = (*EventHandlers[ddd.Event])(nil)

func LogEventHandlerAccess[T ddd.Event](handlers ddd.EventHandler[T], label string, logger zerolog.Logger) EventHandlers[T] {
	return EventHandlers[T]{
		EventHandler: handlers,
		label:        label,
		logger:       logger,
	}
}

func (h EventHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	h.logger.Info().Msgf("--> Baskets.%s.On(%s)", h.label, event.EventName())
	defer func() { h.logger.Info().Err(err).Msgf("<-- Baskets.%s.On(%s)", h.label, event.EventName()) }()
	return h.EventHandler.HandleEvent(ctx, event)
}

type DomainEventHandlers struct {
	application.DomainEventHandlers
	logger zerolog.Logger
}

var _ application.DomainEventHandlers = (*DomainEventHandlers)(nil)

func LogDomainEventHandlerAccess(handlers application.DomainEventHandlers, logger zerolog.Logger) DomainEventHandlers {
	return DomainEventHandlers{
		DomainEventHandlers: handlers,
		logger:              logger,
	}
}

func (h DomainEventHandlers) OnBasketStarted(ctx context.Context, event ddd.Event) (err error) {
	h.logger.Info().Msg("--> Baskets.OnBasketStarted")
	defer func() { h.logger.Info().Err(err).Msg("<-- Baskets.OnBasketStarted") }()
	return h.DomainEventHandlers.OnBasketStarted(ctx, event)
}

func (h DomainEventHandlers) OnBasketItemAdded(ctx context.Context, event ddd.Event) (err error) {
	h.logger.Info().Msg("--> Baskets.OnBasketItemAdded")
	defer func() { h.logger.Info().Err(err).Msg("<-- Baskets.OnBasketItemAdded") }()
	return h.DomainEventHandlers.OnBasketItemAdded(ctx, event)
}

func (h DomainEventHandlers) OnBasketItemRemoved(ctx context.Context, event ddd.Event) (err error) {
	h.logger.Info().Msg("--> Baskets.OnBasketItemRemoved")
	defer func() { h.logger.Info().Err(err).Msg("<-- Baskets.OnBasketItemRemoved") }()
	return h.DomainEventHandlers.OnBasketItemRemoved(ctx, event)
}

func (h DomainEventHandlers) OnBasketCanceled(ctx context.Context, event ddd.Event) (err error) {
	h.logger.Info().Msg("--> Baskets.OnBasketCanceled")
	defer func() { h.logger.Info().Err(err).Msg("<-- Baskets.OnBasketCanceled") }()
	return h.DomainEventHandlers.OnBasketCanceled(ctx, event)
}

func (h DomainEventHandlers) OnBasketCheckedOut(ctx context.Context, event ddd.Event) (err error) {
	h.logger.Info().Msg("--> Baskets.OnBasketCheckedOut")
	defer func() { h.logger.Info().Err(err).Msg("<-- Baskets.OnBasketCheckedOut") }()
	return h.DomainEventHandlers.OnBasketCheckedOut(ctx, event)
}
