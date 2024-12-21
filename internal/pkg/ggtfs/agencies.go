package ggtfs

import (
	"fmt"
)

type Agency struct {
	Id         *ID           // agency_id
	Name       Text          // agency_name
	URL        URL           // agency_url
	Timezone   Timezone      // agency_timezone
	Lang       *LanguageCode // agency_lang
	Phone      *PhoneNumber  // agency_phone
	FareURL    *URL          // agency_fare_url
	Email      *Email        // agency_email
	LineNumber int
}

func (a Agency) Validate() []error {
	var validationErrors []error

	// The 'agency_id' field is conditionally required, check it in the ValidateAgencies function.

	fields := []struct {
		fieldName string
		field     ValidAndPresentField
	}{
		{"agency_name", &a.Name},
		{"agency_url", &a.URL},
		{"agency_timezone", &a.Timezone},
	}

	for _, f := range fields {
		validationErrors = append(validationErrors, validateFieldIsPresentAndValid(f.field, f.fieldName, a.LineNumber, AgenciesFileName)...)
	}

	// Checking the underlying value of the field in ValidAndPresentField for nil would require reflection
	// v := reflect.ValueOf(i)
	// v.Kind() == reflect.Ptr && v.IsNil()
	// which is slow, so we can't use the above mechanism to check optional fields, since they might be nil (pointer field's default value is nil)
	// since CreateTrip might have not processed the field (if its header is missing from the csv).

	if a.Id != nil && !a.Id.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(AgenciesFileName, a.LineNumber, createInvalidFieldString("agency_id")))
	}

	if a.Lang != nil && !a.Lang.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(AgenciesFileName, a.LineNumber, createInvalidFieldString("agency_lang")))
	}

	if a.Phone != nil && !a.Phone.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(AgenciesFileName, a.LineNumber, createInvalidFieldString("agency_phone")))
	}

	if a.FareURL != nil && !a.FareURL.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(AgenciesFileName, a.LineNumber, createInvalidFieldString("agency_fare_url")))
	}

	if a.Email != nil && !a.Email.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(AgenciesFileName, a.LineNumber, createInvalidFieldString("agency_email")))
	}

	return validationErrors
}

func CreateAgency(row []string, headers map[string]int, lineNumber int) *Agency {
	agency := Agency{
		LineNumber: lineNumber,
	}

	for hName, hPos := range headers {
		switch hName {
		case "agency_id":
			agency.Id = NewOptionalID(&row[hPos])
		case "agency_name":
			agency.Name = NewText(&row[hPos])
		case "agency_url":
			agency.URL = NewURL(&row[hPos])
		case "agency_timezone":
			agency.Timezone = NewTimezone(&row[hPos])
		case "agency_lang":
			agency.Lang = NewOptionalLanguageCode(&row[hPos])
		case "agency_phone":
			agency.Phone = NewOptionalPhoneNumber(&row[hPos])
		case "agency_fare_url":
			agency.FareURL = NewOptionalURL(&row[hPos])
		case "agency_email":
			agency.Email = NewOptionalEmail(&row[hPos])
		}
	}

	return &agency
}

func ValidateAgencies(agencies []*Agency) ([]error, []string) {
	var validationErrors []error
	var recommendations []string

	aLength := len(agencies)

	if aLength == 0 || (aLength == 1 && agencies[0] == nil) {
		return validationErrors, recommendations
	}

	if aLength == 1 && !agencies[0].Id.IsValid() {
		validationErrors = append(validationErrors, agencies[0].Validate()...)
		recommendations = append(recommendations, createFileRowRecommendation(AgenciesFileName, agencies[0].LineNumber, "it is recommended that agency_id is specified even when there is only one agency"))
		return validationErrors, recommendations
	}

	if aLength == 1 {
		validationErrors = append(validationErrors, agencies[0].Validate()...)
		return validationErrors, recommendations
	}

	usedIds := make(map[string]bool)
	for _, a := range agencies {
		if a == nil {
			continue
		}

		validationErrors = append(validationErrors, a.Validate()...)

		if !a.Id.IsValid() {
			validationErrors = append(validationErrors, createFileRowError(AgenciesFileName, a.LineNumber, "a valid agency_id must be specified when multiple agencies are declared"))
			continue
		}

		if usedIds[a.Id.String()] {
			validationErrors = append(validationErrors, createFileRowError(AgenciesFileName, a.LineNumber, fmt.Sprintf("agency_id is not unique within the file")))
		} else {
			usedIds[a.Id.String()] = true
		}
	}

	return validationErrors, recommendations
}
