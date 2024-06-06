package validator

import "github.com/go-playground/validator/v10"

func Validate(input interface{}) map[string]string {
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return map[string]string{"error": "validasi input gagal: " + err.Error()}
		}

		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors[err.Field()] = err.Tag()
		}
		return validationErrors
	}
	return nil
}
