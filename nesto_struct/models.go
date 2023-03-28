package nesto_struct

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/volatiletech/null/v9"

	"github.com/nestoca/pkg/addresses/regions"
)

//
// nesto models are the same as v10 models except have no validation tags assign at struct level
// and implementation should be added as required to validate the struct
//

// SIN custom type
type SIN string

// Address struct containing typical address related data
type Address struct {
	Street      string
	City        string
	CountryCode regions.RegionCode
	PostalCode  string
}

// Applicant typical struct with fields
type Applicant struct {
	SocialInsuranceNUmber *SIN
	Email                 null.String
	Phone                 string
	Address
}

type ApplicationApplicants map[int]*Applicant

// Application represents main root struct of various cases
type Application struct {
	Applicants ApplicationApplicants
}

// Validate method to demonstrate validator
func (a Application) Validate() []string {
	validate := validator.New()
	validate.RegisterAlias("canadian_postal_code", "postcode_iso3166_alpha2=CA")
	_ = validate.RegisterValidation("phone", ValidatePhone)
	validate.RegisterCustomTypeFunc(ValidateValuer, null.String{}, null.Int{}, null.Bool{}, null.Float64{}, null.Time{})

	// Decorate can be used for both default and tenant aware validation
	validate.RegisterStructValidation(decorateStructValidation(DefaultValidation), Application{})
	var fields []string
	err := validate.Struct(a)
	if err != nil {
		for _, vErr := range err.(validator.ValidationErrors) {
			fields = append(fields, vErr.Namespace())
		}
	}
	return fields
}

// DecorateStructValidation returns a decorated struct validation function
func decorateStructValidation(customValidation ...validator.StructLevelFunc) validator.StructLevelFunc {
	return func(sl validator.StructLevel) {
		for _, f := range customValidation {
			f(sl)
		}
	}
}

// DefaultValidation sets struct validation that will be shared between all tenants
func DefaultValidation(sl validator.StructLevel) {
	application := sl.Current().Interface().(Application)
	ValidateFieldWithTag(sl, application, application.Applicants, "Applicants", "omitempty,dive,required")
	for _, applicant := range application.Applicants {
		if applicant != nil {
			address := applicant.Address
			if address.CountryCode == "CA" {
				ValidateFieldWithTag(sl, applicant, applicant.SocialInsuranceNUmber, "SocialInsuranceNUmber", "required")
			}
			ValidateFieldWithTag(sl, applicant, applicant.Email, "Email", "required,max=20")
			ValidateFieldWithTag(sl, applicant, applicant.Phone, "Phone", "required,phone")
			ValidateFieldWithTag(sl, address, address.Street, "Street", "omitempty,min=10")
			ValidateFieldWithTag(sl, address, address.City, "City", "omitempty,oneof=Toronto Calgary")
			ValidateFieldWithTag(sl, address, address.CountryCode, "CountryCode", "country_code")
			ValidateFieldWithTag(sl, address, address.PostalCode, "PostalCode", "required,canadian_postal_code")
		}
	}
}

func fieldName(text ...string) string {
	var result string
	for _, t := range text {
		if len(result) == 0 {
			result = t
		} else {
			result = fmt.Sprintf("%s.%s", result, t)
		}
	}
	return result
}

func ValidateFieldWithTag(sl validator.StructLevel, s, field interface{}, fieldName, tag string) {
	fieldValue := field
	if reflect.TypeOf(field).Kind() == reflect.Ptr {
		value := reflect.ValueOf(field)
		if !value.IsNil() {
			fieldValue = reflect.Indirect(value).Interface()
		}
	}
	err := sl.Validator().Var(fieldValue, tag)
	if err != nil {
		fieldTag := extractJSONTag(s, fieldName)
		sl.ReportError(field, fieldTag, fieldTag, tag, "")
	}
}

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
