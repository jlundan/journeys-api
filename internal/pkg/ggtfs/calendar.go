package ggtfs

import (
	"encoding/csv"
	"fmt"
)

// CalendarItem represents the GTFS calendar file structure.
type CalendarItem struct {
	ServiceId  string // service_id
	Monday     string // monday
	Tuesday    string // tuesday
	Wednesday  string // wednesday
	Thursday   string // thursday
	Friday     string // friday
	Saturday   string // saturday
	Sunday     string // sunday
	StartDate  string // start_date
	EndDate    string // end_date
	LineNumber int    // CSV row number
}

// List of valid GTFS calendar headers.
var validCalendarHeaders = []string{
	"service_id", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday", "start_date", "end_date",
}

// LoadCalendarItems reads and parses GTFS calendar data from the provided CSV reader.
func LoadCalendarItems(csvReader *csv.Reader) ([]*CalendarItem, []error) {
	entities, errs := loadEntities(csvReader, validCalendarHeaders, CreateCalendarItem, CalendarFileName)

	calendarItems := make([]*CalendarItem, 0, len(entities))
	for _, entity := range entities {
		if calendarItem, ok := entity.(*CalendarItem); ok {
			calendarItems = append(calendarItems, calendarItem)
		}
	}

	return calendarItems, errs
}

// CreateCalendarItem creates a CalendarItem from a CSV row, using the provided headers.
func CreateCalendarItem(row []string, headers map[string]int, lineNumber int) interface{} {
	var validationErrors []error

	// Create the CalendarItem struct and populate fields dynamically
	calendarItem := &CalendarItem{
		LineNumber: lineNumber,
	}

	// Populate the CalendarItem fields using the headers
	for hName, hPos := range headers {
		switch hName {
		case "service_id":
			calendarItem.ServiceId = getField(row, hName, hPos, &validationErrors, lineNumber, CalendarFileName)
		case "monday":
			calendarItem.Monday = getField(row, hName, hPos, &validationErrors, lineNumber, CalendarFileName)
		case "tuesday":
			calendarItem.Tuesday = getField(row, hName, hPos, &validationErrors, lineNumber, CalendarFileName)
		case "wednesday":
			calendarItem.Wednesday = getField(row, hName, hPos, &validationErrors, lineNumber, CalendarFileName)
		case "thursday":
			calendarItem.Thursday = getField(row, hName, hPos, &validationErrors, lineNumber, CalendarFileName)
		case "friday":
			calendarItem.Friday = getField(row, hName, hPos, &validationErrors, lineNumber, CalendarFileName)
		case "saturday":
			calendarItem.Saturday = getField(row, hName, hPos, &validationErrors, lineNumber, CalendarFileName)
		case "sunday":
			calendarItem.Sunday = getField(row, hName, hPos, &validationErrors, lineNumber, CalendarFileName)
		case "start_date":
			calendarItem.StartDate = getField(row, hName, hPos, &validationErrors, lineNumber, CalendarFileName)
		case "end_date":
			calendarItem.EndDate = getField(row, hName, hPos, &validationErrors, lineNumber, CalendarFileName)
		}
	}

	return calendarItem
}

