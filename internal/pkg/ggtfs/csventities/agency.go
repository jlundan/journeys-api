package ggtfs

import (
	"github.com/jlundan/journeys-api/internal/pkg/ggtfs/types"
)

type CsvAgency struct {
	Id         *string // agency_id 		(conditionally required)
	Name       *string // agency_name 		(required)
	URL        *string // agency_url 		(required)
	Timezone   *string // agency_timezone 	(required)
	Lang       *string // agency_lang 		(optional)
	Phone      *string // agency_phone 		(optional)
	FareURL    *string // agency_fare_url 	(optional)
	Email      *string // agency_email 		(optional)
	LineNumber int
}

func CreateCsvAgency(row []string, headers map[string]int, lineNumber int) *CsvAgency {
	agency := CsvAgency{
		LineNumber: lineNumber,
	}

	for hName := range headers {
		v := getRowValueForHeaderName(row, headers, hName)

		switch hName {
		case "agency_id":
			agency.Id = v
		case "agency_name":
			agency.Name = v
		case "agency_url":
			agency.URL = v
		case "agency_timezone":
			agency.Timezone = v
		case "agency_lang":
			agency.Lang = v
		case "agency_phone":
			agency.Phone = v
		case "agency_fare_url":
			agency.FareURL = v
		case "agency_email":
			agency.Email = v
		}
	}

	return &agency
}

type ValidatableCsvAgency struct {
	CsvAgency
}

func (a ValidatableCsvAgency) GetID() types.ID {
	return types.NewID(a.Id)
}

func (a ValidatableCsvAgency) GetName() types.Text {
	return types.NewText(a.Name)
}

func (a ValidatableCsvAgency) GetURL() types.URL {
	return types.NewURL(a.URL)
}

func (a ValidatableCsvAgency) GetTimezone() types.Timezone {
	return types.NewTimezone(a.Timezone)
}

func (a ValidatableCsvAgency) GetLang() types.LanguageCode {
	return types.NewLanguageCode(a.Lang)
}

func (a ValidatableCsvAgency) GetPhone() types.PhoneNumber {
	return types.NewPhoneNumber(a.Phone)
}

func (a ValidatableCsvAgency) GetFareURL() types.URL {
	return types.NewURL(a.FareURL)
}

func (a ValidatableCsvAgency) GetEmail() types.Email {
	return types.NewEmail(a.Email)
}

func CreateValidatableCsvAgency(row []string, headers map[string]int, lineNumber int) *ValidatableCsvAgency {
	return &ValidatableCsvAgency{*CreateCsvAgency(row, headers, lineNumber)}
}

func getRowValueForHeaderName(row []string, headers map[string]int, headerName string) *string {
	pos, ok := headers[headerName]

	if !ok {
		pos = -1
	}

	if pos < 0 || pos >= len(row) {
		return nil
	}

	return &row[pos]
}
