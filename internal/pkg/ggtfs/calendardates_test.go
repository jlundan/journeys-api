//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

var validCalendarDateHeaders = []string{"service_id", "date", "exception_type"}

func TestShouldReturnEmptyCalendarDateArrayOnEmptyString(t *testing.T) {
	agencies, errors := LoadEntities[*CalendarDate](csv.NewReader(strings.NewReader("")), validCalendarDateHeaders, CreateCalendarDate, CalendarDatesFileName)
	if len(errors) > 0 {
		t.Error(errors)
	}
	if len(agencies) != 0 {
		t.Error("expected zero calendar items")
	}
}

func TestCalendarDateParsing(t *testing.T) {
	loadCalendarDatesFunc := func(reader *csv.Reader) ([]interface{}, []error) {
		calendarDates, errs := LoadEntities[*CalendarDate](reader, validCalendarDateHeaders, CreateCalendarDate, CalendarDatesFileName)
		entities := make([]interface{}, len(calendarDates))
		for i, calendarItem := range calendarDates {
			entities[i] = calendarItem
		}
		return entities, errs
	}

	validateCalendarItemsFunc := func(entities []interface{}, fixtures map[string][]interface{}) ([]error, []string) {
		calendarDates := make([]*CalendarDate, len(entities))
		for i, entity := range entities {
			if calendarDate, ok := entity.(*CalendarDate); ok {
				calendarDates[i] = calendarDate
			}
		}

		ciCount := len(fixtures["calendarItems"])

		if ciCount == 0 {
			return ValidateCalendarDates(calendarDates, nil)
		}

		calendarItems := make([]*CalendarItem, ciCount)
		for i, fixture := range fixtures["calendarItems"] {
			if calendarItem, ok := fixture.(*CalendarItem); ok {
				calendarItems[i] = calendarItem
			} else {
				t.Error(fmt.Sprintf("test setup error: cannot convert %v to CalendarItem pointer. maybe you used value instead of pointer when setting fixtures", fixture))
			}
		}

		return ValidateCalendarDates(calendarDates, calendarItems)
	}

	runGenericGTFSParseTest(t, "CalendarDateNOKTestcases", loadCalendarDatesFunc, validateCalendarItemsFunc, false, getCalendarDateNOKTestcases())
	runGenericGTFSParseTest(t, "CalendarDateOKTestcases", loadCalendarDatesFunc, validateCalendarItemsFunc, false, getCalendarDateOKTestcases())
}

func getCalendarDateOKTestcases() map[string]ggtfsTestCase {
	expected1 := CalendarDate{
		ServiceId:     NewID(stringPtr("1")),
		Date:          NewDate(stringPtr("20200101")),
		ExceptionType: NewExceptionTypeEnum(stringPtr("1")),
	}

	testCases := make(map[string]ggtfsTestCase)
	testCases["1"] = ggtfsTestCase{
		csvRows: [][]string{
			{"service_id", "date", "exception_type"},
			{"1", "20200101", "1"},
		},
		expectedStructs: []interface{}{&expected1},
	}

	return testCases
}

func getCalendarDateNOKTestcases() map[string]ggtfsTestCase {
	testCases := make(map[string]ggtfsTestCase)
	testCases["1"] = ggtfsTestCase{
		csvRows: [][]string{
			{"service_id", "date", "exception_type"},
			{","},
			{" ", " ", "1"},
			{"1000", "20201011", "not an int"},
			{"1001", "20201011", "10"},
		},
		expectedErrors: []string{
			"calendar_dates.txt: record on line 2: wrong number of fields",
			"calendar_dates.txt:1: invalid field: date",
			"calendar_dates.txt:1: invalid field: service_id",
			"calendar_dates.txt:2: invalid field: exception_type",
			"calendar_dates.txt:3: invalid field: exception_type",
		},
	}

	testCases["2"] = ggtfsTestCase{
		csvRows: [][]string{
			{"service_id", "date", "exception_type"},
			{"1000", "20201011", strconv.Itoa(ServiceAddedForCalendarDate)},
			{"1001", "20201011", strconv.Itoa(ServiceRemovedForCalendarDate)},
		},
		expectedErrors: []string{
			"calendar_dates.txt:1: referenced service_id '1001' not found in calendar.txt",
		},
		fixtures: map[string][]interface{}{
			"calendarItems": {
				&CalendarItem{
					ServiceId:  NewID(stringPtr("1000")),
					Monday:     NewAvailableForWeekdayInfo(stringPtr(CalendarAvailableForWeekday)),
					Tuesday:    NewAvailableForWeekdayInfo(stringPtr(CalendarAvailableForWeekday)),
					Wednesday:  NewAvailableForWeekdayInfo(stringPtr(CalendarAvailableForWeekday)),
					Thursday:   NewAvailableForWeekdayInfo(stringPtr(CalendarAvailableForWeekday)),
					Friday:     NewAvailableForWeekdayInfo(stringPtr(CalendarAvailableForWeekday)),
					Saturday:   NewAvailableForWeekdayInfo(stringPtr(CalendarNotAvailableForWeekday)),
					Sunday:     NewAvailableForWeekdayInfo(stringPtr(CalendarNotAvailableForWeekday)),
					StartDate:  NewDate(stringPtr("20201011")),
					EndDate:    NewDate(stringPtr("20201011")),
					LineNumber: 0,
				},
			},
		},
	}

	return testCases
}

func TestValidateCalendarDateReferencesReturnsNoErrorsOnNilValues(t *testing.T) {
	calendarDates := []*CalendarDate{
		nil,
	}
	calendarItems := []*CalendarItem{
		nil,
	}
	var validationErrors *[]error

	validateCalendarDateReferences(calendarDates, calendarItems, validationErrors)
	if validationErrors != nil {
		t.Error("Expected no errors")
	}
}

func TestNewExceptionTypeEnumReturnsEmptyOnNilArgument(t *testing.T) {
	nete := NewExceptionTypeEnum(nil)
	if !nete.IsEmpty() {
		t.Error("Expected empty ExceptionTypeEnum")
	}
}
