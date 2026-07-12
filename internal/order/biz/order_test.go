package biz_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/libtnb/fiber-skeleton/internal/order/biz"
	"github.com/libtnb/fiber-skeleton/internal/pkg/event"
	mocksbiz "github.com/libtnb/fiber-skeleton/mocks/biz"
)

// fakeBus records what the usecase publishes, without a real bus.
type fakeBus struct {
	published []event.Event
}

func (b *fakeBus) Subscribe(string, event.Handler) {}
func (b *fakeBus) Publish(_ context.Context, e event.Event) error {
	b.published = append(b.published, e)
	return nil
}

func TestOrderUsecase_Place(t *testing.T) {
	repo := mocksbiz.NewOrderRepo(t)
	users := mocksbiz.NewUsers(t)
	bus := &fakeBus{}

	users.EXPECT().Exists(mock.Anything, uint(1)).Return(true, nil)
	repo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(o *biz.Order) bool {
		return o.UserID == 1 && o.Amount == 500
	})).Return(nil)

	order, err := biz.NewOrderUsecase(repo, users, bus).Place(context.Background(), 1, 500)

	require.NoError(t, err)
	require.Equal(t, uint(1), order.UserID)

	require.Len(t, bus.published, 1)
	placed, ok := bus.published[0].(biz.OrderPlaced)
	require.True(t, ok)
	require.Equal(t, uint(1), placed.UserID)
	require.EqualValues(t, 500, placed.Amount)
}

func TestOrderUsecase_Place_UnknownUser(t *testing.T) {
	repo := mocksbiz.NewOrderRepo(t) // no Create expectation: it must not be called
	users := mocksbiz.NewUsers(t)
	bus := &fakeBus{}

	users.EXPECT().Exists(mock.Anything, uint(9)).Return(false, nil)

	_, err := biz.NewOrderUsecase(repo, users, bus).Place(context.Background(), 9, 500)

	require.ErrorIs(t, err, biz.ErrUserNotFound)
	require.Empty(t, bus.published)
}
