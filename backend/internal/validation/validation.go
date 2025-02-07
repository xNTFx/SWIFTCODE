package validation

import (
	"regexp"

	"backend/internal/models"
)

// Regular expression definitions for input validation.
var (
	AlnumSpaceRegex = regexp.MustCompile(`^[A-Za-z0-9\s]+$`)     // letters, numbers and spaces
	CountryIsoRegex = regexp.MustCompile(`^[A-Za-z]{2}$`)        // exactly 2 letters
	SwiftCodeRegex  = regexp.MustCompile(`^[A-Za-z0-9]{11}$`)    // exactly 11 letters/digits
	AddressRegex    = regexp.MustCompile(`^[A-Za-z0-9\s,.-/]+$`) // allows commas, periods, dashes and slashes
)

func ValidateSwiftCodeBranch(input models.SwiftCodeBranch) []string {
	var errors []string

	if !SwiftCodeRegex.MatchString(input.SwiftCode) {
		errors = append(errors, "SWIFT code must be exactly 11 alphanumeric characters")
	}
	if len(input.BankName) > 255 || !AlnumSpaceRegex.MatchString(input.BankName) {
		errors = append(errors, "Bank name must be at most 255 characters and contain only letters, numbers, and spaces")
	}
	if !CountryIsoRegex.MatchString(input.CountryISO2) {
		errors = append(errors, "Country ISO2 must be exactly 2 letters")
	}
	if len(input.CountryName) > 100 || !AlnumSpaceRegex.MatchString(input.CountryName) {
		errors = append(errors, "Country name must be at most 100 characters and contain only letters, numbers, and spaces")
	}
	if len(input.Address) > 255 || !AddressRegex.MatchString(input.Address) {
		errors = append(errors, "Address must be at most 255 characters and contain only letters, numbers, spaces, commas, periods, dashes, and slashes")
	}

	return errors
}

func ValidateSwiftCodeFields(body models.SwiftCodeBranch) []string {
	missingFields := []string{}
	fields := map[string]*string{
		"swift_code":   &body.SwiftCode,
		"bank_name":    &body.BankName,
		"country_iso2": &body.CountryISO2,
		"country_name": &body.CountryName,
		"address":      &body.Address,
	}
	for key, value := range fields {
		if value == nil || *value == "" {
			missingFields = append(missingFields, key)
		}
	}
	if body.IsHeadquarter == nil {
		missingFields = append(missingFields, "is_headquarter")
	}
	return missingFields
}
