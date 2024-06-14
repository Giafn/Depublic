package validator

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

func Validate(input interface{}) map[string]string {
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return map[string]string{"error": "validasi input gagal: " + err.Error()}
		}

		validationErrors := make(map[string]string)
		inputVal := reflect.ValueOf(input)
		for _, err := range err.(validator.ValidationErrors) {
			field, _ := inputVal.Type().FieldByName(err.Field())
			jsonTag := field.Tag.Get("json")
			if jsonTag == "" {
				jsonTag = err.Field() // fallback to struct field name if json tag is missing
			}
			validationErrors[jsonTag] = err.Tag()
		}
		return validationErrors
	}
	return nil
}
