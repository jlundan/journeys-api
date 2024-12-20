package ggtfs

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
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

func createMissingMandatoryFieldString(fieldName string) string {
	return fmt.Sprintf("missing mandatory field: %s", fieldName)
}

type GtfsEntity interface {
	*Shape | *Stop | *Agency | *CalendarItem | *CalendarDate | *Route | *StopTime | *Trip
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

func validateFieldIsPresentAndValid(field ValidAndPresentField, fieldName string, lineNumber int, fileName string) []error {
	var validationErrors []error

	if !field.IsPresent() || field.IsEmpty() {
		validationErrors = append(validationErrors, createFileRowError(fileName, lineNumber, createMissingMandatoryFieldString(fieldName)))
	} else if !field.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(fileName, lineNumber, createInvalidFieldString(fieldName)))
	}

	return validationErrors
}

type ValidAndPresentField interface {
	IsValid() bool
	IsPresent() bool
	IsEmpty() bool
}

func getRowValue(row []string, position int) *string {
	if position < 0 || position >= len(row) {
		return nil
	}

	return &row[position]
}
