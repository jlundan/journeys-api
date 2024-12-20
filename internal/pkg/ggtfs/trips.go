package ggtfs

import (
	"fmt"
	"strconv"
)

// Trip struct with fields as strings and optional fields as string pointers.
type Trip struct {
	Id                   ID                    // trip_id
	RouteId              ID                    // route_id
	ServiceId            ID                    // service_id
	HeadSign             *Text                 // trip_headsign (optional)
	ShortName            *Text                 // trip_short_name (optional)
	DirectionId          *DirectionId          // direction_id (optional)
	BlockId              *ID                   // block_id (optional)
	ShapeId              *ID                   // shape_id (optional)
	WheelchairAccessible *WheelchairAccessible // wheelchair_accessible (optional)
	BikesAllowed         *BikesAllowed         // bikes_allowed (optional)
	LineNumber           int                   // Line number in the CSV file for error reporting
}

func (t Trip) Validate() []error {
	var validationErrors []error

	// shape_id is handled in the ValidateStopTimes function since it is conditionally required

	fields := []struct {
		fieldName string
		field     ValidAndPresentField
	}{
		{"route_id", &t.RouteId},
		{"service_id", &t.ServiceId},
		{"trip_id", &t.Id},
	}
	for _, f := range fields {
		validationErrors = append(validationErrors, validateFieldIsPresentAndValid(f.field, f.fieldName, t.LineNumber, TripsFileName)...)
	}

	// Checking the underlying value of the field in ValidAndPresentField for nil would require reflection
	// v := reflect.ValueOf(i)
	// v.Kind() == reflect.Ptr && v.IsNil()
	// which is slow, so we can't use the above mechanism to check optional fields, since they might be nil (pointer field's default value is nil)
	// since CreateTrip might have not processed the field (if its header is missing from the csv).

	if t.HeadSign != nil && !t.HeadSign.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(TripsFileName, t.LineNumber, createInvalidFieldString("trip_headsign")))
	}
	if t.ShortName != nil && !t.ShortName.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(TripsFileName, t.LineNumber, createInvalidFieldString("trip_short_name")))
	}
	if t.DirectionId != nil && !t.DirectionId.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(TripsFileName, t.LineNumber, createInvalidFieldString("direction_id")))
	}
	if t.BlockId != nil && !t.BlockId.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(TripsFileName, t.LineNumber, createInvalidFieldString("block_id")))
	}
	if t.BikesAllowed != nil && !t.BikesAllowed.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(TripsFileName, t.LineNumber, createInvalidFieldString("wheelchair_accessible")))
	}

	return validationErrors
}

// CreateTrip creates and validates a Trip instance from the CSV row data.
func CreateTrip(row []string, headers map[string]int, lineNumber int) *Trip {
	var parseErrors []error

	trip := Trip{
		LineNumber: lineNumber,
	}

	for hName, hPos := range headers {
		switch hName {
		case "trip_id":
			trip.Id = NewID(getRowValue(row, hPos))
		case "route_id":
			trip.RouteId = NewID(getRowValue(row, hPos))
		case "service_id":
			trip.ServiceId = NewID(getRowValue(row, hPos))
		case "trip_headsign":
			trip.HeadSign = NewOptionalText(getRowValue(row, hPos))
		case "trip_short_name":
			trip.ShortName = NewOptionalText(getRowValue(row, hPos))
		case "direction_id":
			trip.DirectionId = NewOptionalDirectionId(getRowValue(row, hPos))
		case "block_id":
			trip.BlockId = NewOptionalID(getRowValue(row, hPos))
		case "shape_id":
			trip.ShapeId = NewOptionalID(getRowValue(row, hPos))
		case "wheelchair_accessible":
			trip.WheelchairAccessible = NewOptionalWheelchairAccessible(getRowValue(row, hPos))
		case "bikes_allowed":
			trip.BikesAllowed = NewOptionalBikesAllowed(getRowValue(row, hPos))
		}
	}

	if len(parseErrors) > 0 {
		return &trip
	}
	return &trip
}

