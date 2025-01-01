package csventities

import (
	"github.com/jlundan/journeys-api/internal/pkg/ggtfs/types"
)

type RawCsvAgency struct {
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

func CreateRawCsvAgency(row []string, headers map[string]int, lineNumber int) *RawCsvAgency {
	agency := RawCsvAgency{
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

type CsvAgency struct {
	RawCsvAgency
}

func (a CsvAgency) GetID() types.ID {
	return types.NewID(a.Id)
}

func (a CsvAgency) GetName() types.Text {
	return types.NewText(a.Name)
}

func (a CsvAgency) GetURL() types.URL {
	return types.NewURL(a.URL)
}

func (a CsvAgency) GetTimezone() types.Timezone {
	return types.NewTimezone(a.Timezone)
}

func (a CsvAgency) GetLang() types.LanguageCode {
	return types.NewLanguageCode(a.Lang)
}

func (a CsvAgency) GetPhone() types.PhoneNumber {
	return types.NewPhoneNumber(a.Phone)
}

func (a CsvAgency) GetFareURL() types.URL {
	return types.NewURL(a.FareURL)
}

func (a CsvAgency) GetEmail() types.Email {
	return types.NewEmail(a.Email)
}

func CreateCsvAgency(row []string, headers map[string]int, lineNumber int) *CsvAgency {
	return &CsvAgency{*CreateRawCsvAgency(row, headers, lineNumber)}
}
