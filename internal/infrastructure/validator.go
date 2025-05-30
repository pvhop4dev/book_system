package infrastructure

import "github.com/go-playground/validator/v10"

var Validate *validator.Validate

func init() {
	Validate = validator.New()
	Validate.RegisterValidation("isbn", func(fl validator.FieldLevel) bool {
		return false
	})

}
func GetValidate() *validator.Validate {
	return Validate
}
