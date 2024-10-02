package ggtfs

import (
	"encoding/csv"
	"fmt"
)

type Agency struct {
	Id         ID            // agency_id
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

	// These should not be implemented with the above method, since checking nil values with Golang interfaces, would
	// require using reflection, which is way too slow.

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

var validAgencyHeaders = []string{"agency_id", "agency_name", "agency_url", "agency_timezone",
	"agency_lang", "agency_phone", "agency_fare_url", "agency_email"}

func LoadAgencies(csvReader *csv.Reader) ([]*Agency, []error) {
	entities, errs := loadEntities(csvReader, validAgencyHeaders, CreateAgency, AgenciesFileName)

	agencies := make([]*Agency, 0, len(entities))

	for _, entity := range entities {
		if agency, ok := entity.(*Agency); ok {
			agencies = append(agencies, agency)
		}
	}

	return agencies, errs
}

func CreateAgency(row []string, headers map[string]int, lineNumber int) interface{} {
	agency := Agency{
		LineNumber: lineNumber,
	}

	for hName, hPos := range headers {
		switch hName {
		case "agency_id":
			agency.Id = NewID(&row[hPos])
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

	if len(agencies) > 1 {
		usedIds := make(map[string]bool)

		for _, a := range agencies {
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
	} else if len(agencies) == 1 && agencies[0].Id.IsEmpty() {
		recommendations = append(recommendations, createFileRowRecommendation(AgenciesFileName, agencies[0].LineNumber, "it is recommended that agency_id is specified even when there is only one agency"))
		validationErrors = append(validationErrors, agencies[0].Validate()...)
	} else if len(agencies) == 1 {
		validationErrors = append(validationErrors, agencies[0].Validate()...)
	}

	return validationErrors, recommendations
}
