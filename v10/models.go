package v10

import (
	"github.com/go-playground/validator/v10"
	"github.com/volatiletech/null/v9"

	"github.com/nestoca/pkg/addresses/regions"
)

//
// v10 models are the same as nesto models except they onyl have v10 validation tags and few custom methods to demonstrate validation
//

// SIN custom type
type SIN string

// Address struct containing typical address related data
type Address struct {
	Street      string             `validate:"omitempty,min=10"` // optional
	City        string             `validate:"omitempty,oneof=Toronto Calgary"`
	CountryCode regions.RegionCode `validate:"country_code"`                  // not required but triggers validation on string
	PostalCode  string             `validate:"required,canadian_postal_code"` // required and additional custom validation
}

// Applicant typical struct with fields
type Applicant struct {
	SocialInsuranceNUmber *SIN        `validate:"required_if=Address.CountryCode CA"`
	Email                 null.String `validate:"required,max=20"`
	Phone                 string      `validate:"required,phone"`
	Address
}

type ApplicationApplicants map[int]*Applicant

// Application represents main root struct of various cases
type Application struct {
	Applicants ApplicationApplicants `validate:"omitempty,dive,required"`
}

// Validate method to demonstrate validator
func (a Application) Validate() []string {
	validate := validator.New()
	validate.RegisterAlias("canadian_postal_code", "postcode_iso3166_alpha2=CA")
	_ = validate.RegisterValidation("phone", ValidatePhone)
	validate.RegisterCustomTypeFunc(ValidateValuer, null.String{}, null.Int{}, null.Bool{}, null.Float64{}, null.Time{})

	var fields []string
	err := validate.Struct(a)
	if err != nil {
		for _, vErr := range err.(validator.ValidationErrors) {
			fields = append(fields, vErr.Namespace())
		}
	}

	// could be error, but for demo purposes just return list of invalid fields
	return fields
}
