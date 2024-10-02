//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"encoding/csv"
	"strings"
	"testing"
)

func TestShouldReturnEmptyCalendarDateArrayOnEmptyString(t *testing.T) {
	agencies, errors := LoadCalendarDates(csv.NewReader(strings.NewReader("")))
	if len(errors) > 0 {
		t.Error(errors)
	}
	if len(agencies) != 0 {
		t.Error("expected zero calendar items")
	}
}

func TestCalendarDateParsing(t *testing.T) {
	loadCalendarDatesFunc := func(reader *csv.Reader) ([]interface{}, []error) {
		calendarDates, errs := LoadCalendarDates(reader)
		entities := make([]interface{}, len(calendarDates))
		for i, calendarItem := range calendarDates {
			entities[i] = calendarItem
		}
		return entities, errs
	}

	validateCalendarItemsFunc := func(entities []interface{}) ([]error, []string) {
		calendarDates := make([]*CalendarDate, len(entities))
		for i, entity := range entities {
			if calendarDate, ok := entity.(*CalendarDate); ok {
				calendarDates[i] = calendarDate
			}
		}

		return ValidateCalendarDates(calendarDates, nil)
	}

	runGenericGTFSParseTest(t, "CalendarDateNOKTestcases", loadCalendarDatesFunc, validateCalendarItemsFunc, false, getCalendarDateNOKTestcases())
	runGenericGTFSParseTest(t, "CalendarDateOKTestcases", loadCalendarDatesFunc, validateCalendarItemsFunc, false, getCalendarDateOKTestcases())
}

func getCalendarDateOKTestcases() map[string]ggtfsTestCase {
	expected1 := CalendarDate{
		ServiceId:     "1",
		Date:          "20200101",
		ExceptionType: "1",
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
			"calendar_dates.txt:1: invalid date format:  ",
			//"calendar_dates.txt:1: service_id must not be empty",
			"calendar_dates.txt:2: exception_type must be '1' or '2', found: not an int",
			"calendar_dates.txt:3: exception_type must be '1' or '2', found: 10",
			// TODO: Fix these
			//"calendar_dates.txt: record on line 2: wrong number of fields",
			//"calendar_dates.txt:1: service_id must be specified",
			//"calendar_dates.txt:1: service_id: empty value not allowed",
			//"calendar_dates.txt:1: date must be specified",
			//"calendar_dates.txt:1: date: empty value not allowed",
			//"calendar_dates.txt:2: exception_type must be specified",
			//"calendar_dates.txt:2: exception_type: strconv.ParseInt: parsing \"not an int\": invalid syntax",
			//"calendar_dates.txt:3: exception_type: invalid value",
			//"calendar_dates.txt:3: exception_type must be specified",
		},
	}

	return testCases
}

// TODO: Integrate these to the test cases run by runGenericGTFSParseTest
//func TestValidateCalendarDates(t *testing.T) {
//	testCases := []struct {
//		calendarDates  []*CalendarDate
//		calendarItems  []*CalendarItem
//		expectedErrors []string
//	}{
//		{
//			calendarDates: []*CalendarDate{
//				{ServiceId: "1000", LineNumber: 0},
//			},
//			calendarItems: []*CalendarItem{
//				{ServiceId: "1000", LineNumber: 0},
//				{ServiceId: "1001", LineNumber: 1},
//			},
//			expectedErrors: []string{},
//		},
//		{
//			calendarDates:  nil,
//			expectedErrors: []string{},
//		},
//		{
//			calendarDates: []*CalendarDate{nil},
//			calendarItems: []*CalendarItem{
//				{ServiceId: "1002", LineNumber: 0},
//				{ServiceId: "1001", LineNumber: 1},
//			},
//			expectedErrors: []string{},
//		},
//		{
//			calendarDates: []*CalendarDate{
//				{ServiceId: "1000", LineNumber: 0},
//			},
//			calendarItems: []*CalendarItem{nil},
//			expectedErrors: []string{
//				"calendar_dates.txt:0: referenced service_id not found in calendar.txt",
//			},
//		},
//		{
//			calendarDates: []*CalendarDate{
//				{ServiceId: "1000", LineNumber: 0},
//			},
//			calendarItems: []*CalendarItem{
//				{ServiceId: "1002", LineNumber: 0},
//				{ServiceId: "1001", LineNumber: 1},
//			},
//			expectedErrors: []string{
//				"calendar_dates.txt:0: referenced service_id not found in calendar.txt",
//			},
//		},
//	}
//
//	for _, tc := range testCases {
//		err := ValidateCalendarDates(tc.calendarDates, tc.calendarItems)
//		checkErrors(tc.expectedErrors, err, t)
//	}
//}
//
//func calendarDatesMatch(a CalendarDate, b CalendarDate) bool {
//	return a.ServiceId == b.ServiceId && a.Date == b.Date && a.ExceptionType == b.ExceptionType
//}
