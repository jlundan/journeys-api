package ggtfs

import (
	"fmt"
	"strconv"
)

type Route struct {
	Id                ID                    // route_id 			(required, unique)
	AgencyId          ID                    // agency_id 			(conditionally required)
	ShortName         Text                  // route_short_name 	(conditionally required)
	LongName          Text                  // route_long_name 		(conditionally required)
	Desc              Text                  // route_desc 			(optional)
	Type              RouteType             // route_type 			(required)
	URL               URL                   // route_url 			(optional)
	Color             Color                 // route_color 			(optional)
	TextColor         Color                 // route_text_color 	(optional)
	SortOrder         Integer               // route_sort_order 	(optional)
	ContinuousPickup  ContinuousPickupType  // continuous_pickup 	(conditionally forbidden)
	ContinuousDropOff ContinuousDropOffType // continuous_drop_off 	(conditionally forbidden)
	NetworkId         ID                    // network_id 			(conditionally forbidden)
	LineNumber        int
}

func (r Route) Validate() ([]error, []string) {
	var validationErrors []error
	var recommendations []string

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
		// agency_id is required only if there are multiple agencies in agencies.txt, recommended if there is only one. This is handled in ValidateAgencies func.
		// route_short_name is required if route_long_name is empty or not present. This is handled below.
		// route_long_name is required if route_short_name is empty or not present. This is handled below.
		{&r.Desc, "route_desc"},
		{&r.URL, "route_url"},
		{&r.Color, "route_color"},
		{&r.TextColor, "route_text_color"},
		{&r.SortOrder, "route_sort_order"},
		{&r.ContinuousPickup, "continuous_pickup"},
		{&r.ContinuousDropOff, "continuous_drop_off"},
		// network id is forbidden if the route_networks.txt file is present. This is handled in ValidateRoutes func.
	}

	for _, field := range optionalFields {
		if field.field != nil && field.field.IsPresent() && !field.field.IsValid() {
			validationErrors = append(validationErrors, createFileRowError(RoutesFileName, r.LineNumber, createInvalidFieldString(field.fieldName)))
		}
	}

	if r.ShortName.IsEmpty() && !r.LongName.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(RoutesFileName, r.LineNumber, "route_short_name must be specified when route_long_name is empty or not present"))
	}

	if r.LongName.IsEmpty() && !r.ShortName.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(RoutesFileName, r.LineNumber, "route_long_name must be specified when route_short_name is empty or not present"))
	}

	if r.ShortName.Length() >= 12 {
		recommendations = append(recommendations, createFileRowRecommendation(RoutesFileName, r.LineNumber, "route_short_name should be less than 12 characters"))
	}

	if r.Desc.IsValid() && (r.Desc == r.ShortName || r.Desc == r.LongName) {
		recommendations = append(recommendations, createFileRowRecommendation(RoutesFileName, r.LineNumber, "route_desc should not be the same as route_short_name or route_long_name"))
	}

	if r.SortOrder.IsValid() && r.SortOrder.Int() < 0 {
		validationErrors = append(validationErrors, createFileRowError(RoutesFileName, r.LineNumber, createInvalidFieldString("sort_order")))
	}

	// TODO: VALIDATION: route_color: The color difference between route_color and route_text_color should provide sufficient contrast when viewed on a black and white screen.
	// Implement this if there is a way to check the contrast between two colors

	// TODO: VALIDATION: route_text_color: The color difference between route_color and route_text_color should provide sufficient contrast when viewed on a black and white screen.
	// Implement this if there is a way to check the contrast between two colors

	return validationErrors, recommendations
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
			route.Desc = NewText(v)
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

	// TODO: IMPLEMENTATION: route_color: Route color designation that matches public facing material. Defaults to white (FFFFFF) when omitted or left empty.
	// Not implemented currently

	// TODO: IMPLEMENTATION: route_text_color: Legible color to use for text drawn against a background of route_color. Defaults to black (000000) when omitted or left empty.
	// Not implemented currently

	return &route
}

func ValidateRoutes(routes []*Route, agencies []*Agency) ([]error, []string) {
	var validationErrors []error
	var validationRecommendations []string

	// Count the number of agencies in agencies.txt, this is used to determine if agency_id is required or recommended later on.
	numAgencies := 0
	if agencies != nil {
		existingAgencyIds := make(map[string]bool)
		for _, agency := range agencies {
			if agency == nil || !agency.Id.IsValid() {
				continue
			}

			if !existingAgencyIds[agency.Id.String()] {
				numAgencies++
				existingAgencyIds[agency.Id.String()] = true
			}
		}
	}

	usedIds := make(map[string]bool)
	for _, route := range routes {
		if route == nil {
			continue
		}

		// basic validation for route
		vErr, vRec := route.Validate()
		if len(vRec) > 0 {
			validationRecommendations = append(validationRecommendations, vRec...)
		}
		if len(vErr) > 0 {
			validationErrors = append(validationErrors, vErr...)
		}

		// agency_id is required only if there are multiple agencies in agencies.txt, recommended otherwise.
		if numAgencies > 1 && !route.AgencyId.IsValid() {
			validationErrors = append(validationErrors, createFileRowError(RoutesFileName, route.LineNumber, "agency_id is required when there are multiple agencies in agencies.txt"))
		} else if !route.AgencyId.IsValid() {
			validationRecommendations = append(validationRecommendations, createFileRowRecommendation(RoutesFileName, route.LineNumber, "specifying agency_id is recommended"))
		}

		// route_id must be unique within the routes.txt file
		if usedIds[route.Id.String()] {
			validationErrors = append(validationErrors, createFileRowError(RoutesFileName, route.LineNumber, fmt.Sprintf("route_id '%s' is not unique within the file", route.Id.String())))
		} else {
			usedIds[route.Id.String()] = true
		}

		if agencies == nil || !route.AgencyId.IsValid() {
			continue
		}

		// agency_id must be a valid agency_id from agencies.txt
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

	// TODO: VALIDATION: route_url: URL of a web page about the particular route. Should be different from the agency.agency_url value.
	// Implement this when the route entity has a method to show the agency it is connected to

	// TODO: VALIDATION: network_id: Forbidden it the route_networks.txt file is present.
	// Implement this when we have support for route_networks.txt

	return validationErrors, validationRecommendations
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
