// Package biz holds the user module's transport-independent business logic and
// its persistence boundary. Other modules depend on this package's interfaces,
// never on its data layer or tables.
package biz

import (
	"context"
	"errors"
	"time"

	"github.com/samber/oops"
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
	List(ctx context.Context, page, limit uint) ([]*User, int64, error)
	Get(ctx context.Context, id uint) (*User, error)
	ExistsName(ctx context.Context, name string) (bool, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, id uint) error
}

// UserUsecase holds transport-independent business logic shared by HTTP,
// CLI and jobs; methods take domain parameters, not request DTOs.
type UserUsecase struct {
	repo UserRepo
}

func NewUserUsecase(repo UserRepo) *UserUsecase {
	return &UserUsecase{
		repo: repo,
	}
}

func (uc *UserUsecase) List(ctx context.Context, page, limit uint) ([]*User, int64, error) {
	return uc.repo.List(ctx, page, limit)
}

func (uc *UserUsecase) Get(ctx context.Context, id uint) (*User, error) {
	return uc.repo.Get(ctx, id)
}

func (uc *UserUsecase) Create(ctx context.Context, name string) (*User, error) {
	// uniqueness is a business invariant, enforced here rather than by a
	// database-backed validation rule that would reach across the data boundary
	taken, err := uc.repo.ExistsName(ctx, name)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, oops.In("user").Code("user.name_taken").Public("name already taken").Wrap(ErrNameTaken)
	}

	user := &User{Name: name}
	if err := uc.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUsecase) Update(ctx context.Context, id uint, name string) (*User, error) {
	return uc.repo.Update(ctx, &User{ID: id, Name: name})
}

func (uc *UserUsecase) Delete(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}
