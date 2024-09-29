//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"testing"
)

func TestCalendarDatesCSVParsing(t *testing.T) {
	items, errors := LoadCalendarDates(csv.NewReader(strings.NewReader("")))
	if len(errors) > 0 {
		t.Error(errors)
	}
	if len(items) != 0 {
		t.Error("expected zero items")
	}

	reader := csv.NewReader(strings.NewReader("foo,bar\n1,2"))
	reader.Comma = ','
	reader.Comment = ','
	_, errors = LoadCalendarDates(reader)
	if len(errors) == 0 {
		t.Error("expected to throw error")
	}
}

func TestCalendarDatesParsingOK(t *testing.T) {
	d, err := parseDate("20200101", false)
	if err != nil {
		t.Error(err)
	}

	expected1 := CalendarDate{
		ServiceId:     "1",
		Date:          d,
		ExceptionType: 1,
	}

	testCases := []struct {
		rows     [][]string
		expected CalendarDate
	}{
		{
			rows: [][]string{
				{"service_id", "date", "exception_type"},
				{"1", "20200101", "1"},
			},
			expected: expected1,
		},
	}

	for _, tc := range testCases {
		stops, err := LoadCalendarDates(csv.NewReader(strings.NewReader(tableToString(tc.rows))))
		if err != nil && len(err) > 0 {
			t.Error(err)
			continue
		}

		if len(stops) != 1 {
			t.Error("expected one row")
			continue
		}

		if !calendarDatesMatch(tc.expected, *stops[0]) {
			s1, err := json.Marshal(tc.expected)
			if err != nil {
				t.Error(err)
			}
			s2, err := json.Marshal(*stops[0])
			if err != nil {
				t.Error(err)
			}
			t.Error(fmt.Sprintf("expected %v, got %v", string(s1), string(s2)))
		}
	}
}

func TestCalendarDatesParsingNOK(t *testing.T) {
	testCases := []struct {
		rows     [][]string
		expected []string
	}{
		{
			rows: [][]string{
				{"service_id", "date", "exception_type"},
				{","},
				{" ", " ", "1"},
				{"1000", "20201011", "not an int"},
				{"1001", "20201011", "10"},
			},
			expected: []string{
				"calendar_dates.txt: record on line 2: wrong number of fields",
				"calendar_dates.txt:1: service_id must be specified",
				"calendar_dates.txt:1: service_id: empty value not allowed",
				"calendar_dates.txt:1: date must be specified",
				"calendar_dates.txt:1: date: empty value not allowed",
				"calendar_dates.txt:2: exception_type must be specified",
				"calendar_dates.txt:2: exception_type: strconv.ParseInt: parsing \"not an int\": invalid syntax",
				"calendar_dates.txt:3: exception_type: invalid value",
				"calendar_dates.txt:3: exception_type must be specified",
			},
		},
	}

	for _, tc := range testCases {
		_, err := LoadCalendarDates(csv.NewReader(strings.NewReader(tableToString(tc.rows))))

		sort.Slice(err, func(x, y int) bool {
			return err[x].Error() < err[y].Error()
		})

		sort.Slice(tc.expected, func(x, y int) bool {
			return tc.expected[x] < tc.expected[y]
		})

		if len(err) == 0 {
			t.Error("expected to throw an error")
			continue
		}

		if len(err) != len(tc.expected) {
			t.Error(fmt.Sprintf("expected %v errors, got %v", len(tc.expected), len(err)))
			for _, e := range err {
				fmt.Println(e)
			}
			continue
		}

		for i, e := range err {
			if e.Error() != tc.expected[i] {
				t.Error(fmt.Sprintf("expected error %s, got %s", tc.expected[i], e.Error()))
			}
		}
	}
}

func TestValidateCalendarDates(t *testing.T) {
	testCases := []struct {
		calendarDates  []*CalendarDate
		calendarItems  []*CalendarItem
		expectedErrors []string
	}{
		{
			calendarDates: []*CalendarDate{
				{ServiceId: "1000", LineNumber: 0},
			},
			calendarItems: []*CalendarItem{
				{ServiceId: "1000", LineNumber: 0},
				{ServiceId: "1001", LineNumber: 1},
			},
			expectedErrors: []string{},
		},
		{
			calendarDates:  nil,
			expectedErrors: []string{},
		},
		{
			calendarDates: []*CalendarDate{nil},
			calendarItems: []*CalendarItem{
				{ServiceId: "1002", LineNumber: 0},
				{ServiceId: "1001", LineNumber: 1},
			},
			expectedErrors: []string{},
		},
		{
			calendarDates: []*CalendarDate{
				{ServiceId: "1000", LineNumber: 0},
			},
			calendarItems: []*CalendarItem{nil},
			expectedErrors: []string{
				"calendar_dates.txt:0: referenced service_id not found in calendar.txt",
			},
		},
		{
			calendarDates: []*CalendarDate{
				{ServiceId: "1000", LineNumber: 0},
			},
			calendarItems: []*CalendarItem{
				{ServiceId: "1002", LineNumber: 0},
				{ServiceId: "1001", LineNumber: 1},
			},
			expectedErrors: []string{
				"calendar_dates.txt:0: referenced service_id not found in calendar.txt",
			},
		},
	}

	for _, tc := range testCases {
		err := ValidateCalendarDates(tc.calendarDates, tc.calendarItems)
		checkErrors(tc.expectedErrors, err, t)
	}
}

func calendarDatesMatch(a CalendarDate, b CalendarDate) bool {
	return a.ServiceId == b.ServiceId && a.Date.Unix() == b.Date.Unix() && a.ExceptionType == b.ExceptionType
}
