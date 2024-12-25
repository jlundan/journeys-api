package ggtfs

import (
	"fmt"
	"strconv"
)

type CalendarDate struct {
	ServiceId     ID                // service_id
	Date          Date              // date
	ExceptionType ExceptionTypeEnum // exception_type
	LineNumber    int               // CSV row number
}

func (cd CalendarDate) Validate() []error {
	var validationErrors []error

	requiredFields := []struct {
		fieldName string
		field     ValidAndPresentField
	}{
		{"service_id", &cd.ServiceId},
		{"date", &cd.Date},
		{"exception_type", &cd.ExceptionType},
	}

	for _, f := range requiredFields {
		if !f.field.IsValid() {
			validationErrors = append(validationErrors, createFileRowError(CalendarDatesFileName, cd.LineNumber, createInvalidRequiredFieldString(f.fieldName)))
		}
	}

	return validationErrors
}

// CreateCalendarDate creates a CalendarDate from a CSV row using the provided headers.
func CreateCalendarDate(row []string, headers map[string]int, lineNumber int) *CalendarDate {
	calendarDate := &CalendarDate{
		LineNumber: lineNumber,
	}

	for hName, hPos := range headers {
		switch hName {
		case "service_id":
			calendarDate.ServiceId = NewID(getRowValue(row, hPos))
		case "date":
			calendarDate.Date = NewDate(getRowValue(row, hPos))
		case "exception_type":
			calendarDate.ExceptionType = NewExceptionTypeEnum(getRowValue(row, hPos))
		}
	}

	return calendarDate
}

// ValidateCalendarDates validates the parsed CalendarDate structs for logical consistency.
func ValidateCalendarDates(calendarDates []*CalendarDate, calendarItems []*CalendarItem) ([]error, []string) {
	var validationErrors []error

	for _, calendarDate := range calendarDates {
		validationErrors = append(validationErrors, calendarDate.Validate()...)
	}

	// Cross-validate service_id references with CalendarItems (if provided)
	if calendarItems != nil {
		validateCalendarDateReferences(calendarDates, calendarItems, &validationErrors)
	}

	return validationErrors, nil
}

// validateCalendarDateReferences checks if the service_id in CalendarDate has a matching entry in CalendarItem.
func validateCalendarDateReferences(calendarDates []*CalendarDate, calendarItems []*CalendarItem, validationErrors *[]error) {
	// Create a map of CalendarItem service_ids for quick lookup
	serviceIDMap := make(map[string]struct{})
	for _, item := range calendarItems {
		if item != nil && !item.ServiceId.IsEmpty() {
			serviceIDMap[item.ServiceId.String()] = struct{}{}
		}
	}

	// Check if each CalendarDate service_id is present in CalendarItems
	for _, calendarDate := range calendarDates {
		if calendarDate == nil || calendarDate.ServiceId.IsEmpty() {
			continue
		}
		if _, found := serviceIDMap[calendarDate.ServiceId.String()]; !found {
			*validationErrors = append(*validationErrors,
				createFileRowError(CalendarDatesFileName, calendarDate.LineNumber,
					fmt.Sprintf("referenced service_id '%s' not found in %s", calendarDate.ServiceId.String(), CalendarFileName)))
		}
	}
}

const (
	ServiceAddedForCalendarDate   int = 1
	ServiceRemovedForCalendarDate int = 2
)

type ExceptionTypeEnum struct {
	Integer
}

func (ete ExceptionTypeEnum) IsValid() bool {
	val, err := strconv.Atoi(ete.Integer.base.raw)
	if err != nil {
		return false
	}

	return val == ServiceAddedForCalendarDate || val == ServiceRemovedForCalendarDate
}

func NewExceptionTypeEnum(raw *string) ExceptionTypeEnum {
	if raw == nil {
		return ExceptionTypeEnum{
			Integer{base: base{raw: ""}}}
	}
	return ExceptionTypeEnum{Integer{base: base{raw: *raw, isPresent: true}}}
}
