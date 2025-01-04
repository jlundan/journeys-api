//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"fmt"
	"testing"
)

/*
	TripId                   *string // trip_id                      (required)
	ArrivalTime              *string // arrival_time                 (conditionally required)
	DepartureTime            *string // departure_time               (conditionally required)
	StopId                   *string // stop_id                      (conditionally required)
	LocationGroupId          *string // location_group_id            (conditionally forbidden)
	LocationId               *string // location_id                  (conditionally forbidden)
	StopSequence             *string // stop_sequence                (required)
	StopHeadSign             *string // stop_headsign                (optional)
	StartPickupDropOffWindow *string // start_pickup_drop_off_window (conditionally required)
	EndPickupDropOffWindow   *string // end_pickup_drop_off_window   (conditionally required)
	PickupType               *string // pickup_type                  (conditionally required)
	DropOffType              *string // drop_off_type                (conditionally required)
	ContinuousPickup         *string // continuous_pickup            (conditionally required)
	ContinuousDropOff        *string // continuous_drop_off          (conditionally required)
	ShapeDistTraveled        *string // shape_dist_traveled          (optional)
	Timepoint                *string // timepoint                    (optional)
	PickupBookingRuleId      *string // pickup_booking_rule_id       (optional)
	DropOffBookingRuleId     *string // drop_off_booking_rule_id     (optional)
*/

func TestCreateStopTime(t *testing.T) {
	headerMap := map[string]int{"trip_id": 0, "arrival_time": 1, "departure_time": 2, "stop_id": 3, "location_group_id": 4,
		"location_id": 5, "stop_sequence": 6, "stop_headsign": 7, "start_pickup_drop_off_window": 8,
		"end_pickup_drop_off_window": 9, "pickup_type": 10, "drop_off_type": 11, "continuous_pickup": 12,
		"continuous_drop_off": 13, "shape_dist_traveled": 14, "timepoint": 15, "pickup_booking_rule_id": 16, "drop_off_booking_rule_id": 17}

	tests := map[string]struct {
		headers    map[string]int
		rows       [][]string
		lineNumber int
		expected   []*StopTime
	}{
		"empty-row": {
			headers: headerMap,
			rows:    [][]string{{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""}},
			expected: []*StopTime{{
				TripId:                   stringPtr(""),
				ArrivalTime:              stringPtr(""),
				DepartureTime:            stringPtr(""),
				StopId:                   stringPtr(""),
				LocationGroupId:          stringPtr(""),
				LocationId:               stringPtr(""),
				StopSequence:             stringPtr(""),
				StopHeadSign:             stringPtr(""),
				StartPickupDropOffWindow: stringPtr(""),
				EndPickupDropOffWindow:   stringPtr(""),
				PickupType:               stringPtr(""),
				DropOffType:              stringPtr(""),
				ContinuousPickup:         stringPtr(""),
				ContinuousDropOff:        stringPtr(""),
				ShapeDistTraveled:        stringPtr(""),
				Timepoint:                stringPtr(""),
				PickupBookingRuleId:      stringPtr(""),
				DropOffBookingRuleId:     stringPtr(""),
				LineNumber:               0,
			}},
		},
		"nil-values": {
			headers: headerMap,
			rows:    [][]string{nil},
			expected: []*StopTime{{
				TripId:                   nil,
				ArrivalTime:              nil,
				DepartureTime:            nil,
				StopId:                   nil,
				LocationGroupId:          nil,
				LocationId:               nil,
				StopSequence:             nil,
				StopHeadSign:             nil,
				StartPickupDropOffWindow: nil,
				EndPickupDropOffWindow:   nil,
				PickupType:               nil,
				DropOffType:              nil,
				ContinuousPickup:         nil,
				ContinuousDropOff:        nil,
				ShapeDistTraveled:        nil,
				Timepoint:                nil,
				PickupBookingRuleId:      nil,
				DropOffBookingRuleId:     nil,
				LineNumber:               0,
			}},
		},
		"OK": {
			headers: headerMap,
			rows: [][]string{
				{"1", "00:10:00", "00:11:00", "0001", "LG1", "L1", "1", "Headsign", "00:11:00", "00:12:00", "1", "2", "2", "2", "100.00", "1", "pickup booking rule", "drop-off booking rule"},
			},
			expected: []*StopTime{{
				TripId:                   stringPtr("1"),
				ArrivalTime:              stringPtr("00:10:00"),
				DepartureTime:            stringPtr("00:11:00"),
				StopId:                   stringPtr("0001"),
				LocationGroupId:          stringPtr("LG1"),
				LocationId:               stringPtr("L1"),
				StopSequence:             stringPtr("1"),
				StopHeadSign:             stringPtr("Headsign"),
				StartPickupDropOffWindow: stringPtr("00:11:00"),
				EndPickupDropOffWindow:   stringPtr("00:12:00"),
				PickupType:               stringPtr("1"),
				DropOffType:              stringPtr("2"),
				ContinuousPickup:         stringPtr("2"),
				ContinuousDropOff:        stringPtr("2"),
				ShapeDistTraveled:        stringPtr("100.00"),
				Timepoint:                stringPtr("1"),
				PickupBookingRuleId:      stringPtr("pickup booking rule"),
				DropOffBookingRuleId:     stringPtr("drop-off booking rule"),
				LineNumber:               0,
			}},
		},
	}

	for name, tt := range tests {
		t.Run(fmt.Sprintf("%s", name), func(t *testing.T) {
			var actual []*StopTime
			for i, row := range tt.rows {
				actual = append(actual, CreateStopTime(row, tt.headers, i))
			}
			handleEntityCreateResults(t, tt.expected, actual)
		})
	}
}

