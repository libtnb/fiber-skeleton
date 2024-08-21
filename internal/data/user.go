package data

import (
	"github.com/TheTNB/go-web-skeleton/internal/app"
	"github.com/TheTNB/go-web-skeleton/internal/biz"
)

type userRepo struct{}

func NewUserRepo() biz.UserRepo {
	return &userRepo{}
}

func (r *userRepo) List(page, limit uint) ([]*biz.User, int64, error) {
	var total int64
	if err := app.Orm.Model(&biz.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []*biz.User
	if err := app.Orm.Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *userRepo) Get(id uint) (*biz.User, error) {
	user := new(biz.User)
	if err := app.Orm.First(user, id).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepo) Save(user *biz.User) error {
	return app.Orm.Save(user).Error
}

func (r *userRepo) Delete(id uint) error {
	return app.Orm.Delete(&biz.User{}, id).Error
}
