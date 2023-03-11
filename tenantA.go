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

//func (v *TenantAUserValidator) ValidateUser(ctx context.Context, user POCUser) error {
//	validate := validator.New()
//	// Register Map Rules to take advantage of built-in validations
//	rulesAddress := map[string]string{
//		"Phone": "e164",
//	}
//	rulesUser := map[string]string{
//		"Email": "required,email",
//	}
//	validate.SetTagName("validateTenantA")
//	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
//		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
//		// skip if tag key says it should be ignored
//		if name == "-" {
//			return ""
//		}
//		return name
//	})
//	validate.RegisterStructValidationMapRules(rulesAddress, Address{})
//	validate.RegisterStructValidationMapRules(rulesUser, POCUser{})
//	validate.RegisterStructValidation(TenantAUserStructLevelValidation, POCUser{})
//
//	err := validate.Struct(user)
//	if err != nil {
//		return err
//	}
//	//if err != nil {
//	//	if ferr, ok := err.(validator.ValidationErrors); ok {
//	//		var errstrings []string
//	//		errstrings = append(errstrings, "failed to validate: ")
//	//		for _, fieldError := range ferr {
//	//			errstrings = append(errstrings, fieldError.Field())
//	//		}
//	//		return fmt.Errorf(strings.Join(errstrings, ", "))
//	//	} else {
//	//		// ... deal with non-flags.Error case, if that's possible.
//	//	}
//	//}
//	return nil
//}
//
//// TenantAUserStructLevelValidation sets struct validation only required for Tenant A
//func TenantAUserStructLevelValidation(sl validator.StructLevel) {
//	user := sl.Current().Interface().(POCUser)
//	// Name has to start with "S"
//	if user.FirstName[0:1] != "S" {
//		sl.ReportError(user.FirstName, "first name", "FirstName", "namestartswiths", "")
//	}
//
//	// Custom email validation
//	err := sl.Validator().Var(user.Email, "required")
//	// Validate Email
//	if err != nil {
//		sl.ReportError(user, "email", "Email", "required", "")
//	}
//}
