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

//func (v *TenantBUserValidator) ValidateUser(ctx context.Context, user POCUser) error {
//	validate := validator.New()
//	validate.RegisterStructValidation(TenantBUserStructLevelValidation, POCUser{})
//	err := validate.Struct(user)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// TenantBUserStructLevelValidation sets struct validation only required for Tenant B
//func TenantBUserStructLevelValidation(sl validator.StructLevel) {
//	user := sl.Current().Interface().(POCUser)
//
//	// Age between 18 and 40
//	if user.Age < 18 || user.Age > 40 {
//		sl.ReportError(user.Age, "age", "Age", "agebetween18and40", "")
//	}
//
//	// Email required
//	if len(user.Email) == 0 {
//		sl.ReportError(user.Email, "email", "Email", "required", "")
//	}
//}
