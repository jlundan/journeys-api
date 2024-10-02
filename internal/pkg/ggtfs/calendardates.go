package ggtfs

import (
	"encoding/csv"
	"fmt"
	"time"
)

// CalendarDate represents the GTFS calendar dates structure.
type CalendarDate struct {
	ServiceId     string // service_id
	Date          string // date
	ExceptionType string // exception_type
	LineNumber    int    // CSV row number
}

// List of valid GTFS calendar date headers.
var validCalendarDateHeaders = []string{"service_id", "date", "exception_type"}

// LoadCalendarDates reads and parses GTFS calendar dates data from the provided CSV reader.
func LoadCalendarDates(csvReader *csv.Reader) ([]*CalendarDate, []error) {
	entities, errs := loadEntities(csvReader, validCalendarDateHeaders, CreateCalendarDate, CalendarDatesFileName)

	calendarDates := make([]*CalendarDate, 0, len(entities))
	for _, entity := range entities {
		if calendarDate, ok := entity.(*CalendarDate); ok {
			calendarDates = append(calendarDates, calendarDate)
		}
	}

	return calendarDates, errs
}

// CreateCalendarDate creates a CalendarDate from a CSV row using the provided headers.
func CreateCalendarDate(row []string, headers map[string]int, lineNumber int) interface{} {
	var validationErrors []error

	// Create the CalendarDate struct and populate fields dynamically
	calendarDate := &CalendarDate{
		LineNumber: lineNumber,
	}

	// Populate the CalendarDate fields using the headers
	for hName, hPos := range headers {
		switch hName {
		case "service_id":
			calendarDate.ServiceId = getField(row, hName, hPos, &validationErrors, lineNumber, CalendarDatesFileName)
		case "date":
			calendarDate.Date = getField(row, hName, hPos, &validationErrors, lineNumber, CalendarDatesFileName)
		case "exception_type":
			calendarDate.ExceptionType = getField(row, hName, hPos, &validationErrors, lineNumber, CalendarDatesFileName)
		}
	}

	return calendarDate
}

// ValidateCalendarDates validates the parsed CalendarDate structs for logical consistency.
func ValidateCalendarDates(calendarDates []*CalendarDate, calendarItems []*CalendarItem) ([]error, []string) {
	var validationErrors []error
	var recommendations []string

	// Validate each calendar date
	for _, date := range calendarDates {
		if date == nil {
			continue
		}

		// Validate that service_id is not empty
		if date.ServiceId == "" {
			validationErrors = append(validationErrors, createFileRowError(CalendarDatesFileName, date.LineNumber, "service_id must not be empty"))
		}

		// Validate that date is not empty and is in correct format
		if date.Date == "" {
			validationErrors = append(validationErrors, createFileRowError(CalendarDatesFileName, date.LineNumber, "date must not be empty"))
		} else if _, err := parseCalendarDate(date.Date); err != nil {
			validationErrors = append(validationErrors, createFileRowError(CalendarDatesFileName, date.LineNumber, fmt.Sprintf("invalid date format: %s", date.Date)))
		}

		// Validate that exception_type is not empty and is either "1" or "2"
		if date.ExceptionType == "" {
			validationErrors = append(validationErrors, createFileRowError(CalendarDatesFileName, date.LineNumber, "exception_type must not be empty"))
		} else if date.ExceptionType != "1" && date.ExceptionType != "2" {
			validationErrors = append(validationErrors, createFileRowError(CalendarDatesFileName, date.LineNumber, fmt.Sprintf("exception_type must be '1' or '2', found: %s", date.ExceptionType)))
		}
	}

	// Cross-validate service_id references with CalendarItems (if provided)
	if calendarItems != nil {
		validateCalendarDateReferences(calendarDates, calendarItems, &validationErrors)
	}

	return validationErrors, recommendations
}

// validateCalendarDateReferences checks if the service_id in CalendarDate has a matching entry in CalendarItem.
func validateCalendarDateReferences(calendarDates []*CalendarDate, calendarItems []*CalendarItem, validationErrors *[]error) {
	// Create a map of CalendarItem service_ids for quick lookup
	serviceIDMap := make(map[string]struct{})
	//for _, item := range calendarItems {
	//	if item != nil && item.ServiceId != "" {
	//		serviceIDMap[item.ServiceId] = struct{}{}
	//	}
	//}

	// Check if each CalendarDate service_id is present in CalendarItems
	for _, calendarDate := range calendarDates {
		if calendarDate == nil || calendarDate.ServiceId == "" {
			continue
		}
		if _, found := serviceIDMap[calendarDate.ServiceId]; !found {
			*validationErrors = append(*validationErrors, createFileRowError(CalendarDatesFileName, calendarDate.LineNumber, fmt.Sprintf("referenced service_id '%s' not found in %s", calendarDate.ServiceId, CalendarFileName)))
		}
	}
}

