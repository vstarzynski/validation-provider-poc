package nesto_map

import (
	"fmt"

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
	addressRules := decorateRules(composeDefaultAddressRules())
	applicantRules := decorateRules(composeDefaultApplicantRules())
	applicationRules := decorateRules(composeDefaultApplicationRules())
	validate.RegisterStructValidationMapRules(addressRules, Address{})
	validate.RegisterStructValidationMapRules(applicantRules, Applicant{})
	validate.RegisterStructValidationMapRules(applicationRules, Application{})

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

func decorateRules(rules ...map[string]string) map[string]string {
	decoratedRules := make(map[string]string)
	for _, r := range rules {
		for k, v := range r {
			appendRule(k, v, decoratedRules)
		}
	}
	return decoratedRules
}

func composeDefaultAddressRules() map[string]string {
	rules := make(map[string]string)
	appendRule("Street", "omitempty,min=10", rules)
	appendRule("City", "omitempty,oneof=Toronto Calgary", rules)
	appendRule("CountryCode", "country_code", rules)
	appendRule("PostalCode", "required,canadian_postal_code", rules)
	return rules
}

func composeDefaultApplicantRules() map[string]string {
	rules := make(map[string]string)
	appendRule("SocialInsuranceNUmber", "required_if=Address.CountryCode CA", rules)
	appendRule("Email", "required,max=20", rules)
	appendRule("Phone", "required,phone", rules)
	return rules
}

func composeDefaultApplicationRules() map[string]string {
	rules := make(map[string]string)
	appendRule("Applicants", "omitempty,dive,required", rules)
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
