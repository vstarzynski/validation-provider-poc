package main

import (
	"context"

	"github.com/go-playground/validator/v10"
)

// TenantBUserValidator uses additional validation
type TenantBUserValidator struct{}

func NewTenantBUserValidator() *TenantBUserValidator {
	return &TenantBUserValidator{}
}

func (v *TenantBUserValidator) ValidateUser(ctx context.Context, user POCUser) error {
	validate := validator.New()
	validate.RegisterStructValidation(TenantBUserStructLevelValidation, POCUser{})
	err := validate.Struct(user)
	if err != nil {
		return err
	}

	return nil
}

// TenantBUserStructLevelValidation sets struct validation only required for Tenant B
func TenantBUserStructLevelValidation(sl validator.StructLevel) {
	user := sl.Current().Interface().(POCUser)

	// Age between 18 and 40
	if user.Age < 18 || user.Age > 40 {
		sl.ReportError(user.Age, "age", "Age", "agebetween18and40", "")
	}

	// Email required
	if len(user.Email) == 0 {
		sl.ReportError(user.Email, "email", "Email", "required", "")
	}
}
