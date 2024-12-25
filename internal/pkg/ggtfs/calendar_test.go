//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"encoding/csv"
	"strings"
	"testing"
)

var validCalendarHeaders = []string{
	"service_id", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday", "start_date", "end_date",
}

func TestShouldReturnEmptyCalendarItemArrayOnEmptyString(t *testing.T) {
	agencies, errors := LoadEntities[*CalendarItem](csv.NewReader(strings.NewReader("")), validCalendarHeaders, CreateCalendarItem, CalendarFileName)
	if len(errors) > 0 {
		t.Error(errors)
	}
	if len(agencies) != 0 {
		t.Error("expected zero calendar items")
	}
}

func TestNewWeekdayEnumReturnsEmptyOnNil(t *testing.T) {
	we := NewAvailableForWeekdayInfo(nil)

	if we.raw != "" {
		t.Error("expected empty AvailableForWeekdayInfo")
	}
}

func TestCalendarItemParsing(t *testing.T) {
	loadCalendarItemsFunc := func(reader *csv.Reader) ([]interface{}, []error) {
		calendarItems, errs := LoadEntities[*CalendarItem](reader, validCalendarHeaders, CreateCalendarItem, CalendarFileName)
		entities := make([]interface{}, len(calendarItems))
		for i, calendarItem := range calendarItems {
			entities[i] = calendarItem
		}
		return entities, errs
	}

	validateCalendarItemsFunc := func(entities []interface{}, _fixtures map[string][]interface{}) ([]error, []string) {
		calendarItems := make([]*CalendarItem, len(entities))
		for i, entity := range entities {
			if calendarItem, ok := entity.(*CalendarItem); ok {
				calendarItems[i] = calendarItem
			}
		}
		return ValidateCalendarItems(calendarItems)
	}

	runGenericGTFSParseTest(t, "CalendarItemNOKTestcases", loadCalendarItemsFunc, validateCalendarItemsFunc, false, getCalendarItemNOKTestcases())
	runGenericGTFSParseTest(t, "CalendarItemOKTestcases", loadCalendarItemsFunc, validateCalendarItemsFunc, false, getCalendarItemOKTestcases())
}

func getCalendarItemNOKTestcases() map[string]ggtfsTestCase {
	testCases := make(map[string]ggtfsTestCase)

	testCases["parse-failures"] = ggtfsTestCase{
		csvRows: [][]string{
			{"service_id", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday", "start_date", "end_date"},
			{" "},
			{","},
			{"", ""},
			{" ", " "},
			{"", "", ""},
			{" ", " ", " "},
			{"", "", "", ""},
			{" ", " ", " ", " "},
			{"", "", "", "", ""},
			{" ", " ", " ", " ", " "},
			{"", "", "", "", "", ""},
			{" ", " ", " ", " ", " ", " "},
			{"", "", "", "", "", "", ""},
			{" ", " ", " ", " ", " ", " ", " "},
			{"", "", "", "", "", "", "", ""},
			{" ", " ", " ", " ", " ", " ", " ", " "},
			{" ", " ", " ", " ", " ", " ", " ", " ", " "},
		},
		expectedErrors: []string{
			"calendar.txt: record on line 10: wrong number of fields",
			"calendar.txt: record on line 11: wrong number of fields",
			"calendar.txt: record on line 12: wrong number of fields",
			"calendar.txt: record on line 13: wrong number of fields",
			"calendar.txt: record on line 14: wrong number of fields",
			"calendar.txt: record on line 15: wrong number of fields",
			"calendar.txt: record on line 16: wrong number of fields",
			"calendar.txt: record on line 17: wrong number of fields",
			"calendar.txt: record on line 18: wrong number of fields",
			"calendar.txt: record on line 2: wrong number of fields",
			"calendar.txt: record on line 3: wrong number of fields",
			"calendar.txt: record on line 4: wrong number of fields",
			"calendar.txt: record on line 5: wrong number of fields",
			"calendar.txt: record on line 6: wrong number of fields",
			"calendar.txt: record on line 7: wrong number of fields",
			"calendar.txt: record on line 8: wrong number of fields",
			"calendar.txt: record on line 9: wrong number of fields",
		},
	}

	testCases["invalid-fields"] = ggtfsTestCase{
		csvRows: [][]string{
			{"service_id", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday", "start_date", "end_date"},
			{" ", " ", " ", " ", " ", " ", " ", " ", " ", " "},
			{"SERVICE0", "3", "not a day", "-1", "2", "2", "2", "2", "not a date", "20181313"},
		},
		expectedErrors: []string{
			"calendar.txt:0: invalid mandatory field: end_date",
			"calendar.txt:0: invalid mandatory field: friday",
			"calendar.txt:0: invalid mandatory field: monday",
			"calendar.txt:0: invalid mandatory field: saturday",
			"calendar.txt:0: invalid mandatory field: service_id",
			"calendar.txt:0: invalid mandatory field: start_date",
			"calendar.txt:0: invalid mandatory field: sunday",
			"calendar.txt:0: invalid mandatory field: thursday",
			"calendar.txt:0: invalid mandatory field: tuesday",
			"calendar.txt:0: invalid mandatory field: wednesday",
			"calendar.txt:1: invalid mandatory field: end_date",
			"calendar.txt:1: invalid mandatory field: friday",
			"calendar.txt:1: invalid mandatory field: monday",
			"calendar.txt:1: invalid mandatory field: saturday",
			"calendar.txt:1: invalid mandatory field: start_date",
			"calendar.txt:1: invalid mandatory field: sunday",
			"calendar.txt:1: invalid mandatory field: thursday",
			"calendar.txt:1: invalid mandatory field: tuesday",
			"calendar.txt:1: invalid mandatory field: wednesday",
		},
	}

	return testCases
}

func getCalendarItemOKTestcases() map[string]ggtfsTestCase {
	expected1 := CalendarItem{
		ServiceId: NewID(stringPtr("111")),
		Monday:    NewAvailableForWeekdayInfo(stringPtr("1")),
		Tuesday:   NewAvailableForWeekdayInfo(stringPtr("1")),
		Wednesday: NewAvailableForWeekdayInfo(stringPtr("1")),
		Thursday:  NewAvailableForWeekdayInfo(stringPtr("1")),
		Friday:    NewAvailableForWeekdayInfo(stringPtr("1")),
		Saturday:  NewAvailableForWeekdayInfo(stringPtr("1")),
		Sunday:    NewAvailableForWeekdayInfo(stringPtr("1")),
		StartDate: NewDate(stringPtr("20200101")),
		EndDate:   NewDate(stringPtr("20200102")),
	}

	testCases := make(map[string]ggtfsTestCase)
	testCases["1"] = ggtfsTestCase{
		csvRows: [][]string{
			{"service_id", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday", "start_date", "end_date"},
			{"111", "1", "1", "1", "1", "1", "1", "1", "20200101", "20200102"},
		},
		expectedStructs: []interface{}{&expected1},
	}

	return testCases
}
