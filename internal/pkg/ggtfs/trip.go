package ggtfs

type Trip struct {
	RouteId              *string // route_id                (required)
	ServiceId            *string // service_id              (required)
	Id                   *string // trip_id                 (required)
	HeadSign             *string // trip_headsign           (optional)
	ShortName            *string // trip_short_name         (optional)
	DirectionId          *string // direction_id            (optional)
	BlockId              *string // block_id                (optional)
	ShapeId              *string // shape_id                (conditionally required)
	WheelchairAccessible *string // wheelchair_accessible   (optional)
	BikesAllowed         *string // bikes_allowed           (optional)
	LineNumber           int
}

func CreateTrip(row []string, headers map[string]int, lineNumber int) *Trip {
	var parseErrors []error

	trip := Trip{
		LineNumber: lineNumber,
	}

	for hName := range headers {
		v := getRowValueForHeaderName(row, headers, hName)

		switch hName {
		case "trip_id":
			trip.Id = v
		case "route_id":
			trip.RouteId = v
		case "service_id":
			trip.ServiceId = v
		case "trip_headsign":
			trip.HeadSign = v
		case "trip_short_name":
			trip.ShortName = v
		case "direction_id":
			trip.DirectionId = v
		case "block_id":
			trip.BlockId = v
		case "shape_id":
			trip.ShapeId = v
		case "wheelchair_accessible":
			trip.WheelchairAccessible = v
		case "bikes_allowed":
			trip.BikesAllowed = v
		}
	}

	if len(parseErrors) > 0 {
		return &trip
	}
	return &trip
}

func ValidateTrip(t Trip) []Result {
	var validationResults []Result

	fields := []struct {
		fieldType FieldType
		name      string
		value     *string
		required  bool
	}{
		{FieldTypeID, "trip_id", t.Id, true},
		{FieldTypeID, "route_id", t.RouteId, true},
		{FieldTypeID, "service_id", t.ServiceId, true},
		{FieldTypeText, "trip_headsign", t.HeadSign, false},
		{FieldTypeText, "trip_short_name", t.ShortName, false},
		{FieldTypeDirectionId, "direction_id", t.DirectionId, false},
		{FieldTypeID, "block_id", t.BlockId, false},
		{FieldTypeID, "shape_id", t.ShapeId, false},
		{FieldTypeWheelchairAccessible, "wheelchair_accessible", t.WheelchairAccessible, false},
		{FieldTypeBikesAllowed, "bikes_allowed", t.BikesAllowed, false},
	}

	for _, field := range fields {
		validationResults = append(validationResults, validateField(field.fieldType, field.name, field.value, field.required, FileNameTrips, t.LineNumber)...)
	}

	return validationResults
}

func ValidateTrips(trips []*Trip, routes []*Route, calendarItems []*CalendarItem, shapes []*Shape) []Result {
	var validationResults []Result

	if trips == nil {
		return validationResults
	}

	for _, trip := range trips {
		if trip == nil {
			continue
		}

		validationResults = append(validationResults, ValidateTrip(*trip)...)

		if routes != nil {
			routeFound := false
			for _, route := range routes {
				if route == nil {
					continue
				}
				if *trip.RouteId == *route.Id {
					routeFound = true
					break
				}
			}
			if !routeFound {
				validationResults = append(validationResults, ForeignKeyViolationResult{
					ReferencingFileName:  FileNameTrips,
					ReferencingFieldName: "route_id",
					ReferencedFieldName:  FileNameRoutes,
					ReferencedFileName:   "route_id",
					OffendingValue:       *trip.RouteId,
					ReferencedAtRow:      trip.LineNumber,
				})
			}
		}

		if calendarItems != nil {
			serviceFound := false
			for _, calendarItem := range calendarItems {
				if calendarItem == nil {
					continue
				}
				if *trip.ServiceId == *calendarItem.ServiceId {
					serviceFound = true
					break
				}
			}
			if !serviceFound {
				validationResults = append(validationResults, ForeignKeyViolationResult{
					ReferencingFileName:  FileNameTrips,
					ReferencingFieldName: "service_id",
					ReferencedFieldName:  FileNameCalendar,
					ReferencedFileName:   "service_id",
					OffendingValue:       *trip.ServiceId,
					ReferencedAtRow:      trip.LineNumber,
				})
			}
		}

		if shapes != nil {
			if !StringIsNilOrEmpty(trip.ShapeId) {
				shapeFound := false
				for _, shape := range shapes {
					if shape == nil {
						continue
					}
					// TODO nil check
					if *trip.ShapeId == *shape.Id {
						shapeFound = true
						break
					}
				}
				if !shapeFound {
					validationResults = append(validationResults, ForeignKeyViolationResult{
						ReferencingFileName:  FileNameTrips,
						ReferencingFieldName: "shape_id",
						ReferencedFieldName:  FileNameShapes,
						ReferencedFileName:   "shape_id",
						OffendingValue:       *trip.ShapeId,
						ReferencedAtRow:      trip.LineNumber,
					})
				}
			}
		}
	}

	return validationResults
}

//const (
//	BikesAllowedNoInformation     = 0
//	BikesAllowedAtLeastOneBicycle = 1
//	BikesAllowedNone              = 2
//)
//
//type BikesAllowed struct {
//	Integer
//}
//
//func (ba BikesAllowed) IsValid() bool {
//	val, err := strconv.Atoi(ba.Integer.base.raw)
//	if err != nil {
//		return false
//	}
//
//	return val == BikesAllowedNoInformation || val == BikesAllowedAtLeastOneBicycle ||
//		val == BikesAllowedNone
//}
//
//func NewBikesAllowed(raw *string) BikesAllowed {
//	if raw == nil {
//		return BikesAllowed{
//			Integer{base: base{raw: ""}}}
//	}
//	return BikesAllowed{Integer{base: base{raw: *raw, isPresent: true}}}
//}
//
//const (
//	WheelchairAccessibleNoInformation   = 0
//	WheelchairAccessibleAtLeastOneRider = 1
//	WheelchairAccessibleNoAccommodation = 2
//)
//
//type WheelchairAccessible struct {
//	Integer
//}
//
//func (wa WheelchairAccessible) IsValid() bool {
//	val, err := strconv.Atoi(wa.Integer.base.raw)
//	if err != nil {
//		return false
//	}
//
//	return val == WheelchairAccessibleNoInformation || val == WheelchairAccessibleAtLeastOneRider ||
//		val == WheelchairAccessibleNoAccommodation
//}
//
//func NewWheelchairAccessible(raw *string) WheelchairAccessible {
//	if raw == nil {
//		return WheelchairAccessible{
//			Integer{base: base{raw: ""}}}
//	}
//	return WheelchairAccessible{Integer{base: base{raw: *raw, isPresent: true}}}
//}
//
//const (
//	TripTravelInOneDirection      = 0
//	TripTravelInOppositeDirection = 1
//)
//
//type DirectionId struct {
//	Integer
//}
//
//func (di DirectionId) IsValid() bool {
//	val, err := strconv.Atoi(di.Integer.base.raw)
//	if err != nil {
//		return false
//	}
//
//	return val == TripTravelInOneDirection || val == TripTravelInOppositeDirection
//}
//
//func NewDirectionId(raw *string) DirectionId {
//	if raw == nil {
//		return DirectionId{
//			Integer{base: base{raw: ""}}}
//	}
//	return DirectionId{Integer{base: base{raw: *raw, isPresent: true}}}
//}
