package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateCreateUser(data interface{}) error {
	return validate.Struct(data)
}

func ValidateUpdateUser(data interface{}) error {
	return validate.Struct(data)
}
