package service

import (
	"context"
	"log/slog"

	"github.com/libtnb/fiber-skeleton/internal/order/biz"
	"github.com/libtnb/fiber-skeleton/internal/shared/event"
)

// NewOrderPlacedLogger logs placed orders; replace with a real subscriber.
func NewOrderPlacedLogger(bus event.Bus, log *slog.Logger) event.Subscription {
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

	return event.Subscription{}
}
