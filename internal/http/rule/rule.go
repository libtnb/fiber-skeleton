package rule

import (
	"github.com/gookit/validate"
	"gorm.io/gorm"
)

func GlobalRules(db *gorm.DB) {
	validate.AddValidators(validate.M{
		"exists":     NewExists(db).Passes,
		"not_exists": NewNotExists(db).Passes,
	})
	validate.AddGlobalMessages(map[string]string{
		"exists":     "{field} 不存在",
		"not_exists": "{field} 已存在",
	})
}
