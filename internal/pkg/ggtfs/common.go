package ggtfs

import (
	"encoding/csv"
	"errors"
	"fmt"
)

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

func createFieldError(fileName string, fieldName string, index int, err error) error {
	return errors.New(fmt.Sprintf("%s:%v: %s: %s", fileName, index, fieldName, err.Error()))
}
func createFileRowError(fileName string, row int, err string) error {
	return errors.New(fmt.Sprintf("%s:%v: %s", fileName, row, err))
}

func createFileError(fileName string, err string) error {
	return errors.New(fmt.Sprintf("%s: %s", fileName, err))
}

// getField and getOptionalField
// There are four cases for returning values from the row:
// For a MANDATORY (from the GTFS spec point of view) field:
//	- if the field is missing (it has a header, but no value on the row), an empty string is returned.
//  - if the field is present, a string containing the field value is returned, empty or not.
// For an OPTIONAL field:
//	- if the field is missing (it has a header, but no value on the row), a nil pointer is returned.
//  - if the field is present, a pointer to a string containing the field value is returned, the string being empty or not.
//
// In other words, there is a distinct difference between an empty value ("") on the row when its header is present, and
// a missing value (nil) when the header is present but the value is not. This is why we have two functions to handle them.
// You cannot make getOptionalField call getField and then check if the returned value was empty, because you cannot
// distinguish between a missing value and an empty value in getOptionalField by using getField's result (it will return
// empty string when the field is missing and when the field is empty). Stop trying to do that. =)

func getField(row []string, headerName string, headerPosition uint8, errs *[]error, lineNumber int, fileName string) string {
	if len(row) <= int(headerPosition) {
		*errs = append(*errs, createFileRowError(fileName, lineNumber, fmt.Sprintf("missing value for field: %s", headerName)))
		return ""
	}

	return row[headerPosition]
}

func getOptionalField(row []string, headerName string, headerPosition uint8, errs *[]error, lineNumber int, fileName string) *string {
	if len(row) <= int(headerPosition) {
		*errs = append(*errs, createFileRowError(fileName, lineNumber, fmt.Sprintf("missing value for field: %s", headerName)))
		return nil
	}

	return &row[headerPosition]
}

type entityCreator func(row []string, headers map[string]uint8, lineNumber int) (interface{}, []error)

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
func loadEntities(csvReader *csv.Reader, validHeaders []string, entityCreator entityCreator, fileName string) ([]interface{}, []error) {
	entities := make([]interface{}, 0)
	errs := make([]error, 0)

	headers, err := ReadHeaderRow(csvReader, validHeaders)
	if err != nil {
		errs = append(errs, createFileError(fileName, fmt.Sprintf("read error: %v", err.Error())))
		return entities, errs
	}
	if headers == nil {
		return entities, errs
	}

	lineNumber := 0
	for {
		row, err := ReadDataRow(csvReader)
		if err != nil {
			errs = append(errs, createFileError(fileName, fmt.Sprintf("%v", err.Error())))
			lineNumber++
			continue
		}
		if row == nil {
			break
		}
		if len(row) == 0 {
			continue
		}

		var entityCreateErrors []error
		entity, entityCreateErrors := entityCreator(row, headers, lineNumber)
		if len(entityCreateErrors) > 0 {
			errs = append(errs, entityCreateErrors...)
		}
		entities = append(entities, entity)

		lineNumber++
	}

	return entities, errs
}
