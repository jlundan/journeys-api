package ggtfs

import (
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
)

type Trip struct {
	Id                   string
	RouteId              string
	ServiceId            string
	HeadSign             *string
	ShortName            *string
	DirectionId          *int
	BlockId              *string
	ShapeId              *string
	WheelchairAccessible *int
	BikesAllowed         *int
	lineNumber           int
}

func ExtractTripsByRoutes(input *csv.Reader, output *csv.Writer, routes map[string]struct{}) (map[string]struct{}, map[string]struct{}, map[string]struct{}, []error) {
	errs := make([]error, 0)

	headers, err := ReadHeaderRow(input)
	if err != nil {
		errs = append(errs, createFileError(TripsFileName, fmt.Sprintf("read error: %v", err.Error())))
		return nil, nil, nil, errs
	}

	if headers == nil { // EOF
		return nil, nil, nil, nil
	}

	var idHeaderPos uint8
	if routeIdPos, hasRouteId := headers["route_id"]; hasRouteId {
		idHeaderPos = routeIdPos
	} else {
		errs = append(errs, createFileError(TripsFileName, "cannot extract trips without route_id column"))
		return nil, nil, nil, errs
	}

	var tripIdPos uint8
	if pos, hasTripId := headers["trip_id"]; hasTripId {
		tripIdPos = pos
	} else {
		errs = append(errs, createFileError(TripsFileName, "cannot extract trips without trip_id column"))
		return nil, nil, nil, errs
	}

	var serviceIdPos *uint8
	if pos, hasServiceId := headers["service_id"]; hasServiceId {
		serviceIdPos = &pos
	} else {
		errs = append(errs, createFileError(TripsFileName, "no service_id column, services will not be extracted"))
		serviceIdPos = nil
	}

	var shapeIdPos *uint8
	if pos, hasShapeId := headers["shape_id"]; hasShapeId {
		shapeIdPos = &pos
	} else {
		errs = append(errs, createFileError(TripsFileName, "no shape_id column, shapes will not be extracted"))
		shapeIdPos = nil
	}

	err = writeHeaderRow(headers, output)
	if err != nil {
		errs = append(errs, err)
		return nil, nil, nil, errs
	}

	tripIdMap := make(map[string]struct{}, 0)
	serviceIdMap := make(map[string]struct{}, 0)
	shapeIdMap := make(map[string]struct{}, 0)
	for {
		row, rErr := ReadDataRow(input)
		if rErr != nil {
			errs = append(errs, createFileError(TripsFileName, fmt.Sprintf("%v", rErr.Error())))
			continue
		}

		if row == nil { // EOF
			break
		}

		if _, shouldBeExtracted := routes[row[idHeaderPos]]; shouldBeExtracted {
			wErr := writeDataRow(row, output)
			if wErr != nil {
				errs = append(errs, wErr)
				return nil, nil, nil, errs
			}

			tripId := row[tripIdPos]
			if _, tripAlreadyExists := tripIdMap[tripId]; !tripAlreadyExists {
				tripIdMap[tripId] = struct{}{}
			}

			if serviceIdPos == nil {
				continue
			}
			serviceId := row[*serviceIdPos]
			if _, serviceAlreadyExists := serviceIdMap[serviceId]; !serviceAlreadyExists {
				serviceIdMap[serviceId] = struct{}{}
			}

			if shapeIdPos == nil {
				continue
			}
			shapeId := row[*shapeIdPos]
			if _, alreadyExists := shapeIdMap[shapeId]; !alreadyExists {
				shapeIdMap[shapeId] = struct{}{}
			}
		}
	}

	return tripIdMap, serviceIdMap, shapeIdMap, nil
}

