package main

import "github.com/go-playground/validator/v10"

// custom validation functions - this can be placed in a shared location (such as pkg.validation)

// ProvinceCode
func isProvinceCode(fl validator.FieldLevel) bool {
	switch fl.Field().String() {
	case "AB", "BC", "MB", "NB", "NS", "ON", "PE", "QC", "SK":
		return true
	default:
		return false
	}
}

// ProvinceName
func isProvinceName(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	switch value {
	case "Alberta", "British Columbia", "Manitoba", "New Brunswick", "Nova Scotia", "Ontario", "Prince Edward", "Quebec", "Saskatchewan":
		return true
	default:
		return false
	}
}
