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

// jlundan 6.4.2024 commenting code below out for now, since we have no use for it (yet). This was used in the prototype
// to extract agencies from the GTFS file, but that functionality is not imported yet.

//func ExtractAgencies(input *csv.Reader, output *csv.Writer, agencies map[string]struct{}) []error {
//	errs := make([]error, 0)
//
//	headers, err := ReadHeaderRow(input)
//	if err != nil {
//		errs = append(errs, createFileError(AgenciesFileName, fmt.Sprintf("read error: %v", err.Error())))
//		return errs
//	}
//
//	if headers == nil { // EOF
//		// This means empty file -> no errors and nothing to extract
//		return nil
//	}
//
//	err = writeHeaderRow(headers, output)
//	if err != nil {
//		errs = append(errs, err)
//		return errs
//	}
//
//	// GTFS states that agency_id should be present if multiple agencies are used in the data set. If the agency_id
//	// is not present, we either have invalid data set or data set with single agency. In any case, we cannot proceed
//	// since anything we extract would be based on the agency_id field.
//	var idHeaderPos uint8
//	if agencyIdPos, hasAgencyId := headers["agency_id"]; hasAgencyId {
//		idHeaderPos = agencyIdPos
//	} else {
//		errs = append(errs, createFileError(AgenciesFileName, "cannot extract agency without agency_id column"))
//		return errs
//	}
//
//	// Before extracting the agency, we should check that we actually have multiple agencies declared in the file.
//	// If we do not, then it does not make sense to extract anything, since the extract result would be exactly the
//	// same data set. To do this, we keep track of the agencies we read, and compare each to the previous agency_id
//	// we read and check if it is the same. When we have finished the file, we know if any of the agency_ids were
//	// different from the previous one (did not have single agency id). While we are checking this, we also store
//	// the rows (for efficiency), so we can write those out if we actually had multiple agencies in the file.
//	hadSingleAgency := true
//	previousAgency := ""
//	rows := make([][]string, 0)
//	for {
//		row, rErr := ReadDataRow(input)
//		if rErr != nil {
//			errs = append(errs, createFileError(AgenciesFileName, fmt.Sprintf("%v", rErr.Error())))
//			continue
//		}
//
//		if row == nil { // EOF
//			break
//		}
//
//		if previousAgency != "" && previousAgency != row[idHeaderPos] {
//			hadSingleAgency = false
//		}
//		previousAgency = row[idHeaderPos]
//
//		if _, shouldBeExtracted := agencies[row[idHeaderPos]]; shouldBeExtracted {
//			rows = append(rows, row)
//		}
//	}
//
//	if hadSingleAgency {
//		errs = append(errs, createFileError(AgenciesFileName, "only one agency specified - no need to extract anything"))
//		return errs
//	}
//
//	for _, row := range rows {
//		wErr := output.Write(row)
//		if wErr != nil {
//			output.Flush()
//			errs = append(errs, wErr)
//			return errs
//		}
//		output.Flush()
//	}
//
//	return nil
//}

func LoadAgencies(csvReader *csv.Reader) ([]*Agency, []error) {
	agencies := make([]*Agency, 0)
	errs := make([]error, 0)

	headers, err := ReadHeaderRow(csvReader)
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
