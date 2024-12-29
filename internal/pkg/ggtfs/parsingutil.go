package ggtfs

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

func LoadEntitiesFromCSV[T GtfsEntity](csvReader *csv.Reader, validHeaders []string, entityCreator entityCreator[T], fileName string) ([]T, []error) {
	var errs []error

	headers, indexingErrors := getHeaderIndex(csvReader, validHeaders)

	if len(indexingErrors) > 0 {
		for _, e := range indexingErrors {
			errs = append(errs, createFileError(fileName, fmt.Sprintf("%v", e.Error())))
		}
	}

	if len(headers) == 0 {
		return []T{}, errs
	}

	var entities []T

	lineNumber := 2
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			errs = append(errs, createFileError(fileName, fmt.Sprintf("%v", err.Error())))
			lineNumber++
			continue
		}

		entities = append(entities, entityCreator(row, headers, lineNumber))

		lineNumber++
	}

	return entities, errs
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

func getHeaderIndex(r *csv.Reader, validHeaderList []string) (map[string]int, []error) {
	headerRow, err := r.Read()
	if err == io.EOF {
		return map[string]int{}, []error{}
	}

	if err != nil {
		return map[string]int{}, []error{err}
	}

	var readErrors []error
	headerIndex := map[string]int{}
	encounteredHeaders := map[string]bool{}

	validHeaders := toSet(validHeaderList)

	for index, header := range headerRow {
		header = strings.TrimSpace(header)

		if encounteredHeaders[header] {
			readErrors = append(readErrors, fmt.Errorf("duplicate header name: %s", header))
			continue
		}

		if _, found := validHeaders[header]; !found {
			headerIndex[header] = -1
			continue
		}

		encounteredHeaders[header] = true
		headerIndex[header] = index
	}
	return headerIndex, readErrors
}

func toSet[T comparable](slice []T) map[T]struct{} {
	set := make(map[T]struct{}, len(slice))
	for _, item := range slice {
		set[item] = struct{}{}
	}
	return set
}
