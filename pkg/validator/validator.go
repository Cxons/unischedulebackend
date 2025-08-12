package validator

import "github.com/go-playground/validator/v10"


var validate = validator.New()

func ValidateStruct(payload interface{}) error{
	return validate.Struct(payload)
}