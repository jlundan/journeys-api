package ggtfs

import (
	"encoding/csv"
	"fmt"
)

type FareAttributes struct {
	Id               string
	Price            float64
	CurrencyType     string
	PaymentMethod    int
	Transfers        int
	AgencyId         *string
	TransferDuration *uint
	LineNumber       int
}

func ExtractFareAttributesByAgencies(input *csv.Reader, output *csv.Writer, agencyIds map[string]struct{}) (map[string]struct{}, []error) {
	errs := make([]error, 0)

	headers, err := ReadHeaderRow(input)
	if err != nil {
		errs = append(errs, createFileError(FareAttributesFileName, fmt.Sprintf("read error: %v", err.Error())))
		return nil, errs
	}

	if headers == nil { // EOF
		return nil, nil
	}

	err = writeHeaderRow(headers, output)
	if err != nil {
		errs = append(errs, err)
		return nil, errs
	}

	var idHeaderPos uint8
	if pos, columnExists := headers["agency_id"]; columnExists {
		idHeaderPos = pos
	} else {
		errs = append(errs, createFileError(FareAttributesFileName, "cannot extract agencies without agency_id column"))
		return nil, errs
	}

	var fareIdPos uint8
	if pos, columnExists := headers["fare_id"]; columnExists {
		fareIdPos = pos
	} else {
		errs = append(errs, createFileError(FareAttributesFileName, "cannot extract fares without fare_id column"))
		return nil, errs
	}

	fareIdMap := make(map[string]struct{})
	for {
		row, rErr := ReadDataRow(input)
		if rErr != nil {
			errs = append(errs, createFileError(FareAttributesFileName, fmt.Sprintf("%v", rErr.Error())))
			continue
		}

		if row == nil { // EOF
			break
		}

		if _, shouldBeExtracted := agencyIds[row[idHeaderPos]]; shouldBeExtracted {
			wErr := writeDataRow(row, output)
			if wErr != nil {
				errs = append(errs, wErr)
				return nil, errs
			}

			fareId := row[fareIdPos]
			if _, tripAlreadyExists := fareIdMap[fareId]; !tripAlreadyExists {
				fareIdMap[fareId] = struct{}{}
			}
		}
	}

	return fareIdMap, nil
}
