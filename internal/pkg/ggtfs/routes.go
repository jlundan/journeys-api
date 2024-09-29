package ggtfs

import (
	"encoding/csv"
	"fmt"
)

// Route struct with optional fields as string pointers and required fields as strings.
type Route struct {
	Id                string
	AgencyId          string
	ShortName         *string
	LongName          *string
	Desc              *string
	Type              string
	Url               *string
	Color             *string
	TextColor         *string
	SortOrder         *string
	ContinuousPickup  *string
	ContinuousDropOff *string
	LineNumber        int
}

var validRouteHeaders = []string{"route_id", "agency_id", "route_short_name", "route_long_name", "route_desc",
	"route_type", "route_url", "route_color", "route_text_color", "route_sort_order", "continuous_pickup",
	"continuous_drop_off", "network_id"}

// LoadRoutes loads routes from a CSV reader and returns them along with any errors.
func LoadRoutes(csvReader *csv.Reader) ([]*Route, []error) {
	entities, errs := loadEntities(csvReader, validRouteHeaders, CreateRoute, RoutesFileName)

	routes := make([]*Route, 0, len(entities))
	for _, entity := range entities {
		if route, ok := entity.(*Route); ok {
			routes = append(routes, route)
		}
	}

	return routes, errs
}

// CreateRoute creates and validates a Route instance from the CSV row data.
func CreateRoute(row []string, headers map[string]uint8, lineNumber int) (interface{}, []error) {
	var parseErrors []error

	route := Route{
		LineNumber: lineNumber,
	}

	for hName, hPos := range headers {
		switch hName {
		case "route_id":
			route.Id = getField(row, hName, hPos, &parseErrors, lineNumber, RoutesFileName)
		case "agency_id":
			route.AgencyId = getField(row, hName, hPos, &parseErrors, lineNumber, RoutesFileName)
		case "route_short_name":
			route.ShortName = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, RoutesFileName)
		case "route_long_name":
			route.LongName = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, RoutesFileName)
		case "route_desc":
			route.Desc = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, RoutesFileName)
		case "route_url":
			route.Url = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, RoutesFileName)
		case "route_color":
			route.Color = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, RoutesFileName)
		case "route_text_color":
			route.TextColor = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, RoutesFileName)
		case "route_sort_order":
			route.SortOrder = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, RoutesFileName)
		case "continuous_pickup":
			route.ContinuousPickup = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, RoutesFileName)
		case "continuous_drop_off":
			route.ContinuousDropOff = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, RoutesFileName)
		case "route_type":
			route.Type = getField(row, hName, hPos, &parseErrors, lineNumber, RoutesFileName)
		}
	}

	if len(parseErrors) > 0 {
		return &route, parseErrors
	}
	return &route, nil
}

// ValidateRoutes performs additional validation for a list of Route instances.
func ValidateRoutes(routes []*Route, agencies []*Agency) []error {
	var validationErrors []error

	if routes == nil {
		return validationErrors
	}

	for _, route := range routes {
		// Additional required field checks for individual Route.
		if route.Id == "" {
			validationErrors = append(validationErrors, createFileRowError(RoutesFileName, route.LineNumber, "route_id must be specified"))
		}
		if route.Type == "" {
			validationErrors = append(validationErrors, createFileRowError(RoutesFileName, route.LineNumber, "route_type must be specified"))
		}
		if route.ShortName == nil && route.LongName == nil {
			validationErrors = append(validationErrors, createFileRowError(RoutesFileName, route.LineNumber, "either route_short_name or route_long_name must be specified"))
		}
	}

	if agencies == nil {
		return validationErrors
	}

	// Check for valid agency_id references in the routes.
	for _, route := range routes {
		if route == nil {
			continue
		}
		agencyFound := false
		for _, agency := range agencies {
			if agency == nil {
				continue
			}
			if route.AgencyId == agency.Id {
				agencyFound = true
				break
			}
		}
		if !agencyFound {
			validationErrors = append(validationErrors, createFileRowError(RoutesFileName, route.LineNumber, fmt.Sprintf("referenced agency_id '%s' not found in %s", route.AgencyId, AgenciesFileName)))
		}
	}

	// Check for duplicate route_id entries.
	usedIds := make(map[string]bool)
	for _, route := range routes {
		if route == nil {
			continue
		}
		if usedIds[route.Id] {
			validationErrors = append(validationErrors, createFileRowError(RoutesFileName, route.LineNumber, fmt.Sprintf("route_id '%s' is not unique within the file", route.Id)))
		} else {
			usedIds[route.Id] = true
		}
	}

	return validationErrors
}

