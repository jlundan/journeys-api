package ggtfs

import (
	"encoding/csv"
	"fmt"
)

type Level struct {
	Id         string
	LevelIndex float64
	LevelName  string
	LineNumber int
}

func ExtractLevels(input *csv.Reader, output *csv.Writer, levelIds map[string]struct{}) []error {
	errs := make([]error, 0)

	headers, err := ReadHeaderRow(input)
	if err != nil {
		errs = append(errs, createFileError(LevelsFileName, fmt.Sprintf("read error: %v", err.Error())))
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
	if pos, columnExists := headers["level_id"]; columnExists {
		idHeaderPos = pos
	} else {
		errs = append(errs, createFileError(LevelsFileName, "cannot extract levels without level_id column"))
		return errs
	}

	for {
		row, rErr := ReadDataRow(input)
		if rErr != nil {
			errs = append(errs, createFileError(LevelsFileName, fmt.Sprintf("%v", rErr.Error())))
			continue
		}

		if row == nil { // EOF
			break
		}

		if _, shouldBeExtracted := levelIds[row[idHeaderPos]]; shouldBeExtracted {
			wErr := writeDataRow(row, output)
			if wErr != nil {
				errs = append(errs, wErr)
				return errs
			}
		}
	}

	return nil
}
