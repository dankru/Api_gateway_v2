package validation

import "github.com/go-playground/validator/v10"

func Validate(entity interface{}) error {
	validate := validator.New()
	return validate.Struct(entity)
}
