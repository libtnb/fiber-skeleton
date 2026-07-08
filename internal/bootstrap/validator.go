package bootstrap

import (
	"reflect"
	"strings"

	"github.com/libtnb/validator"
	"github.com/libtnb/validator/contrib/gormrules"
	"github.com/libtnb/validator/translations"
	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/config"
	"github.com/libtnb/fiber-skeleton/internal/data"
)

// NewValidator builds the validator behind service.Bind: localized messages,
// wire-name field reporting and the database rules.
func NewValidator(i do.Injector) (*validator.Validator, error) {
	conf := do.MustInvoke[*config.Config](i)

	opts := []validator.Option{
		validator.WithTagNameFunc(fieldName),
	}
	if messages := localeMessages(conf.App.Locale); messages != nil {
		opts = append(opts, validator.WithTranslation(messages))
	}

	v := validator.NewValidator(opts...)
	gormrules.Register(v, do.MustInvoke[*data.Data](i).DB)

	return v, nil
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
