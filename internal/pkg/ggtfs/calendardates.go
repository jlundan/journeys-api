package ggtfs

import (
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type CalendarDate struct {
	ServiceId     string
	Date          time.Time
	ExceptionType int
	LineNumber    int
}

func ExtractCalendarDates(input *csv.Reader, output *csv.Writer, serviceIds map[string]struct{}) []error {
	errs := make([]error, 0)

	headers, err := ReadHeaderRow(input)
	if err != nil {
		errs = append(errs, createFileError(CalendarDatesFileName, fmt.Sprintf("read error: %v", err.Error())))
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
	if pos, columnExists := headers["service_id"]; columnExists {
		idHeaderPos = pos
	} else {
		errs = append(errs, createFileError(CalendarDatesFileName, "cannot extract calendar dates without service_id column"))
		return errs
	}

	for {
		row, rErr := ReadDataRow(input)
		if rErr != nil {
			errs = append(errs, createFileError(CalendarDatesFileName, fmt.Sprintf("%v", rErr.Error())))
			continue
		}

		if row == nil { // EOF
			break
		}

		if _, shouldBeExtracted := serviceIds[row[idHeaderPos]]; shouldBeExtracted {
			wErr := writeDataRow(row, output)
			if wErr != nil {
				errs = append(errs, wErr)
				return errs
			}
		}
	}

	return nil
}

func LoadCalendarDates(csvReader *csv.Reader) ([]*CalendarDate, []error) {
	calendarDates := make([]*CalendarDate, 0)
	errs := make([]error, 0)

	headers, err := ReadHeaderRow(csvReader)
	if err != nil {
		errs = append(errs, createFileError(CalendarDatesFileName, fmt.Sprintf("read error: %v", err.Error())))
		return calendarDates, errs
	}
	if headers == nil {
		return calendarDates, errs
	}

	index := 0
	for {
		row, err := ReadDataRow(csvReader)
		if err != nil {
			errs = append(errs, createFileError(CalendarDatesFileName, fmt.Sprintf("%v", err.Error())))
			index++
			continue
		}
		if row == nil {
			break
		}

		rowErrs := make([]error, 0)
		calendarDate := CalendarDate{
			LineNumber: index,
		}

		var (
			serviceId     *string
			date          *time.Time
			exceptionType *int
		)

		for name, column := range headers {
			switch name {
			case "service_id":
				serviceId = handleIDField(row[column], CalendarDatesFileName, name, index, &rowErrs)
			case "date":
				date = handleDateField(row[column], CalendarDatesFileName, name, index, false, &rowErrs)
			case "exception_type":
				exceptionType = handleExceptionTypeField(row[column], CalendarDatesFileName, name, index, &rowErrs)
			}
		}

		if serviceId == nil {
			rowErrs = append(rowErrs, createFileRowError(CalendarDatesFileName, calendarDate.LineNumber, "service_id must be specified"))
		} else {
			calendarDate.ServiceId = *serviceId
		}

		if date == nil {
			rowErrs = append(rowErrs, createFileRowError(CalendarDatesFileName, calendarDate.LineNumber, "date must be specified"))
		} else {
			calendarDate.Date = *date
		}

		if exceptionType == nil {
			rowErrs = append(rowErrs, createFileRowError(CalendarDatesFileName, calendarDate.LineNumber, "exception_type must be specified"))
		} else {
			calendarDate.ExceptionType = *exceptionType
		}

		if len(rowErrs) > 0 {
			errs = append(errs, rowErrs...)
		} else {
			calendarDates = append(calendarDates, &calendarDate)
		}

		index++
	}

	return calendarDates, errs
}

func ValidateCalendarDates(calendarDates []*CalendarDate, calendarItems []*CalendarItem) []error {
	var validationErrors []error

	if calendarDates == nil {
		return validationErrors
	}

	if calendarItems != nil {
		for _, calendarDate := range calendarDates {
			if calendarDate == nil {
				continue
			}
			notFound := true
			for _, calendarItem := range calendarItems {
				if calendarItem == nil {
					continue
				}
				if calendarDate.ServiceId == calendarItem.ServiceId {
					notFound = false
					break
				}
			}
			if notFound {
				validationErrors = append(validationErrors, createFileRowError(CalendarDatesFileName, calendarDate.LineNumber, fmt.Sprintf("referenced service_id not found in %s", CalendarFileName)))
			}
		}
	}

	return validationErrors
}

func handleExceptionTypeField(str string, fileName string, fieldName string, index int, errs *[]error) *int {
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, err))
		return nil
	}

	if v := int(n); v >= 1 && v <= 2 {
		return &v
	} else {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, errors.New(invalidValue)))
		return nil
	}
}
