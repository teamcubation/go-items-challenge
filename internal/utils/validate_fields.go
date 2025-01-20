package utils

import "github.com/go-playground/validator/v10"

func ValidateStruct(data interface{}) error {
	validate := validator.New()
	return validate.Struct(data)
}
