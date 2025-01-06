//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"fmt"
	"testing"
)

func TestCreateCalendarItem(t *testing.T) {
	headerMap := map[string]int{"service_id": 0, "monday": 1, "tuesday": 2, "wednesday": 3,
		"thursday": 4, "friday": 5, "saturday": 6, "sunday": 7, "start_date": 8, "end_date": 9,
	}

	tests := map[string]struct {
		headers    map[string]int
		rows       [][]string
		lineNumber int
		expected   []*CalendarItem
	}{
		"empty-row": {
			headers: headerMap,
			rows:    [][]string{{"", "", "", "", "", "", "", "", "", ""}},
			expected: []*CalendarItem{{
				ServiceId:  stringPtr(""),
				Monday:     stringPtr(""),
				Tuesday:    stringPtr(""),
				Wednesday:  stringPtr(""),
				Thursday:   stringPtr(""),
				Friday:     stringPtr(""),
				Saturday:   stringPtr(""),
				Sunday:     stringPtr(""),
				StartDate:  stringPtr(""),
				EndDate:    stringPtr(""),
				LineNumber: 0,
			}},
		},
		"nil-values": {
			headers: headerMap,
			rows:    [][]string{nil},
			expected: []*CalendarItem{{
				ServiceId:  nil,
				Monday:     nil,
				Tuesday:    nil,
				Wednesday:  nil,
				Thursday:   nil,
				Friday:     nil,
				Saturday:   nil,
				Sunday:     nil,
				StartDate:  nil,
				EndDate:    nil,
				LineNumber: 0,
			}},
		},
		"OK": {
			headers: headerMap,
			rows: [][]string{
				{"111", "1", "1", "1", "1", "1", "0", "0", "20200101", "20200102"},
				{"112", "0", "0", "0", "0", "0", "1", "1", "20200101", "20200102"},
			},
			expected: []*CalendarItem{{
				ServiceId:  stringPtr("111"),
				Monday:     stringPtr("1"),
				Tuesday:    stringPtr("1"),
				Wednesday:  stringPtr("1"),
				Thursday:   stringPtr("1"),
				Friday:     stringPtr("1"),
				Saturday:   stringPtr("0"),
				Sunday:     stringPtr("0"),
				StartDate:  stringPtr("20200101"),
				EndDate:    stringPtr("20200102"),
				LineNumber: 0,
			}, {
				ServiceId:  stringPtr("112"),
				Monday:     stringPtr("0"),
				Tuesday:    stringPtr("0"),
				Wednesday:  stringPtr("0"),
				Thursday:   stringPtr("0"),
				Friday:     stringPtr("0"),
				Saturday:   stringPtr("1"),
				Sunday:     stringPtr("1"),
				StartDate:  stringPtr("20200101"),
				EndDate:    stringPtr("20200102"),
				LineNumber: 1,
			}},
		},
	}

	for name, tt := range tests {
		t.Run(fmt.Sprintf("%s", name), func(t *testing.T) {
			var actual []*CalendarItem
			for i, row := range tt.rows {
				actual = append(actual, CreateCalendarItem(row, tt.headers, i))
			}
			handleEntityCreateResults(t, tt.expected, actual)
		})
	}
}

func TestValidateCalendarItems(t *testing.T) {
	tests := map[string]struct {
		actualEntities  []*CalendarItem
		expectedResults []ValidationNotice
	}{
		"nil-slice": {
			actualEntities:  nil,
			expectedResults: []ValidationNotice{},
		},
		"nil-slice-items": {
			actualEntities:  []*CalendarItem{nil},
			expectedResults: []ValidationNotice{},
		},
		"invalid-fields": {
			actualEntities: []*CalendarItem{
				{
					ServiceId: stringPtr("111"), // avoid missing required field
					Monday:    stringPtr("-1"),
					Tuesday:   stringPtr("2"),
					Wednesday: stringPtr("0"),        // avoid missing required field
					Thursday:  stringPtr("0"),        // avoid missing required field
					Friday:    stringPtr("0"),        // avoid missing required field
					Saturday:  stringPtr("1"),        // avoid missing required field
					Sunday:    stringPtr("1"),        // avoid missing required field
					StartDate: stringPtr("20200101"), // avoid missing required field
					EndDate:   stringPtr("20200102"), // avoid missing required field
				},
			},
			expectedResults: []ValidationNotice{
				InvalidCalendarDayNotice{SingleLineNotice{FileName: "calendar.txt", FieldName: "monday"}},
				InvalidCalendarDayNotice{SingleLineNotice{FileName: "calendar.txt", FieldName: "tuesday"}},
			},
		},
	}

	for name, tt := range tests {
		t.Run(fmt.Sprintf("%s", name), func(t *testing.T) {
			handleValidationResults(t, ValidateCalendarItems(tt.actualEntities), tt.expectedResults)
		})
	}
}
