//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"fmt"
	"testing"
)

func TestCreateCalendarDate(t *testing.T) {
	headerMap := map[string]int{"service_id": 0, "date": 1, "exception_type": 2}

	tests := map[string]struct {
		headers    map[string]int
		rows       [][]string
		lineNumber int
		expected   []*CalendarDate
	}{
		"empty-row": {
			headers: headerMap,
			rows:    [][]string{{"", "", ""}},
			expected: []*CalendarDate{{
				ServiceId:     stringPtr(""),
				Date:          stringPtr(""),
				ExceptionType: stringPtr(""),
				LineNumber:    0,
			}},
		},
		"nil-values": {
			headers: headerMap,
			rows:    [][]string{nil},
			expected: []*CalendarDate{{
				ServiceId:     nil,
				Date:          nil,
				ExceptionType: nil,
				LineNumber:    0,
			}},
		},
		"OK": {
			headers: headerMap,
			rows: [][]string{
				{"111", "20200101", "1"},
				{"111", "20201201", "2"},
			},
			expected: []*CalendarDate{{
				ServiceId:     stringPtr("111"),
				Date:          stringPtr("20200101"),
				ExceptionType: stringPtr("1"),
				LineNumber:    0,
			}, {
				ServiceId:     stringPtr("111"),
				Date:          stringPtr("20201201"),
				ExceptionType: stringPtr("2"),
				LineNumber:    1,
			}},
		},
	}

	for name, tt := range tests {
		t.Run(fmt.Sprintf("%s", name), func(t *testing.T) {
			var actual []*CalendarDate
			for i, row := range tt.rows {
				actual = append(actual, CreateCalendarDate(row, tt.headers, i))
			}
			handleEntityCreateResults(t, tt.expected, actual)
		})
	}
}

func TestValidateCalendarDates(t *testing.T) {
	tests := map[string]struct {
		actualEntities  []*CalendarDate
		expectedResults []Result
		calendarItems   []*CalendarItem
	}{
		"nil-slice": {
			actualEntities:  nil,
			expectedResults: []Result{},
		},
		"nil-slice-items": {
			actualEntities:  []*CalendarDate{nil},
			expectedResults: []Result{},
		},
		"invalid-fields": {
			actualEntities: []*CalendarDate{
				{
					ServiceId:     stringPtr("111"), // avoid missing required field
					Date:          stringPtr("Not a date"),
					ExceptionType: stringPtr("3"),
				},
			},
			expectedResults: []Result{
				InvalidDateResult{SingleLineResult{FileName: "calendar_dates.txt", FieldName: "date"}},
				InvalidCalendarExceptionResult{SingleLineResult{FileName: "calendar_dates.txt", FieldName: "exception_type"}},
			},
		},
		"empty-calendar-item-slice": {
			actualEntities: []*CalendarDate{
				{
					ServiceId:     stringPtr("111"), // avoid missing required field
					Date:          stringPtr("20201201"),
					ExceptionType: stringPtr("1"),
				},
			},
			calendarItems: []*CalendarItem{},
			expectedResults: []Result{
				ForeignKeyViolationResult{
					ReferencingFileName:  "calendar_dates.txt",
					ReferencingFieldName: "service_id",
					ReferencedFieldName:  "calendar.txt",
					ReferencedFileName:   "service_id",
					OffendingValue:       "111",
					ReferencedAtRow:      0,
				},
			},
		},
		"empty-calendar-item-slice-item": {
			actualEntities: []*CalendarDate{
				{
					ServiceId:     stringPtr("111"), // avoid missing required field
					Date:          stringPtr("20201201"),
					ExceptionType: stringPtr("1"),
				},
			},
			calendarItems: []*CalendarItem{nil},
			expectedResults: []Result{
				ForeignKeyViolationResult{
					ReferencingFileName:  "calendar_dates.txt",
					ReferencingFieldName: "service_id",
					ReferencedFieldName:  "calendar.txt",
					ReferencedFileName:   "service_id",
					OffendingValue:       "111",
					ReferencedAtRow:      0,
				},
			},
		},
		"missing-calendar-item": {
			actualEntities: []*CalendarDate{
				{
					ServiceId:     stringPtr("111"), // avoid missing required field
					Date:          stringPtr("20201201"),
					ExceptionType: stringPtr("1"),
				},
			},
			calendarItems: []*CalendarItem{{ServiceId: stringPtr("112")}},
			expectedResults: []Result{
				ForeignKeyViolationResult{
					ReferencingFileName:  "calendar_dates.txt",
					ReferencingFieldName: "service_id",
					ReferencedFieldName:  "calendar.txt",
					ReferencedFileName:   "service_id",
					OffendingValue:       "111",
					ReferencedAtRow:      0,
				},
			},
		},
		"missing-calendar-date-with-calendar-items": {
			actualEntities:  []*CalendarDate{nil},
			calendarItems:   []*CalendarItem{{ServiceId: stringPtr("112")}},
			expectedResults: []Result{},
		},
	}

	for name, tt := range tests {
		t.Run(fmt.Sprintf("%s", name), func(t *testing.T) {
			handleValidationResults(t, ValidateCalendarDates(tt.actualEntities, tt.calendarItems), tt.expectedResults)
		})
	}
}
