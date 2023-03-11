package main

import (
	"context"

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
	validate.RegisterStructValidation(decorateStructValidation(vp.DefaultUserValidation, vp.tenantValidators[tenantID].UserValidation), POCUser{})
	err := validate.Struct(user)
	if err != nil {
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

	// Validate First Name - Max Length is 10
	err := sl.Validator().Var(user.FirstName, "max=10")
	if err != nil {
		sl.ReportError(user, "firstName", "FirstName", "len=10", "")
	}

	// Validate Age - 18+
	err = sl.Validator().Var(user.Age, "min=18")
	if err != nil {
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
