//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"encoding/csv"
	"fmt"
	"strings"
	"testing"
)

var validStopTimeHeaders = []string{"trip_id", "arrival_time", "departure_time", "stop_id", "stop_sequence",
	"stop_headsign", "pickup_type", "drop_off_type", "continuous_pickup", "continuous_drop_off",
	"shape_dist_traveled", "timepoint"}

func TestShouldReturnEmptyStopTimeArrayOnEmptyString(t *testing.T) {
	stopTimes, errors := LoadEntitiesFromCSV[*StopTime](csv.NewReader(strings.NewReader("")), validStopTimeHeaders, CreateStopTime, StopTimesFileName)
	if len(errors) > 0 {
		t.Error(errors)
	}
	if len(stopTimes) != 0 {
		t.Error("expected zero calendar items")
	}
}

func TestStopTimeParsing(t *testing.T) {
	loadStopTimesFunc := func(reader *csv.Reader) ([]interface{}, []error) {
		stopTimes, errs := LoadEntitiesFromCSV[*StopTime](reader, validStopTimeHeaders, CreateStopTime, StopTimesFileName)
		entities := make([]interface{}, len(stopTimes))
		for i, stopTime := range stopTimes {
			entities[i] = stopTime
		}
		return entities, errs
	}

	validateStopTimesFunc := func(entities []interface{}, fixtures map[string][]interface{}) ([]error, []string) {
		stopTimes := make([]*StopTime, len(entities))
		for i, entity := range entities {
			if stopTime, ok := entity.(*StopTime); ok {
				stopTimes[i] = stopTime
			}
		}

		stoPCount := len(fixtures["stops"])

		if stoPCount == 0 {
			return ValidateStopTimes(stopTimes, nil)
		}

		stops := make([]*Stop, stoPCount)
		for i, fixture := range fixtures["stops"] {
			if stop, ok := fixture.(*Stop); ok {
				stops[i] = stop
			} else {
				t.Error(fmt.Sprintf("test setup error: cannot convert %v to Stop pointer. maybe you used value instead of pointer when setting fixtures", fixture))
			}
		}

		return ValidateStopTimes(stopTimes, stops)
	}

	runGenericGTFSParseTest(t, "StopTimeOKTestcases", loadStopTimesFunc, validateStopTimesFunc, false, getStopTimeOKTestcases())
	runGenericGTFSParseTest(t, "StopTimeNOKTestcases", loadStopTimesFunc, validateStopTimesFunc, false, getStopTimeNOKTestcases())
}

func getStopTimeOKTestcases() map[string]ggtfsTestCase {
	expected1 := StopTime{
		TripId:            NewID(stringPtr("1")),
		ArrivalTime:       NewTime(stringPtr("00:10:00")),
		DepartureTime:     NewTime(stringPtr("00:11:00")),
		StopId:            NewID(stringPtr("0001")),
		StopSequence:      NewInteger(stringPtr("1")),
		StopHeadSign:      NewText(stringPtr("Foo city")),
		PickupType:        NewPickupType(stringPtr("1")),
		DropOffType:       NewDropOffType(stringPtr("0")),
		ContinuousPickup:  NewContinuousPickupType(stringPtr("2")),
		ContinuousDropOff: NewContinuousDropOffType(stringPtr("3")),
		ShapeDistTraveled: NewFloat(stringPtr("100")),
		Timepoint:         NewTimePoint(stringPtr("1")),
		LineNumber:        2,
	}

	expected2 := StopTime{
		TripId:            NewID(stringPtr("1")),
		ArrivalTime:       NewTime(stringPtr("00:20:00")),
		DepartureTime:     NewTime(stringPtr("00:21:00")),
		StopId:            NewID(stringPtr("0002")),
		StopSequence:      NewInteger(stringPtr("2")),
		StopHeadSign:      NewText(stringPtr("Bar city")),
		PickupType:        NewPickupType(stringPtr("1")),
		DropOffType:       NewDropOffType(stringPtr("0")),
		ContinuousPickup:  NewContinuousPickupType(stringPtr("2")),
		ContinuousDropOff: NewContinuousDropOffType(stringPtr("3")),
		ShapeDistTraveled: NewFloat(stringPtr("100")),
		Timepoint:         NewTimePoint(stringPtr("1")),
		LineNumber:        3,
	}

	testCases := make(map[string]ggtfsTestCase)
	testCases["1"] = ggtfsTestCase{
		csvRows: [][]string{
			{"trip_id", "arrival_time", "departure_time", "stop_id", "stop_sequence", "stop_headsign", "pickup_type", "drop_off_type", "continuous_pickup", "continuous_drop_off", "shape_dist_traveled", "timepoint"},
			{"1", "00:10:00", "00:11:00", "0001", "1", "Foo city", "1", "0", "2", "3", "100", "1"},
			{"1", "00:20:00", "00:21:00", "0002", "2", "Bar city", "1", "0", "2", "3", "100", "1"},
		},
		expectedErrors:  []string{},
		expectedStructs: []interface{}{&expected1, &expected2},
		fixtures: map[string][]interface{}{
			"stops": {
				&Stop{
					Id: stringPtr("0001"),
				},
				&Stop{
					Id: stringPtr("0002"),
				},
			},
		},
	}

	return testCases
}

func getStopTimeNOKTestcases() map[string]ggtfsTestCase {
	testCases := make(map[string]ggtfsTestCase)
	testCases["invalid-fields-must-error-out"] = ggtfsTestCase{
		csvRows: [][]string{
			{"trip_id", "arrival_time", "departure_time", "stop_id", "stop_sequence",
				"stop_headsign", "pickup_type", "drop_off_type", "continuous_pickup", "continuous_drop_off",
				"shape_dist_traveled", "timepoint"},
			{"", "", "", "", "", "", "", "", "", "", "", ""},
			{"ID1", "00:01:00", "50:00:00", "STOP1", "1", "STOP 1", "1", "1", "1", "1", "100", "1"},
		},
		expectedErrors: []string{
			"stop_times.txt:2: invalid mandatory field: stop_sequence",
			"stop_times.txt:2: invalid mandatory field: trip_id",
			"stop_times.txt:2: stop_id () references to an unknown stop",
			"stop_times.txt:3: invalid field: departure_time",
			"stop_times.txt:3: stop_id (STOP1) references to an unknown stop",
		},
	}

	return testCases
}
