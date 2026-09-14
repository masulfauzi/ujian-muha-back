package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateRegister(data interface{}) error {
	return validate.Struct(data)
}

func ValidateLogin(data interface{}) error {
	return validate.Struct(data)
}
