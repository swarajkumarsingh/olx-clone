package validators

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidateStruct - validates the struct and return error if they occur
func ValidateStruct(s interface{}) error {
	v := validator.New()
	err := v.Struct(s)

	if err != nil {
		var invalidArgs []string

		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, validationError := range validationErrors {
				invalidArgs = append(invalidArgs, strings.ToLower(validationError.Field()))
			}
		}
		errorMsg := "Missing or invalid arguments: " + strings.Join(invalidArgs, ", ")
		return errors.New(errorMsg)
	}

	return err
}

// IsEmpty - checks if a value is empty.
func IsEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		return value.(string) == ""
	case reflect.Slice, reflect.Map, reflect.Array, reflect.Chan:
		return reflect.ValueOf(value).Len() == 0
	case reflect.Ptr:
		return reflect.ValueOf(value).IsNil()
	default:
		return false
	}
}
