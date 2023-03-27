package nesto_struct

import (
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
	validate.RegisterStructValidation(decorateStructValidation(DefaultValidation))
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
	user := sl.Current().Interface().(Application)

}
