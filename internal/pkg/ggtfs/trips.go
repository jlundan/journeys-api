package ggtfs

import (
	"fmt"
	"strconv"
)

type Trip struct {
	RouteId              ID                   // route_id 				(required)
	ServiceId            ID                   // service_id 			(required)
	Id                   ID                   // trip_id 				(required)
	HeadSign             Text                 // trip_headsign 			(optional)
	ShortName            Text                 // trip_short_name 		(optional)
	DirectionId          DirectionId          // direction_id 			(optional)
	BlockId              ID                   // block_id 				(optional)
	ShapeId              ID                   // shape_id 				(conditionally required)
	WheelchairAccessible WheelchairAccessible // wheelchair_accessible 	(optional)
	BikesAllowed         BikesAllowed         // bikes_allowed 			(optional)
	LineNumber           int
}

func (t Trip) Validate() []error {
	var validationErrors []error

	requiredFields := []struct {
		fieldName string
		field     ValidAndPresentField
	}{
		{"route_id", &t.RouteId},
		{"service_id", &t.ServiceId},
		{"trip_id", &t.Id},
	}
	for _, f := range requiredFields {
		if !f.field.IsValid() {
			validationErrors = append(validationErrors, createFileRowError(TripsFileName, t.LineNumber, createInvalidRequiredFieldString(f.fieldName)))
		}
	}

	optionalFields := []struct {
		field     ValidAndPresentField
		fieldName string
	}{
		{&t.HeadSign, "trip_headsign"},
		{&t.ShortName, "trip_short_name"},
		{&t.DirectionId, "direction_id"},
		{&t.BlockId, "block_id"},
		{&t.ShapeId, "shape_id"},
		{&t.WheelchairAccessible, "wheelchair_accessible"},
		{&t.BikesAllowed, "bikes_allowed"},
	}

	for _, field := range optionalFields {
		if field.field != nil && field.field.IsPresent() && !field.field.IsValid() {
			validationErrors = append(validationErrors, createFileRowError(TripsFileName, t.LineNumber, createInvalidFieldString(field.fieldName)))
		}
	}

	return validationErrors
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
			trip.Id = NewID(v)
		case "route_id":
			trip.RouteId = NewID(v)
		case "service_id":
			trip.ServiceId = NewID(v)
		case "trip_headsign":
			trip.HeadSign = NewText(v)
		case "trip_short_name":
			trip.ShortName = NewText(v)
		case "direction_id":
			trip.DirectionId = NewDirectionId(v)
		case "block_id":
			trip.BlockId = NewID(v)
		case "shape_id":
			trip.ShapeId = NewID(v)
		case "wheelchair_accessible":
			trip.WheelchairAccessible = NewWheelchairAccessible(v)
		case "bikes_allowed":
			trip.BikesAllowed = NewBikesAllowed(v)
		}
	}

	if len(parseErrors) > 0 {
		return &trip
	}
	return &trip
}

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
			if trip.ShapeId.IsValid() {
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

	// TODO: VALIDATION: route_id: Foreign ID referencing routes.route_id (must refer to an existing route)
	// TODO: VALIDATION: service_id: Foreign ID referencing calendar.service_id or calendar_dates.service_id. (calendar_dates service_id relation is not checked)
	// TODO: VALIDATION: shape_id: Foreign ID referencing shapes.shape_id (must refer to an existing shape)
	// TODO: VALIDATION: shape_id: Required if the trip has a continuous pickup or drop-off behavior defined either in routes.txt or in stop_times.txt. Otherwise, optional.

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

func NewBikesAllowed(raw *string) BikesAllowed {
	if raw == nil {
		return BikesAllowed{
			Integer{base: base{raw: ""}}}
	}
	return BikesAllowed{Integer{base: base{raw: *raw, isPresent: true}}}
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

func NewWheelchairAccessible(raw *string) WheelchairAccessible {
	if raw == nil {
		return WheelchairAccessible{
			Integer{base: base{raw: ""}}}
	}
	return WheelchairAccessible{Integer{base: base{raw: *raw, isPresent: true}}}
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

func NewDirectionId(raw *string) DirectionId {
	if raw == nil {
		return DirectionId{
			Integer{base: base{raw: ""}}}
	}
	return DirectionId{Integer{base: base{raw: *raw, isPresent: true}}}
}
