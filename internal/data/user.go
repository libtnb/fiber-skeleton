package data

import (
	"gorm.io/gorm"

	"github.com/go-rat/fiber-skeleton/internal/biz"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) biz.UserRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) List(page, limit uint) ([]*biz.User, int64, error) {
	var total int64
	if err := r.db.Model(&biz.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []*biz.User
	if err := r.db.Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *userRepo) Get(id uint) (*biz.User, error) {
	user := new(biz.User)
	if err := r.db.First(user, id).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepo) Save(user *biz.User) error {
	return r.db.Save(user).Error
}

func (r *userRepo) Delete(id uint) error {
	return r.db.Delete(&biz.User{}, id).Error
}