// ValidateTrips performs additional validation for a list of Trip instances.
func ValidateTrips(trips []*Trip, routes []*Route, calendarItems []*CalendarItem, shapes []*Shape) ([]error, []string) {
	var validationErrors []error
	var recommendations []string

	if trips == nil {
		return validationErrors, recommendations
	}

	for _, trip := range trips {
		if trip == nil {
			continue
		}

		validationErrors = append(validationErrors, trip.Validate()...)

		if routes != nil {
			routeFound := false
			for _, route := range routes {
				if route == nil {
					continue
				}
				if trip.RouteId.String() == route.Id.String() {
					routeFound = true
					break
				}
			}
			if !routeFound {
				validationErrors = append(validationErrors, createFileRowError(TripsFileName, trip.LineNumber, fmt.Sprintf("referenced route_id not found in %s", RoutesFileName)))
			}
		}

		if calendarItems != nil {
			serviceFound := false
			for _, calendarItem := range calendarItems {
				if calendarItem == nil {
					continue
				}
				if trip.ServiceId.String() == calendarItem.ServiceId.String() {
					serviceFound = true
					break
				}
			}
			if !serviceFound {
				validationErrors = append(validationErrors, createFileRowError(TripsFileName, trip.LineNumber, fmt.Sprintf("referenced service_id not found in %s", CalendarFileName)))
			}
		}

		if shapes != nil {
			if trip.ShapeId != nil {
				shapeFound := false
				for _, shape := range shapes {
					if shape == nil {
						continue
					}
					if trip.ShapeId.String() == shape.Id.String() {
						shapeFound = true
						break
					}
				}
				if !shapeFound {
					validationErrors = append(validationErrors, createFileRowError(TripsFileName, trip.LineNumber, fmt.Sprintf("referenced shape_id not found in %s", ShapesFileName)))
				}
			}
		}
	}

	return validationErrors, recommendations
}

const (
	BikesAllowedNoInformation     = 0
	BikesAllowedAtLeastOneBicycle = 1
	BikesAllowedNone              = 2
)

type BikesAllowed struct {
	Integer
}

func (ba BikesAllowed) IsValid() bool {
	val, err := strconv.Atoi(ba.Integer.base.raw)
	if err != nil {
		return false
	}

	return val == BikesAllowedNoInformation || val == BikesAllowedAtLeastOneBicycle ||
		val == BikesAllowedNone
}

func NewOptionalBikesAllowed(raw *string) *BikesAllowed {
	if raw == nil {
		return &BikesAllowed{
			Integer{base: base{raw: ""}}}
	}
	return &BikesAllowed{Integer{base: base{raw: *raw, isPresent: true}}}
}

const (
	WheelchairAccessibleNoInformation   = 0
	WheelchairAccessibleAtLeastOneRider = 1
	WheelchairAccessibleNoAccommodation = 2
)

type WheelchairAccessible struct {
	Integer
}

func (wa WheelchairAccessible) IsValid() bool {
	val, err := strconv.Atoi(wa.Integer.base.raw)
	if err != nil {
		return false
	}

	return val == WheelchairAccessibleNoInformation || val == WheelchairAccessibleAtLeastOneRider ||
		val == WheelchairAccessibleNoAccommodation
}

func NewOptionalWheelchairAccessible(raw *string) *WheelchairAccessible {
	if raw == nil {
		return &WheelchairAccessible{
			Integer{base: base{raw: ""}}}
	}
	return &WheelchairAccessible{Integer{base: base{raw: *raw, isPresent: true}}}
}

const (
	TripTravelInOneDirection      = 0
	TripTravelInOppositeDirection = 1
)

type DirectionId struct {
	Integer
}

func (di DirectionId) IsValid() bool {
	val, err := strconv.Atoi(di.Integer.base.raw)
	if err != nil {
		return false
	}

	return val == TripTravelInOneDirection || val == TripTravelInOppositeDirection
}

func NewOptionalDirectionId(raw *string) *DirectionId {
	if raw == nil {
		return &DirectionId{
			Integer{base: base{raw: ""}}}
	}
	return &DirectionId{Integer{base: base{raw: *raw, isPresent: true}}}
}
