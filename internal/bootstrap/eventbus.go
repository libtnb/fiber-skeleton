package bootstrap

import (
	"context"
	"log/slog"
	"sync"

	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/pkg/event"
)

// inProcessBus dispatches synchronously in the caller's goroutine. Handler
// errors are logged and do not abort the remaining handlers or fail Publish: a
// domain event is a fact that already happened, not a request that can be
// rejected. Swap this implementation for a Kafka/NATS one — publishers and
// subscribers depend on event.Bus, not on this type.
type inProcessBus struct {
	mu       sync.RWMutex
	log      *slog.Logger
	handlers map[string][]event.Handler
}

func NewBus(i do.Injector) (event.Bus, error) {
	return &inProcessBus{
		log:      do.MustInvoke[*slog.Logger](i),
		handlers: make(map[string][]event.Handler),
	}, nil
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
