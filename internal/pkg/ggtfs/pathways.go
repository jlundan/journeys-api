package ggtfs

import (
	"encoding/csv"
	"fmt"
)

type Pathway struct {
	Id                  string
	FromStopId          string
	ToStopId            string
	PathwayMode         uint
	IsBidirectional     uint
	Length              float64
	TraversalTime       uint
	StairCount          uint
	MaxSlope            float64
	MinWidth            float64
	SignpostedAs        string
	ReverseSignpostedAs string
	LineNumber          int
}

func ExtractPathways(input *csv.Reader, output *csv.Writer, stopIds map[string]struct{}) []error {
	errs := make([]error, 0)

	headers, err := ReadHeaderRow(input)
	if err != nil {
		errs = append(errs, createFileError(PathwaysFileName, fmt.Sprintf("read error: %v", err.Error())))
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

	var fromStopIdHeaderPos uint8
	if pos, columnExists := headers["from_stop_id"]; columnExists {
		fromStopIdHeaderPos = pos
	} else {
		errs = append(errs, createFileError(PathwaysFileName, "cannot extract pathways without from_stop_id column"))
		return errs
	}

	var toStopIdHeaderPos uint8
	if pos, columnExists := headers["to_stop_id"]; columnExists {
		toStopIdHeaderPos = pos
	} else {
		errs = append(errs, createFileError(PathwaysFileName, "cannot extract pathways without to_stop_id column"))
		return errs
	}

	for {
		row, rErr := ReadDataRow(input)
		if rErr != nil {
			errs = append(errs, createFileError(PathwaysFileName, fmt.Sprintf("%v", rErr.Error())))
			continue
		}

		if row == nil { // EOF
			break
		}

		_, fromStopExists := stopIds[row[fromStopIdHeaderPos]]
		_, toStopExists := stopIds[row[toStopIdHeaderPos]]
		if fromStopExists && toStopExists {
			wErr := writeDataRow(row, output)
			if wErr != nil {
				errs = append(errs, wErr)
				return errs
			}
		}
	}

	return nil
}
