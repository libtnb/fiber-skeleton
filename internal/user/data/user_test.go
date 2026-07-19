package data

import (
	"path/filepath"
	"testing"

	"github.com/go-rio/migrate"
	"github.com/go-rio/rio"
	"github.com/go-rio/sqlite"
	"github.com/stretchr/testify/require"

	"github.com/libtnb/fiber-skeleton/internal/user/biz"
)

// newTestRepo binds the repo to a throwaway, fully migrated SQLite database.
func newTestRepo(t *testing.T) *userRepo {
	t.Helper()

	db, err := sqlite.Open("file:" + filepath.Join(t.TempDir(), "test.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	m, err := migrate.New(db.Unwrap(), migrate.SQLite)
	require.NoError(t, err)
	require.NoError(t, m.Up(t.Context()))

	return &userRepo{db: db}
}

func TestUserRepo_CRUD(t *testing.T) {
	repo := newTestRepo(t)
	ctx := t.Context()

	user := &biz.User{Name: "alice"}
	require.NoError(t, repo.Create(ctx, user))
	require.NotZero(t, user.ID)
	require.False(t, user.CreatedAt.IsZero())

	got, err := repo.Get(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, "alice", got.Name)

	list, total, err := repo.List(ctx, 1, 10)
	require.NoError(t, err)
	require.Equal(t, int64(1), total)
	require.Len(t, list, 1)

	// Update changes the name but keeps CreatedAt.
	updated, err := repo.Update(ctx, &biz.User{ID: user.ID, Name: "bob"})
	require.NoError(t, err)
	require.Equal(t, "bob", updated.Name)
	require.Equal(t, user.CreatedAt.Unix(), updated.CreatedAt.Unix())
}

func TestUserRepo_Create_DuplicateName(t *testing.T) {
	repo := newTestRepo(t)
	ctx := t.Context()

	require.NoError(t, repo.Create(ctx, &biz.User{Name: "alice"}))
	require.ErrorIs(t, repo.Create(ctx, &biz.User{Name: "alice"}), biz.ErrNameTaken)
}

func TestUserRepo_Get_NotFound(t *testing.T) {
	repo := newTestRepo(t)

	_, err := repo.Get(t.Context(), 404)
	require.ErrorIs(t, err, rio.ErrNotFound)
}

func TestUserRepo_Delete_SoftDeletesAndReports(t *testing.T) {
	repo := newTestRepo(t)
	ctx := t.Context()

	user := &biz.User{Name: "carol"}
	require.NoError(t, repo.Create(ctx, user))

	// Delete soft-deletes: the row is gone from default reads.
	require.NoError(t, repo.Delete(ctx, user.ID))
	_, err := repo.Get(ctx, user.ID)
	require.ErrorIs(t, err, rio.ErrNotFound)

	require.ErrorIs(t, repo.Delete(ctx, user.ID), rio.ErrNotFound)

	// soft-deleting released the name
	require.NoError(t, repo.Create(ctx, &biz.User{Name: "carol"}))
}
