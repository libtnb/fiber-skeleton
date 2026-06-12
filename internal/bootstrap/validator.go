package bootstrap

import (
	"github.com/gookit/validate/v2"
	"github.com/gookit/validate/v2/locales/zhcn"
	"gorm.io/gorm"

	"github.com/libtnb/fiber-skeleton/internal/http/rule"
)

// NewValidator just for register global rules
func NewValidator(db *gorm.DB) *validate.Validation {
	zhcn.RegisterGlobal()
	validate.Config(func(opt *validate.GlobalOption) {
		opt.StopOnError = false
		opt.SkipOnEmpty = true
		opt.FieldTag = "form"
	})

	// register global rules
	rule.GlobalRules(db)

	return validate.NewEmpty()
}
