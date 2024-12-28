package rule

import (
	"github.com/gookit/validate"
	"gorm.io/gorm"
)

func GlobalRules(db *gorm.DB) {
	validate.AddValidators(validate.M{
		"exists":    NewExists(db).Passes,
		"notExists": NewNotExists(db).Passes,
	})
	validate.AddGlobalMessages(map[string]string{
		"exists":    "{field} 不存在",
		"notExists": "{field} 已存在",
	})
}
