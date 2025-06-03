package infrastructure

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()
	Validate.RegisterValidation("phone_number", validator.Func(func(fl validator.FieldLevel) bool {
		reStartWith0 := regexp.MustCompile(`^0\d{9}$`)
		reStartWithPlus := regexp.MustCompile(`^\+\d{11}$`)
		return reStartWith0.MatchString(fl.Field().String()) || reStartWithPlus.MatchString(fl.Field().String())
	}))
}

// func GetValidate() *validator.Validate {
// 	return Validate
// }
