package bootstrap

import (
	"fmt"
	"reflect"

	"github.com/go-playground/locales/zh_Hans_CN"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/zh"

	"github.com/TheTNB/go-web-skeleton/internal/app"
)

func initValidator() {
	translator := zh_Hans_CN.New()
	uni := ut.New(translator, translator)
	trans, _ := uni.GetTranslator("zh_Hans_CN")

	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("comment")
	})
	if err := zh.RegisterDefaultTranslations(validate, trans); err != nil {
		panic(fmt.Sprintf("failed to register validator translations: %v", err))
	}

	app.Translator = &trans
	app.Validator = validate
}
