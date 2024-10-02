package ggtfs

import (
	"encoding/csv"
	"fmt"
)

// Trip struct with fields as strings and optional fields as string pointers.
type Trip struct {
	Id                   string  // trip_id
	RouteId              string  // route_id
	ServiceId            string  // service_id
	HeadSign             *string // trip_headsign (optional)
	ShortName            *string // trip_short_name (optional)
	DirectionId          *string // direction_id (optional)
	BlockId              *string // block_id (optional)
	ShapeId              *string // shape_id (optional)
	WheelchairAccessible *string // wheelchair_accessible (optional)
	BikesAllowed         *string // bikes_allowed (optional)
	LineNumber           int     // Line number in the CSV file for error reporting
}

// ValidTripHeaders defines the headers expected in the trips CSV file.
var validTripHeaders = []string{"route_id", "service_id", "trip_id", "trip_headsign", "trip_short_name",
	"direction_id", "block_id", "shape_id", "wheelchair_accessible", "bikes_allowed"}

// LoadTrips loads trips from a CSV reader and returns them along with any errors.
func LoadTrips(csvReader *csv.Reader) ([]*Trip, []error) {
	entities, errs := loadEntities(csvReader, validTripHeaders, CreateTrip, TripsFileName)

	trips := make([]*Trip, 0, len(entities))
	for _, entity := range entities {
		if trip, ok := entity.(*Trip); ok {
			trips = append(trips, trip)
		}
	}

	return trips, errs
}

// CreateTrip creates and validates a Trip instance from the CSV row data.
func CreateTrip(row []string, headers map[string]int, lineNumber int) (interface{}, []error) {
	var parseErrors []error

	trip := Trip{
		LineNumber: lineNumber,
	}

	for hName, hPos := range headers {
		switch hName {
		case "trip_id":
			trip.Id = getField(row, hName, hPos, &parseErrors, lineNumber, TripsFileName)
		case "route_id":
			trip.RouteId = getField(row, hName, hPos, &parseErrors, lineNumber, TripsFileName)
		case "service_id":
			trip.ServiceId = getField(row, hName, hPos, &parseErrors, lineNumber, TripsFileName)
		case "trip_headsign":
			trip.HeadSign = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, TripsFileName)
		case "trip_short_name":
			trip.ShortName = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, TripsFileName)
		case "direction_id":
			trip.DirectionId = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, TripsFileName)
		case "block_id":
			trip.BlockId = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, TripsFileName)
		case "shape_id":
			trip.ShapeId = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, TripsFileName)
		case "wheelchair_accessible":
			trip.WheelchairAccessible = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, TripsFileName)
		case "bikes_allowed":
			trip.BikesAllowed = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, TripsFileName)
		}
	}

	if len(parseErrors) > 0 {
		return &trip, parseErrors
	}
	return &trip, nil
}

// ValidateTrips performs additional validation for a list of Trip instances.
func ValidateTrips(trips []*Trip, routes []*Route, calendarItems []*CalendarItem, shapes []*Shape) ([]error, []string) {
	var validationErrors []error
	var recommendations []string

	if trips == nil {
		return validationErrors, recommendations
	}

	for _, trip := range trips {
		// Additional required field checks for individual Trip.
		if trip.Id == "" {
			validationErrors = append(validationErrors, createFileRowError(TripsFileName, trip.LineNumber, "trip_id must be specified"))
		}
		if trip.RouteId == "" {
			validationErrors = append(validationErrors, createFileRowError(TripsFileName, trip.LineNumber, "route_id must be specified"))
		}
		if trip.ServiceId == "" {
			validationErrors = append(validationErrors, createFileRowError(TripsFileName, trip.LineNumber, "service_id must be specified"))
		}
	}

	for _, trip := range trips {
		if trip == nil {
			continue
		}

		if routes != nil {
			routeFound := false
			for _, route := range routes {
				if route == nil {
					continue
				}
				if trip.RouteId == route.Id {
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
				if trip.ServiceId == calendarItem.ServiceId {
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
					if *trip.ShapeId == shape.Id {
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
