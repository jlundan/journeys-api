package ggtfs

import (
	"encoding/csv"
	"fmt"
)

type Agency struct {
	Id         string // agency_id
	Name       string // agency_name
	Url        string // agency_url
	Timezone   string // agency_timezone
	Lang       string // agency_lang
	Phone      string // agency_phone
	FareURL    string // agency_fare_url
	Email      string // agency_email
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

		if len(row) == 0 {
			continue
		}

		agency := Agency{
			lineNumber: index,
			Id:         getRowValue(row, headers, "agency_id", errs, index, AgenciesFileName),
			Name:       getRowValue(row, headers, "agency_name", errs, index, AgenciesFileName),
			Url:        getRowValue(row, headers, "agency_url", errs, index, AgenciesFileName),
			Timezone:   getRowValue(row, headers, "agency_timezone", errs, index, AgenciesFileName),
			Lang:       getRowValue(row, headers, "agency_lang", errs, index, AgenciesFileName),
			Phone:      getRowValue(row, headers, "agency_phone", errs, index, AgenciesFileName),
			FareURL:    getRowValue(row, headers, "agency_fare_url", errs, index, AgenciesFileName),
			Email:      getRowValue(row, headers, "agency_email", errs, index, AgenciesFileName),
		}

		if agency.Url == "" {
			errs = append(errs, createFileRowError(AgenciesFileName, agency.lineNumber, "agency_url must be specified"))
		}

		if agency.Timezone == "" {
			errs = append(errs, createFileRowError(AgenciesFileName, agency.lineNumber, "agency_timezone must be specified"))
		}

		agencies = append(agencies, &agency)
		index++
	}

	if len(agencies) > 1 {
		usedIds := make([]string, 0)

		for _, a := range agencies {
			if a.Id == "" {
				errs = append(errs, createFileRowError(AgenciesFileName, a.lineNumber, "agency id must be specified when multiple agencies are declared"))
				continue
			}

			if StringArrayContainsItem(usedIds, a.Id) {
				errs = append(errs, createFileRowError(AgenciesFileName, index, fmt.Sprintf("%s: agency_id", nonUniqueId)))
			} else {
				usedIds = append(usedIds, a.Id)
			}
		}
	}

	return agencies, errs
}
