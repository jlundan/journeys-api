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
	testCases["1"] = ggtfsTestCase{
		csvRows: [][]string{
			{"service_id"},
			{","},
		},
		expectedErrors: []string{
			"calendar.txt: record on line 2: wrong number of fields",
		},
	}
	testCases["2"] = ggtfsTestCase{
		csvRows: [][]string{
			{"service_id"},
			{" "},
		},
		expectedErrors: []string{
			"calendar.txt:0: invalid field: service_id",
			"calendar.txt:0: missing mandatory field: end_date",
			"calendar.txt:0: missing mandatory field: friday",
			"calendar.txt:0: missing mandatory field: monday",
			"calendar.txt:0: missing mandatory field: saturday",
			"calendar.txt:0: missing mandatory field: start_date",
			"calendar.txt:0: missing mandatory field: sunday",
			"calendar.txt:0: missing mandatory field: thursday",
			"calendar.txt:0: missing mandatory field: tuesday",
			"calendar.txt:0: missing mandatory field: wednesday",
		},
	}

	testCases["3"] = ggtfsTestCase{
		csvRows: [][]string{
			{"service_id", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday", "start_date", "end_date"},
			{"1000", "not an int", "10", "1", "1", "1", "1", "1", "20210101", "20201011"},
			{"1000", "1", "1", "1", "1", "1", "1", "1", "20210101", "20201011"},
			{"1002", "1", "1", "1", "1", "1", "1", "1", "202x0101", "20201011"},
			{"1003", "1", "1", "1", "1", "1", "1", "1", "2021x101", "20201011"},
			{"1004", "1", "1", "1", "1", "1", "1", "1", "20210x01", "20201011"},
			{"1005", "1", "1", "1", "1", "1", "1", "1", "202101x1", "20201011"},
			{"1006", "1", "1", "1", "1", "1", "1", "1", "2021011", "20201011"},
		},
		expectedErrors: []string{
			"calendar.txt:0: invalid field: monday",
			"calendar.txt:0: invalid field: tuesday",
			"calendar.txt:2: invalid field: start_date",
			"calendar.txt:3: invalid field: start_date",
			"calendar.txt:4: invalid field: start_date",
			"calendar.txt:5: invalid field: start_date",
			"calendar.txt:6: invalid field: start_date",
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