// parseDate is a helper function that validates and parses the date in YYYYMMDD format.
func parseCalendarDate(dateStr string) (time.Time, error) {
	return time.Parse("20060102", dateStr)
}

//package ggtfs
//
//import (
//	"encoding/csv"
//	"errors"
//	"fmt"
//	"strconv"
//	"time"
//)
//
//type CalendarDate struct {
//	ServiceId     string
//	Date          time.Time
//	ExceptionType int
//	LineNumber    int
//}
//
//var validCalendarDateHeaders = []string{"service_id", "date", "exception_type"}
//
//func LoadCalendarDates(csvReader *csv.Reader) ([]*CalendarDate, []error) {
//	calendarDates := make([]*CalendarDate, 0)
//	errs := make([]error, 0)
//
//	headers, err := ReadHeaderRow(csvReader, validCalendarDateHeaders)
//	if err != nil {
//		errs = append(errs, createFileError(CalendarDatesFileName, fmt.Sprintf("read error: %v", err.Error())))
//		return calendarDates, errs
//	}
//	if headers == nil {
//		return calendarDates, errs
//	}
//
//	index := 0
//	for {
//		row, err := ReadDataRow(csvReader)
//		if err != nil {
//			errs = append(errs, createFileError(CalendarDatesFileName, fmt.Sprintf("%v", err.Error())))
//			index++
//			continue
//		}
//		if row == nil {
//			break
//		}
//
//		rowErrs := make([]error, 0)
//		calendarDate := CalendarDate{
//			LineNumber: index,
//		}
//
//		var (
//			serviceId     *string
//			date          *time.Time
//			exceptionType *int
//		)
//
//		for name, column := range headers {
//			switch name {
//			case "service_id":
//				serviceId = handleIDField(row[column], CalendarDatesFileName, name, index, &rowErrs)
//			case "date":
//				date = handleDateField(row[column], CalendarDatesFileName, name, index, false, &rowErrs)
//			case "exception_type":
//				exceptionType = handleExceptionTypeField(row[column], CalendarDatesFileName, name, index, &rowErrs)
//			}
//		}
//
//		if serviceId == nil {
//			rowErrs = append(rowErrs, createFileRowError(CalendarDatesFileName, calendarDate.LineNumber, "service_id must be specified"))
//		} else {
//			calendarDate.ServiceId = *serviceId
//		}
//
//		if date == nil {
//			rowErrs = append(rowErrs, createFileRowError(CalendarDatesFileName, calendarDate.LineNumber, "date must be specified"))
//		} else {
//			calendarDate.Date = *date
//		}
//
//		if exceptionType == nil {
//			rowErrs = append(rowErrs, createFileRowError(CalendarDatesFileName, calendarDate.LineNumber, "exception_type must be specified"))
//		} else {
//			calendarDate.ExceptionType = *exceptionType
//		}
//
//		if len(rowErrs) > 0 {
//			errs = append(errs, rowErrs...)
//		} else {
//			calendarDates = append(calendarDates, &calendarDate)
//		}
//
//		index++
//	}
//
//	return calendarDates, errs
//}
//
//func ValidateCalendarDates(calendarDates []*CalendarDate, calendarItems []*CalendarItem) []error {
//	var validationErrors []error
//
//	if calendarDates == nil {
//		return validationErrors
//	}
//
//	if calendarItems != nil {
//		for _, calendarDate := range calendarDates {
//			if calendarDate == nil {
//				continue
//			}
//			notFound := true
//			for _, calendarItem := range calendarItems {
//				if calendarItem == nil {
//					continue
//				}
//				if calendarDate.ServiceId == calendarItem.ServiceId {
//					notFound = false
//					break
//				}
//			}
//			if notFound {
//				validationErrors = append(validationErrors, createFileRowError(CalendarDatesFileName, calendarDate.LineNumber, fmt.Sprintf("referenced service_id not found in %s", CalendarFileName)))
//			}
//		}
//	}
//
//	return validationErrors
//}
//
//func handleExceptionTypeField(str string, fileName string, fieldName string, index int, errs *[]error) *int {
//	n, err := strconv.ParseInt(str, 10, 64)
//	if err != nil {
//		*errs = append(*errs, createFieldError(fileName, fieldName, index, err))
//		return nil
//	}
//
//	if v := int(n); v >= 1 && v <= 2 {
//		return &v
//	} else {
//		*errs = append(*errs, createFieldError(fileName, fieldName, index, errors.New(invalidValue)))
//		return nil
//	}
//}
