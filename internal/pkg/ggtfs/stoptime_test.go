//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"fmt"
	"testing"
)

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
		expectedResults []ValidationNotice
		stops           []*Stop
	}{
		"nil-slice": {
			actualEntities:  nil,
			expectedResults: []ValidationNotice{},
		},
		"nil-slice-items": {
			actualEntities:  []*StopTime{nil},
			expectedResults: []ValidationNotice{},
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
			expectedResults: []ValidationNotice{
				InvalidTimeNotice{SingleLineNotice{FileName: "stop_times.txt", FieldName: "arrival_time"}},
				InvalidTimeNotice{SingleLineNotice{FileName: "stop_times.txt", FieldName: "departure_time"}},

				InvalidIntegerNotice{SingleLineNotice{FileName: "stop_times.txt", FieldName: "stop_sequence"}},
				InvalidTimeNotice{SingleLineNotice{FileName: "stop_times.txt", FieldName: "start_pickup_drop_off_window"}},
				InvalidTimeNotice{SingleLineNotice{FileName: "stop_times.txt", FieldName: "end_pickup_drop_off_window"}},

				InvalidPickupTypeNotice{SingleLineNotice{FileName: "stop_times.txt", FieldName: "pickup_type"}},
				InvalidDropOffTypeNotice{SingleLineNotice{FileName: "stop_times.txt", FieldName: "drop_off_type"}},

				InvalidContinuousPickupNotice{SingleLineNotice{FileName: "stop_times.txt", FieldName: "continuous_pickup"}},
				InvalidContinuousDropOffNotice{SingleLineNotice{FileName: "stop_times.txt", FieldName: "continuous_drop_off"}},
				InvalidFloatNotice{SingleLineNotice{FileName: "stop_times.txt", FieldName: "shape_dist_traveled"}},
				InvalidTimepointNotice{SingleLineNotice{FileName: "stop_times.txt", FieldName: "timepoint"}},
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
			expectedResults: []ValidationNotice{
				ForeignKeyViolationNotice{
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
