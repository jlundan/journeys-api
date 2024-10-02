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

	if !a.Name.IsPresent() {
		validationErrors = append(validationErrors, createFileRowError(AgenciesFileName, a.LineNumber, createMissingMandatoryFieldString("agency_name")))
	} else if !a.Name.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(AgenciesFileName, a.LineNumber, createInvalidFieldString("agency_name")))
	}

	if !a.URL.IsPresent() {
		validationErrors = append(validationErrors, createFileRowError(AgenciesFileName, a.LineNumber, createMissingMandatoryFieldString("agency_url")))
	} else if !a.URL.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(AgenciesFileName, a.LineNumber, createInvalidFieldString("agency_url")))
	}

	if !a.Timezone.IsPresent() {
		validationErrors = append(validationErrors, createFileRowError(AgenciesFileName, a.LineNumber, createMissingMandatoryFieldString("agency_timezone")))
	} else if !a.Timezone.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(AgenciesFileName, a.LineNumber, createInvalidFieldString("agency_timezone")))
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

func CreateAgency(row []string, headers map[string]int, lineNumber int) (interface{}, []error) {
	var parseErrors []error

	agency := Agency{
		LineNumber: lineNumber,
	}

	for hName, hPos := range headers {
		switch hName {
		case "agency_id":
			agency.Id = NewID(getOptionalField(row, hName, hPos, &parseErrors, lineNumber, AgenciesFileName))
		case "agency_name":
			agency.Name = NewText(getOptionalField(row, hName, hPos, &parseErrors, lineNumber, AgenciesFileName))
		case "agency_url":
			agency.URL = NewURL(getOptionalField(row, hName, hPos, &parseErrors, lineNumber, AgenciesFileName))
		case "agency_timezone":
			agency.Timezone = NewTimezone(getOptionalField(row, hName, hPos, &parseErrors, lineNumber, AgenciesFileName))
		case "agency_lang":
			agency.Lang = NewOptionalLanguageCode(getOptionalField(row, hName, hPos, &parseErrors, lineNumber, AgenciesFileName))
		case "agency_phone":
			agency.Phone = NewOptionalPhoneNumber(getOptionalField(row, hName, hPos, &parseErrors, lineNumber, AgenciesFileName))
		case "agency_fare_url":
			agency.FareURL = NewOptionalURL(getOptionalField(row, hName, hPos, &parseErrors, lineNumber, AgenciesFileName))
		case "agency_email":
			agency.Email = NewOptionalEmail(getOptionalField(row, hName, hPos, &parseErrors, lineNumber, AgenciesFileName))
		}
	}

	if len(parseErrors) > 0 {
		return &agency, parseErrors
	}

	return &agency, nil
}

func ValidateAgencies(agencies []*Agency) ([]error, []string) {
	var validationErrors []error
	var recommendations []string

	if len(agencies) > 1 {
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
	} else if len(agencies) == 1 && agencies[0].Id.IsEmpty() {
		recommendations = append(recommendations, createFileRowRecommendation(AgenciesFileName, agencies[0].LineNumber, "it is recommended that agency_id is specified even when there is only one agency"))
		validationErrors = append(validationErrors, agencies[0].Validate()...)
	} else if len(agencies) == 1 {
		validationErrors = append(validationErrors, agencies[0].Validate()...)
	}

	return validationErrors, recommendations
}
