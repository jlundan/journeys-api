package csventities

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

func NewReader(r *csv.Reader) *GtfsCsvReader {
	return &GtfsCsvReader{
		csvReader: r,

		FailOnHeaderErrors: true,
		SkipRowsWithErrors: true,

		AgencyHeaders: defaultAgencyHeaders,
	}
}

func LoadAgencies(reader *GtfsCsvReader) ([]*CsvAgency, []error) {
	return loadCsvEntities[*CsvAgency]("agency", reader, CreateCsvAgency)
}

func loadCsvEntities[T CsvEntity](entityType string, reader *GtfsCsvReader, entityCreator csvEntityCreator[T]) ([]T, []error) {
	var errs []error

	headerNames, hErr := getEntityHeaders(entityType, reader)
	if hErr != nil {
		return []T{}, []error{hErr}
	}

	headers, indexingErrors := getHeaderIndex(reader.csvReader, headerNames)

	if len(indexingErrors) > 0 {
		for _, err := range indexingErrors {
			errs = append(errs, fmt.Errorf("line 1: %v", err.Error()))
		}

		if reader.FailOnHeaderErrors {
			return []T{}, errs
		}
	}

	if len(headers) == 0 {
		return []T{}, errs
	}

	var entities []T

	lineNumber := 2
	for {
		row, rErr := reader.csvReader.Read()
		if rErr == io.EOF {
			break
		}

		if rErr != nil {
			errs = append(errs, fmt.Errorf("line %d: %v", lineNumber, rErr.Error()))

			if reader.SkipRowsWithErrors {
				lineNumber++
				continue
			}
		}

		entities = append(entities, entityCreator(row, headers, lineNumber))

		lineNumber++
	}

	return entities, errs
}

func getEntityHeaders(entityType string, reader *GtfsCsvReader) ([]string, error) {
	switch entityType {
	case "agency":
		return reader.AgencyHeaders, nil
	default:
		return []string{}, fmt.Errorf("unknown entity: %s", entityType)
	}
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

type GtfsCsvReader struct {
	csvReader          *csv.Reader
	FailOnHeaderErrors bool
	SkipRowsWithErrors bool
	AgencyHeaders      []string
}

type csvEntityCreator[T CsvEntity] func(row []string, headers map[string]int, lineNumber int) T

type CsvEntity interface {
	*CsvAgency | any
}

var defaultAgencyHeaders = []string{"agency_id", "agency_name", "agency_url", "agency_timezone",
	"agency_lang", "agency_phone", "agency_fare_url", "agency_email"}
