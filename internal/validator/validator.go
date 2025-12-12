package validator

import (
	"errors"
	"fmt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	govalidator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	validate *govalidator.Validate
	trans    ut.Translator
}

func New() *Validator {
	en := en.New()
	uni := ut.New(en, en)

	trans, found := uni.GetTranslator("en")
	if !found {
		// The validator is call on config app first, so is ok.
		panic("validator translator not found")
	}

	validate := govalidator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)

	v := &Validator{
		validate: validate,
		trans:    trans,
	}
	v.registerTranslation()
	return v
}

func (v *Validator) registerTranslation() {
	v.validate.RegisterTranslation("port", v.trans,
		func(ut ut.Translator) error {
			return ut.Add("port", "{0} has invalid value of {1}", true)
		},
		func(ut ut.Translator, fe govalidator.FieldError) string {
			t, _ := ut.T("port", fe.Field(), fmt.Sprintf("%v", fe.Value()))
			return t
		},
	)
}

func (v *Validator) Struct(input any) error {
	err := v.validate.Struct(input)
	if err != nil {
		var validationErrs govalidator.ValidationErrors

		switch {
		case errors.As(err, &validationErrs):
			errs := err.(govalidator.ValidationErrors)
			return ValidationErrorMap(errs.Translate(v.trans))
		default:
			return err
		}

	}
	return nil
}

// Implement validator interface from echo
func (v *Validator) Validate(i any) error {
	return v.Struct(i)
}
