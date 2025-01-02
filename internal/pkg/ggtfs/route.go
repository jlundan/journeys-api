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
	SortOrder         PositiveInteger       // route_sort_order 	(optional)
	ContinuousPickup  ContinuousPickupType  // continuous_pickup 	(conditionally forbidden)
	ContinuousDropOff ContinuousDropOffType // continuous_drop_off 	(conditionally forbidden)
	NetworkId         ID                    // network_id 			(conditionally forbidden)
	LineNumber        int
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
			route.SortOrder = NewPositiveInteger(v)
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

func ValidateRoute(r Route) ([]error, []string) {
	var validationErrors []error
	var recommendations []string

	requiredFields := map[string]FieldTobeValidated{
		"route_id":   &r.Id,
		"route_type": &r.Type,
	}
	validateRequiredFields(requiredFields, &validationErrors, r.LineNumber, RoutesFileName)

	optionalFields := map[string]FieldTobeValidated{
		"route_desc":          &r.Desc,
		"route_url":           &r.URL,
		"route_color":         &r.Color,
		"route_text_color":    &r.TextColor,
		"route_sort_order":    &r.SortOrder,
		"continuous_pickup":   &r.ContinuousPickup,
		"continuous_drop_off": &r.ContinuousDropOff,
	}
	validateOptionalFields(optionalFields, &validationErrors, r.LineNumber, RoutesFileName)

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

	return validationErrors, recommendations
}

func ValidateRoutes(routes []*Route, agencies []*Agency) ([]error, []string) {
	var validationErrors []error
	var validationRecommendations []string

	// Count the number of agencies in agencies.txt, this is used to determine if agency_id is required or recommended later on.
	numAgencies := 0
	if agencies != nil {
		existingAgencyIds := make(map[string]bool)
		for _, agency := range agencies {
			if agency == nil || StringIsNilOrEmpty(agency.Id) {
				continue
			}

			if !existingAgencyIds[*agency.Id] {
				numAgencies++
				existingAgencyIds[*agency.Id] = true
			}
		}
	}

	usedIds := make(map[string]bool)
	for _, route := range routes {
		if route == nil {
			continue
		}

		// basic validation for route
		vErr, vRec := ValidateRoute(*route)
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
		if usedIds[route.Id.Raw()] {
			validationErrors = append(validationErrors, createFileRowError(RoutesFileName, route.LineNumber, fmt.Sprintf("route_id '%s' is not unique within the file", route.Id.Raw())))
		} else {
			usedIds[route.Id.Raw()] = true
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

			if route.AgencyId.Raw() == *agency.Id {
				matchingAgencyFound = true
				break
			}
		}

		if !matchingAgencyFound {
			validationErrors = append(validationErrors, createFileRowError(RoutesFileName, route.LineNumber, fmt.Sprintf("referenced agency_id '%s' not found in %s", route.AgencyId.Raw(), AgenciesFileName)))
		}
	}

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
