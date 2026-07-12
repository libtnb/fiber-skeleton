package service

import (
	"context"
	"log/slog"

	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/order/biz"
	"github.com/libtnb/fiber-skeleton/internal/pkg/event"
)

// NewOrderPlacedLogger subscribes a handler that records placed orders — a
// stand-in for sending a confirmation or updating analytics, decoupled from
// the order flow. Building it (the app collects subscribers at startup)
// registers the handler on the bus.
func NewOrderPlacedLogger(i do.Injector) (event.Subscription, error) {
	bus := do.MustInvoke[event.Bus](i)
	log := do.MustInvoke[*slog.Logger](i)

	bus.Subscribe(biz.OrderPlaced{}.Name(), func(ctx context.Context, e event.Event) error {
		placed, ok := e.(biz.OrderPlaced)
		if !ok {
			return nil
		}
		log.InfoContext(ctx, "order placed",
			slog.Uint64("order_id", uint64(placed.OrderID)),
			slog.Uint64("user_id", uint64(placed.UserID)),
			slog.Int64("amount", placed.Amount),
		)
		return nil
	})

	return event.Subscription{}, nil
}
