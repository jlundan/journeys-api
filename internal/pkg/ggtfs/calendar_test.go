//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestCalendarCSVParsing(t *testing.T) {
	calendarItems, errs := LoadCalendarItems(csv.NewReader(strings.NewReader("")))
	if len(errs) > 0 {
		t.Error(errs)
	}
	if len(calendarItems) != 0 {
		t.Error("expected zero calendarItems")
	}

	reader := csv.NewReader(strings.NewReader("foo,bar\n1,2"))
	reader.Comma = ','
	reader.Comment = ','
	_, errs = LoadCalendarItems(reader)
	if len(errs) == 0 {
		t.Error("expected to throw error")
	}
}

func TestSCalendarParsingOK(t *testing.T) {
	s, err := parseDate("20200101", false)
	if err != nil {
		t.Error(err)
	}
	e, err := parseDate("20200102", true)
	if err != nil {
		t.Error(err)
	}

	expected1 := CalendarItem{
		ServiceId: "111",
		Monday:    1,
		Tuesday:   1,
		Wednesday: 1,
		Thursday:  1,
		Friday:    1,
		Saturday:  1,
		Sunday:    1,
		Start:     s,
		End:       e,
	}

	testCases := []struct {
		rows     [][]string
		expected CalendarItem
	}{
		{
			rows: [][]string{
				{"service_id", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday", "start_date", "end_date"},
				{"111", "1", "1", "1", "1", "1", "1", "1", "20200101", "20200102"},
			},
			expected: expected1,
		},
	}

	for _, tc := range testCases {
		calendarItems, err := LoadCalendarItems(csv.NewReader(strings.NewReader(tableToString(tc.rows))))
		if len(err) > 0 {
			t.Error(err)
			continue
		}

		if len(calendarItems) != 1 {
			t.Error("expected one row")
			continue
		}

		if !calendarItemsMatch(tc.expected, *calendarItems[0]) {
			c1, err := json.Marshal(tc.expected)
			if err != nil {
				t.Error(err)
			}
			c2, err := json.Marshal(*calendarItems[0])
			if err != nil {
				t.Error(err)
			}
			t.Error(fmt.Sprintf("expected %v, got %v", string(c1), string(c2)))
		}
	}
}

func TestCalendarParsingNOK(t *testing.T) {
	testCases := []struct {
		rows     [][]string
		expected []string
	}{
		{
			rows: [][]string{
				{"service_id"},
				{","},
			},
			expected: []string{
				"calendar.txt: record on line 2: wrong number of fields",
			},
		},
		{
			rows: [][]string{
				{"service_id"},
				{" "},
			},
			expected: []string{
				"calendar.txt:0: service_id must be specified",
				"calendar.txt:0: service_id: empty value not allowed",
				"calendar.txt:0: monday must be specified",
				"calendar.txt:0: tuesday must be specified",
				"calendar.txt:0: wednesday must be specified",
				"calendar.txt:0: thursday must be specified",
				"calendar.txt:0: friday must be specified",
				"calendar.txt:0: saturday must be specified",
				"calendar.txt:0: sunday must be specified",
				"calendar.txt:0: start_date must be specified",
				"calendar.txt:0: end_date must be specified",
			},
		},
		{
			rows: [][]string{
				{"service_id", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday", "start_date", "end_date"},
				{"1000", "not an int", "10", "1", "1", "1", "1", "1", "20210101", "20201011"},
				{"1000", "1", "1", "1", "1", "1", "1", "1", "20210101", "20201011"},
				{"1002", "1", "1", "1", "1", "1", "1", "1", "202x0101", "20201011"},
				{"1003", "1", "1", "1", "1", "1", "1", "1", "2021x101", "20201011"},
				{"1004", "1", "1", "1", "1", "1", "1", "1", "20210x01", "20201011"},
				{"1005", "1", "1", "1", "1", "1", "1", "1", "202101x1", "20201011"},
				{"1006", "1", "1", "1", "1", "1", "1", "1", "2021011", "20201011"},
			},
			expected: []string{
				"calendar.txt:0: monday must be specified",
				"calendar.txt:0: monday: strconv.ParseInt: parsing \"not an int\": invalid syntax",
				"calendar.txt:0: tuesday: invalid value",
				"calendar.txt:0: tuesday must be specified",
				"calendar.txt:1: non-unique id: service_id",
				"calendar.txt:2: start_date must be specified",
				"calendar.txt:2: start_date: strconv.ParseInt: parsing \"202x\": invalid syntax",
				"calendar.txt:3: start_date must be specified",
				"calendar.txt:3: start_date: strconv.ParseInt: parsing \"x1\": invalid syntax",
				"calendar.txt:4: start_date must be specified",
				"calendar.txt:4: start_date: strconv.ParseInt: parsing \"0x\": invalid syntax",
				"calendar.txt:5: start_date must be specified",
				"calendar.txt:5: start_date: strconv.ParseInt: parsing \"x1\": invalid syntax",
				"calendar.txt:6: start_date: invalid date format",
				"calendar.txt:6: start_date must be specified",
			},
		},
	}

	for _, tc := range testCases {
		_, err := LoadCalendarItems(csv.NewReader(strings.NewReader(tableToString(tc.rows))))

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

func calendarItemsMatch(a CalendarItem, b CalendarItem) bool {
	return a.ServiceId == b.ServiceId && a.Monday == b.Monday && a.Tuesday == b.Tuesday && a.Wednesday == b.Wednesday &&
		a.Thursday == b.Thursday && a.Friday == b.Friday && a.Saturday == b.Saturday && a.Sunday == b.Sunday &&
		a.Start.Unix() == b.Start.Unix() && a.End.Unix() == b.End.Unix()
}

func parseDate(str string, fillEnd bool) (time.Time, error) {
	year, err := strconv.ParseInt(str[:4], 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	month, err := strconv.ParseInt(str[4:6], 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	day, err := strconv.ParseInt(str[6:8], 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	if fillEnd {
		return time.Date(int(year), time.Month(int(month)), int(day), 23, 59, 59, 0, time.FixedZone("UTC+2", 2*60*60)), nil
	}
	return time.Date(int(year), time.Month(int(month)), int(day), 0, 0, 0, 0, time.FixedZone("UTC+2", 2*60*60)), nil
}
