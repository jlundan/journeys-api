package ggtfs

type StopTime struct {
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
	LineNumber               int
}

func CreateStopTime(row []string, headers map[string]int, lineNumber int) *StopTime {
	stopTime := StopTime{
		LineNumber: lineNumber,
	}

	for hName := range headers {
		v := getRowValueForHeaderName(row, headers, hName)

		switch hName {
		case "trip_id":
			stopTime.TripId = v
		case "arrival_time":
			stopTime.ArrivalTime = v
		case "departure_time":
			stopTime.DepartureTime = v
		case "stop_id":
			stopTime.StopId = v
		case "location_group_id":
			stopTime.LocationGroupId = v
		case "location_id":
			stopTime.LocationId = v
		case "stop_sequence":
			stopTime.StopSequence = v
		case "stop_headsign":
			stopTime.StopHeadSign = v
		case "start_pickup_drop_off_window":
			stopTime.StartPickupDropOffWindow = v
		case "end_pickup_drop_off_window":
			stopTime.EndPickupDropOffWindow = v
		case "pickup_type":
			stopTime.PickupType = v
		case "drop_off_type":
			stopTime.DropOffType = v
		case "continuous_pickup":
			stopTime.ContinuousPickup = v
		case "continuous_drop_off":
			stopTime.ContinuousDropOff = v
		case "shape_dist_traveled":
			stopTime.ShapeDistTraveled = v
		case "timepoint":
			stopTime.Timepoint = v
		case "pickup_booking_rule_id":
			stopTime.PickupBookingRuleId = v
		case "drop_off_booking_rule_id":
			stopTime.DropOffBookingRuleId = v
		}

	}

	return &stopTime
}

func ValidateStopTime(st StopTime) []Result {
	var validationResults []Result

	fields := []struct {
		fieldType FieldType
		name      string
		value     *string
		required  bool
	}{
		{FieldTypeID, "trip_id", st.TripId, true},
		{FieldTypeTime, "arrival_time", st.ArrivalTime, false},
		{FieldTypeTime, "departure_time", st.DepartureTime, false},
		{FieldTypeID, "stop_id", st.StopId, false},
		{FieldTypeID, "location_group_id", st.LocationGroupId, false},
		{FieldTypeID, "location_id", st.LocationId, false},
		{FieldTypeInteger, "stop_sequence", st.StopSequence, true},
		{FieldTypeText, "stop_headsign", st.StopHeadSign, false},
		{FieldTypeTime, "start_pickup_drop_off_window", st.StartPickupDropOffWindow, false},
		{FieldTypeTime, "end_pickup_drop_off_window", st.EndPickupDropOffWindow, false},
		{FieldTypePickupType, "pickup_type", st.PickupType, false},
		{FieldTypeDropOffType, "drop_off_type", st.DropOffType, false},
		{FieldTypeContinuousPickup, "continuous_pickup", st.ContinuousPickup, false},
		{FieldTypeContinuousDropOff, "continuous_drop_off", st.ContinuousDropOff, false},
		{FieldTypeFloat, "shape_dist_traveled", st.ShapeDistTraveled, false},
		{FieldTypeTimepoint, "timepoint", st.Timepoint, false},
		{FieldTypeID, "pickup_booking_rule_id", st.PickupBookingRuleId, false},
		{FieldTypeID, "drop_off_booking_rule_id", st.DropOffBookingRuleId, false},
	}

	for _, field := range fields {
		validationResults = append(validationResults, validateField(field.fieldType, field.name, field.value, field.required, FileNameStopTimes, st.LineNumber)...)
	}

	return validationResults
}

func ValidateStopTimes(stopTimes []*StopTime, stops []*Stop) []Result {
	var validationResults []Result

	if stopTimes == nil {
		return validationResults
	}

	for _, stopTimeItem := range stopTimes {
		if stopTimeItem == nil {
			continue
		}

		vRes := ValidateStopTime(*stopTimeItem)
		if len(vRes) > 0 {
			validationResults = append(validationResults, vRes...)
		}

		stopFound := false
		if stops != nil {
			for _, stop := range stops {
				if stop == nil {
					continue
				}
				// TODO: nil check
				if *stopTimeItem.StopId == *stop.Id {
					stopFound = true
					break
				}
			}
		}
		if !stopFound {
			validationResults = append(validationResults, ForeignKeyViolationResult{
				ReferencingFileName:  FileNameStopTimes,
				ReferencingFieldName: "stop_id",
				ReferencedFieldName:  FileNameStops,
				ReferencedFileName:   "stop_id",
				OffendingValue:       *stopTimeItem.StopId,
				ReferencedAtRow:      stopTimeItem.LineNumber,
			})
		}
	}

	return validationResults
}

