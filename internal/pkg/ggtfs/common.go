package ggtfs

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"
)

//goland:noinspection GoUnusedConst
const (
	AgenciesFileName      = "agency.txt"
	RoutesFileName        = "routes.txt"
	StopsFileName         = "stops.txt"
	TripsFileName         = "trips.txt"
	StopTimesFileName     = "stop_times.txt"
	CalendarFileName      = "calendar.txt"
	CalendarDatesFileName = "calendar_dates.txt"
	ShapesFileName        = "shapes.txt"
	//These are not used yet, but are part of the spec
	//FareAttributesFileName = "fare_attributes.txt"
	//FareRulesFileName      = "fare_rules.txt"
	//FrequenciesFileName    = "frequencies.txt"
	//TransfersFileName      = "transfers.txt"
	//PathwaysFileName       = "pathways.txt"
	//LevelsFileName         = "levels.txt"

	emptyValueNotAllowed = "empty value not allowed"
	invalidValue         = "invalid value"
	nonUniqueId          = "non-unique id"
)

func createFileRowError(fileName string, row int, err string) error {
	return errors.New(fmt.Sprintf("%s:%v: %s", fileName, row, err))
}

func createFileRowRecommendation(fileName string, row int, err string) string {
	return fmt.Sprintf("%s:%v: %s", fileName, row, err)
}

func createFileError(fileName string, err string) error {
	return errors.New(fmt.Sprintf("%s: %s", fileName, err))
}

func createInvalidFieldString(fieldName string) string {
	return fmt.Sprintf("invalid field: %s", fieldName)
}

func createInvalidRequiredFieldString(fieldName string) string {
	return fmt.Sprintf("invalid mandatory field: %s", fieldName)
}

type GtfsEntity interface {
	*Shape | *Stop | *Agency | *CalendarItem | *CalendarDate | *Route | *StopTime | *Trip | any
}

type entityCreator[T GtfsEntity] func(row []string, headers map[string]int, lineNumber int) T

// LoadEntities is a generic function for loading entities from a CSV file using a provided entity creation callback.
// This function reads each row from the given CSV reader, creates entities using the provided callback function, and
// collects any errors encountered during the loading process.
//
// Parameters:
//   - csvReader (*csv.Reader): The CSV reader from which to read the data.
//   - validHeaders ([]string): A list of expected column headers to validate against the CSV header row.
//   - entityCreator (EntityCreator): A callback function that is called for each data row to create an entity.
//     The function has the following signature:
//     func(row []string, headers map[string]int, lineNumber int) (interface{}, error)
//   - row: Represents a single row of data from the CSV.
//   - headers: A map of header names to their corresponding index positions in the row.
//   - lineNumber: The current row number being processed.
//     The callback should return the created entity and an error, if any.
//   - fileName (string): The name of the file being processed. Used in error reporting.
//
// Returns:
//   - ([]interface{}): A slice of entities created from the CSV file. The entities are returned as generic interfaces,
//     and should be type-asserted by the caller to their concrete types.
//   - ([]error): A slice of errors encountered during the loading process, including errors from reading the CSV,
//     missing or invalid headers, and errors returned from the entity creation callback.
func LoadEntities[T GtfsEntity](csvReader *csv.Reader, validHeaders []string, entityCreator entityCreator[T], fileName string) ([]T, []error) {
	entities := make([]T, 0)
	errs := make([]error, 0)

	headers, err := ReadHeaderRow(csvReader, validHeaders)
	if err == io.EOF {
		return entities, errs
	}

	if err != nil {
		errs = append(errs, createFileError(fileName, fmt.Sprintf("read error: %v", err.Error())))
		return entities, errs
	}
	if headers == nil {
		return entities, errs
	}

	lineNumber := 0
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

type ValidAndPresentField interface {
	IsValid() bool
	IsPresent() bool
}

func getRowValue(row []string, position int) *string {
	if position < 0 || position >= len(row) {
		return nil
	}

	return &row[position]
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
