package pkg

import (
	"github.com/go-playground/validator/v10"
)

type ValidatorWrapper struct {
	Validator *validator.Validate
}

func NewValidator() *validator.Validate {
	MyValidator := validator.New()
	return MyValidator
}

func NewValidatorWrapper() *ValidatorWrapper {
	return &ValidatorWrapper{
		Validator: NewValidator(),
	}
}

func (v *ValidatorWrapper) Struct(s any) error {
	return v.Validator.Struct(s)
}
