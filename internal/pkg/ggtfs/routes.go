package ggtfs

import (
	"encoding/csv"
	"fmt"
	"strconv"
)

type Route struct {
	Id                ID                     // route_id
	AgencyId          *ID                    // agency_id
	ShortName         *Text                  // route_short_name
	LongName          *Text                  // route_long_name
	Description       *Text                  // route_desc
	Type              RouteType              // route_type
	Url               *URL                   // route_url
	Color             *Color                 // route_color
	TextColor         *Color                 // route_text_color
	SortOrder         *Integer               // route_sort_order
	ContinuousPickup  *ContinuousPickupType  // continuous_pickup
	ContinuousDropOff *ContinuousDropOffType // continuous_drop_off
	NetworkId         *ID                    // network_id
	LineNumber        int
}

var validRouteHeaders = []string{"route_id", "agency_id", "route_short_name", "route_long_name", "route_desc",
	"route_type", "route_url", "route_color", "route_text_color", "route_sort_order", "continuous_pickup",
	"continuous_drop_off", "network_id"}

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

func CreateRoute(row []string, headers map[string]int, lineNumber int) interface{} {
	route := Route{
		LineNumber: lineNumber,
	}

	for hName, hPos := range headers {
		switch hName {
		case "route_id":
			route.Id = NewID(getRowValue(row, hPos))
		case "agency_id":
			route.AgencyId = NewOptionalID(getRowValue(row, hPos))
		case "route_short_name":
			route.ShortName = NewOptionalText(getRowValue(row, hPos))
		case "route_long_name":
			route.LongName = NewOptionalText(getRowValue(row, hPos))
		case "route_desc":
			route.Description = NewOptionalText(getRowValue(row, hPos))
		case "route_type":
			route.Type = NewRouteType(getRowValue(row, hPos))
		case "route_url":
			route.Url = NewOptionalURL(getRowValue(row, hPos))
		case "route_color":
			route.Color = NewOptionalColor(getRowValue(row, hPos))
		case "route_text_color":
			route.TextColor = NewOptionalColor(getRowValue(row, hPos))
		case "route_sort_order":
			route.SortOrder = NewOptionalInteger(getRowValue(row, hPos))
		case "continuous_pickup":
			route.ContinuousPickup = NewOptionalContinuousPickupType(getRowValue(row, hPos))
		case "continuous_drop_off":
			route.ContinuousDropOff = NewOptionalContinuousDropOffType(getRowValue(row, hPos))

		case "network_id":
			route.NetworkId = NewOptionalID(getRowValue(row, hPos))
		}
	}

	return &route
}

func ValidateRoutes(routes []*Route, agencies []*Agency) ([]error, []string) {
	var validationErrors []error
	var recommendations []string

	if routes == nil {
		return validationErrors, recommendations
	}

	for _, route := range routes {
		// Additional required field checks for individual Route.
		if route.Id.String() == "" {
			validationErrors = append(validationErrors, createFileRowError(RoutesFileName, route.LineNumber, "route_id must be specified"))
		}
		if route.Type.String() == "" {
			validationErrors = append(validationErrors, createFileRowError(RoutesFileName, route.LineNumber, "route_type must be specified"))
		}
		if route.ShortName == nil && route.LongName == nil {
			validationErrors = append(validationErrors, createFileRowError(RoutesFileName, route.LineNumber, "either route_short_name or route_long_name must be specified"))
		}
	}

	if agencies == nil {
		return validationErrors, recommendations
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

			if route.AgencyId != nil && route.AgencyId.String() == agency.Id.String() {
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
		if usedIds[route.Id.String()] {
			validationErrors = append(validationErrors, createFileRowError(RoutesFileName, route.LineNumber, fmt.Sprintf("route_id '%s' is not unique within the file", route.Id.String())))
		} else {
			usedIds[route.Id.String()] = true
		}
	}

	return validationErrors, recommendations
}

const (
	TramStreetCarLightRailRouteType   = 0
	SubwayMetroRouteType              = 1
	RailRouteType                     = 2
	BusRouteType                      = 3
	FerryRouteType                    = 4
	CableTramRouteType                = 5
	AerialLiftSuspendedCableRouteType = 6
	FunicularRouteType                = 7
	TrolleyBusRouteType               = 11
	MonorailRouteType                 = 12
)

type RouteType struct {
	Integer
}

func (ete RouteType) IsValid() bool {
	val, err := strconv.Atoi(ete.Integer.base.raw)
	if err != nil {
		return false
	}

	return val == TramStreetCarLightRailRouteType || val == SubwayMetroRouteType || val == RailRouteType || val == BusRouteType ||
		val == FerryRouteType || val == CableTramRouteType || val == AerialLiftSuspendedCableRouteType || val == FunicularRouteType ||
		val == TrolleyBusRouteType || val == MonorailRouteType
}

func NewRouteType(raw *string) RouteType {
	if raw == nil {
		return RouteType{
			Integer{base: base{raw: ""}}}
	}
	return RouteType{Integer{base: base{raw: *raw, isPresent: true}}}
}

const (
	ContinuousStoppingPickupType       = 0
	NoContinuousStoppingPickupType     = 1
	MustPhoneAgencyPickupType          = 2
	MustCoordinateWithDriverPickupType = 3
)

type ContinuousPickupType struct {
	Integer
}

func (cpt ContinuousPickupType) IsValid() bool {
	// Spec says "1 or empty - No continuous stopping drop off."
	// Empty = valid
	if cpt.Integer.base.IsEmpty() {
		return true
	}

	val, err := strconv.Atoi(cpt.Integer.base.raw)
	if err != nil {
		return false
	}

	return val == ContinuousStoppingPickupType || val == NoContinuousStoppingPickupType ||
		val == MustPhoneAgencyPickupType || val == MustCoordinateWithDriverPickupType
}

func NewOptionalContinuousPickupType(raw *string) *ContinuousPickupType {
	if raw == nil {
		return &ContinuousPickupType{
			Integer{base: base{raw: ""}}}
	}
	return &ContinuousPickupType{Integer{base: base{raw: *raw}}}
}

const (
	ContinuousStoppingDropOffType       = 0
	NoContinuousStoppingDropOffType     = 1
	MustPhoneAgencyDropOffType          = 2
	MustCoordinateWithDriverDropOffType = 3
)

type ContinuousDropOffType struct {
	Integer
}

func (cpt ContinuousDropOffType) IsValid() bool {
	// Spec says "1 or empty - No continuous stopping drop off."
	// Empty = valid
	if cpt.Integer.base.IsEmpty() {
		return true
	}

	val, err := strconv.Atoi(cpt.Integer.base.raw)
	if err != nil {
		return false
	}

	return val == ContinuousStoppingDropOffType || val == NoContinuousStoppingDropOffType ||
		val == MustPhoneAgencyDropOffType || val == MustCoordinateWithDriverDropOffType
}

func NewOptionalContinuousDropOffType(raw *string) *ContinuousDropOffType {
	if raw == nil {
		return &ContinuousDropOffType{
			Integer{base: base{raw: ""}}}
	}
	return &ContinuousDropOffType{Integer{base: base{raw: *raw}}}
}
