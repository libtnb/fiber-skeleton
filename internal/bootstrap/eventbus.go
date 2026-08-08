package bootstrap

import (
	"context"
	"log/slog"
	"sync"

	"github.com/libtnb/fiber-skeleton/internal/pkg/event"
)

// inProcessBus dispatches synchronously; handler errors are logged, never
// returned — an event is a fact, not a rejectable request.
type inProcessBus struct {
	mu       sync.RWMutex
	log      *slog.Logger
	handlers map[string][]event.Handler
}

func NewBus(log *slog.Logger) event.Bus {
	return &inProcessBus{
		log:      log,
		handlers: make(map[string][]event.Handler),
	}
}

func (b *inProcessBus) Subscribe(name string, h event.Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[name] = append(b.handlers[name], h)
}

func (b *inProcessBus) Publish(ctx context.Context, e event.Event) error {
	b.mu.RLock()
	handlers := b.handlers[e.Name()]
	b.mu.RUnlock()

	for _, h := range handlers {
		if err := h(ctx, e); err != nil {
			b.log.ErrorContext(ctx, "event handler failed",
				slog.String("event", e.Name()),
				slog.Any("err", err),
			)
		}
	}

	return nil
}
