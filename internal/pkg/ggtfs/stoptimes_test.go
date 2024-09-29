//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"encoding/csv"
	"strings"
	"testing"
)

func TestShouldReturnEmptyStopTimeArrayOnEmptyString(t *testing.T) {
	stopTimes, errors := LoadStopTimes(csv.NewReader(strings.NewReader("")))
	if len(errors) > 0 {
		t.Error(errors)
	}
	if len(stopTimes) != 0 {
		t.Error("expected zero calendar items")
	}
}

func TestStopTimeParsing(t *testing.T) {
	loadStopTimesFunc := func(reader *csv.Reader) ([]interface{}, []error) {
		stopTimes, errs := LoadStopTimes(reader)
		entities := make([]interface{}, len(stopTimes))
		for i, stopTime := range stopTimes {
			entities[i] = stopTime
		}
		return entities, errs
	}

	validateStopTimesFunc := func(entities []interface{}) []error {
		stopTimes := make([]*StopTime, len(entities))
		for i, entity := range entities {
			if stopTime, ok := entity.(*StopTime); ok {
				stopTimes[i] = stopTime
			}
		}

		return ValidateStopTimes(stopTimes, nil)
	}

	runGenericGTFSParseTest(t, "StopTimeOKTestcases", loadStopTimesFunc, validateStopTimesFunc, false, getStopTimeOKTestcases())
	runGenericGTFSParseTest(t, "StopTimeNOKTestcases", loadStopTimesFunc, validateStopTimesFunc, false, getStopTimeNOKTestcases())
}

func getStopTimeOKTestcases() map[string]ggtfsTestCase {
	headSign := "Foo city"
	pickup := "1"
	dropOff := "0"
	cp := "2"
	cd := "3"
	sd := "100"
	tp := "1"

	expected1 := StopTime{
		TripId:            "1",
		ArrivalTime:       "00:00",
		DepartureTime:     "02:00",
		StopId:            "0001",
		StopSequence:      "1",
		StopHeadSign:      &headSign,
		PickupType:        &pickup,
		DropOffType:       &dropOff,
		ContinuousPickup:  &cp,
		ContinuousDropOff: &cd,
		ShapeDistTraveled: &sd,
		Timepoint:         &tp,
		LineNumber:        0,
	}

	testCases := make(map[string]ggtfsTestCase)
	testCases["1"] = ggtfsTestCase{
		csvRows: [][]string{
			{"trip_id", "arrival_time", "departure_time", "stop_id", "stop_sequence", "stop_headsign", "pickup_type", "drop_off_type", "continuous_pickup", "continuous_drop_off", "shape_dist_traveled", "timepoint"},
			{"1", "00:00", "02:00", "0001", "1", "Foo city", "1", "0", "2", "3", "100", "1"},
		},
		//FIXME: this is not right
		expectedErrors: []string{
			"stop_times.txt: trip (1) has less than two defined stop times",
			"stop_times.txt: trip (1) references to an unknown stop_id (0001)",
		},
		expectedStructs: []interface{}{&expected1},
	}

	return testCases
}

func getStopTimeNOKTestcases() map[string]ggtfsTestCase {
	testCases := make(map[string]ggtfsTestCase)
	testCases["1"] = ggtfsTestCase{
		csvRows: [][]string{
			{"trip_id"},
			{","},
			{" "},
		},
		expectedErrors: []string{
			"stop_times.txt: record on line 2: wrong number of fields",
			"stop_times.txt: trip () has less than two defined stop times",
			"stop_times.txt: trip () references to an unknown stop_id ()",
			"stop_times.txt:1: arrival_time must be specified",
			"stop_times.txt:1: departure_time must be specified",
			"stop_times.txt:1: stop_id must be specified",
			"stop_times.txt:1: stop_sequence must be specified",
			"stop_times.txt:1: trip_id must be specified",
		},
	}
	testCases["2"] = ggtfsTestCase{
		csvRows: [][]string{
			{"trip_id", "stop_sequence", "stop_id", "departure_time", "arrival_time", "pickup_type", "drop_off_type", "timepoint"},
			{"0000", "0", "0000", "11:11", "22:22", "invalid", "invalid", "invalid"},
			{"0001", "0", "0001", "11:11", "22:22", "5", "5", "5"},
		},
		expectedErrors: []string{
			"stop_times.txt: trip (0000) has less than two defined stop times",
			"stop_times.txt: trip (0000) references to an unknown stop_id (0000)",
			"stop_times.txt: trip (0001) has less than two defined stop times",
			"stop_times.txt: trip (0001) references to an unknown stop_id (0001)",
			//"stop_times.txt:0: drop_off_type: strconv.ParseInt: parsing \"invalid\": invalid syntax",
			//"stop_times.txt:0: pickup_type: strconv.ParseInt: parsing \"invalid\": invalid syntax",
			//"stop_times.txt:0: timepoint: strconv.ParseInt: parsing \"invalid\": invalid syntax",
			//"stop_times.txt:1: drop_off_type: invalid value",
			//"stop_times.txt:1: pickup_type: invalid value",
			//"stop_times.txt:1: timepoint: invalid value",
		},
	}

	return testCases
}

// TODO: Integrate these to the tests run by runGenericGTFSParseTest
//func TestValidateStoptimes(t *testing.T) {
//	testCases := []struct {
//		stopTimes      []*StopTime
//		stops          []*Stop
//		expectedErrors []string
//	}{
//		{
//			stopTimes:      nil,
//			expectedErrors: []string{},
//		},
//		{
//			stopTimes:      []*StopTime{nil},
//			expectedErrors: []string{},
//		},
//		{
//			stopTimes: []*StopTime{
//				{TripId: "0001", StopId: "1000", lineNumber: 0},
//				{TripId: "0001", StopId: "1001", lineNumber: 1},
//				{TripId: "0001", StopId: "1005", lineNumber: 2},
//				{TripId: "0002", StopId: "1000", lineNumber: 3},
//			},
//			stops: []*Stop{
//				nil,
//				{Id: "1001", LineNumber: 1},
//				{Id: "1002", LineNumber: 2},
//			},
//			expectedErrors: []string{
//				"stop_times.txt: trip (0001) references to an unknown stop_id (1000)",
//				"stop_times.txt: trip (0001) references to an unknown stop_id (1005)",
//				"stop_times.txt: trip (0001) has less than two defined stop times",
//				"stop_times.txt: trip (0002) references to an unknown stop_id (1000)",
//				"stop_times.txt: trip (0002) has less than two defined stop times",
//			},
//		},
//		{
//			stopTimes: []*StopTime{
//				{TripId: "0001", StopId: "1000", lineNumber: 0},
//				{TripId: "0001", StopId: "1001", lineNumber: 1},
//				{TripId: "0001", StopId: "1005", lineNumber: 2},
//				{TripId: "0002", StopId: "1000", lineNumber: 3},
//			},
//			stops: []*Stop{
//				{Id: "1000", LineNumber: 0},
//				{Id: "1001", LineNumber: 1},
//				{Id: "1002", LineNumber: 2},
//			},
//			expectedErrors: []string{
//				"stop_times.txt: trip (0001) references to an unknown stop_id (1005)",
//				"stop_times.txt: trip (0002) has less than two defined stop times",
//			},
//		},
//	}
//
//	for _, tc := range testCases {
//		err := ValidateStoptimes(tc.stopTimes, tc.stops)
//		checkErrors(tc.expectedErrors, err, t)
//	}
//}
