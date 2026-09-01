package biz_test

import (
	"context"
	"testing"

	"github.com/go-rio/rio"
	"github.com/libtnb/assert/must"

	mocksbiz "github.com/libtnb/fiber-skeleton/internal/mocks/user/biz"
	"github.com/libtnb/fiber-skeleton/internal/user/biz"
)

func TestUserUsecase_Create(t *testing.T) {
	repo := &mocksbiz.UserRepoMock{
		ExistsNameFunc: func(context.Context, string) (bool, error) { return false, nil },
		CreateFunc:     func(context.Context, *biz.User) error { return nil },
	}

	user, err := biz.NewUserUsecase(repo).Create(t.Context(), "alice")

	must.NoError(t, err)
	must.Equal(t, user.Name, "alice")
	created := repo.CreateCalls()
	must.Len(t, created, 1)
	must.Equal(t, created[0].User.Name, "alice")
}

func TestUserUsecase_Create_NameTaken(t *testing.T) {
	// CreateFunc stays nil: a Create call would panic the test
	repo := &mocksbiz.UserRepoMock{
		ExistsNameFunc: func(context.Context, string) (bool, error) { return true, nil },
	}

	_, err := biz.NewUserUsecase(repo).Create(t.Context(), "alice")

	must.ErrorIs(t, err, biz.ErrNameTaken)
}

func TestUserUsecase_Get_NotFound(t *testing.T) {
	repo := &mocksbiz.UserRepoMock{
		GetFunc: func(context.Context, uint) (*biz.User, error) { return nil, rio.ErrNotFound },
	}

	_, err := biz.NewUserUsecase(repo).Get(t.Context(), 9)

	must.ErrorIs(t, err, rio.ErrNotFound)
}

func TestUserUsecase_Update(t *testing.T) {
	repo := &mocksbiz.UserRepoMock{
		UpdateFunc: func(_ context.Context, u *biz.User) (*biz.User, error) { return u, nil },
	}

	user, err := biz.NewUserUsecase(repo).Update(t.Context(), 1, "bob")

	must.NoError(t, err)
	must.Equal(t, user.Name, "bob")
	updated := repo.UpdateCalls()
	must.Len(t, updated, 1)
	must.Equal(t, updated[0].User.ID, uint(1))
}
