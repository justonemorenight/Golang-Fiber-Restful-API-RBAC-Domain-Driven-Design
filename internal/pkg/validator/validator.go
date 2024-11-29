package validator

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func GetValidator() *validator.Validate {
	return validate
}

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}