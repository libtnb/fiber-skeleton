package data

import (
	"path/filepath"
	"testing"

	"github.com/go-rio/migrate"
	"github.com/go-rio/rio"
	"github.com/go-rio/sqlite"
	"github.com/libtnb/assert/must"

	"github.com/libtnb/fiber-skeleton/internal/migrations"
	"github.com/libtnb/fiber-skeleton/internal/user/biz"
)

// newTestRepo binds the repo to a throwaway, fully migrated SQLite database.
func newTestRepo(t *testing.T) *userRepo {
	t.Helper()

	db, err := sqlite.Open("file:" + filepath.Join(t.TempDir(), "test.db"))
	must.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	m, err := migrate.New(db.Unwrap(), migrate.SQLite, migrate.WithCollection(migrations.Collection()))
	must.NoError(t, err)
	must.NoError(t, m.Up(t.Context()))

	return &userRepo{db: db}
}

func TestUserRepo_CRUD(t *testing.T) {
	repo := newTestRepo(t)
	ctx := t.Context()

	user := &biz.User{Name: "alice"}
	must.NoError(t, repo.Create(ctx, user))
	must.NotZero(t, user.ID)
	must.False(t, user.CreatedAt.IsZero())

	got, err := repo.Get(ctx, user.ID)
	must.NoError(t, err)
	must.Equal(t, got.Name, "alice")

	list, total, err := repo.List(ctx, 1, 10)
	must.NoError(t, err)
	must.Equal(t, total, int64(1))
	must.Len(t, list, 1)

	// Update changes the name but keeps CreatedAt.
	updated, err := repo.Update(ctx, &biz.User{ID: user.ID, Name: "bob"})
	must.NoError(t, err)
	must.Equal(t, updated.Name, "bob")
	must.Equal(t, updated.CreatedAt.Unix(), user.CreatedAt.Unix())
}

func TestUserRepo_Create_DuplicateName(t *testing.T) {
	repo := newTestRepo(t)
	ctx := t.Context()

	must.NoError(t, repo.Create(ctx, &biz.User{Name: "alice"}))
	must.ErrorIs(t, repo.Create(ctx, &biz.User{Name: "alice"}), biz.ErrNameTaken)
}

func TestUserRepo_Get_NotFound(t *testing.T) {
	repo := newTestRepo(t)

	_, err := repo.Get(t.Context(), 404)
	must.ErrorIs(t, err, rio.ErrNotFound)
}

func TestUserRepo_PredefinedQueriesUsePerCallArguments(t *testing.T) {
	repo := newTestRepo(t)
	ctx := t.Context()

	alice := &biz.User{Name: "alice"}
	bob := &biz.User{Name: "bob"}
	must.NoError(t, repo.Create(ctx, alice))
	must.NoError(t, repo.Create(ctx, bob))

	exists, err := repo.ExistsName(ctx, "alice")
	must.NoError(t, err)
	must.True(t, exists)
	exists, err = repo.ExistsName(ctx, "bob")
	must.NoError(t, err)
	must.True(t, exists)
	exists, err = repo.ExistsName(ctx, "missing")
	must.NoError(t, err)
	must.False(t, exists)

	must.NoError(t, repo.Delete(ctx, bob.ID))
	_, err = repo.Get(ctx, bob.ID)
	must.ErrorIs(t, err, rio.ErrNotFound)
	_, err = repo.Get(ctx, alice.ID)
	must.NoError(t, err)
}

func TestUserRepo_Delete_SoftDeletesAndReports(t *testing.T) {
	repo := newTestRepo(t)
	ctx := t.Context()

	user := &biz.User{Name: "carol"}
	must.NoError(t, repo.Create(ctx, user))

	// Delete soft-deletes: the row is gone from default reads.
	must.NoError(t, repo.Delete(ctx, user.ID))
	_, err := repo.Get(ctx, user.ID)
	must.ErrorIs(t, err, rio.ErrNotFound)

	must.ErrorIs(t, repo.Delete(ctx, user.ID), rio.ErrNotFound)

	// soft-deleting released the name
	must.NoError(t, repo.Create(ctx, &biz.User{Name: "carol"}))
}
