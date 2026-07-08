package biz

import (
	"context"

	"github.com/dromara/carbon/v2"
	"gorm.io/gorm"
)

type User struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Name      string          `json:"name"`
	CreatedAt carbon.DateTime `json:"created_at"`
	UpdatedAt carbon.DateTime `json:"updated_at"`
	DeletedAt gorm.DeletedAt  `json:"-"`
}

// UserRepo is the persistence boundary; a missing row is ErrNotFound.
type UserRepo interface {
	List(ctx context.Context, page, limit uint) ([]*User, int64, error)
	Get(ctx context.Context, id uint) (*User, error)
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
