package ggtfs

import (
	"encoding/csv"
	"fmt"
)

type FareRule struct {
	Id            string
	RouteId       *string
	OriginId      *string
	DestinationId *string
	ContainsId    *string
	LineNumber    int
}

func ExtractFareRulesByFareIds(input *csv.Reader, output *csv.Writer, fareIds map[string]struct{}) []error {
	errs := make([]error, 0)

	headers, err := ReadHeaderRow(input)
	if err != nil {
		errs = append(errs, createFileError(FareRulesFileName, fmt.Sprintf("read error: %v", err.Error())))
		return errs
	}

	if headers == nil { // EOF
		return nil
	}

	err = writeHeaderRow(headers, output)
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	var idHeaderPos uint8
	if pos, columnExists := headers["fare_id"]; columnExists {
		idHeaderPos = pos
	} else {
		errs = append(errs, createFileError(FareRulesFileName, "cannot extract fares without fare_id column"))
		return errs
	}

	for {
		row, rErr := ReadDataRow(input)
		if rErr != nil {
			errs = append(errs, createFileError(FareRulesFileName, fmt.Sprintf("%v", rErr.Error())))
			continue
		}

		if row == nil { // EOF
			break
		}

		if _, shouldBeExtracted := fareIds[row[idHeaderPos]]; shouldBeExtracted {
			wErr := writeDataRow(row, output)
			if wErr != nil {
				errs = append(errs, wErr)
				return errs
			}
		}
	}

	return nil
}
