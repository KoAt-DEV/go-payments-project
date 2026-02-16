package bootstrap

import "github.com/go-playground/validator/v10"

func InitValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	return v
}
