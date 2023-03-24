package main

import "github.com/go-playground/validator/v10"

// TenantBUserValidator uses additional validation
type TenantBUserValidator struct{}

func NewTenantBUserValidator() *TenantBUserValidator {
	return &TenantBUserValidator{}
}

// UserValidation sets struct validation only required for Tenant A
func (v *TenantBUserValidator) UserValidation(sl validator.StructLevel) {
	user := sl.Current().Interface().(POCUser)

	// Maximum age is 40
	if user.Age < 20 || user.Age > 40 {
		sl.ReportError(user.Age, "age", "Age", "agenotinbetween20and40", "")
	}

	// Address province is province name
	addresses := user.Addresses
	for _, a := range addresses {
		sl.Validator().RegisterValidation("isprovincecode", isProvinceCode)
		err := sl.Validator().Var(a.Province, "isprovincecode")
		if err != nil {
			sl.ReportError(a, "province", "Province", "isprovincecode", "")
		}
	}
}

func (v *TenantBUserValidator) UserValidationRules() map[string]string {
	userRules := make(map[string]string)
	return userRules
}
