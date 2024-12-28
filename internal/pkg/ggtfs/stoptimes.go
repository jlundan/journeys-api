package ggtfs

import (
	"fmt"
	"strconv"
)

type StopTime struct {
	TripId                   ID                    // trip_id 						(required)
	ArrivalTime              Time                  // arrival_time 					(conditionally required)
	DepartureTime            Time                  // departure_time 				(conditionally required)
	StopId                   ID                    // stop_id 						(conditionally required)
	LocationGroupId          ID                    // location_group_id 			(conditionally forbidden)
	LocationId               ID                    // location_id 					(conditionally forbidden)
	StopSequence             Integer               // stop_sequence 				(required)
	StopHeadSign             Text                  // stop_headsign 				(optional)
	StartPickupDropOffWindow Time                  // start_pickup_drop_off_window 	(conditionally required)
	EndPickupDropOffWindow   Time                  // end_pickup_drop_off_window 	(conditionally required)
	PickupType               PickupType            // pickup_type 					(conditionally required)
	DropOffType              DropOffType           // drop_off_type 				(conditionally required)
	ContinuousPickup         ContinuousPickupType  // continuous_pickup 			(conditionally required)
	ContinuousDropOff        ContinuousDropOffType // continuous_drop_off 			(conditionally required)
	ShapeDistTraveled        Float                 // shape_dist_traveled 			(optional)
	Timepoint                TimePoint             // timepoint 					(optional)
	PickupBookingRuleId      ID                    // pickup_booking_rule_id 		(optional)
	DropOffBookingRuleId     ID                    // drop_off_booking_rule_id 		(optional)
	LineNumber               int
}

func (st StopTime) Validate() []error {
	var validationErrors []error

	requiredFields := []struct {
		fieldName string
		field     ValidAndPresentField
	}{
		{"trip_id", &st.TripId},
		{"stop_sequence", &st.StopSequence},
	}
	for _, f := range requiredFields {
		if !f.field.IsValid() {
			validationErrors = append(validationErrors, createFileRowError(StopTimesFileName, st.LineNumber, createInvalidRequiredFieldString(f.fieldName)))
		}
	}

	optionalFields := []struct {
		field     ValidAndPresentField
		fieldName string
	}{
		{&st.ArrivalTime, "arrival_time"},
		{&st.DepartureTime, "departure_time"},
		{&st.StopId, "stop_id"},
		{&st.LocationGroupId, "location_group_id"},
		{&st.LocationId, "location_id"},
		{&st.StopHeadSign, "stop_headsign"},
		{&st.StartPickupDropOffWindow, "start_pickup_drop_off_window"},
		{&st.EndPickupDropOffWindow, "end_pickup_drop_off_window"},
		{&st.PickupType, "pickup_type"},
		{&st.DropOffType, "drop_off_type"},
		{&st.ContinuousPickup, "continuous_pickup"},
		{&st.ContinuousDropOff, "continuous_drop_off"},
		{&st.ShapeDistTraveled, "shape_dist_traveled"},
		{&st.Timepoint, "timepoint"},
		{&st.PickupBookingRuleId, "pickup_booking_rule_id"},
		{&st.DropOffBookingRuleId, "drop_off_booking_rule_id"},
	}

	for _, field := range optionalFields {
		if field.field != nil && field.field.IsPresent() && !field.field.IsValid() {
			validationErrors = append(validationErrors, createFileRowError(StopTimesFileName, st.LineNumber, createInvalidFieldString(field.fieldName)))
		}
	}

	return validationErrors
}

func CreateStopTime(row []string, headers map[string]int, lineNumber int) *StopTime {
	var parseErrors []error

	stopTime := StopTime{
		LineNumber: lineNumber,
	}

	for hName := range headers {
		v := getRowValueForHeaderName(row, headers, hName)

		switch hName {
		case "trip_id":
			stopTime.TripId = NewID(v)
		case "arrival_time":
			stopTime.ArrivalTime = NewTime(v)
		case "departure_time":
			stopTime.DepartureTime = NewTime(v)
		case "stop_id":
			stopTime.StopId = NewID(v)
		case "location_group_id":
			stopTime.LocationGroupId = NewID(v)
		case "location_id":
			stopTime.LocationId = NewID(v)
		case "stop_sequence":
			stopTime.StopSequence = NewInteger(v)
		case "stop_headsign":
			stopTime.StopHeadSign = NewText(v)
		case "start_pickup_drop_off_window":
			stopTime.StartPickupDropOffWindow = NewTime(v)
		case "end_pickup_drop_off_window":
			stopTime.EndPickupDropOffWindow = NewTime(v)
		case "pickup_type":
			stopTime.PickupType = NewPickupType(v)
		case "drop_off_type":
			stopTime.DropOffType = NewDropOffType(v)
		case "continuous_pickup":
			stopTime.ContinuousPickup = NewContinuousPickupType(v)
		case "continuous_drop_off":
			stopTime.ContinuousDropOff = NewContinuousDropOffType(v)
		case "shape_dist_traveled":
			stopTime.ShapeDistTraveled = NewFloat(v)
		case "timepoint":
			stopTime.Timepoint = NewTimePoint(v)
		case "pickup_booking_rule_id":
			stopTime.PickupBookingRuleId = NewID(v)
		case "drop_off_booking_rule_id":
			stopTime.DropOffBookingRuleId = NewID(v)
		}

	}

	if len(parseErrors) > 0 {
		return &stopTime
	}
	return &stopTime
}

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

		vErr := stopTimeItem.Validate()
		if len(vErr) > 0 {
			validationErrors = append(validationErrors, vErr...)
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

func NewPickupType(raw *string) PickupType {
	if raw == nil {
		return PickupType{
			Integer{base: base{raw: ""}}}
	}
	return PickupType{Integer{base: base{raw: *raw, isPresent: true}}}
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

func NewDropOffType(raw *string) DropOffType {
	if raw == nil {
		return DropOffType{
			Integer{base: base{raw: ""}}}
	}
	return DropOffType{Integer{base: base{raw: *raw, isPresent: true}}}
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

func NewTimePoint(raw *string) TimePoint {
	if raw == nil {
		return TimePoint{
			Integer{base: base{raw: ""}}}
	}
	return TimePoint{Integer{base: base{raw: *raw, isPresent: true}}}
}
