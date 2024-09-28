package ggtfs

import (
	"encoding/csv"
	"fmt"
)

type Agency struct {
	Id         string  // agency_id
	Name       string  // agency_name
	Url        string  // agency_url
	Timezone   string  // agency_timezone
	Lang       *string // agency_lang
	Phone      *string // agency_phone
	FareURL    *string // agency_fare_url
	Email      *string // agency_email
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
	var validationErrors []error

	agency := Agency{
		LineNumber: lineNumber,
	}

	for hName, hPos := range headers {
		switch hName {
		case "agency_id":
			agency.Id = getField(row, hName, hPos, &validationErrors, lineNumber, AgenciesFileName)
		case "agency_name":
			agency.Name = getField(row, hName, hPos, &validationErrors, lineNumber, AgenciesFileName)
		case "agency_url":
			agency.Url = getField(row, hName, hPos, &validationErrors, lineNumber, AgenciesFileName)
		case "agency_timezone":
			agency.Timezone = getField(row, hName, hPos, &validationErrors, lineNumber, AgenciesFileName)
		case "agency_lang":
			agency.Lang = getOptionalField(row, hName, hPos, &validationErrors, lineNumber, AgenciesFileName)
		case "agency_phone":
			agency.Phone = getOptionalField(row, hName, hPos, &validationErrors, lineNumber, AgenciesFileName)
		case "agency_fare_url":
			agency.FareURL = getOptionalField(row, hName, hPos, &validationErrors, lineNumber, AgenciesFileName)
		case "agency_email":
			agency.Email = getOptionalField(row, hName, hPos, &validationErrors, lineNumber, AgenciesFileName)
		}
	}

	if len(validationErrors) > 0 {
		return &agency, validationErrors
	}
	return &agency, nil
}

func ValidateAgencies(agencies []*Agency) []error {
	var validationErrors []error

	for _, agency := range agencies {
		if agency == nil {
			continue
		}

		if agency.Url == "" {
			validationErrors = append(validationErrors, createFileRowError(AgenciesFileName, agency.LineNumber, "agency_url must be specified"))
		}

		if agency.Timezone == "" {
			validationErrors = append(validationErrors, createFileRowError(AgenciesFileName, agency.LineNumber, "agency_timezone must be specified"))
		}
	}

	if len(agencies) > 1 {
		usedIds := make(map[string]bool)

		for _, a := range agencies {
			if a == nil {
				continue
			}

			if a.Id == "" {
				validationErrors = append(validationErrors, createFileRowError(AgenciesFileName, a.LineNumber, "agency_id must be specified when multiple agencies are declared"))
				continue
			}

			if usedIds[a.Id] {
				validationErrors = append(validationErrors, createFileRowError(AgenciesFileName, a.LineNumber, fmt.Sprintf("agency_id is not unique within the file")))
			} else {
				usedIds[a.Id] = true
			}
		}
	}

	return validationErrors
}
