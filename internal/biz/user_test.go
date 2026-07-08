package biz_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/libtnb/fiber-skeleton/internal/biz"
	mocksbiz "github.com/libtnb/fiber-skeleton/mocks/biz"
)

// Usecase tests talk to a mocked repo directly — no HTTP, no binding — which
// is where business logic tests belong once it grows past CRUD.

func TestUserUsecase_Create(t *testing.T) {
	repo := mocksbiz.NewUserRepo(t)
	repo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(u *biz.User) bool {
		return u.Name == "alice"
	})).Return(nil)

	user, err := biz.NewUserUsecase(repo).Create(context.Background(), "alice")

	assert.NoError(t, err)
	assert.Equal(t, "alice", user.Name)
}

func TestUserUsecase_Get_NotFound(t *testing.T) {
	repo := mocksbiz.NewUserRepo(t)
	repo.EXPECT().Get(mock.Anything, uint(9)).Return(nil, biz.ErrNotFound)

	_, err := biz.NewUserUsecase(repo).Get(context.Background(), 9)

	assert.ErrorIs(t, err, biz.ErrNotFound)
}

func TestUserUsecase_Update(t *testing.T) {
	repo := mocksbiz.NewUserRepo(t)
	repo.EXPECT().Update(mock.Anything, mock.MatchedBy(func(u *biz.User) bool {
		return u.ID == 1 && u.Name == "bob"
	})).Return(&biz.User{ID: 1, Name: "bob"}, nil)

	user, err := biz.NewUserUsecase(repo).Update(context.Background(), 1, "bob")

	assert.NoError(t, err)
	assert.Equal(t, "bob", user.Name)
}