// ValidateCalendarItems validates the parsed CalendarItem structs for logical consistency.
func ValidateCalendarItems(calendarItems []*CalendarItem) ([]error, []string) {
	var validationErrors []error
	var recommendations []string

	for _, item := range calendarItems {
		if item == nil {
			continue
		}

		// Check for required fields' logical content, e.g., days of the week should be 0 or 1.
		if item.ServiceId == "" {
			validationErrors = append(validationErrors, createFileRowError(CalendarFileName, item.LineNumber, "service_id must not be empty"))
		}

		// Weekdays must be "0" or "1"
		for _, day := range []struct {
			fieldName string
			value     string
		}{
			{"monday", item.Monday},
			{"tuesday", item.Tuesday},
			{"wednesday", item.Wednesday},
			{"thursday", item.Thursday},
			{"friday", item.Friday},
			{"saturday", item.Saturday},
			{"sunday", item.Sunday},
		} {
			if day.value != "0" && day.value != "1" {
				validationErrors = append(validationErrors, createFileRowError(CalendarFileName, item.LineNumber, fmt.Sprintf("%s must be '0' or '1'", day.fieldName)))
			}
		}

		// Ensure start and end dates are not empty
		if item.StartDate == "" {
			validationErrors = append(validationErrors, createFileRowError(CalendarFileName, item.LineNumber, "start_date must not be empty"))
		}

		if item.EndDate == "" {
			validationErrors = append(validationErrors, createFileRowError(CalendarFileName, item.LineNumber, "end_date must not be empty"))
		}
	}

	return validationErrors, recommendations
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
//type CalendarItem struct {
//	ServiceId  string
//	Monday     int
//	Tuesday    int
//	Wednesday  int
//	Thursday   int
//	Friday     int
//	Saturday   int
//	Sunday     int
//	Start      time.Time
//	End        time.Time
//	lineNumber int
//}
//
//var validCalendarHeaders = []string{"service_id", "monday", "tuesday", "wednesday",
//	"thursday", "friday", "saturday", "sunday", "start_date", "end_date"}
//
//func LoadCalendarItems(csvReader *csv.Reader) ([]*CalendarItem, []error) {
//	calendarItems := make([]*CalendarItem, 0)
//	errs := make([]error, 0)
//
//	headers, err := ReadHeaderRow(csvReader, validCalendarHeaders)
//	if err != nil {
//		errs = append(errs, createFileError(CalendarFileName, fmt.Sprintf("read error: %v", err.Error())))
//		return calendarItems, errs
//	}
//	if headers == nil {
//		return calendarItems, errs
//	}
//
//	usedIds := make([]string, 0)
//	index := 0
//	for {
//		row, err := ReadDataRow(csvReader)
//		if err != nil {
//			errs = append(errs, createFileError(CalendarFileName, fmt.Sprintf("%v", err.Error())))
//			index++
//			break
//		}
//		if row == nil {
//			break
//		}
//
//		rowErrs := make([]error, 0)
//		calendarItem := CalendarItem{
//			lineNumber: index,
//		}
//
//		var (
//			mon, tue, wed, thu, fri, sat, sun *int
//			service                           *string
//			start, end                        *time.Time
//		)
//
//		for name, column := range headers {
//			switch name {
//			case "service_id":
//				service = handleIDField(row[column], CalendarFileName, name, index, &rowErrs)
//			case "monday":
//				mon = handleWeekdayField(row[column], CalendarFileName, name, index, &rowErrs)
//			case "tuesday":
//				tue = handleWeekdayField(row[column], CalendarFileName, name, index, &rowErrs)
//			case "wednesday":
//				wed = handleWeekdayField(row[column], CalendarFileName, name, index, &rowErrs)
//			case "thursday":
//				thu = handleWeekdayField(row[column], CalendarFileName, name, index, &rowErrs)
//			case "friday":
//				fri = handleWeekdayField(row[column], CalendarFileName, name, index, &rowErrs)
//			case "saturday":
//				sat = handleWeekdayField(row[column], CalendarFileName, name, index, &rowErrs)
//			case "sunday":
//				sun = handleWeekdayField(row[column], CalendarFileName, name, index, &rowErrs)
//			case "start_date":
//				start = handleDateField(row[column], CalendarFileName, name, index, false, &rowErrs)
//			case "end_date":
//				end = handleDateField(row[column], CalendarFileName, name, index, true, &rowErrs)
//			}
//		}
//
//		if service == nil {
//			rowErrs = append(rowErrs, createFileRowError(CalendarFileName, calendarItem.lineNumber, "service_id must be specified"))
//		} else {
//			calendarItem.ServiceId = *service
//			if StringArrayContainsItem(usedIds, *service) {
//				errs = append(errs, createFileRowError(CalendarFileName, index, fmt.Sprintf("%s: service_id", nonUniqueId)))
//			} else {
//				usedIds = append(usedIds, *service)
//			}
//		}
//
//		if mon == nil {
//			rowErrs = append(rowErrs, createFileRowError(CalendarFileName, calendarItem.lineNumber, "monday must be specified"))
//		} else {
//			calendarItem.Monday = *mon
//		}
//
//		if tue == nil {
//			rowErrs = append(rowErrs, createFileRowError(CalendarFileName, calendarItem.lineNumber, "tuesday must be specified"))
//		} else {
//			calendarItem.Tuesday = *tue
//		}
//
//		if wed == nil {
//			rowErrs = append(rowErrs, createFileRowError(CalendarFileName, calendarItem.lineNumber, "wednesday must be specified"))
//		} else {
//			calendarItem.Wednesday = *wed
//		}
//
//		if thu == nil {
//			rowErrs = append(rowErrs, createFileRowError(CalendarFileName, calendarItem.lineNumber, "thursday must be specified"))
//		} else {
//			calendarItem.Thursday = *thu
//		}
//
//		if fri == nil {
//			rowErrs = append(rowErrs, createFileRowError(CalendarFileName, calendarItem.lineNumber, "friday must be specified"))
//		} else {
//			calendarItem.Friday = *fri
//		}
//
//		if sat == nil {
//			rowErrs = append(rowErrs, createFileRowError(CalendarFileName, calendarItem.lineNumber, "saturday must be specified"))
//		} else {
//			calendarItem.Saturday = *sat
//		}
//
//		if sun == nil {
//			rowErrs = append(rowErrs, createFileRowError(CalendarFileName, calendarItem.lineNumber, "sunday must be specified"))
//		} else {
//			calendarItem.Sunday = *sun
//		}
//
//		if start == nil {
//			rowErrs = append(rowErrs, createFileRowError(CalendarFileName, calendarItem.lineNumber, "start_date must be specified"))
//		} else {
//			calendarItem.Start = *start
//		}
//
//		if end == nil {
//			rowErrs = append(rowErrs, createFileRowError(CalendarFileName, calendarItem.lineNumber, "end_date must be specified"))
//		} else {
//			calendarItem.End = *end
//		}
//
//		if len(rowErrs) > 0 {
//			errs = append(errs, rowErrs...)
//		} else {
//			calendarItems = append(calendarItems, &calendarItem)
//		}
//
//		index++
//	}
//
//	return calendarItems, errs
//}
//
//func handleWeekdayField(str string, fileName string, fieldName string, index int, errs *[]error) *int {
//	n, err := strconv.ParseInt(str, 10, 64)
//	if err != nil {
//		*errs = append(*errs, createFieldError(fileName, fieldName, index, err))
//		return nil
//	}
//
//	if v := int(n); v <= 1 {
//		return &v
//	} else {
//		*errs = append(*errs, createFieldError(fileName, fieldName, index, errors.New(invalidValue)))
//		return nil
//	}
//}
