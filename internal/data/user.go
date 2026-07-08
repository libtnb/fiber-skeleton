package data

import (
	"context"

	"github.com/samber/do/v2"
	"gorm.io/gorm"

	"github.com/libtnb/fiber-skeleton/internal/biz"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(i do.Injector) (biz.UserRepo, error) {
	return &userRepo{
		db: do.MustInvoke[*Data](i).DB,
	}, nil
}

func (r *userRepo) List(ctx context.Context, page, limit uint) ([]*biz.User, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&biz.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []*biz.User
	if err := r.db.WithContext(ctx).Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *userRepo) Get(ctx context.Context, id uint) (*biz.User, error) {
	user := new(biz.User)
	if err := r.db.WithContext(ctx).First(user, id).Error; err != nil {
		return nil, wrapErr(err)
	}

	return user, nil
}

func (r *userRepo) Create(ctx context.Context, user *biz.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Update is read-modify-write: it keeps CreatedAt intact and reports a
// missing row instead of upserting, unlike Save on a partial struct.
func (r *userRepo) Update(ctx context.Context, user *biz.User) (*biz.User, error) {
	existing := new(biz.User)
	if err := r.db.WithContext(ctx).First(existing, user.ID).Error; err != nil {
		return nil, wrapErr(err)
	}

	existing.Name = user.Name
	if err := r.db.WithContext(ctx).Save(existing).Error; err != nil {
		return nil, err
	}

	return existing, nil
}

func (r *userRepo) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&biz.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return biz.ErrNotFound
	}

	return nil
}
