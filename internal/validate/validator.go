package validate

import (
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	v10 "github.com/go-playground/validator/v10"
	zh_translation "github.com/go-playground/validator/v10/translations/zh"
	"github.com/pot-code/gobit/pkg/validate"
	"golang.org/x/text/language"
)

var Validator *validate.Validator

func Init() {
	Validator = validate.New().RegisterLocale(language.SimplifiedChinese, zh.New(), func(v *v10.Validate, t ut.Translator) error {
		return zh_translation.RegisterDefaultTranslations(v, t)
	}).DefaultLocale(language.SimplifiedChinese).Build()
}