func TestValidateStopTimes(t *testing.T) {
	tests := map[string]struct {
		actualEntities  []*StopTime
		expectedResults []Result
		stops           []*Stop
	}{
		"nil-slice": {
			actualEntities:  nil,
			expectedResults: []Result{},
		},
		"nil-slice-items": {
			actualEntities:  []*StopTime{nil},
			expectedResults: []Result{},
		},
		"invalid-fields": {
			actualEntities: []*StopTime{
				{
					TripId:                   stringPtr("1"),
					ArrivalTime:              stringPtr("Not a time"),
					DepartureTime:            stringPtr("Not a time"),
					StopId:                   stringPtr("0001"),
					LocationGroupId:          stringPtr("LG1"),
					LocationId:               stringPtr("L1"),
					StopSequence:             stringPtr("Not an integer"),
					StopHeadSign:             stringPtr("Headsign"),
					StartPickupDropOffWindow: stringPtr("Not a time"),
					EndPickupDropOffWindow:   stringPtr("Not a time"),
					PickupType:               stringPtr("4"),
					DropOffType:              stringPtr("5"),
					ContinuousPickup:         stringPtr("6"),
					ContinuousDropOff:        stringPtr("7"),
					ShapeDistTraveled:        stringPtr("Not a float"),
					Timepoint:                stringPtr("3"),
					PickupBookingRuleId:      stringPtr("pickup booking rule"),
					DropOffBookingRuleId:     stringPtr("drop-off booking rule"),
				},
			},
			stops: []*Stop{{Id: stringPtr("0001")}},
			expectedResults: []Result{
				InvalidTimeResult{SingleLineResult{FileName: "stop_times.txt", FieldName: "arrival_time"}},
				InvalidTimeResult{SingleLineResult{FileName: "stop_times.txt", FieldName: "departure_time"}},

				InvalidIntegerResult{SingleLineResult{FileName: "stop_times.txt", FieldName: "stop_sequence"}},
				InvalidTimeResult{SingleLineResult{FileName: "stop_times.txt", FieldName: "start_pickup_drop_off_window"}},
				InvalidTimeResult{SingleLineResult{FileName: "stop_times.txt", FieldName: "end_pickup_drop_off_window"}},

				InvalidPickupTypeResult{SingleLineResult{FileName: "stop_times.txt", FieldName: "pickup_type"}},
				InvalidDropOffTypeResult{SingleLineResult{FileName: "stop_times.txt", FieldName: "drop_off_type"}},

				InvalidContinuousPickupResult{SingleLineResult{FileName: "stop_times.txt", FieldName: "continuous_pickup"}},
				InvalidContinuousDropOffResult{SingleLineResult{FileName: "stop_times.txt", FieldName: "continuous_drop_off"}},
				InvalidFloatResult{SingleLineResult{FileName: "stop_times.txt", FieldName: "shape_dist_traveled"}},
				InvalidTimepointResult{SingleLineResult{FileName: "stop_times.txt", FieldName: "timepoint"}},
			},
		},
		"missing-stop": {
			actualEntities: []*StopTime{
				{
					StopId:       stringPtr("1000"),
					TripId:       stringPtr("1"),
					StopSequence: stringPtr("1"),
				},
			},
			stops: []*Stop{
				{Id: stringPtr("0002")},
				nil, // For coverage
			},
			expectedResults: []Result{
				ForeignKeyViolationResult{
					ReferencingFileName:  "stop_times.txt",
					ReferencingFieldName: "stop_id",
					ReferencedFieldName:  "stops.txt",
					ReferencedFileName:   "stop_id",
					OffendingValue:       "1000",
					ReferencedAtRow:      0,
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(fmt.Sprintf("%s", name), func(t *testing.T) {
			handleValidationResults(t, ValidateStopTimes(tt.actualEntities, tt.stops), tt.expectedResults)
		})
	}
}

//var validStopTimeHeaders = []string{"trip_id", "arrival_time", "departure_time", "stop_id", "stop_sequence",
//	"stop_headsign", "pickup_type", "drop_off_type", "continuous_pickup", "continuous_drop_off",
//	"shape_dist_traveled", "timepoint"}
//
//func TestShouldReturnEmptyStopTimeArrayOnEmptyString(t *testing.T) {
//	stopTimes, errors := LoadEntitiesFromCSV[*StopTime](csv.NewReader(strings.NewReader("")), validStopTimeHeaders, CreateStopTime, StopTimesFileName)
//	if len(errors) > 0 {
//		t.Error(errors)
//	}
//	if len(stopTimes) != 0 {
//		t.Error("expected zero calendar items")
//	}
//}
//
//func TestStopTimeParsing(t *testing.T) {
//	loadStopTimesFunc := func(reader *csv.Reader) ([]interface{}, []error) {
//		stopTimes, errs := LoadEntitiesFromCSV[*StopTime](reader, validStopTimeHeaders, CreateStopTime, StopTimesFileName)
//		entities := make([]interface{}, len(stopTimes))
//		for i, stopTime := range stopTimes {
//			entities[i] = stopTime
//		}
//		return entities, errs
//	}
//
//	validateStopTimesFunc := func(entities []interface{}, fixtures map[string][]interface{}) ([]error, []string) {
//		stopTimes := make([]*StopTime, len(entities))
//		for i, entity := range entities {
//			if stopTime, ok := entity.(*StopTime); ok {
//				stopTimes[i] = stopTime
//			}
//		}
//
//		stoPCount := len(fixtures["stops"])
//
//		if stoPCount == 0 {
//			return ValidateStopTimes(stopTimes, nil)
//		}
//
//		stops := make([]*Stop, stoPCount)
//		for i, fixture := range fixtures["stops"] {
//			if stop, ok := fixture.(*Stop); ok {
//				stops[i] = stop
//			} else {
//				t.Error(fmt.Sprintf("test setup error: cannot convert %v to Stop pointer. maybe you used value instead of pointer when setting fixtures", fixture))
//			}
//		}
//
//		return ValidateStopTimes(stopTimes, stops)
//	}
//
//	runGenericGTFSParseTest(t, "StopTimeOKTestcases", loadStopTimesFunc, validateStopTimesFunc, false, getStopTimeOKTestcases())
//	runGenericGTFSParseTest(t, "StopTimeNOKTestcases", loadStopTimesFunc, validateStopTimesFunc, false, getStopTimeNOKTestcases())
//}
//
//func getStopTimeOKTestcases() map[string]ggtfsTestCase {
//	expected1 := StopTime{
//		TripId:            NewID(stringPtr("1")),
//		ArrivalTime:       NewTime(stringPtr("00:10:00")),
//		DepartureTime:     NewTime(stringPtr("00:11:00")),
//		StopId:            NewID(stringPtr("0001")),
//		StopSequence:      NewInteger(stringPtr("1")),
//		StopHeadSign:      NewText(stringPtr("Foo city")),
//		PickupType:        NewPickupType(stringPtr("1")),
//		DropOffType:       NewDropOffType(stringPtr("0")),
//		ContinuousPickup:  NewContinuousPickupType(stringPtr("2")),
//		ContinuousDropOff: NewContinuousDropOffType(stringPtr("3")),
//		ShapeDistTraveled: NewFloat(stringPtr("100")),
//		Timepoint:         NewTimePoint(stringPtr("1")),
//		LineNumber:        2,
//	}
//
//	expected2 := StopTime{
//		TripId:            NewID(stringPtr("1")),
//		ArrivalTime:       NewTime(stringPtr("00:20:00")),
//		DepartureTime:     NewTime(stringPtr("00:21:00")),
//		StopId:            NewID(stringPtr("0002")),
//		StopSequence:      NewInteger(stringPtr("2")),
//		StopHeadSign:      NewText(stringPtr("Bar city")),
//		PickupType:        NewPickupType(stringPtr("1")),
//		DropOffType:       NewDropOffType(stringPtr("0")),
//		ContinuousPickup:  NewContinuousPickupType(stringPtr("2")),
//		ContinuousDropOff: NewContinuousDropOffType(stringPtr("3")),
//		ShapeDistTraveled: NewFloat(stringPtr("100")),
//		Timepoint:         NewTimePoint(stringPtr("1")),
//		LineNumber:        3,
//	}
//
//	testCases := make(map[string]ggtfsTestCase)
//	testCases["1"] = ggtfsTestCase{
//		csvRows: [][]string{
//			{"trip_id", "arrival_time", "departure_time", "stop_id", "stop_sequence", "stop_headsign", "pickup_type", "drop_off_type", "continuous_pickup", "continuous_drop_off", "shape_dist_traveled", "timepoint"},
//			{"1", "00:10:00", "00:11:00", "0001", "1", "Foo city", "1", "0", "2", "3", "100", "1"},
//			{"1", "00:20:00", "00:21:00", "0002", "2", "Bar city", "1", "0", "2", "3", "100", "1"},
//		},
//		expectedErrors:  []string{},
//		expectedStructs: []interface{}{&expected1, &expected2},
//		fixtures: map[string][]interface{}{
//			"stops": {
//				&Stop{
//					Id: stringPtr("0001"),
//				},
//				&Stop{
//					Id: stringPtr("0002"),
//				},
//			},
//		},
//	}
//
//	return testCases
//}
//
//func getStopTimeNOKTestcases() map[string]ggtfsTestCase {
//	testCases := make(map[string]ggtfsTestCase)
//	testCases["invalid-fields-must-error-out"] = ggtfsTestCase{
//		csvRows: [][]string{
//			{"trip_id", "arrival_time", "departure_time", "stop_id", "stop_sequence",
//				"stop_headsign", "pickup_type", "drop_off_type", "continuous_pickup", "continuous_drop_off",
//				"shape_dist_traveled", "timepoint"},
//			{"", "", "", "", "", "", "", "", "", "", "", ""},
//			{"ID1", "00:01:00", "50:00:00", "STOP1", "1", "STOP 1", "1", "1", "1", "1", "100", "1"},
//		},
//		expectedErrors: []string{
//			"stop_times.txt:2: invalid mandatory field: stop_sequence",
//			"stop_times.txt:2: invalid mandatory field: trip_id",
//			"stop_times.txt:2: stop_id () references to an unknown stop",
//			"stop_times.txt:3: invalid field: departure_time",
//			"stop_times.txt:3: stop_id (STOP1) references to an unknown stop",
//		},
//	}
//
//	return testCases
//}
