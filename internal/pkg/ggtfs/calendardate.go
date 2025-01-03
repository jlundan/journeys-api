package ggtfs

import (
	"fmt"
	"strconv"
)

type CalendarDate struct {
	ServiceId     ID                // service_id 		(required)
	Date          Date              // date 			(required)
	ExceptionType ExceptionTypeEnum // exception_type 	(required)
	LineNumber    int
}

func CreateCalendarDate(row []string, headers map[string]int, lineNumber int) *CalendarDate {
	calendarDate := &CalendarDate{
		LineNumber: lineNumber,
	}

	for hName := range headers {
		v := getRowValueForHeaderName(row, headers, hName)
		switch hName {
		case "service_id":
			calendarDate.ServiceId = NewID(v)
		case "date":
			calendarDate.Date = NewDate(v)
		case "exception_type":
			calendarDate.ExceptionType = NewExceptionTypeEnum(v)
		}
	}

	return calendarDate
}

func ValidateCalendarDate(cd CalendarDate) []error {
	var validationErrors []error

	requiredFields := map[string]FieldTobeValidated{
		"service_id":     &cd.ServiceId,
		"date":           &cd.Date,
		"exception_type": &cd.ExceptionType,
	}
	validateRequiredFields(requiredFields, &validationErrors, cd.LineNumber, CalendarDatesFileName)

	return validationErrors
}

func ValidateCalendarDates(calendarDates []*CalendarDate, calendarItems []*CalendarItem) ([]error, []string) {
	var validationErrors []error

	for _, calendarDate := range calendarDates {
		validationErrors = append(validationErrors, ValidateCalendarDate(*calendarDate)...)
	}

	if calendarItems != nil {
		validateCalendarDateReferences(calendarDates, calendarItems, &validationErrors)
	}

	return validationErrors, nil
}

func validateCalendarDateReferences(calendarDates []*CalendarDate, calendarItems []*CalendarItem, validationErrors *[]error) {
	serviceIDMap := make(map[string]struct{})
	for _, item := range calendarItems {
		if item != nil && !StringIsNilOrEmpty(item.ServiceId) {
			serviceIDMap[*item.ServiceId] = struct{}{}
		}
	}

	for _, calendarDate := range calendarDates {
		if calendarDate == nil || calendarDate.ServiceId.IsEmpty() {
			continue
		}
		if _, found := serviceIDMap[calendarDate.ServiceId.Raw()]; !found {
			*validationErrors = append(*validationErrors,
				createFileRowError(CalendarDatesFileName, calendarDate.LineNumber,
					fmt.Sprintf("referenced service_id '%s' not found in %s", calendarDate.ServiceId.Raw(), CalendarFileName)))
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
