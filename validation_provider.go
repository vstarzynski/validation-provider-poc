package main

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type POCValidator interface {
	UserValidation(sl validator.StructLevel)
}

// POCValidationProvider provides ways to validate fields.
type POCValidationProvider interface {
	ValidateUser(ctx context.Context, user POCUser) error
}

// POCDefaultValidationProvider is the default validation provider.
// It has an embedded sanitizer that should be used to sanitize data before validation is executed.
type POCDefaultValidationProvider struct {
	tenantValidators map[int]POCValidator // allows multi tenancy validation
}

// NewPOCDefaultValidationProvider returns a new POCDefaultValidationProvider
func NewPOCDefaultValidationProvider() *POCDefaultValidationProvider {
	return &POCDefaultValidationProvider{
		tenantValidators: make(map[int]POCValidator),
	}
}

func (vp *POCDefaultValidationProvider) SetTenantValidator(tenantID int, validator POCValidator) {
	vp.tenantValidators[tenantID] = validator
}

func (vp *POCDefaultValidationProvider) ValidateUser(ctx context.Context, user POCUser) error {
	// validation that is applied to all tenants
	tenantID := ctx.Value("tenant").(int)
	validate := validator.New()
	// Register function to get tag name from json tags by default, then field names
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return fld.Name
		}
		return name
	})
	validate.RegisterStructValidation(decorateStructValidation(vp.DefaultUserValidation, vp.tenantValidators[tenantID].UserValidation), POCUser{})
	err := validate.Struct(user)
	if err != nil {
		test := err.(validator.ValidationErrors)
		test = test
		return err
	}
	return nil
}

// DecorateStructValidation returns a decorated struct validation function
func decorateStructValidation(customValidation ...validator.StructLevelFunc) validator.StructLevelFunc {
	return func(sl validator.StructLevel) {
		for _, f := range customValidation {
			f(sl)
		}
	}
}

// DefaultUserValidation sets struct validation that will be shared between all tenants
func (vp *POCDefaultValidationProvider) DefaultUserValidation(sl validator.StructLevel) {
	user := sl.Current().Interface().(POCUser)

	//field, _ := reflect.TypeOf(user).FieldByName("FirstName")
	//name := field.Tag.Get("json")
	//vp.validateField(sl, user.FirstName, name, "max=10")

	// Validate Age - 18+
	err := sl.Validator().Var(user.Age, "min=18")
	if err != nil {
		test := err.(validator.ValidationErrors)
		test = test
		sl.ReportError(user, "age", "Age", "min=18", "")
	}

	// Validate Email
	err = sl.Validator().Var(user.Email, "required,email")
	if err != nil {
		sl.ReportError(user, "email", "Email", "required,email", "")
	}

	// Validate Addresses
	address := user.Addresses
	for _, a := range address {
		err = sl.Validator().Var(a.ZipCode, "required")
		if err != nil {
			sl.ReportError(a, "zipcode", "ZipCode", "required", "")
		}
	}

	// Validate Account
	account := user.Account

	err = sl.Validator().Var(account.ID, "required")
	if err != nil {
		sl.ReportError(account, "id", "ID", "required", "")
	}
}

func (vp *POCDefaultValidationProvider) validateField(sl validator.StructLevel, field interface{}, name, tag string) {

	// Validate First Name - Max Length is 10
	//err := sl.Validator().Var(field, tag).(validator.ValidationErrors)
	//if err != nil {
	//	sl.ReportError(field, name, name, tag, "")
	//}
}

func extractJSONTag(T any, name string) string {
	if field, ok := reflect.TypeOf(T).FieldByName(name); ok {
		tagName := field.Tag.Get("json")
		if len(tagName) > 0 {
			return tagName
		}
	}
	return name
}

func tempFunc(s interface{}, f interface{}) *reflect.StructField {
	if reflect.ValueOf(s).Type().Kind() != reflect.Ptr {

	}
	temp := findStructField(s, f)
	fmt.Println(temp.Name)
	return temp
}

// findStructField looks for a field (f) in the given struct (s).
// This function receives:
// s: pointer to the struct.
// f: pointer to the field being looked for and should be a pointer to the actual struct field.
//
// If found, the field info will be returned. Otherwise, nil will be returned.
func findStructField(s interface{}, f interface{}) *reflect.StructField {
	// Check if s (struct) is a pointer to an interface
	var structValue reflect.Value
	if reflect.ValueOf(s).Type().Kind() == reflect.Ptr {
		structValue = reflect.ValueOf(s).Elem()
	} else {
		return nil
	}
	// Check if f (field) is a pointer to an interface
	var fieldValue reflect.Value
	if reflect.ValueOf(f).Type().Kind() == reflect.Ptr {
		fieldValue = reflect.ValueOf(f)
	} else {
		return nil
	}
	// Set field pointer and type
	fieldPointer := fieldValue.Pointer()
	fieldType := fieldValue.Elem().Type()
	for i := structValue.NumField() - 1; i >= 0; i-- {
		structField := structValue.Type().Field(i)
		// Compare if field and struct field are the same type
		if structValue.Field(i).Type().Kind() == fieldType.Kind() {
			// Compare if field and struct field are the same
			if structValue.Field(i).CanAddr() {
				if fieldPointer == structValue.Field(i).Addr().Pointer() {
					return &structField
				}
			} else {
				p := structValue.Field(i).Pointer()
				if fieldPointer == p {
					return &structField
				}
			}

		}
	}
	return nil
}
