package ggtfs

import (
	"encoding/csv"
	"fmt"
)

type Frequency struct {
	TripId      string
	StartTime   string
	EndTime     string
	HeadwaySecs uint
	ExactTimes  *int
	LineNumber  int
}

func ExtractFrequencies(input *csv.Reader, output *csv.Writer, tripIds map[string]struct{}) []error {
	errs := make([]error, 0)

	headers, err := ReadHeaderRow(input)
	if err != nil {
		errs = append(errs, createFileError(FrequenciesFileName, fmt.Sprintf("read error: %v", err.Error())))
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
	if pos, columnExists := headers["trip_id"]; columnExists {
		idHeaderPos = pos
	} else {
		errs = append(errs, createFileError(FrequenciesFileName, "cannot extract frequencies without trip_id column"))
		return errs
	}

	for {
		row, rErr := ReadDataRow(input)
		if rErr != nil {
			errs = append(errs, createFileError(FrequenciesFileName, fmt.Sprintf("%v", rErr.Error())))
			continue
		}

		if row == nil { // EOF
			break
		}

		if _, shouldBeExtracted := tripIds[row[idHeaderPos]]; shouldBeExtracted {
			wErr := writeDataRow(row, output)
			if wErr != nil {
				errs = append(errs, wErr)
				return errs
			}
		}
	}

	return nil
}
