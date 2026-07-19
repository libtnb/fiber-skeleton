package biz_test

import (
	"testing"

	"github.com/go-rio/rio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/libtnb/fiber-skeleton/internal/user/biz"
	mocksbiz "github.com/libtnb/fiber-skeleton/mocks/user/biz"
)

func TestUserUsecase_Create(t *testing.T) {
	repo := mocksbiz.NewUserRepo(t)
	repo.EXPECT().ExistsName(mock.Anything, "alice").Return(false, nil)
	repo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(u *biz.User) bool {
		return u.Name == "alice"
	})).Return(nil)

	user, err := biz.NewUserUsecase(repo).Create(t.Context(), "alice")

	require.NoError(t, err)
	assert.Equal(t, "alice", user.Name)
}

func TestUserUsecase_Create_NameTaken(t *testing.T) {
	repo := mocksbiz.NewUserRepo(t) // no Create expectation: it must not be called
	repo.EXPECT().ExistsName(mock.Anything, "alice").Return(true, nil)

	_, err := biz.NewUserUsecase(repo).Create(t.Context(), "alice")

	assert.ErrorIs(t, err, biz.ErrNameTaken)
}

func TestUserUsecase_Get_NotFound(t *testing.T) {
	repo := mocksbiz.NewUserRepo(t)
	repo.EXPECT().Get(mock.Anything, uint(9)).Return(nil, rio.ErrNotFound)

	_, err := biz.NewUserUsecase(repo).Get(t.Context(), 9)

	assert.ErrorIs(t, err, rio.ErrNotFound)
}

func TestUserUsecase_Update(t *testing.T) {
	repo := mocksbiz.NewUserRepo(t)
	repo.EXPECT().Update(mock.Anything, mock.MatchedBy(func(u *biz.User) bool {
		return u.ID == 1 && u.Name == "bob"
	})).Return(&biz.User{ID: 1, Name: "bob"}, nil)

	user, err := biz.NewUserUsecase(repo).Update(t.Context(), 1, "bob")

	require.NoError(t, err)
	assert.Equal(t, "bob", user.Name)
}
