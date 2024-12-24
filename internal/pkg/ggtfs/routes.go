package ggtfs

import (
	"fmt"
	"strconv"
)

type Route struct {
	Id                ID                    // route_id, required
	AgencyId          ID                    // agency_id, conditionally required
	ShortName         Text                  // route_short_name, conditionally required
	LongName          Text                  // route_long_name, conditionally required
	Description       Text                  // route_desc, optional
	Type              RouteType             // route_type, required
	URL               URL                   // route_url
	Color             Color                 // route_color
	TextColor         Color                 // route_text_color
	SortOrder         Integer               // route_sort_order
	ContinuousPickup  ContinuousPickupType  // continuous_pickup
	ContinuousDropOff ContinuousDropOffType // continuous_drop_off
	NetworkId         ID                    // network_id
	LineNumber        int
}

func (r Route) Validate() []error {
	var validationErrors []error

	requiredFields := []struct {
		fieldName string
		field     ValidAndPresentField
	}{
		{"route_id", &r.Id},
		{"route_type", &r.Type},
	}
	for _, f := range requiredFields {
		if !f.field.IsValid() {
			validationErrors = append(validationErrors, createFileRowError(RoutesFileName, r.LineNumber, createInvalidRequiredFieldString(f.fieldName)))
		}
	}

	optionalFields := []struct {
		field     ValidAndPresentField
		fieldName string
	}{
		{&r.Id, "agency_id"},
		{&r.ShortName, "route_short_name"},
		{&r.LongName, "route_long_name"},
		{&r.Description, "route_desc"},
		{&r.URL, "route_url"},
		{&r.Color, "route_color"},
		{&r.TextColor, "route_text_color"},
		{&r.SortOrder, "route_sort_order"},
		{&r.ContinuousPickup, "continuous_pickup"},
		{&r.ContinuousDropOff, "continuous_drop_off"},
		{&r.NetworkId, "network_id"},
	}

	for _, field := range optionalFields {
		if field.field != nil && field.field.IsPresent() && !field.field.IsValid() {
			validationErrors = append(validationErrors, createFileRowError(RoutesFileName, r.LineNumber, createInvalidFieldString(field.fieldName)))
		}
	}

	return validationErrors
}

func CreateRoute(row []string, headers map[string]int, lineNumber int) *Route {
	route := Route{
		LineNumber: lineNumber,
	}

	for hName := range headers {
		v := getRowValueForHeaderName(row, headers, hName)

		switch hName {
		case "route_id":
			route.Id = NewID(v)
		case "agency_id":
			route.AgencyId = NewID(v)
		case "route_short_name":
			route.ShortName = NewText(v)
		case "route_long_name":
			route.LongName = NewText(v)
		case "route_desc":
			route.Description = NewText(v)
		case "route_type":
			route.Type = NewRouteType(v)
		case "route_url":
			route.URL = NewURL(v)
		case "route_color":
			route.Color = NewColor(v)
		case "route_text_color":
			route.TextColor = NewColor(v)
		case "route_sort_order":
			route.SortOrder = NewInteger(v)
		case "continuous_pickup":
			route.ContinuousPickup = NewContinuousPickupType(v)
		case "continuous_drop_off":
			route.ContinuousDropOff = NewContinuousDropOffType(v)
		case "network_id":
			route.NetworkId = NewID(v)
		}
	}

	return &route
}

func ValidateRoutes(routes []*Route, agencies []*Agency) ([]error, []string) {
	var validationErrors []error

	usedIds := make(map[string]bool)
	for _, route := range routes {
		if route == nil {
			continue
		}

		vErr := route.Validate()
		if len(vErr) > 0 {
			validationErrors = append(validationErrors, vErr...)
			continue
		}

		if usedIds[route.Id.String()] {
			validationErrors = append(validationErrors, createFileRowError(RoutesFileName, route.LineNumber, fmt.Sprintf("route_id '%s' is not unique within the file", route.Id.String())))
		} else {
			usedIds[route.Id.String()] = true
		}
	}

	if agencies == nil {
		return validationErrors, nil
	}

	for _, route := range routes {
		if route == nil || !route.AgencyId.IsValid() {
			continue
		}

		matchingAgencyFound := false
		for _, agency := range agencies {
			if agency == nil {
				continue
			}

			if route.AgencyId.String() == agency.Id.String() {
				matchingAgencyFound = true
				break
			}
		}

		if !matchingAgencyFound {
			validationErrors = append(validationErrors, createFileRowError(RoutesFileName, route.LineNumber, fmt.Sprintf("referenced agency_id '%s' not found in %s", route.AgencyId.String(), AgenciesFileName)))
		}
	}

	return validationErrors, nil
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
	// Spec says values "1 or empty mean No continuous stopping drop off."
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

func NewContinuousPickupType(raw *string) ContinuousPickupType {
	if raw == nil {
		return ContinuousPickupType{
			Integer{base: base{raw: ""}}}
	}
	return ContinuousPickupType{Integer{base: base{raw: *raw, isPresent: true}}}
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

func NewContinuousDropOffType(raw *string) ContinuousDropOffType {
	if raw == nil {
		return ContinuousDropOffType{
			Integer{base: base{raw: ""}}}
	}
	return ContinuousDropOffType{Integer{base: base{raw: *raw, isPresent: true}}}
}
