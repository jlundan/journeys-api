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
	trips, errors := LoadEntitiesFromCSV[*Trip](csv.NewReader(strings.NewReader("")), validTripHeaders, CreateTrip, TripsFileName)
	if len(errors) > 0 {
		t.Error(errors)
	}
	if len(trips) != 0 {
		t.Error("expected zero calendar items")
	}
}

func TestTripParsing(t *testing.T) {
	loadTripsFunc := func(reader *csv.Reader) ([]interface{}, []error) {
		trips, errs := LoadEntitiesFromCSV[*Trip](reader, validTripHeaders, CreateTrip, TripsFileName)
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
		HeadSign:             NewText(stringPtr("Trip 1")),
		ShortName:            NewText(stringPtr("trip1")),
		DirectionId:          NewDirectionId(stringPtr("0")),
		BlockId:              NewID(stringPtr("4")),
		ShapeId:              NewID(stringPtr("5")),
		WheelchairAccessible: NewWheelchairAccessible(stringPtr("1")),
		BikesAllowed:         NewBikesAllowed(stringPtr("2")),
		LineNumber:           2,
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
	testCases["invalid-fields-must-error-out"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "service_id", "trip_id", "trip_headsign", "trip_short_name",
				"direction_id", "block_id", "shape_id", "wheelchair_accessible", "bikes_allowed"},
			{"", "", "", "", "", "", "", "", "", ""},
		},
		expectedErrors: []string{
			"trips.txt:2: invalid field: bikes_allowed",
			"trips.txt:2: invalid field: block_id",
			"trips.txt:2: invalid field: direction_id",
			"trips.txt:2: invalid field: shape_id",
			"trips.txt:2: invalid field: trip_headsign",
			"trips.txt:2: invalid field: trip_short_name",
			"trips.txt:2: invalid field: wheelchair_accessible",
			"trips.txt:2: invalid mandatory field: route_id",
			"trips.txt:2: invalid mandatory field: service_id",
			"trips.txt:2: invalid mandatory field: trip_id",
		},
	}

	return testCases
}