//package ggtfs
//
//import (
//	"encoding/csv"
//	"errors"
//	"fmt"
//	"strconv"
//)
//
//type Route struct {
//	Id                string
//	AgencyId          string
//	ShortName         *string
//	LongName          *string
//	Desc              *string
//	Type              int
//	Url               *string
//	Color             *string
//	TextColor         *string
//	SortOrder         *int
//	ContinuousPickup  *int
//	ContinuousDropOff *int
//	lineNumber        int
//}
//
//var validRouteHeaders = []string{"route_id", "agency_id", "route_short_name", "route_long_name", "route_desc",
//	"route_type", "route_url", "route_color", "route_text_color", "route_sort_order", "continuous_pickup",
//	"continuous_drop_off", "network_id"}
//
//func LoadRoutes(csvReader *csv.Reader) ([]*Route, []error) {
//	routes := make([]*Route, 0)
//	errs := make([]error, 0)
//
//	headers, err := ReadHeaderRow(csvReader, validRouteHeaders)
//	if err != nil {
//		errs = append(errs, createFileError(RoutesFileName, fmt.Sprintf("read error: %v", err.Error())))
//		return routes, errs
//	}
//	if headers == nil {
//		return routes, errs
//	}
//
//	usedIds := make([]string, 0)
//	index := 0
//	for {
//		row, err := ReadDataRow(csvReader)
//		if err != nil {
//			errs = append(errs, createFileError(RoutesFileName, fmt.Sprintf("%v", err.Error())))
//			index++
//			continue
//		}
//		if row == nil {
//			break
//		}
//
//		rowErrs := make([]error, 0)
//		route := Route{
//			lineNumber: index,
//		}
//
//		var (
//			routeId   *string
//			routeType *int
//		)
//		for name, column := range headers {
//			switch name {
//			case "route_id":
//				routeId = handleStringField(row[column], RoutesFileName, name, index, &rowErrs)
//			case "agency_id":
//				route.AgencyId = row[column] //handleStringField(row[column], RoutesFileName, name, index, &rowErrs)
//			case "route_short_name":
//				route.ShortName = handleStringField(row[column], RoutesFileName, name, index, &rowErrs)
//			case "route_long_name":
//				//route.LongName = handleStringField(row[column], RoutesFileName, name, index, &rowErrs)
//				route.LongName = &row[column]
//			case "route_desc":
//				route.Desc = handleStringField(row[column], RoutesFileName, name, index, &rowErrs)
//			case "route_url":
//				route.Url = handleURLField(row[column], RoutesFileName, name, index, &rowErrs)
//			case "route_color":
//				route.Color = handleColorField(row[column], RoutesFileName, name, index, &rowErrs)
//			case "route_text_color":
//				route.TextColor = handleColorField(row[column], RoutesFileName, name, index, &rowErrs)
//			case "route_sort_order":
//				route.SortOrder = handleIntField(row[column], RoutesFileName, name, index, &rowErrs)
//			case "continuous_pickup":
//				route.ContinuousPickup = handleContinuousPickupField(row[column], RoutesFileName, name, index, &rowErrs)
//			case "continuous_drop_off":
//				route.ContinuousDropOff = handleContinuousDropOffField(row[column], RoutesFileName, name, index, &rowErrs)
//			case "route_type":
//				routeType = handleRouteTypeField(row[column], RoutesFileName, name, index, &rowErrs)
//			}
//		}
//
//		if routeId == nil {
//			rowErrs = append(rowErrs, createFileRowError(RoutesFileName, index, "route_id must be specified"))
//		} else {
//			route.Id = *routeId
//			if StringArrayContainsItem(usedIds, *routeId) {
//				errs = append(errs, createFileRowError(RoutesFileName, index, fmt.Sprintf("%s: route_id", nonUniqueId)))
//			} else {
//				usedIds = append(usedIds, *routeId)
//			}
//		}
//
//		if routeType == nil {
//			rowErrs = append(rowErrs, createFileRowError(RoutesFileName, index, "route_type must be specified"))
//		} else {
//			route.Type = *routeType
//		}
//
//		if route.ShortName == nil && route.LongName == nil {
//			rowErrs = append(rowErrs, createFileRowError(RoutesFileName, index, "either route_short_name or route_long_name must be specified"))
//		}
//
//		if len(rowErrs) > 0 {
//			errs = append(errs, rowErrs...)
//		} else {
//			routes = append(routes, &route)
//		}
//
//		index++
//	}
//
//	return routes, errs
//}
//
//func ValidateRoutes(routes []*Route, agencies []*Agency) []error {
//	var validationErrors []error
//
//	if routes == nil || agencies == nil {
//		return validationErrors
//	}
//
//	for _, route := range routes {
//		if route == nil {
//			continue
//		}
//		notFound := true
//		for _, agency := range agencies {
//			if agency == nil {
//				continue
//			}
//			if route.AgencyId == agency.Id {
//				notFound = false
//				break
//			}
//		}
//		if notFound {
//			validationErrors = append(validationErrors, createFileRowError(RoutesFileName, route.lineNumber, fmt.Sprintf("referenced agency_id not found in %s", AgenciesFileName)))
//		}
//	}
//
//	return validationErrors
//}
//
//func handleRouteTypeField(str string, fileName string, fieldName string, index int, errs *[]error) *int {
//	n, err := strconv.ParseInt(str, 10, 64)
//	if err != nil {
//		*errs = append(*errs, createFieldError(fileName, fieldName, index, err))
//	}
//
//	if rt := int(n); rt <= 7 || (rt >= 11 && rt <= 12) {
//		return &rt
//	} else {
//		*errs = append(*errs, createFieldError(fileName, fieldName, index, errors.New(invalidValue)))
//		return nil
//	}
//}
//
////const (
////	TramStreetcarLightRail RouteType = iota
////	SubwayMetro
////	Rail
////	Bus
////	Ferry
////	CableTram
////	AerialLiftSuspendedCableCar
////	Funicular
////	TrolleyBus
////	Monorail
////)
////
////const (
////	ContinuousPickup ContinuousPickupType = iota
////	NoContinuousPickup
////	PhoneAgencyPickup
////	CoordinateWithDriverPickup
////)
////
////const (
////	ContinuousDropOff ContinuousDropOffType = iota
////	NoContinuousDropOff
////	PhoneAgencyDropOff
////	CoordinateWithDriverDropOff
////)
