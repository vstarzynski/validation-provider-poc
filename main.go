package main

import (
	"context"
	"fmt"
)

// POCUser contains POC user information
type POC struct {
	BaseUser
	FirstName string `json:"FIRSTNAME"`
	Age       uint8  `json:"myAge"`
	Email     string
	Phone     string
	Account   *Account `json:"account"`
}

// POCUser contains POC user information
type POCUser struct {
	BaseUser
	FirstName string `json:"FIRSTNAME"`
	Age       uint8  `json:"myAge,string"`
	Email     string
	Phone     string
	Addresses []*Address
	Account   *Account `json:"account"`
}

type BaseUser struct {
	LastName string
}

// Address info
type Address struct {
	ZipCode  string
	Province string
}

// Account info
type Account struct {
	ID      string `json:"anID"`
	Balance float64
}

func main() {

	vp := NewPOCDefaultValidationProvider()
	tav := NewTenantAUserValidator()
	tbv := NewTenantBUserValidator()

	vp.SetTenantValidator(1, tav) // tenant 1
	vp.SetTenantValidator(2, tbv) // tenant 2

	pocUser := POCUser{
		BaseUser: BaseUser{
			LastName: "Smith",
		},
		FirstName: "Pam with a very very very long name",
		Age:       17,
		Email:     "jdoe@mail.com",
		Phone:     "+16175551212",
		Addresses: []*Address{
			//{
			//	ZipCode:  "zip",
			//	Province: "Quebec",
			//},
			//{
			//	ZipCode:  "another zip",
			//	Province: "Ontario",
			//},
		},
		Account: &Account{
			ID:      "anuuid",
			Balance: -4.78,
		},
	}

	// test entity map composition
	vp.validationEntities = ComposeEntityFieldsMap(POCUser{})

	ctx := context.WithValue(context.Background(), "tenant", 1) // 1 - nesto | 2 - ig
	err := vp.ValidateUser(ctx, &pocUser)
	if err != nil {
		fmt.Println("Validation Provider failed...")
		fmt.Println(err)
	}
}