func LoadTrips(csvReader *csv.Reader) ([]*Trip, []error) {
	stops := make([]*Trip, 0)
	errs := make([]error, 0)

	headers, err := ReadHeaderRow(csvReader)
	if err != nil {
		errs = append(errs, createFileError(TripsFileName, fmt.Sprintf("read error: %v", err.Error())))
		return stops, errs
	}
	if headers == nil {
		return stops, errs
	}

	usedIds := make([]string, 0)
	index := 0
	for {
		row, err := ReadDataRow(csvReader)
		if err != nil {
			errs = append(errs, createFileError(TripsFileName, fmt.Sprintf("%v", err.Error())))
			index++
			continue
		}
		if row == nil {
			break
		}

		rowErrs := make([]error, 0)
		trip := Trip{
			lineNumber: index,
		}

		var (
			tripId    *string
			routeId   *string
			serviceId *string
		)

		for name, column := range headers {
			switch name {
			case "route_id":
				routeId = handleIDField(row[column], TripsFileName, name, index, &rowErrs)
			case "service_id":
				serviceId = handleIDField(row[column], TripsFileName, name, index, &rowErrs)
			case "trip_id":
				tripId = handleIDField(row[column], TripsFileName, name, index, &rowErrs)
			case "trip_headsign":
				trip.HeadSign = handleTextField(row[column], TripsFileName, name, index, &rowErrs)
			case "trip_short_name":
				trip.ShortName = handleTextField(row[column], TripsFileName, name, index, &rowErrs)
			case "direction_id":
				trip.DirectionId = handleDirectionField(row[column], TripsFileName, name, index, &rowErrs)
			case "block_id":
				trip.BlockId = handleIDField(row[column], TripsFileName, name, index, &rowErrs)
			case "shape_id":
				trip.ShapeId = handleIDField(row[column], TripsFileName, name, index, &rowErrs)
			case "wheelchair_accessible":
				trip.WheelchairAccessible = handleWheelchairAccessibleField(row[column], TripsFileName, name, index, &rowErrs)
			case "bikes_allowed":
				trip.BikesAllowed = handleBikesAllowedField(row[column], TripsFileName, name, index, &rowErrs)
			}
		}

		if tripId == nil {
			rowErrs = append(rowErrs, createFileRowError(TripsFileName, trip.lineNumber, "trip_id must be specified"))
		} else {
			trip.Id = *tripId
			if StringArrayContainsItem(usedIds, *tripId) {
				errs = append(errs, createFileRowError(TripsFileName, index, fmt.Sprintf("%s: trip_id", nonUniqueId)))
			} else {
				usedIds = append(usedIds, *tripId)
			}
		}

		if serviceId == nil {
			rowErrs = append(rowErrs, createFileRowError(TripsFileName, trip.lineNumber, "service_id must be specified"))
		} else {
			trip.ServiceId = *serviceId
		}

		if routeId == nil {
			rowErrs = append(rowErrs, createFileRowError(TripsFileName, trip.lineNumber, "route_id must be specified"))
		} else {
			trip.RouteId = *routeId
		}

		if len(rowErrs) > 0 {
			errs = append(errs, rowErrs...)
		} else {
			stops = append(stops, &trip)
		}

		index++
	}

	return stops, errs
}

func ValidateTrips(trips []*Trip, routes []*Route, calendarItems []*CalendarItem, shapes []*Shape) []error {
	var validationErrors []error

	if trips == nil {
		return validationErrors
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
				validationErrors = append(validationErrors, createFileRowError(TripsFileName, trip.lineNumber, fmt.Sprintf("referenced route_id not found in %s", RoutesFileName)))
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
				validationErrors = append(validationErrors, createFileRowError(TripsFileName, trip.lineNumber, fmt.Sprintf("referenced service_id not found in %s", CalendarFileName)))
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
					validationErrors = append(validationErrors, createFileRowError(TripsFileName, trip.lineNumber, fmt.Sprintf("referenced shape_id not found in %s", ShapesFileName)))
				}
			}
		}
	}

	return validationErrors
}

func handleDirectionField(str string, fileName string, fieldName string, index int, errs *[]error) *int {
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, err))
		return nil
	}

	if v := int(n); v <= 1 {
		return &v
	} else {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, errors.New(invalidValue)))
		return nil
	}
}

func handleWheelchairAccessibleField(str string, fileName string, fieldName string, index int, errs *[]error) *int {
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, err))
		return nil
	}

	if v := int(n); v <= 2 {
		return &v
	} else {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, errors.New(invalidValue)))
		return nil
	}
}

func handleBikesAllowedField(str string, fileName string, fieldName string, index int, errs *[]error) *int {
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, err))
		return nil
	}

	if v := int(n); v <= 2 {
		return &v
	} else {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, errors.New(invalidValue)))
		return nil
	}
}
