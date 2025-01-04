//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"fmt"
	"testing"
)

/*
	RouteId              *string // route_id                (required)
	ServiceId            *string // service_id              (required)
	Id                   *string // trip_id                 (required)
	HeadSign             *string // trip_headsign           (optional)
	ShortName            *string // trip_short_name         (optional)
	DirectionId          *string // direction_id            (optional)
	BlockId              *string // block_id                (optional)
	ShapeId              *string // shape_id                (conditionally required)
	WheelchairAccessible *string // wheelchair_accessible   (optional)
	BikesAllowed         *string // bikes_allowed           (optional)
*/

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
		expectedResults []Result
		routes          []*Route
		calendarItems   []*CalendarItem
		shapes          []*Shape
	}{
		"nil-slice": {
			actualEntities:  nil,
			expectedResults: []Result{},
		},
		"nil-slice-items": {
			actualEntities:  []*Trip{nil},
			expectedResults: []Result{},
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
			expectedResults: []Result{
				InvalidDirectionIdResult{SingleLineResult{FileName: "trips.txt", FieldName: "direction_id"}},
				InvalidWheelchairAccessibleResult{SingleLineResult{FileName: "trips.txt", FieldName: "wheelchair_accessible"}},
				InvalidBikesAllowedResult{SingleLineResult{FileName: "trips.txt", FieldName: "bikes_allowed"}},
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
			expectedResults: []Result{
				ForeignKeyViolationResult{
					ReferencingFileName:  "trips.txt",
					ReferencingFieldName: "route_id",
					ReferencedFieldName:  "routes.txt",
					ReferencedFileName:   "route_id",
					OffendingValue:       "ROUTE_1",
					ReferencedAtRow:      0,
				},
				ForeignKeyViolationResult{
					ReferencingFileName:  "trips.txt",
					ReferencingFieldName: "service_id",
					ReferencedFieldName:  "calendar.txt",
					ReferencedFileName:   "service_id",
					OffendingValue:       "SERVICE_1",
					ReferencedAtRow:      0,
				},
				ForeignKeyViolationResult{
					ReferencingFileName:  "trips.txt",
					ReferencingFieldName: "shape_id",
					ReferencedFieldName:  "shapes.txt",
					ReferencedFileName:   "shape_id",
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

//
//var validTripHeaders = []string{"route_id", "service_id", "trip_id", "trip_headsign", "trip_short_name",
//	"direction_id", "block_id", "shape_id", "wheelchair_accessible", "bikes_allowed"}
//
//func TestShouldReturnEmptyTripArrayOnEmptyString(t *testing.T) {
//	trips, errors := LoadEntitiesFromCSV[*Trip](csv.NewReader(strings.NewReader("")), validTripHeaders, CreateTrip, TripsFileName)
//	if len(errors) > 0 {
//		t.Error(errors)
//	}
//	if len(trips) != 0 {
//		t.Error("expected zero calendar items")
//	}
//}
//
//func TestTripParsing(t *testing.T) {
//	loadTripsFunc := func(reader *csv.Reader) ([]interface{}, []error) {
//		trips, errs := LoadEntitiesFromCSV[*Trip](reader, validTripHeaders, CreateTrip, TripsFileName)
//		entities := make([]interface{}, len(trips))
//		for i, trip := range trips {
//			entities[i] = trip
//		}
//		return entities, errs
//	}
//
//	validateTripsFunc := func(entities []interface{}, _ map[string][]interface{}) ([]error, []string) {
//		trips := make([]*Trip, len(entities))
//		for i, entity := range entities {
//			if trip, ok := entity.(*Trip); ok {
//				trips[i] = trip
//			}
//		}
//
//		return ValidateTrips(trips, nil, nil, nil)
//	}
//
//	runGenericGTFSParseTest(t, "TripOKTestcases", loadTripsFunc, validateTripsFunc, false, getTripOKTestcases())
//	runGenericGTFSParseTest(t, "TripNOKTestcases", loadTripsFunc, validateTripsFunc, false, getTripNOKTestcases())
//}
//
//func getTripOKTestcases() map[string]ggtfsTestCase {
//	expected1 := Trip{
//		Id:                   NewID(stringPtr("1")),
//		RouteId:              NewID(stringPtr("2")),
//		ServiceId:            NewID(stringPtr("3")),
//		HeadSign:             NewText(stringPtr("Trip 1")),
//		ShortName:            NewText(stringPtr("trip1")),
//		DirectionId:          NewDirectionId(stringPtr("0")),
//		BlockId:              NewID(stringPtr("4")),
//		ShapeId:              NewID(stringPtr("5")),
//		WheelchairAccessible: NewWheelchairAccessible(stringPtr("1")),
//		BikesAllowed:         NewBikesAllowed(stringPtr("2")),
//		LineNumber:           2,
//	}
//
//	testCases := make(map[string]ggtfsTestCase)
//	testCases["1"] = ggtfsTestCase{
//		csvRows: [][]string{
//			{"route_id", "service_id", "trip_id", "trip_headsign", "trip_short_name", "direction_id", "block_id", "shape_id", "wheelchair_accessible", "bikes_allowed"},
//			{"2", "3", "1", "Trip 1", "trip1", "0", "4", "5", "1", "2"},
//		},
//		expectedStructs: []interface{}{&expected1},
//	}
//
//	return testCases
//}
//
//func getTripNOKTestcases() map[string]ggtfsTestCase {
//	testCases := make(map[string]ggtfsTestCase)
//	testCases["invalid-fields-must-error-out"] = ggtfsTestCase{
//		csvRows: [][]string{
//			{"route_id", "service_id", "trip_id", "trip_headsign", "trip_short_name",
//				"direction_id", "block_id", "shape_id", "wheelchair_accessible", "bikes_allowed"},
//			{"", "", "", "", "", "", "", "", "", ""},
//		},
//		expectedErrors: []string{
//			"trips.txt:2: invalid mandatory field: route_id",
//			"trips.txt:2: invalid mandatory field: service_id",
//			"trips.txt:2: invalid mandatory field: trip_id",
//		},
//	}
//
//	return testCases
//}
