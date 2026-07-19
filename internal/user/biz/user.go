// Package biz holds the user module's transport-independent business logic;
// other modules depend on it, never on the data layer.
package biz

import (
	"context"
	"errors"
	"time"

	"github.com/libtnb/fiber-skeleton/internal/pkg/apperr"
)

type User struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `rio:",softdelete" json:"-"`
}

var ErrNameTaken = errors.New("name already taken")

// UserRepo is the persistence boundary; a missing row is rio.ErrNotFound.
type UserRepo interface {
	List(ctx context.Context, page, limit int) ([]*User, int64, error)
	Get(ctx context.Context, id uint) (*User, error)
	ExistsName(ctx context.Context, name string) (bool, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, id uint) error
}

// UserUsecase is shared by HTTP, CLI and jobs; methods take domain
// parameters, not request DTOs.
type UserUsecase struct {
	repo UserRepo
}

func NewUserUsecase(repo UserRepo) *UserUsecase {
	return &UserUsecase{
		repo: repo,
	}
}

func (uc *UserUsecase) List(ctx context.Context, page, limit int) ([]*User, int64, error) {
	return uc.repo.List(ctx, page, limit)
}

func (uc *UserUsecase) Get(ctx context.Context, id uint) (*User, error) {
	return uc.repo.Get(ctx, id)
}

func (uc *UserUsecase) Create(ctx context.Context, name string) (*User, error) {
	// the unique index is the real guarantee; the pre-check just answers earlier
	taken, err := uc.repo.ExistsName(ctx, name)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, errNameTaken()
	}

	user := &User{Name: name}
	if err := uc.repo.Create(ctx, user); err != nil {
		// lost the pre-check race
		if errors.Is(err, ErrNameTaken) {
			return nil, errNameTaken()
		}
		return nil, err
	}

	return user, nil
}

func errNameTaken() error {
	return apperr.Conflict("user.name_taken", "name already taken").In("user").Wrap(ErrNameTaken)
}

func (uc *UserUsecase) Update(ctx context.Context, id uint, name string) (*User, error) {
	return uc.repo.Update(ctx, &User{ID: id, Name: name})
}

func (uc *UserUsecase) Delete(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}
