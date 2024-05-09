package ggtfs

import (
	"encoding/csv"
	"fmt"
)

type Agency struct {
	Id         *string // agency_id
	Name       string  // agency_name
	Url        string  // agency_url
	Timezone   string  // agency_timezone
	Lang       *string // agency_lang
	Phone      *string // agency_phone
	FareURL    *string // agency_fare_url
	Email      *string // agency_email
	lineNumber int
}

var validAgencyHeaders = []string{"agency_id", "agency_name", "agency_url", "agency_timezone",
	"agency_lang", "agency_phone", "agency_fare_url", "agency_email"}

func LoadAgencies(csvReader *csv.Reader) ([]*Agency, []error) {
	agencies := make([]*Agency, 0)
	errs := make([]error, 0)

	headers, err := ReadHeaderRow(csvReader, validAgencyHeaders)
	if err != nil {
		errs = append(errs, createFileError(AgenciesFileName, fmt.Sprintf("read error: %v", err.Error())))
		return agencies, errs
	}
	if headers == nil {
		return agencies, errs
	}

	usedIds := make([]string, 0)
	index := 0
	for {
		row, err := ReadDataRow(csvReader)
		if err != nil {
			errs = append(errs, createFileError(AgenciesFileName, fmt.Sprintf("%v", err.Error())))
			index++
			continue
		}
		if row == nil {
			break
		}

		rowErrs := make([]error, 0)
		agency := Agency{
			lineNumber: index,
		}

		var agencyName, agencyUrl, agencyTimezone *string

		for name, column := range headers {
			switch name {
			case "agency_id":
				agency.Id = handleIDField(row[column], AgenciesFileName, name, index, &rowErrs)
			case "agency_name":
				agencyName = handleTextField(row[column], AgenciesFileName, name, index, &rowErrs)
			case "agency_url":
				agencyUrl = handleURLField(row[column], AgenciesFileName, name, index, &rowErrs)
			case "agency_timezone":
				agencyTimezone = handleTimeZoneField(row[column], AgenciesFileName, name, index, &rowErrs)
			case "agency_lang":
				agency.Lang = handleLanguageCodeField(row[column], AgenciesFileName, name, index, &rowErrs)
			case "agency_phone":
				agency.Phone = handlePhoneNumberField(row[column], AgenciesFileName, name, index, &rowErrs)
			case "agency_fare_url":
				agency.FareURL = handleURLField(row[column], AgenciesFileName, name, index, &rowErrs)
			case "agency_email":
				agency.Email = handleEmailField(row[column], AgenciesFileName, name, index, &rowErrs)
			}
		}

		if agency.Id != nil {
			if StringArrayContainsItem(usedIds, *agency.Id) {
				errs = append(errs, createFileRowError(AgenciesFileName, index, fmt.Sprintf("%s: agency_id", nonUniqueId)))
			} else {
				usedIds = append(usedIds, *agency.Id)
			}
		}

		if agencyName == nil {
			rowErrs = append(rowErrs, createFileRowError(AgenciesFileName, agency.lineNumber, "agency_name must be specified"))
		} else {
			agency.Name = *agencyName
		}

		if agencyUrl == nil {
			rowErrs = append(rowErrs, createFileRowError(AgenciesFileName, agency.lineNumber, "agency_url must be specified"))
		} else if *agencyUrl == "" {
			rowErrs = append(rowErrs, createFileRowError(AgenciesFileName, agency.lineNumber, "agency_url must not be empty"))
		} else {
			agency.Url = *agencyUrl
		}

		if agencyTimezone == nil {
			rowErrs = append(rowErrs, createFileRowError(AgenciesFileName, agency.lineNumber, "agency_timezone must be specified"))
		} else {
			agency.Timezone = *agencyTimezone
		}

		if len(rowErrs) > 0 {
			errs = append(errs, rowErrs...)
		} else {
			agencies = append(agencies, &agency)
		}

		index++
	}

	if len(agencies) > 1 {
		for _, a := range agencies {
			if a.Id == nil {
				errs = append(errs, createFileRowError(AgenciesFileName, a.lineNumber, "agency id must be specified when multiple agencies are declared"))
			}
		}
	}

	return agencies, errs
}
