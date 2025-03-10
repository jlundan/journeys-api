//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"fmt"
	"testing"
)

func TestCreateTrip(t *testing.T) {
	headerMap := map[string]int{"route_id": 0, "service_id": 1, "trip_id": 2, "trip_headsign": 3, "trip_short_name": 4,
		"direction_id": 5, "block_id": 6, "shape_id": 7, "wheelchair_accessible": 8, "bikes_allowed": 9}

	tests := map[string]struct {
		headers    map[string]int
		rows       [][]string
		lineNumber int
		expected   []*Trip
	}{
		"empty-row": {
			headers: headerMap,
			rows:    [][]string{{"", "", "", "", "", "", "", "", "", ""}},
			expected: []*Trip{{
				RouteId:              stringPtr(""),
				ServiceId:            stringPtr(""),
				Id:                   stringPtr(""),
				HeadSign:             stringPtr(""),
				ShortName:            stringPtr(""),
				DirectionId:          stringPtr(""),
				BlockId:              stringPtr(""),
				ShapeId:              stringPtr(""),
				WheelchairAccessible: stringPtr(""),
				BikesAllowed:         stringPtr(""),
				LineNumber:           0,
			}},
		},
		"nil-values": {
			headers: headerMap,
			rows:    [][]string{nil},
			expected: []*Trip{{
				RouteId:              nil,
				ServiceId:            nil,
				Id:                   nil,
				HeadSign:             nil,
				ShortName:            nil,
				DirectionId:          nil,
				BlockId:              nil,
				ShapeId:              nil,
				WheelchairAccessible: nil,
				BikesAllowed:         nil,
				LineNumber:           0,
			}},
		},
		"OK": {
			headers: headerMap,
			rows: [][]string{
				{"route id", "service id", "trip id", "headsign", "shortname", "0", "block id", "shape id", "1", "2"},
			},
			expected: []*Trip{{
				RouteId:              stringPtr("route id"),
				ServiceId:            stringPtr("service id"),
				Id:                   stringPtr("trip id"),
				HeadSign:             stringPtr("headsign"),
				ShortName:            stringPtr("shortname"),
				DirectionId:          stringPtr("0"),
				BlockId:              stringPtr("block id"),
				ShapeId:              stringPtr("shape id"),
				WheelchairAccessible: stringPtr("1"),
				BikesAllowed:         stringPtr("2"),
				LineNumber:           0,
			}},
		},
	}

	for name, tt := range tests {
		t.Run(fmt.Sprintf("%s", name), func(t *testing.T) {
			var actual []*Trip
			for i, row := range tt.rows {
				actual = append(actual, CreateTrip(row, tt.headers, i))
			}
			handleEntityCreateResults(t, tt.expected, actual)
		})
	}
}

func TestValidateTrips(t *testing.T) {
	tests := map[string]struct {
		actualEntities  []*Trip
		expectedResults []ValidationNotice
		routes          []*Route
		calendarItems   []*CalendarItem
		shapes          []*Shape
	}{
		"nil-slice": {
			actualEntities:  nil,
			expectedResults: []ValidationNotice{},
		},
		"nil-slice-items": {
			actualEntities:  []*Trip{nil},
			expectedResults: []ValidationNotice{},
		},
		"invalid-fields": {
			actualEntities: []*Trip{
				{
					RouteId:              stringPtr("route id"),
					ServiceId:            stringPtr("service id"),
					Id:                   stringPtr("trip id"),
					DirectionId:          stringPtr("3"),
					WheelchairAccessible: stringPtr("5"),
					BikesAllowed:         stringPtr("5"),
				},
			},
			expectedResults: []ValidationNotice{
				InvalidDirectionIdNotice{SingleLineNotice{FileName: "trips.txt", FieldName: "direction_id"}},
				InvalidWheelchairAccessibleNotice{SingleLineNotice{FileName: "trips.txt", FieldName: "wheelchair_accessible"}},
				InvalidBikesAllowedNotice{SingleLineNotice{FileName: "trips.txt", FieldName: "bikes_allowed"}},
			},
		},
		"missing-foreign-keys": {
			actualEntities: []*Trip{
				{
					RouteId:   stringPtr("ROUTE_1"),
					ServiceId: stringPtr("SERVICE_1"),
					ShapeId:   stringPtr("SHAPE_1"),
					Id:        stringPtr("trip id"),
				},
			},
			calendarItems: []*CalendarItem{},
			routes:        []*Route{},
			shapes:        []*Shape{},
			expectedResults: []ValidationNotice{
				ForeignKeyViolationNotice{
					ReferencingFileName:  "trips.txt",
					ReferencingFieldName: "route_id",
					ReferencedFieldName:  "route_id",
					ReferencedFileName:   "routes.txt",
					OffendingValue:       "ROUTE_1",
					ReferencedAtRow:      0,
				},
				ForeignKeyViolationNotice{
					ReferencingFileName:  "trips.txt",
					ReferencingFieldName: "service_id",
					ReferencedFieldName:  "service_id",
					ReferencedFileName:   "calendar.txt",
					OffendingValue:       "SERVICE_1",
					ReferencedAtRow:      0,
				},
				ForeignKeyViolationNotice{
					ReferencingFileName:  "trips.txt",
					ReferencingFieldName: "shape_id",
					ReferencedFieldName:  "shape_id",
					ReferencedFileName:   "shapes.txt",
					OffendingValue:       "SHAPE_1",
					ReferencedAtRow:      0,
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(fmt.Sprintf("%s", name), func(t *testing.T) {
			handleValidationResults(t, ValidateTrips(tt.actualEntities, tt.routes, tt.calendarItems, tt.shapes), tt.expectedResults)
		})
	}
}
