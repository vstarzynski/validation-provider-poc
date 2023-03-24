package main

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

type POCValidator interface {
	UserValidation(sl validator.StructLevel)
	UserValidationRules() map[string]string
}

// POCValidationProvider provides ways to validate fields.
type POCValidationProvider interface {
	ValidateUser(ctx context.Context, user POCUser) error
}

// POCDefaultValidationProvider is the default validation provider.
// It has an embedded sanitizer that should be used to sanitize data before validation is executed.
type POCDefaultValidationProvider struct {
	tenantValidators   map[int]POCValidator // allows multi tenancy validation
	tenantRules        map[int]map[string]string
	validationEntities map[string]map[string]string
}

// NewPOCDefaultValidationProvider returns a new POCDefaultValidationProvider
func NewPOCDefaultValidationProvider() *POCDefaultValidationProvider {
	return &POCDefaultValidationProvider{
		tenantValidators:   make(map[int]POCValidator),
		validationEntities: make(map[string]map[string]string),
	}
}

func (vp *POCDefaultValidationProvider) SetTenantValidator(tenantID int, validator POCValidator) {
	vp.tenantValidators[tenantID] = validator
}

func (vp *POCDefaultValidationProvider) SetTenantRules(tenantID int, rules map[string]string) {
	vp.tenantRules[tenantID] = rules
}

func (vp *POCDefaultValidationProvider) ValidateUser(ctx context.Context, user POCUser) error {
	// validation that is applied to all tenants
	tenantID := ctx.Value("tenant").(int)
	validate := validator.New()
	validate.RegisterValidation("startswiths", ValidateFieldStartsWithS)
	// Register function to get tag name from json tags by default, then field names

	// Register Struct Validation Pattern
	//validate.RegisterStructValidation(decorateStructValidation(vp.DefaultUserValidation, vp.tenantValidators[tenantID].UserValidation), POCUser{})
	//err := validate.Struct(user)
	//if err != nil {
	//	test := err.(validator.ValidationErrors)
	//	test = test
	//	return err
	//}

	//  RegisterStructValidationMapRules Pattern
	userRules := decorateRules(ComposeDefaultUserRules(), vp.tenantValidators[tenantID].UserValidationRules())
	validate.RegisterStructValidationMapRules(userRules, POCUser{})
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

	ValidateFieldWithTag(sl, user, user.FirstName, "FirstName", "max=10", vp.validationEntities)

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

// extractJSONTag extracts the JSON tag of a field based on its name.
// Function parameters:
// T: Struct to extract the JSON tag from
// name: Field name to extract JSON tag
//
// If a JSON tag does not exist for the specified field the field name is returned instead.
// If field name does not exist in the struct, an empty string is returned.
func extractJSONTag(T any, name string) string {
	// newT represents the actual struct where the field JSON tag will be extracted from
	var s interface{}
	if reflect.TypeOf(T).Kind() == reflect.Ptr {
		// If T is a pointer to an interface, set newT to the value that T points to
		s = reflect.ValueOf(T).Elem().Interface()
	} else {
		s = T
	}
	// Using reflection, extract the field tag name
	if field, ok := reflect.TypeOf(s).FieldByName(name); ok {
		tagName := field.Tag.Get("json")
		if len(tagName) > 0 {
			return tagName
		}
	} else {
		return ""
	}
	return name
}

// findStructField looks for a field (f) in the given struct (s).
// Function parameters:
// s: pointer to the struct. If s is not a pointer to the struct, nil will be returned.
// f: pointer to the field being looked for and should be a pointer to the actual struct field. If f is not a pointer to the field, nil will be returned.
//
// If found, the field info will be returned. Otherwise, nil will be returned.
func findStructField_old(s interface{}, f interface{}) *reflect.StructField {
	// Check if s (struct) is a pointer to an interface
	var structValue reflect.Value
	if reflect.ValueOf(s).Type().Kind() == reflect.Ptr || reflect.ValueOf(s).Type().Kind() == reflect.Interface {
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
			if fieldPointer == structValue.Field(i).Addr().Pointer() {
				return &structField
			}
		}
	}
	return nil
}

func decorateRules(rules ...map[string]string) map[string]string {
	decoratedRules := make(map[string]string)
	for _, r := range rules {
		for k, v := range r {
			appendRule(k, v, decoratedRules)
		}
	}
	return decoratedRules
}

func ComposeDefaultUserRules() map[string]string {
	rules := make(map[string]string)
	appendRule("FirstName", "max=10", rules)
	appendRule("Age", "min=18", rules)
	appendRule("Email", "required,email", rules)
	return rules
}

func ComposeDefaultAddressRules() map[string]string {
	rules := make(map[string]string)
	appendRule("FirstName", "max=10", rules)
	appendRule("Age", "min=18", rules)
	appendRule("Email", "required,email", rules)
	return rules
}

func appendRule(field, rule string, rules map[string]string) {

	// check if field is there
	if val, ok := rules[field]; ok {
		// append
		if len(val) != 0 {
			rules[field] = fmt.Sprintf("%s,%s", rules[field], rule)
		} else {
			rules[field] = rule
		}
	} else {
		// include
		rules[field] = rule
	}
}

// ValidateFieldWithTag validates a field based on the tag provided and reports an error based on the name provided.
// Function parameters:
// sl: StructLevel to be used
// s: A pointer to the struct where validation is applied (only used for reporting error)
// field: A pointer to the field being validated
// tag: String that defined which validation should be executed
//
// For validation purposes, the "field" parameter can be both a pointer or a value, however to properly report errors
// both "s" and "field" should be pointers.
func ValidateFieldWithTag(sl validator.StructLevel, s, field interface{}, fieldName, tag string, structFields map[string]map[string]string) {
	fieldValue := field
	if reflect.TypeOf(field).Kind() == reflect.Ptr {
		fieldValue = reflect.Indirect(reflect.ValueOf(field)).Interface()
	}
	err := sl.Validator().Var(fieldValue, tag)
	if err != nil {
		structValue := reflect.ValueOf(s)
		mapKey := fmt.Sprintf("%s.%s", structValue.Type().PkgPath(), structValue.Type().Name())

		fieldName = structFields[mapKey][fieldName]
		sl.ReportError(field, fieldName, fieldName, tag, "")
	}
}

func ComposeEntityFieldsMap(structs ...interface{}) map[string]map[string]string {
	entityFields := make(map[string]map[string]string)
	for i := 0; i < len(structs); i++ {
		structValue := reflect.ValueOf(structs[i])
		if reflect.ValueOf(structValue).Type().Kind() != reflect.Struct {
			return nil
		}
		fieldsMap := make(map[string]string)
		for i := 0; i < structValue.NumField(); i++ {
			structField := structValue.Type().Field(i)
			fieldName := structField.Name
			fieldJSONTag := extractJSONTag(structValue.Interface(), fieldName)
			fieldsMap[fieldName] = fieldJSONTag
		}
		entityFields[fmt.Sprintf("%s.%s", structValue.Type().PkgPath(), structValue.Type().Name())] = fieldsMap
	}
	return entityFields
}

// ValidateFieldStartsWithS implements validator.Func
func ValidateFieldStartsWithS(fl validator.FieldLevel) bool {
	return fl.Field().String()[0:1] == "S"
}
