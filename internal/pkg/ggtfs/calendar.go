package ggtfs

import (
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type CalendarItem struct {
	ServiceId  string
	Monday     int
	Tuesday    int
	Wednesday  int
	Thursday   int
	Friday     int
	Saturday   int
	Sunday     int
	Start      time.Time
	End        time.Time
	lineNumber int
}

var validCalendarHeaders = []string{"service_id", "monday", "tuesday", "wednesday",
	"thursday", "friday", "saturday", "sunday", "start_date", "end_date"}

func LoadCalendarItems(csvReader *csv.Reader) ([]*CalendarItem, []error) {
	calendarItems := make([]*CalendarItem, 0)
	errs := make([]error, 0)

	headers, err := ReadHeaderRow(csvReader, validCalendarHeaders)
	if err != nil {
		errs = append(errs, createFileError(CalendarFileName, fmt.Sprintf("read error: %v", err.Error())))
		return calendarItems, errs
	}
	if headers == nil {
		return calendarItems, errs
	}

	usedIds := make([]string, 0)
	index := 0
	for {
		row, err := ReadDataRow(csvReader)
		if err != nil {
			errs = append(errs, createFileError(CalendarFileName, fmt.Sprintf("%v", err.Error())))
			index++
			break
		}
		if row == nil {
			break
		}

		rowErrs := make([]error, 0)
		calendarItem := CalendarItem{
			lineNumber: index,
		}

		var (
			mon, tue, wed, thu, fri, sat, sun *int
			service                           *string
			start, end                        *time.Time
		)

		for name, column := range headers {
			switch name {
			case "service_id":
				service = handleIDField(row[column], CalendarFileName, name, index, &rowErrs)
			case "monday":
				mon = handleWeekdayField(row[column], CalendarFileName, name, index, &rowErrs)
			case "tuesday":
				tue = handleWeekdayField(row[column], CalendarFileName, name, index, &rowErrs)
			case "wednesday":
				wed = handleWeekdayField(row[column], CalendarFileName, name, index, &rowErrs)
			case "thursday":
				thu = handleWeekdayField(row[column], CalendarFileName, name, index, &rowErrs)
			case "friday":
				fri = handleWeekdayField(row[column], CalendarFileName, name, index, &rowErrs)
			case "saturday":
				sat = handleWeekdayField(row[column], CalendarFileName, name, index, &rowErrs)
			case "sunday":
				sun = handleWeekdayField(row[column], CalendarFileName, name, index, &rowErrs)
			case "start_date":
				start = handleDateField(row[column], CalendarFileName, name, index, false, &rowErrs)
			case "end_date":
				end = handleDateField(row[column], CalendarFileName, name, index, true, &rowErrs)
			}
		}

		if service == nil {
			rowErrs = append(rowErrs, createFileRowError(CalendarFileName, calendarItem.lineNumber, "service_id must be specified"))
		} else {
			calendarItem.ServiceId = *service
			if StringArrayContainsItem(usedIds, *service) {
				errs = append(errs, createFileRowError(CalendarFileName, index, fmt.Sprintf("%s: service_id", nonUniqueId)))
			} else {
				usedIds = append(usedIds, *service)
			}
		}

		if mon == nil {
			rowErrs = append(rowErrs, createFileRowError(CalendarFileName, calendarItem.lineNumber, "monday must be specified"))
		} else {
			calendarItem.Monday = *mon
		}

		if tue == nil {
			rowErrs = append(rowErrs, createFileRowError(CalendarFileName, calendarItem.lineNumber, "tuesday must be specified"))
		} else {
			calendarItem.Tuesday = *tue
		}

		if wed == nil {
			rowErrs = append(rowErrs, createFileRowError(CalendarFileName, calendarItem.lineNumber, "wednesday must be specified"))
		} else {
			calendarItem.Wednesday = *wed
		}

		if thu == nil {
			rowErrs = append(rowErrs, createFileRowError(CalendarFileName, calendarItem.lineNumber, "thursday must be specified"))
		} else {
			calendarItem.Thursday = *thu
		}

		if fri == nil {
			rowErrs = append(rowErrs, createFileRowError(CalendarFileName, calendarItem.lineNumber, "friday must be specified"))
		} else {
			calendarItem.Friday = *fri
		}

		if sat == nil {
			rowErrs = append(rowErrs, createFileRowError(CalendarFileName, calendarItem.lineNumber, "saturday must be specified"))
		} else {
			calendarItem.Saturday = *sat
		}

		if sun == nil {
			rowErrs = append(rowErrs, createFileRowError(CalendarFileName, calendarItem.lineNumber, "sunday must be specified"))
		} else {
			calendarItem.Sunday = *sun
		}

		if start == nil {
			rowErrs = append(rowErrs, createFileRowError(CalendarFileName, calendarItem.lineNumber, "start_date must be specified"))
		} else {
			calendarItem.Start = *start
		}

		if end == nil {
			rowErrs = append(rowErrs, createFileRowError(CalendarFileName, calendarItem.lineNumber, "end_date must be specified"))
		} else {
			calendarItem.End = *end
		}

		if len(rowErrs) > 0 {
			errs = append(errs, rowErrs...)
		} else {
			calendarItems = append(calendarItems, &calendarItem)
		}

		index++
	}

	return calendarItems, errs
}

func handleWeekdayField(str string, fileName string, fieldName string, index int, errs *[]error) *int {
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, err))
		return nil
	}

	if v := int(n); v <= 1 {
		return &v
	} else {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, errors.New(invalidValue)))
		return nil
	}
}
