package ggtfs

import (
	"fmt"
)

type Agency struct {
	Id         ID           // agency_id 		(conditionally required)
	Name       Text         // agency_name 		(required)
	URL        URL          // agency_url 		(required)
	Timezone   Timezone     // agency_timezone 	(required)
	Lang       LanguageCode // agency_lang 		(optional)
	Phone      PhoneNumber  // agency_phone 	(optional)
	FareURL    URL          // agency_fare_url 	(optional)
	Email      Email        // agency_email 	(optional)
	LineNumber int
}

func (a Agency) Validate() []error {
	var validationErrors []error

	requiredFields := []struct {
		fieldName string
		field     ValidAndPresentField
	}{
		{"agency_name", &a.Name},
		{"agency_url", &a.URL},
		{"agency_timezone", &a.Timezone},
	}
	for _, f := range requiredFields {
		if !f.field.IsValid() {
			validationErrors = append(validationErrors, createFileRowError(AgenciesFileName, a.LineNumber, createInvalidRequiredFieldString(f.fieldName)))
		}
	}

	optionalFields := []struct {
		field     ValidAndPresentField
		fieldName string
	}{
		{&a.Id, "agency_id"},
		{&a.Lang, "agency_lang"},
		{&a.Phone, "agency_phone"},
		{&a.FareURL, "agency_fare_url"},
		{&a.Email, "agency_email"},
	}

	for _, field := range optionalFields {
		if field.field != nil && field.field.IsPresent() && !field.field.IsValid() {
			validationErrors = append(validationErrors, createFileRowError(AgenciesFileName, a.LineNumber, createInvalidFieldString(field.fieldName)))
		}
	}

	return validationErrors
}

func CreateAgency(row []string, headers map[string]int, lineNumber int) *Agency {
	agency := Agency{
		LineNumber: lineNumber,
	}

	for hName := range headers {
		v := getRowValueForHeaderName(row, headers, hName)

		switch hName {
		case "agency_id":
			agency.Id = NewID(v)
		case "agency_name":
			agency.Name = NewText(v)
		case "agency_url":
			agency.URL = NewURL(v)
		case "agency_timezone":
			agency.Timezone = NewTimezone(v)
		case "agency_lang":
			agency.Lang = NewLanguageCode(v)
		case "agency_phone":
			agency.Phone = NewPhoneNumber(v)
		case "agency_fare_url":
			agency.FareURL = NewURL(v)
		case "agency_email":
			agency.Email = NewEmail(v)
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
