package bootstrap

import (
	"reflect"
	"strings"

	"github.com/libtnb/validator"
	"github.com/libtnb/validator/translations"
	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/conf"
)

func NewValidator(i do.Injector) (*validator.Validator, error) {
	config := do.MustInvoke[*conf.Config](i)

	opts := []validator.Option{
		validator.WithTagNameFunc(fieldName),
	}
	if messages := localeMessages(config.App.Locale); messages != nil {
		opts = append(opts, validator.WithTranslation(messages))
	}

	return validator.NewValidator(opts...), nil
}

// fieldName reports fields in error messages by the name the client sent.
func fieldName(field reflect.StructField) string {
	for _, tag := range []string{"form", "json", "query", "uri"} {
		if name, _, _ := strings.Cut(field.Tag.Get(tag), ","); name != "" && name != "-" {
			return name
		}
	}
	return field.Name
}

func localeMessages(locale string) map[string]string {
	switch locale {
	case "zh_Hans", "zh_CN":
		return translations.ZhHans()
	case "zh_Hant", "zh_TW":
		return translations.ZhHant()
	case "ja":
		return translations.Ja()
	case "ko":
		return translations.Ko()
	case "es":
		return translations.Es()
	case "ru":
		return translations.Ru()
	default:
		return nil // built-in English
	}
}
