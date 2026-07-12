package data

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/go-rio/migrate"
	"github.com/go-rio/rio"
	"github.com/go-rio/sqlite"
	"github.com/stretchr/testify/require"

	"github.com/libtnb/fiber-skeleton/internal/order/biz"
)

func newTestRepo(t *testing.T) *orderRepo {
	t.Helper()

	db, err := sqlite.Open("file:" + filepath.Join(t.TempDir(), "test.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	m, err := migrate.New(db.Unwrap(), migrate.SQLite)
	require.NoError(t, err)
	require.NoError(t, m.Up(t.Context()))

	return &orderRepo{db: db}
}

func TestOrderRepo_CRUD(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	order := &biz.Order{UserID: 7, Amount: 1299}
	require.NoError(t, repo.Create(ctx, order))
	require.NotZero(t, order.ID)
	require.False(t, order.CreatedAt.IsZero())

	got, err := repo.Get(ctx, order.ID)
	require.NoError(t, err)
	require.EqualValues(t, 7, got.UserID)
	require.EqualValues(t, 1299, got.Amount)

	list, total, err := repo.List(ctx, 1, 10)
	require.NoError(t, err)
	require.Equal(t, int64(1), total)
	require.Len(t, list, 1)

	require.NoError(t, repo.Delete(ctx, order.ID))
	_, err = repo.Get(ctx, order.ID)
	require.ErrorIs(t, err, rio.ErrNotFound)
}
