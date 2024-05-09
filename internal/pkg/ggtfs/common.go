package ggtfs

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
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

func handleIDField(str string, fileName string, fieldName string, index int, errs *[]error) *string {
	return handleStringField(str, fileName, fieldName, index, errs)
}

func handleTextField(str string, fileName string, fieldName string, index int, errs *[]error) *string {
	return handleStringField(str, fileName, fieldName, index, errs)
}

func handleTimeZoneField(str string, fileName string, fieldName string, index int, errs *[]error) *string {
	return handleStringField(str, fileName, fieldName, index, errs)
}

func handleLanguageCodeField(str string, fileName string, fieldName string, index int, errs *[]error) *string {
	return handleStringField(str, fileName, fieldName, index, errs)
}

func handlePhoneNumberField(str string, fileName string, fieldName string, index int, errs *[]error) *string {
	return handleStringField(str, fileName, fieldName, index, errs)
}

func handleEmailField(str string, fileName string, fieldName string, index int, errs *[]error) *string {
	return handleStringField(str, fileName, fieldName, index, errs)
}

func handleIntField(str string, fileName string, fieldName string, index int, errs *[]error) *int {
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, err))
		return nil
	}
	so := int(val)
	return &so
}

func handleContinuousPickupField(str string, fileName string, fieldName string, index int, errs *[]error) *int {
	if str == "" {
		c := 1
		return &c
	}

	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, err))
		return nil
	}

	if cpt := int(n); cpt == 0 || (cpt >= 1 && cpt <= 3) {
		return &cpt
	} else {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, errors.New(invalidValue)))
		return nil
	}
}

func handleContinuousDropOffField(str string, fileName string, fieldName string, index int, errs *[]error) *int {
	return handleContinuousPickupField(str, fileName, fieldName, index, errs)
}

func handleFloat64Field(str string, fileName string, fieldName string, index int, errs *[]error) *float64 {
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, err))
		return nil
	}
	return &val
}

func handleDateField(str string, fileName string, fieldName string, index int, fillEnd bool, errs *[]error) *time.Time {
	if str == "" {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, errors.New(emptyValueNotAllowed)))
		return nil
	}

	if len(str) < 8 {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, errors.New("invalid date format")))
		return nil
	}

	year, err := strconv.ParseInt(str[:4], 10, 64)
	if err != nil {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, err))
		return nil
	}

	month, err := strconv.ParseInt(str[4:6], 10, 64)
	if err != nil {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, err))
		return nil
	}

	day, err := strconv.ParseInt(str[6:8], 10, 64)
	if err != nil {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, err))
		return nil
	}

	var d time.Time
	if fillEnd {
		d = time.Date(int(year), time.Month(int(month)), int(day), 23, 59, 59, 0, time.FixedZone("UTC+2", 2*60*60))
	} else {
		d = time.Date(int(year), time.Month(int(month)), int(day), 0, 0, 0, 0, time.FixedZone("UTC+2", 2*60*60))
	}
	return &d
}

func handleTimeField(str string, fileName string, fieldName string, index int, errs *[]error) *string {
	return handleStringField(str, fileName, fieldName, index, errs)
}

func handleStringField(str string, fileName string, fieldName string, index int, errs *[]error) *string {
	if str == "" {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, errors.New(emptyValueNotAllowed)))
		return nil
	}

	return &str
}

func handleColorField(str string, fileName string, fieldName string, index int, errs *[]error) *string {
	if str == "" {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, errors.New(emptyValueNotAllowed)))
		return nil
	}

	_, err := hex.DecodeString(fmt.Sprintf("%sFF", str))
	if err != nil {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, err))
		return nil
	}

	return &str
}

func handleURLField(str string, fileName string, fieldName string, index int, errs *[]error) *string {
	_, err := url.Parse(str)
	if err != nil {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, err))
		return nil
	}
	return &str
}

func getRowValue(row []string, headers map[string]uint8, header string, errs []error, rowIndex int, fileName string) string {
	headerPosition, ok := headers[header]

	if !ok {
		// Requested header not specified in the CSV file, return empty string
		errs = append(errs, createFileRowError(fileName, rowIndex, fmt.Sprintf("invalid header: %s", header)))
		return ""
	}

	if len(row) <= int(headerPosition) {
		// Requested header is found but the row does not have enough columns, return empty string
		errs = append(errs, createFileRowError(fileName, rowIndex, fmt.Sprintf("missing required field: %s", header)))
		return ""
	}

	return row[headerPosition]
}
