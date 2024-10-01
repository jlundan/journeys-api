package ggtfs

import (
	"encoding/csv"
	"fmt"
)

type Agency struct {
	Id         ID            // agency_id
	Name       Text          // agency_name
	Url        URL           // agency_url
	Timezone   Timezone      // agency_timezone
	Lang       *LanguageCode // agency_lang
	Phone      *PhoneNumber  // agency_phone
	FareURL    *URL          // agency_fare_url
	Email      *Email        // agency_email
	LineNumber int
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

func CreateAgency(row []string, headers map[string]uint8, lineNumber int) (interface{}, []error) {
	var parseErrors []error

	agency := Agency{
		LineNumber: lineNumber,
	}

	for hName, hPos := range headers {
		switch hName {
		case "agency_id":
			agency.Id = NewID(getField(row, hName, hPos, &parseErrors, lineNumber, AgenciesFileName))
		case "agency_name":
			agency.Name = NewText(getField(row, hName, hPos, &parseErrors, lineNumber, AgenciesFileName))
		case "agency_url":
			agency.Url = NewURL(getField(row, hName, hPos, &parseErrors, lineNumber, AgenciesFileName))
		case "agency_timezone":
			agency.Timezone = NewTimezone(getField(row, hName, hPos, &parseErrors, lineNumber, AgenciesFileName))
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

func ValidateAgencies(agencies []*Agency) []error {
	var validationErrors []error

	for _, agency := range agencies {
		if agency == nil {
			continue
		}

		if agency.Url.IsEmpty() {
			validationErrors = append(validationErrors, createFileRowError(AgenciesFileName, agency.LineNumber, "agency_url must be specified"))
		}

		if agency.Timezone.IsEmpty() {
			validationErrors = append(validationErrors, createFileRowError(AgenciesFileName, agency.LineNumber, "agency_timezone must be specified"))
		}
	}

	if len(agencies) > 1 {
		usedIds := make(map[string]bool)

		for _, a := range agencies {
			if a == nil {
				continue
			}

			if a.Id.IsEmpty() {
				validationErrors = append(validationErrors, createFileRowError(AgenciesFileName, a.LineNumber, "agency_id must be specified when multiple agencies are declared"))
				continue
			}

			if usedIds[a.Id.String()] {
				validationErrors = append(validationErrors, createFileRowError(AgenciesFileName, a.LineNumber, fmt.Sprintf("agency_id is not unique within the file")))
			} else {
				usedIds[a.Id.String()] = true
			}
		}
	}

	return validationErrors
}
