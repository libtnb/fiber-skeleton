package data

import (
	"path/filepath"
	"testing"

	"github.com/go-rio/migrate"
	"github.com/go-rio/rio"
	"github.com/go-rio/sqlite"
	"github.com/libtnb/assert/must"

	"github.com/libtnb/fiber-skeleton/internal/migrations"
	"github.com/libtnb/fiber-skeleton/internal/order/biz"
)

func newTestRepo(t *testing.T) *orderRepo {
	t.Helper()

	db, err := sqlite.Open("file:" + filepath.Join(t.TempDir(), "test.db"))
	must.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	m, err := migrate.New(db.Unwrap(), migrate.SQLite, migrate.WithCollection(migrations.Collection()))
	must.NoError(t, err)
	must.NoError(t, m.Up(t.Context()))

	return &orderRepo{db: db}
}

func TestOrderRepo_CRUD(t *testing.T) {
	repo := newTestRepo(t)
	ctx := t.Context()

	order := &biz.Order{UserID: 7, Amount: 1299}
	must.NoError(t, repo.Create(ctx, order))
	must.NotZero(t, order.ID)
	must.False(t, order.CreatedAt.IsZero())

	got, err := repo.Get(ctx, order.ID)
	must.NoError(t, err)
	must.Equal(t, got.UserID, 7)
	must.Equal(t, got.Amount, 1299)

	list, total, err := repo.List(ctx, 1, 10)
	must.NoError(t, err)
	must.Equal(t, total, int64(1))
	must.Len(t, list, 1)

	must.NoError(t, repo.Delete(ctx, order.ID))
	_, err = repo.Get(ctx, order.ID)
	must.ErrorIs(t, err, rio.ErrNotFound)
}
