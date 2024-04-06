package ggtfs

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"testing"
)

func TestStopTimesCSVParsing(t *testing.T) {
	items, errors := LoadStopTimes(csv.NewReader(strings.NewReader("")))
	if len(errors) > 0 {
		t.Error(errors)
	}
	if len(items) != 0 {
		t.Error("expected zero items")
	}

	reader := csv.NewReader(strings.NewReader("foo,bar\n1,2"))
	reader.Comma = ','
	reader.Comment = ','
	_, errors = LoadStopTimes(reader)
	if len(errors) == 0 {
		t.Error("expected to throw error")
	}
}

func TestStopTimesParsingOK(t *testing.T) {
	headSign := "Foo city"
	pickup := 1
	dropOff := 0
	cp := 2
	cd := 3
	sd := float64(100)
	tp := 1

	expected1 := StopTime{
		TripId:            "1",
		ArrivalTime:       "00:00",
		DepartureTime:     "02:00",
		StopId:            "0001",
		StopSequence:      1,
		StopHeadSign:      &headSign,
		PickupType:        &pickup,
		DropOffType:       &dropOff,
		ContinuousPickup:  &cp,
		ContinuousDropOff: &cd,
		ShapeDistTraveled: &sd,
		Timepoint:         &tp,
		lineNumber:        0,
	}

	testCases := []struct {
		headers  map[string]uint8
		rows     [][]string
		expected StopTime
	}{
		{
			rows: [][]string{
				{"trip_id", "arrival_time", "departure_time", "stop_id", "stop_sequence", "stop_headsign", "pickup_type", "drop_off_type", "continuous_pickup", "continuous_drop_off", "shape_dist_traveled", "timepoint"},
				{"1", "00:00", "02:00", "0001", "1", "Foo city", "1", "0", "2", "3", "100", "1"},
			},
			expected: expected1,
		},
	}

	for _, tc := range testCases {
		stops, err := LoadStopTimes(csv.NewReader(strings.NewReader(tableToString(tc.rows))))
		if err != nil && len(err) > 0 {
			t.Error(err)
			continue
		}

		if len(stops) != 1 {
			t.Error("expected one row")
			continue
		}

		if !stopTimesMatch(tc.expected, *stops[0]) {
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

func TestStopTimesParsingNOK(t *testing.T) {
	testCases := []struct {
		rows     [][]string
		expected []string
	}{
		{
			rows: [][]string{
				{"trip_id"},
				{","},
				{" "},
			},
			expected: []string{
				"stop_times.txt: record on line 2: wrong number of fields",
				"stop_times.txt:1: trip_id must be specified",
				"stop_times.txt:1: trip_id: empty value not allowed",
				"stop_times.txt:1: stop_sequence must be specified",
				"stop_times.txt:1: stop_id must be specified",
				"stop_times.txt:1: departure_time must be specified",
				"stop_times.txt:1: arrival_time must be specified",
			},
		},
		{
			rows: [][]string{
				{"trip_id", "stop_sequence", "stop_id", "departure_time", "arrival_time", "pickup_type", "drop_off_type", "timepoint"},
				{"0000", "0", "0000", "11:11", "22:22", "invalid", "invalid", "invalid"},
				{"0001", "0", "0001", "11:11", "22:22", "5", "5", "5"},
			},
			expected: []string{
				"stop_times.txt:0: drop_off_type: strconv.ParseInt: parsing \"invalid\": invalid syntax",
				"stop_times.txt:0: pickup_type: strconv.ParseInt: parsing \"invalid\": invalid syntax",
				"stop_times.txt:0: timepoint: strconv.ParseInt: parsing \"invalid\": invalid syntax",
				"stop_times.txt:1: drop_off_type: invalid value",
				"stop_times.txt:1: pickup_type: invalid value",
				"stop_times.txt:1: timepoint: invalid value",
			},
		},
	}

	for _, tc := range testCases {
		_, err := LoadStopTimes(csv.NewReader(strings.NewReader(tableToString(tc.rows))))

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

func TestValidateStoptimes(t *testing.T) {
	testCases := []struct {
		stopTimes      []*StopTime
		stops          []*Stop
		expectedErrors []string
	}{
		{
			stopTimes:      nil,
			expectedErrors: []string{},
		},
		{
			stopTimes:      []*StopTime{nil},
			expectedErrors: []string{},
		},
		{
			stopTimes: []*StopTime{
				{TripId: "0001", StopId: "1000", lineNumber: 0},
				{TripId: "0001", StopId: "1001", lineNumber: 1},
				{TripId: "0001", StopId: "1005", lineNumber: 2},
				{TripId: "0002", StopId: "1000", lineNumber: 3},
			},
			stops: []*Stop{
				nil,
				{Id: "1001", lineNumber: 1},
				{Id: "1002", lineNumber: 2},
			},
			expectedErrors: []string{
				"stop_times.txt: trip (0001) references to an unknown stop_id (1000)",
				"stop_times.txt: trip (0001) references to an unknown stop_id (1005)",
				"stop_times.txt: trip (0001) has less than two defined stop times",
				"stop_times.txt: trip (0002) references to an unknown stop_id (1000)",
				"stop_times.txt: trip (0002) has less than two defined stop times",
			},
		},
		{
			stopTimes: []*StopTime{
				{TripId: "0001", StopId: "1000", lineNumber: 0},
				{TripId: "0001", StopId: "1001", lineNumber: 1},
				{TripId: "0001", StopId: "1005", lineNumber: 2},
				{TripId: "0002", StopId: "1000", lineNumber: 3},
			},
			stops: []*Stop{
				{Id: "1000", lineNumber: 0},
				{Id: "1001", lineNumber: 1},
				{Id: "1002", lineNumber: 2},
			},
			expectedErrors: []string{
				"stop_times.txt: trip (0001) references to an unknown stop_id (1005)",
				"stop_times.txt: trip (0002) has less than two defined stop times",
			},
		},
	}

	for _, tc := range testCases {
		err := ValidateStoptimes(tc.stopTimes, tc.stops)
		checkErrors(tc.expectedErrors, err, t)
	}
}

func stopTimesMatch(a StopTime, b StopTime) bool {
	return a.TripId == b.TripId && a.ArrivalTime == b.ArrivalTime && a.DepartureTime == b.DepartureTime && a.StopId == b.StopId && a.StopSequence == b.StopSequence &&
		*a.StopHeadSign == *b.StopHeadSign && *a.PickupType == *b.PickupType && *a.DropOffType == *b.DropOffType && *a.ContinuousPickup == *b.ContinuousPickup &&
		*a.ContinuousDropOff == *b.ContinuousDropOff && *a.ShapeDistTraveled == *b.ShapeDistTraveled && *a.Timepoint == *b.Timepoint
}
