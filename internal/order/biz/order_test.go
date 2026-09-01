package biz_test

import (
	"context"
	"testing"

	"github.com/libtnb/assert/must"

	mocksbiz "github.com/libtnb/fiber-skeleton/internal/mocks/order/biz"
	"github.com/libtnb/fiber-skeleton/internal/order/biz"
	"github.com/libtnb/fiber-skeleton/internal/shared/event"
)

// fakeBus records what the usecase publishes.
type fakeBus struct {
	published []event.Event
}

func (b *fakeBus) Subscribe(string, event.Handler) {}
func (b *fakeBus) Publish(_ context.Context, e event.Event) error {
	b.published = append(b.published, e)
	return nil
}

func TestOrderUsecase_Place(t *testing.T) {
	repo := &mocksbiz.OrderRepoMock{
		CreateFunc: func(context.Context, *biz.Order) error { return nil },
	}
	users := &mocksbiz.UsersMock{
		ExistsFunc: func(context.Context, uint) (bool, error) { return true, nil },
	}
	bus := &fakeBus{}

	order, err := biz.NewOrderUsecase(repo, users, bus).Place(t.Context(), 1, 500)

	must.NoError(t, err)
	must.Equal(t, order.UserID, uint(1))

	created := repo.CreateCalls()
	must.Len(t, created, 1)
	must.Equal(t, created[0].Order.UserID, uint(1))
	must.Equal(t, created[0].Order.Amount, 500)

	must.Len(t, bus.published, 1)
	placed, ok := bus.published[0].(biz.OrderPlaced)
	must.True(t, ok)
	must.Equal(t, placed.UserID, uint(1))
	must.Equal(t, placed.Amount, 500)
}

func TestOrderUsecase_Place_UnknownUser(t *testing.T) {
	// CreateFunc stays nil: a Create call would panic the test
	repo := &mocksbiz.OrderRepoMock{}
	users := &mocksbiz.UsersMock{
		ExistsFunc: func(context.Context, uint) (bool, error) { return false, nil },
	}
	bus := &fakeBus{}

	_, err := biz.NewOrderUsecase(repo, users, bus).Place(t.Context(), 9, 500)

	must.ErrorIs(t, err, biz.ErrUserNotFound)
	must.Empty(t, bus.published)
}
