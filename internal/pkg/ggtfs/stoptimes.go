package ggtfs

import (
	"encoding/csv"
	"fmt"
	"strconv"
)

// StopTime struct with fields as strings and optional fields as string pointers.
type StopTime struct {
	TripId                   ID                     // trip_id
	ArrivalTime              Time                   // arrival_time
	DepartureTime            Time                   // departure_time
	StopId                   ID                     // stop_id
	LocationGroupId          *ID                    // location_group_id (optional)
	LocationId               *ID                    // location_id (optional)
	StopSequence             Integer                // stop_sequence
	StopHeadSign             *Text                  // stop_headsign (optional)
	StartPickupDropOffWindow *Time                  // start_pickup_drop_off_window (optional)
	EndPickupDropOffWindow   *Time                  // end_pickup_drop_off_window (optional)
	PickupType               *PickupType            // pickup_type (optional)
	DropOffType              *DropOffType           // drop_off_type (optional)
	ContinuousPickup         *ContinuousPickupType  // continuous_pickup (optional)
	ContinuousDropOff        *ContinuousDropOffType // continuous_drop_off (optional)
	ShapeDistTraveled        *Float                 // shape_dist_traveled (optional)
	Timepoint                *TimePoint             // timepoint (optional)
	PickupBookingRuleId      *ID                    // pickup_booking_rule_id (optional)
	DropOffBookingRuleId     *ID                    // drop_off_booking_rule_id (optional)
	LineNumber               int
}

func (st StopTime) Validate() []error {
	var validationErrors []error

	// arrival_time is handled in the ValidateStopTimes function since it is conditionally required
	// departure_time is handled in the ValidateStopTimes function since it is conditionally required
	// stop_id is handled in the ValidateStopTimes function since it is conditionally required
	// location_group_id is handled in the ValidateStopTimes function since it is conditionally forbidden
	// location_id is handled in the ValidateStopTimes function since it is conditionally forbidden
	// start_pickup_drop_off_window is handled in the ValidateStopTimes function since it is conditionally required
	// end_pickup_drop_off_window is handled in the ValidateStopTimes function since it is conditionally required
	// pickup_type is handled in the ValidateStopTimes function since it is conditionally forbidden
	// drop_off_type is handled in the ValidateStopTimes function since it is conditionally forbidden
	// continuous_pickup is handled in the ValidateStopTimes function since it is conditionally forbidden
	// continuous_drop_off is handled in the ValidateStopTimes function since it is conditionally forbidden

	fields := []struct {
		fieldName string
		field     ValidAndPresentField
	}{
		{"trip_id", &st.TripId},
		{"arrival_time", &st.ArrivalTime},
		{"departure_time", &st.DepartureTime},
		{"stop_id", &st.StopId},
		{"stop_sequence", &st.StopSequence},
	}
	for _, f := range fields {
		validationErrors = append(validationErrors, validateFieldIsPresentAndValid(f.field, f.fieldName, st.LineNumber, StopTimesFileName)...)
	}

	if st.StopHeadSign != nil && !st.StopHeadSign.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(StopTimesFileName, st.LineNumber, createInvalidFieldString("stop_headsign")))
	}
	if st.ShapeDistTraveled != nil && !st.ShapeDistTraveled.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(StopTimesFileName, st.LineNumber, createInvalidFieldString("shape_dist_traveled")))
	}
	if st.Timepoint != nil && !st.Timepoint.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(StopTimesFileName, st.LineNumber, createInvalidFieldString("timepoint")))
	}
	if st.PickupBookingRuleId != nil && !st.PickupBookingRuleId.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(StopTimesFileName, st.LineNumber, createInvalidFieldString("pickup_booking_rule_id")))
	}
	if st.DropOffBookingRuleId != nil && !st.DropOffBookingRuleId.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(StopTimesFileName, st.LineNumber, createInvalidFieldString("drop_off_booking_rule_id")))
	}

	return validationErrors
}

// ValidStopTimeHeaders defines the headers expected in the stop_times CSV file.
var validStopTimeHeaders = []string{"trip_id", "arrival_time", "departure_time", "stop_id", "stop_sequence",
	"stop_headsign", "pickup_type", "drop_off_type", "continuous_pickup", "continuous_drop_off",
	"shape_dist_traveled", "timepoint"}

// LoadStopTimes loads stop times from a CSV reader and returns them along with any errors.
func LoadStopTimes(csvReader *csv.Reader) ([]*StopTime, []error) {
	entities, errs := loadEntities(csvReader, validStopTimeHeaders, CreateStopTime, StopTimesFileName)

	stopTimes := make([]*StopTime, 0, len(entities))
	for _, entity := range entities {
		if stopTime, ok := entity.(*StopTime); ok {
			stopTimes = append(stopTimes, stopTime)
		}
	}

	return stopTimes, errs
}

