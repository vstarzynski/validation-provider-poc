package nesto_map

import (
	"database/sql/driver"
	"reflect"

	"github.com/go-playground/validator/v10"
)

// ValidateValuer implements validator.CustomTypeFunc for driver.Valuer type
// this si sued for complex types like sql.* or null.* which conveniently implement this common Go interface
func ValidateValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {

		val, err := valuer.Value()
		if err != nil {
			return validator.InvalidValidationError{Type: field.Type()}
		}

		return val
	}

	return nil
}

// ValidatePhone is custom validation method to demonstrate ay custom validation
func ValidatePhone(fl validator.FieldLevel) bool {
	return fl.Field().String() == "555-555-555" // dummy example
}