//const (
//	StopTimePickupTypeRegularlyScheduled       = 0
//	StopTimePickupTypeNoPickup                 = 1
//	StopTimePickupTypeMustPhoneAgency          = 2
//	StopTimePickupTypeMustCoordinateWithDriver = 3
//)
//
//type PickupType struct {
//	Integer
//}
//
//func (pt PickupType) IsValid() bool {
//	val, err := strconv.Atoi(pt.Integer.base.raw)
//	if err != nil {
//		return false
//	}
//
//	return val == StopTimePickupTypeRegularlyScheduled || val == StopTimePickupTypeNoPickup ||
//		val == StopTimePickupTypeMustPhoneAgency || val == StopTimePickupTypeMustCoordinateWithDriver
//}
//
//func NewPickupType(raw *string) PickupType {
//	if raw == nil {
//		return PickupType{
//			Integer{base: base{raw: ""}}}
//	}
//	return PickupType{Integer{base: base{raw: *raw, isPresent: true}}}
//}
//
//const (
//	DropOffTypeRegularlyScheduled        = 0
//	DropOffTypeNoDropOff                 = 1
//	DropOffTypeMustPhoneAgency           = 2
//	DropOffTypeRMustCoordinateWithDriver = 3
//)
//
//type DropOffType struct {
//	Integer
//}
//
//func (dot DropOffType) IsValid() bool {
//	val, err := strconv.Atoi(dot.Integer.base.raw)
//	if err != nil {
//		return false
//	}
//
//	return val == DropOffTypeRegularlyScheduled || val == DropOffTypeNoDropOff ||
//		val == DropOffTypeMustPhoneAgency || val == DropOffTypeRMustCoordinateWithDriver
//}
//
//func NewDropOffType(raw *string) DropOffType {
//	if raw == nil {
//		return DropOffType{
//			Integer{base: base{raw: ""}}}
//	}
//	return DropOffType{Integer{base: base{raw: *raw, isPresent: true}}}
//}
//
//const (
//	TimePointApproximate = 0
//	TimePointExact       = 1
//)
//
//type TimePoint struct {
//	Integer
//}
//
//func (dot TimePoint) IsValid() bool {
//	val, err := strconv.Atoi(dot.Integer.base.raw)
//	if err != nil {
//		return false
//	}
//
//	return val == TimePointApproximate || val == TimePointExact
//}
//
//func NewTimePoint(raw *string) TimePoint {
//	if raw == nil {
//		return TimePoint{
//			Integer{base: base{raw: ""}}}
//	}
//	return TimePoint{Integer{base: base{raw: *raw, isPresent: true}}}
//}
//
//const (
//	ContinuousStoppingPickupType       = 0
//	NoContinuousStoppingPickupType     = 1
//	MustPhoneAgencyPickupType          = 2
//	MustCoordinateWithDriverPickupType = 3
//)
//
//type ContinuousPickupType struct {
//	Integer
//}
//
//func (cpt ContinuousPickupType) IsValid() bool {
//	// Spec says values "1 or empty mean No continuous stopping drop off."
//	// Empty = valid
//	if cpt.Integer.base.IsEmpty() {
//		return true
//	}
//
//	val, err := strconv.Atoi(cpt.Integer.base.raw)
//
//	if err != nil {
//		return false
//	}
//
//	return val == ContinuousStoppingPickupType || val == NoContinuousStoppingPickupType ||
//		val == MustPhoneAgencyPickupType || val == MustCoordinateWithDriverPickupType
//}
//
//func NewContinuousPickupType(raw *string) ContinuousPickupType {
//	if raw == nil {
//		return ContinuousPickupType{
//			Integer{base: base{raw: ""}}}
//	}
//	return ContinuousPickupType{Integer{base: base{raw: *raw, isPresent: true}}}
//}
//
//const (
//	ContinuousStoppingDropOffType       = 0
//	NoContinuousStoppingDropOffType     = 1
//	MustPhoneAgencyDropOffType          = 2
//	MustCoordinateWithDriverDropOffType = 3
//)
//
//type ContinuousDropOffType struct {
//	Integer
//}
//
//func (cpt ContinuousDropOffType) IsValid() bool {
//	// Spec says "1 or empty - No continuous stopping drop off."
//	// Empty = valid
//	if cpt.Integer.base.IsEmpty() {
//		return true
//	}
//
//	val, err := strconv.Atoi(cpt.Integer.base.raw)
//	if err != nil {
//		return false
//	}
//
//	return val == ContinuousStoppingDropOffType || val == NoContinuousStoppingDropOffType ||
//		val == MustPhoneAgencyDropOffType || val == MustCoordinateWithDriverDropOffType
//}
//
//func NewContinuousDropOffType(raw *string) ContinuousDropOffType {
//	if raw == nil {
//		return ContinuousDropOffType{
//			Integer{base: base{raw: ""}}}
//	}
//	return ContinuousDropOffType{Integer{base: base{raw: *raw, isPresent: true}}}
//}
