package main

import (
	"context"

	"github.com/go-playground/validator/v10"
)

// TenantAUserValidator uses the default validate fields and non negative balance
type TenantAUserValidator struct {
}

func NewTenantAUserValidator() *TenantAUserValidator {
	return &TenantAUserValidator{}
}

func (v *TenantAUserValidator) ValidateUser(ctx context.Context, user POCUser) error {
	validate := validator.New()
	// Register Map Rules to take advantage of built-in validations
	rules := map[string]string{
		"Phone": "e164",
	}
	validate.RegisterStructValidationMapRules(rules, Address{})
	validate.RegisterStructValidation(TenantAUserStructLevelValidation, POCUser{})

	err := validate.Struct(user)
	if err != nil {
		return err
	}
	return nil
}

// TenantAUserStructLevelValidation sets struct validation only required for Tenant A
func TenantAUserStructLevelValidation(sl validator.StructLevel) {
	user := sl.Current().Interface().(POCUser)
	// Name has to start with "S"
	if user.FirstName[0:1] != "S" {
		sl.ReportError(user.FirstName, "first name", "FirstName", "namestartswiths", "")
	}
}
