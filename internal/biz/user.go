package biz

import (
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

type UserRepo interface {
	List(page, limit uint) ([]*User, int64, error)
	Get(id uint) (*User, error)
	Save(user *User) error
	Delete(id uint) error
}
