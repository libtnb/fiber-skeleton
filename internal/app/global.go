package app

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/knadh/koanf/v2"
	"gorm.io/gorm"
)

var (
	Conf       *koanf.Koanf
	Http       *fiber.App
	Orm        *gorm.DB
	Validator  *validator.Validate
	Translator *ut.Translator
)