// CreateStopTime creates and validates a StopTime instance from the CSV row data.
func CreateStopTime(row []string, headers map[string]int, lineNumber int) interface{} {
	var parseErrors []error

	stopTime := StopTime{
		LineNumber: lineNumber,
	}

	for hName, hPos := range headers {
		switch hName {
		case "trip_id":
			stopTime.TripId = NewID(getRowValue(row, hPos))
		case "arrival_time":
			stopTime.ArrivalTime = NewTime(getRowValue(row, hPos))
		case "departure_time":
			stopTime.DepartureTime = NewTime(getRowValue(row, hPos))
		case "stop_id":
			stopTime.StopId = NewID(getRowValue(row, hPos))
		case "location_group_id":
			stopTime.LocationGroupId = NewOptionalID(getRowValue(row, hPos))
		case "location_id":
			stopTime.LocationId = NewOptionalID(getRowValue(row, hPos))
		case "stop_sequence":
			stopTime.StopSequence = NewInteger(getRowValue(row, hPos))
		case "stop_headsign":
			stopTime.StopHeadSign = NewOptionalText(getRowValue(row, hPos))
		case "start_pickup_drop_off_window":
			stopTime.StartPickupDropOffWindow = NewOptionalTime(getRowValue(row, hPos))
		case "end_pickup_drop_off_window":
			stopTime.EndPickupDropOffWindow = NewOptionalTime(getRowValue(row, hPos))
		case "pickup_type":
			stopTime.PickupType = NewOptionalPickupType(getRowValue(row, hPos))
		case "drop_off_type":
			stopTime.DropOffType = NewOptionalDropOffType(getRowValue(row, hPos))
		case "continuous_pickup":
			stopTime.ContinuousPickup = NewOptionalContinuousPickupType(getRowValue(row, hPos))
		case "continuous_drop_off":
			stopTime.ContinuousDropOff = NewOptionalContinuousDropOffType(getRowValue(row, hPos))
		case "shape_dist_traveled":
			stopTime.ShapeDistTraveled = NewOptionalFloat(getRowValue(row, hPos))
		case "timepoint":
			stopTime.Timepoint = NewOptionalTimePoint(getRowValue(row, hPos))
		case "pickup_booking_rule_id":
			stopTime.PickupBookingRuleId = NewOptionalID(getRowValue(row, hPos))
		case "drop_off_booking_rule_id":
			stopTime.DropOffBookingRuleId = NewOptionalID(getRowValue(row, hPos))
		}

	}

	if len(parseErrors) > 0 {
		return &stopTime
	}
	return &stopTime
}

// ValidateStopTimes performs additional validation for a list of StopTime instances.
func ValidateStopTimes(stopTimes []*StopTime, stops []*Stop) ([]error, []string) {
	var validationErrors []error
	var recommendations []string

	if stopTimes == nil {
		return validationErrors, recommendations
	}

	// Group stop times by trip_id to ensure each trip has at least two stops.
	tripIdToStopCount := make(map[string]int)
	for _, stopTimeItem := range stopTimes {
		if stopTimeItem == nil {
			continue
		}
		tripIdToStopCount[stopTimeItem.TripId.String()] = 0
	}

	for _, stopTimeItem := range stopTimes {
		if stopTimeItem == nil {
			continue
		}
		stopFound := false
		if stops != nil {
			for _, stop := range stops {
				if stop == nil {
					continue
				}
				if stopTimeItem.StopId.String() == stop.Id.String() {
					stopFound = true
					break
				}
			}
		}
		if !stopFound {
			validationErrors = append(validationErrors, createFileError(StopTimesFileName, fmt.Sprintf("trip (%v) references to an unknown stop_id (%s)", stopTimeItem.TripId.String(), stopTimeItem.StopId.String())))
		} else {
			tripIdToStopCount[stopTimeItem.TripId.String()]++
		}
	}

	for tripId, stopCount := range tripIdToStopCount {
		if stopCount < 2 {
			validationErrors = append(validationErrors, createFileError(StopTimesFileName, fmt.Sprintf("trip (%v) has less than two defined stop times", tripId)))
		}
	}

	return validationErrors, recommendations
}

const (
	StopTimePickupTypeRegularlyScheduled       = 0
	StopTimePickupTypeNoPickup                 = 1
	StopTimePickupTypeMustPhoneAgency          = 2
	StopTimePickupTypeMustCoordinateWithDriver = 3
)

type PickupType struct {
	Integer
}

func (pt PickupType) IsValid() bool {
	val, err := strconv.Atoi(pt.Integer.base.raw)
	if err != nil {
		return false
	}

	return val == StopTimePickupTypeRegularlyScheduled || val == StopTimePickupTypeNoPickup ||
		val == StopTimePickupTypeMustPhoneAgency || val == StopTimePickupTypeMustCoordinateWithDriver
}

func NewOptionalPickupType(raw *string) *PickupType {
	if raw == nil {
		return &PickupType{
			Integer{base: base{raw: ""}}}
	}
	return &PickupType{Integer{base: base{raw: *raw}}}
}

const (
	DropOffTypeRegularlyScheduled        = 0
	DropOffTypeNoDropOff                 = 1
	DropOffTypeMustPhoneAgency           = 2
	DropOffTypeRMustCoordinateWithDriver = 3
)

type DropOffType struct {
	Integer
}

func (dot DropOffType) IsValid() bool {
	val, err := strconv.Atoi(dot.Integer.base.raw)
	if err != nil {
		return false
	}

	return val == DropOffTypeRegularlyScheduled || val == DropOffTypeNoDropOff ||
		val == DropOffTypeMustPhoneAgency || val == DropOffTypeRMustCoordinateWithDriver
}

func NewOptionalDropOffType(raw *string) *DropOffType {
	if raw == nil {
		return &DropOffType{
			Integer{base: base{raw: ""}}}
	}
	return &DropOffType{Integer{base: base{raw: *raw}}}
}

const (
	TimePointApproximate = 0
	TimePointExact       = 1
)

type TimePoint struct {
	Integer
}

func (dot TimePoint) IsValid() bool {
	val, err := strconv.Atoi(dot.Integer.base.raw)
	if err != nil {
		return false
	}

	return val == TimePointApproximate || val == TimePointExact
}

func NewOptionalTimePoint(raw *string) *TimePoint {
	if raw == nil {
		return &TimePoint{
			Integer{base: base{raw: ""}}}
	}
	return &TimePoint{Integer{base: base{raw: *raw}}}
}
