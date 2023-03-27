package nesto_struct

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v9"

	"github.com/nestoca/pkg/addresses/regions"
)

// variable to store results of benchmark to prevent any compiler optimizations
var record []string

type testCase struct {
	name string
	app  func() Application
	errs []string
}

// unit test for the nesto common cases
// PLEASE DO NOT MODIFY THE TEST
// it's the same as v10 folder and should prove the working solution
func TestStructValidationDefaultValidator(t *testing.T) {
	for _, tc := range provideDefaultTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			app := tc.app()
			errors := app.Validate()

			assert.Equal(t, tc.errs, errors)
		})
	}
}

// benchmark test written to measure performance of a validator on single struct
func BenchmarkSingleStructValidation(b *testing.B) {
	var r []string
	app := provideValidStruct()

	b.ResetTimer() // to eliminate prep time spoil the results

	// run the benchmark function b.N times
	for n := 0; n < b.N; n++ {
		r = app.Validate()
	}

	record = r // this is here just to avoid any go compiler optimization
}

// benchmark test written to measure performance of a validator on batch of applications
func Benchmark1000StructValidations(b *testing.B) {
	var r []string
	applications := make([]Application, 1000)

	for i := 0; i < 1000; i++ {
		applications = append(applications, provideValidStruct())
	}

	b.ResetTimer() // to eliminate prep time spoil the results

	// run the benchmark function b.N times
	for n := 0; n < b.N; n++ {
		for _, app := range applications {
			r = app.Validate()
		}
	}

	record = r // this is here just to avoid any go compiler optimization
}

func provideDefaultTestCases() []testCase {
	return []testCase{
		{
			"1/valid",
			provideValidStruct,
			nil,
		},
		{
			"2/valid/omitempty",
			func() Application {
				app := provideValidStruct()
				app.Applicants[123456].Address.Street = ""         // this is valid due to omitempty
				app.Applicants[123456].Address.CountryCode = "US"  // setting to US so it will make SIN not required
				app.Applicants[123456].SocialInsuranceNUmber = nil // because of US does not have social insurance number
				return app
			},
			nil,
		},
		{
			"3/valid/nothing there",
			func() Application {
				app := provideValidStruct()
				app.Applicants = nil // omit empty tag
				return app
			},
			nil,
		},
		{
			"4/valid/nothing there 2",
			func() Application {
				app := provideValidStruct()
				app.Applicants = make(map[int]*Applicant, 2)
				return app
			},
			nil,
		},
		{
			"5/invalid/required dive on applicants",
			func() Application {
				app := provideValidStruct()
				app.Applicants = make(map[int]*Applicant, 2)
				app.Applicants[1111] = nil
				return app
			},
			[]string{"Application.Applicants[1111]"},
		},
		{
			"6/invalid/applicant all address fields",
			func() Application {
				app := provideValidStruct()
				app.Applicants[123456].Address.CountryCode = ""       // missing country code
				app.Applicants[123456].Address.City = "Vancouver"     // not part of enum
				app.Applicants[123456].Address.Street = "Short St"    // street name too short
				app.Applicants[123456].Address.PostalCode = "0805 03" // no Canadian postal code

				return app
			},
			[]string{
				"Application.Applicants[123456].Address.Street",
				"Application.Applicants[123456].Address.City",
				"Application.Applicants[123456].Address.CountryCode",
				"Application.Applicants[123456].Address.PostalCode"},
		},
		{
			"7/invalid/applicant country code",
			func() Application {
				app := provideValidStruct()
				app.Applicants[123456].Address.CountryCode = "ABCG" // not a valid ISO country code

				return app
			},
			[]string{"Application.Applicants[123456].Address.CountryCode"},
		},
		{
			"8/invalid/email too long",
			func() Application {
				app := provideValidStruct()
				app.Applicants[123456].Email.SetValid("TooLongEmailThatIJustCameUpWith@domain.com")

				return app
			},
			[]string{"Application.Applicants[123456].Email"},
		},
		{
			"9/invalid/email missing",
			func() Application {
				app := provideValidStruct()
				app.Applicants[123456].Email.SetValid("")

				return app
			},
			[]string{"Application.Applicants[123456].Email"},
		},
		{
			"10/invalid/email not defined",
			func() Application {
				app := provideValidStruct()
				app.Applicants[123456].Email = null.String{}

				return app
			},
			[]string{"Application.Applicants[123456].Email"},
		},
		{
			"11/invalid/custom validation/phone",
			func() Application {
				app := provideValidStruct()
				app.Applicants[123456].Phone = "403-111-5555" // still valid number but different from what validation expects for demo purposes

				return app
			},
			[]string{"Application.Applicants[123456].Phone"},
		},
		{
			"12/invalid/missing SIN for Canada",
			func() Application {
				app := provideValidStruct()
				app.Applicants[123456].SocialInsuranceNUmber = nil

				return app
			},
			[]string{"Application.Applicants[123456].SocialInsuranceNUmber"},
		},
	}
}

func provideValidStruct() Application {
	sin := SIN("666-666-666")
	return Application{
		Applicants: map[int]*Applicant{
			123456: {
				SocialInsuranceNUmber: &sin,
				Email:                 null.StringFrom("myemail@email.com"),
				Phone:                 "555-555-555",
				Address: Address{
					Street:      "Long St SW",
					City:        "Calgary",
					CountryCode: regions.RegionCodeCA,
					PostalCode:  "T2Y5G1",
				},
			},
		},
	}
}
