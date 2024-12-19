//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"encoding/csv"
	"strings"
	"testing"
)

var validTripHeaders = []string{"route_id", "service_id", "trip_id", "trip_headsign", "trip_short_name",
	"direction_id", "block_id", "shape_id", "wheelchair_accessible", "bikes_allowed"}

func TestShouldReturnEmptyTripArrayOnEmptyString(t *testing.T) {
	trips, errors := LoadEntities[*Trip](csv.NewReader(strings.NewReader("")), validTripHeaders, CreateTrip, TripsFileName)
	if len(errors) > 0 {
		t.Error(errors)
	}
	if len(trips) != 0 {
		t.Error("expected zero calendar items")
	}
}

func TestTripParsing(t *testing.T) {
	loadTripsFunc := func(reader *csv.Reader) ([]interface{}, []error) {
		trips, errs := LoadEntities[*Trip](reader, validTripHeaders, CreateTrip, TripsFileName)
		entities := make([]interface{}, len(trips))
		for i, trip := range trips {
			entities[i] = trip
		}
		return entities, errs
	}

	validateTripsFunc := func(entities []interface{}, _ map[string][]interface{}) ([]error, []string) {
		trips := make([]*Trip, len(entities))
		for i, entity := range entities {
			if trip, ok := entity.(*Trip); ok {
				trips[i] = trip
			}
		}

		return ValidateTrips(trips, nil, nil, nil)
	}

	runGenericGTFSParseTest(t, "TripOKTestcases", loadTripsFunc, validateTripsFunc, false, getTripOKTestcases())
	runGenericGTFSParseTest(t, "TripNOKTestcases", loadTripsFunc, validateTripsFunc, false, getTripNOKTestcases())
}

func getTripOKTestcases() map[string]ggtfsTestCase {
	expected1 := Trip{
		Id:                   NewID(stringPtr("1")),
		RouteId:              NewID(stringPtr("2")),
		ServiceId:            NewID(stringPtr("3")),
		HeadSign:             NewOptionalText(stringPtr("Trip 1")),
		ShortName:            NewOptionalText(stringPtr("trip1")),
		DirectionId:          NewOptionalDirectionId(stringPtr("0")),
		BlockId:              NewOptionalID(stringPtr("4")),
		ShapeId:              NewOptionalID(stringPtr("5")),
		WheelchairAccessible: NewOptionalWheelchairAccessible(stringPtr("1")),
		BikesAllowed:         NewOptionalBikesAllowed(stringPtr("2")),
	}

	testCases := make(map[string]ggtfsTestCase)
	testCases["1"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "service_id", "trip_id", "trip_headsign", "trip_short_name", "direction_id", "block_id", "shape_id", "wheelchair_accessible", "bikes_allowed"},
			{"2", "3", "1", "Trip 1", "trip1", "0", "4", "5", "1", "2"},
		},
		expectedStructs: []interface{}{&expected1},
	}

	return testCases
}

func getTripNOKTestcases() map[string]ggtfsTestCase {
	testCases := make(map[string]ggtfsTestCase)
	testCases["1"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id"},
			{","},
			{" "},
		},
		expectedErrors: []string{
			"trips.txt: record on line 2: wrong number of fields",
		},
	}
	testCases["2"] = ggtfsTestCase{
		csvRows: [][]string{
			{"trip_id", "service_id", "route_id", "direction_id", "wheelchair_accessible", "bikes_allowed"},
			{"001", "002", "001", "invalid", "invalid", "invalid"},
			{"002", "002", "001", "5", "5", "5"},
			{"002", "002", "001", "0", "0", "0"},
		},
		expectedErrors: []string{
			//"trips.txt:0: bikes_allowed: strconv.ParseInt: parsing \"invalid\": invalid syntax",
			//"trips.txt:0: direction_id: strconv.ParseInt: parsing \"invalid\": invalid syntax",
			//"trips.txt:0: wheelchair_accessible: strconv.ParseInt: parsing \"invalid\": invalid syntax",
			//"trips.txt:1: bikes_allowed: invalid value",
			//"trips.txt:1: direction_id: invalid value",
			//"trips.txt:1: wheelchair_accessible: invalid value",
			//"trips.txt:2: non-unique id: trip_id",
		},
	}

	return testCases
}

// TODO: Integrate these to the tests run by runGenericGTFSParseTest
//func TestValidateTrips(t *testing.T) {
//	shapeId1 := "S001"
//	shapeId2 := "S002"
//
//	testCases := []struct {
//		trips          []*Trip
//		routes         []*Route
//		services       []*CalendarItem
//		shapes         []*Shape
//		expectedErrors []string
//	}{
//		{
//			trips: []*Trip{
//				{RouteId: "R001", ServiceId: "C001", ShapeId: &shapeId1, LineNumber: 0},
//				{RouteId: "R002", ServiceId: "C002", ShapeId: &shapeId2, LineNumber: 1},
//			},
//			routes: []*Route{{
//				Id: "R001",
//			}},
//			services: []*CalendarItem{{
//				ServiceId: "C001",
//			}},
//			shapes: []*Shape{{
//				Id: "S001",
//			}},
//			expectedErrors: []string{
//				"trips.txt:1: referenced route_id not found in routes.txt",
//				"trips.txt:1: referenced service_id not found in calendar.txt",
//				"trips.txt:1: referenced shape_id not found in shapes.txt",
//			},
//		},
//		{
//			trips:          nil,
//			expectedErrors: []string{},
//		},
//		{
//			trips: []*Trip{
//				{RouteId: "R001", ServiceId: "C001", ShapeId: &shapeId1, LineNumber: 0},
//				{RouteId: "R002", ServiceId: "C002", ShapeId: &shapeId2, LineNumber: 1},
//			},
//			expectedErrors: []string{},
//		},
//		{
//			trips:          []*Trip{nil},
//			expectedErrors: []string{},
//		},
//		{
//			trips: []*Trip{
//				{RouteId: "R001", ServiceId: "C001", ShapeId: &shapeId1, LineNumber: 0},
//				{RouteId: "R002", ServiceId: "C002", ShapeId: &shapeId2, LineNumber: 1},
//			},
//			routes: []*Route{nil},
//			expectedErrors: []string{
//				"trips.txt:0: referenced route_id not found in routes.txt",
//				"trips.txt:1: referenced route_id not found in routes.txt",
//			},
//		},
//		{
//			trips: []*Trip{
//				{RouteId: "R001", ServiceId: "C001", ShapeId: &shapeId1, LineNumber: 0},
//				{RouteId: "R002", ServiceId: "C002", ShapeId: &shapeId2, LineNumber: 1},
//			},
//			services: []*CalendarItem{nil},
//			expectedErrors: []string{
//				"trips.txt:0: referenced service_id not found in calendar.txt",
//				"trips.txt:1: referenced service_id not found in calendar.txt",
//			},
//		},
//		{
//			trips: []*Trip{
//				{RouteId: "R001", ServiceId: "C001", ShapeId: &shapeId1, LineNumber: 0},
//				{RouteId: "R002", ServiceId: "C002", ShapeId: &shapeId2, LineNumber: 1},
//			},
//			shapes: []*Shape{nil},
//			expectedErrors: []string{
//				"trips.txt:0: referenced shape_id not found in shapes.txt",
//				"trips.txt:1: referenced shape_id not found in shapes.txt",
//			},
//		},
//	}
//
//	for _, tc := range testCases {
//		err := ValidateTrips(tc.trips, tc.routes, tc.services, tc.shapes)
//		checkErrors(tc.expectedErrors, err, t)
//	}
//}
