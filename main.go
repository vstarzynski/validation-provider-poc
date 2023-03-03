package main

import (
	"context"
	"fmt"
)

// POCUser contains POC user information
type POCUser struct {
	FirstName string
	Age       uint8
	Email     string     `validate:"email"`
	Addresses []*Address `validate:"required,dive"`
	Account   *Account
}

// Address info
type Address struct {
	ZipCode string `validate:"required"`
	Phone   string
}

// Account info
type Account struct {
	ID      string `validate:"required"`
	Balance float64
}

type POCUserValidator interface {
	ValidateUser(ctx context.Context, user POCUser) error
}

// UserValidationProvider provides ways to validate fields.
type UserValidationProvider interface {
	Validate(ctx context.Context) error
}

// DefaultUserValidationProvider is the default validation provider.
// It has an embedded sanitizer that should be used to sanitize data before validation is executed.
type DefaultUserValidationProvider struct {
	validators map[int]POCUserValidator // allows multi tenancy validation
}

// NewDefaultUserValidationProvider returns a new DefaultUserValidationProvider with an embedded DefaultSanitizer.
func NewDefaultUserValidationProvider() *DefaultUserValidationProvider {
	return &DefaultUserValidationProvider{
		validators: make(map[int]POCUserValidator),
	}
}

func (vp *DefaultUserValidationProvider) SetTenantValidator(tenantID int, validator POCUserValidator) {
	vp.validators[tenantID] = validator
}

func main() {

	vp := NewDefaultUserValidationProvider()
	tav := NewTenantAUserValidator()
	tbv := NewTenantBUserValidator()

	vp.SetTenantValidator(1, tav) // nesto
	vp.SetTenantValidator(2, tbv) // ig

	pocUser := POCUser{
		FirstName: "John",
		Age:       43,
		Email:     "jdoe@email.com",
		Addresses: []*Address{
			{
				ZipCode: "zip",

				Phone: "+16175551212",
			},
			{
				Phone: "abc",
			},
		},
		Account: &Account{
			ID:      "anuuid",
			Balance: -4.78,
		},
	}

	err := vp.validators[1].ValidateUser(context.Background(), pocUser)
	if err != nil {
		fmt.Println("Validation for Tenant A failed")
		fmt.Println(err)
	}

	fmt.Println("*****")

	err = vp.validators[2].ValidateUser(context.Background(), pocUser)
	if err != nil {
		fmt.Println("Validation for Tenant B failed")
		fmt.Println(err)
	}
}
