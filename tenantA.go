package main

import (
	"github.com/go-playground/validator/v10"
)

// TenantAUserValidator uses the default validate fields and non negative balance
type TenantAUserValidator struct {
}

func NewTenantAUserValidator() *TenantAUserValidator {
	return &TenantAUserValidator{}
}

// UserValidation sets struct validation only required for Tenant A
func (v *TenantAUserValidator) UserValidation(sl validator.StructLevel) {
	user := sl.Current().Interface().(POCUser)

	// Name has to start with "S"
	if user.FirstName[0:1] != "S" {
		sl.ReportError(user.FirstName, "first name", "FirstName", "namestartswiths", "")
	}

	// Phone number has to be valid
	err := sl.Validator().Var(user.Phone, "e164")
	if err != nil {
		sl.ReportError(user, "phone", "Phone", "e164", "")
	}

	// Address province is province name
	addresses := user.Addresses
	for _, a := range addresses {
		sl.Validator().RegisterValidation("isprovincename", isProvinceName)
		err = sl.Validator().Var(a.Province, "isprovincename")
		if err != nil {
			sl.ReportError(addresses, "province", "Province", "isprovincename", "")
		}
	}
}

func (v *TenantAUserValidator) UserValidationRules() map[string]string {
	userRules := make(map[string]string)
	appendRule("FirstName", "startswiths", userRules)
	appendRule("Phone", "e164", userRules)
	return userRules
}
