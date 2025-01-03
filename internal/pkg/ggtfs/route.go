package ggtfs

type Route struct {
	Id                *string // route_id 			(required, unique)
	AgencyId          *string // agency_id 			(conditionally required)
	ShortName         *string // route_short_name 	(conditionally required)
	LongName          *string // route_long_name 		(conditionally required)
	Desc              *string // route_desc 			(optional)
	Type              *string // route_type 			(required)
	URL               *string // route_url 			(optional)
	Color             *string // route_color 			(optional)
	TextColor         *string // route_text_color 	(optional)
	SortOrder         *string // route_sort_order 	(optional)
	ContinuousPickup  *string // continuous_pickup 	(conditionally forbidden)
	ContinuousDropOff *string // continuous_drop_off 	(conditionally forbidden)
	NetworkId         *string // network_id 			(conditionally forbidden)
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
			route.Id = v
		case "agency_id":
			route.AgencyId = v
		case "route_short_name":
			route.ShortName = v
		case "route_long_name":
			route.LongName = v
		case "route_desc":
			route.Desc = v
		case "route_type":
			route.Type = v
		case "route_url":
			route.URL = v
		case "route_color":
			route.Color = v
		case "route_text_color":
			route.TextColor = v
		case "route_sort_order":
			route.SortOrder = v
		case "continuous_pickup":
			route.ContinuousPickup = v
		case "continuous_drop_off":
			route.ContinuousDropOff = v
		case "network_id":
			route.NetworkId = v
		}
	}

	return &route
}

func ValidateRoute(r Route) []Result {
	var validationResults []Result

	fields := []struct {
		fieldType FieldType
		name      string
		value     *string
		required  bool
	}{
		{FieldTypeID, "route_id", r.Id, true},
		{FieldTypeID, "agency_id", r.AgencyId, false},
		{FieldTypeText, "route_short_name", r.ShortName, false},
		{FieldTypeText, "route_long_name", r.LongName, false},
		{FieldTypeText, "route_desc", r.Desc, false},
		{FieldTypeRouteType, "route_type", r.Type, true},
		{FieldTypeURL, "route_url", r.URL, false},
		{FieldTypeColor, "route_color", r.Color, false},
		{FieldTypeColor, "route_text_color", r.TextColor, false},
		{FieldTypeInteger, "route_sort_order", r.SortOrder, false},
		{FieldTypeContinuousPickup, "continuous_pickup", r.ContinuousPickup, false},
		{FieldTypeContinuousDropOff, "continuous_drop_off", r.ContinuousDropOff, false},
		{FieldTypeID, "network_id", r.NetworkId, false},
	}

	for _, field := range fields {
		validationResults = append(validationResults, validateField(field.fieldType, field.name, field.value, field.required, FileNameRoutes, r.LineNumber)...)
	}

	if StringIsNilOrEmpty(r.ShortName) && StringIsNilOrEmpty(r.LongName) {
		validationResults = append(validationResults, MissingRouteShortNameWhenLongNameIsNotPresentResult{SingleLineResult{
			FileName:  FileNameRoutes,
			FieldName: "route_short_name",
			Line:      r.LineNumber,
		}})
		validationResults = append(validationResults, MissingRouteLongNameWhenShortNameIsNotPresentResult{SingleLineResult{
			FileName:  FileNameRoutes,
			FieldName: "route_long_name",
			Line:      r.LineNumber,
		}})
	}

	if r.ShortName != nil && len(*r.ShortName) >= 12 {
		validationResults = append(validationResults, TooLongRouteShortNameResult{SingleLineResult{
			FileName:  FileNameRoutes,
			FieldName: "route_short_name",
			Line:      r.LineNumber,
		}})
	}

	if r.Desc != nil && r.ShortName != nil && *r.Desc == *r.ShortName {
		validationResults = append(validationResults, DescriptionDuplicatesRouteNameResult{SingleLineResult{
			FileName:  FileNameRoutes,
			FieldName: "route_desc",
			Line:      r.LineNumber,
		}, "route_short_name"})
	}

	if r.Desc != nil && r.LongName != nil && *r.Desc == *r.LongName {
		validationResults = append(validationResults, DescriptionDuplicatesRouteNameResult{SingleLineResult{
			FileName:  FileNameRoutes,
			FieldName: "route_desc",
			Line:      r.LineNumber,
		}, "route_long_name"})
	}

	return validationResults
}

func ValidateRoutes(routes []*Route, agencies []*Agency) []Result {
	var validationResults []Result

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

		results := ValidateRoute(*route)
		if len(results) > 0 {
			validationResults = append(validationResults, results...)
		}

		// agency_id is required only if there are multiple agencies in agencies.txt, recommended otherwise.
		if numAgencies > 1 && StringIsNilOrEmpty(route.AgencyId) {
			validationResults = append(validationResults, AgencyIdRequiredForRouteWhenMultipleAgenciesResult{SingleLineResult{
				FileName:  FileNameRoutes,
				FieldName: "agency_id",
				Line:      route.LineNumber,
			}})
		} else if StringIsNilOrEmpty(route.AgencyId) {
			validationResults = append(validationResults, AgencyIdRecommendedForRouteResult{SingleLineResult{
				FileName:  FileNameRoutes,
				FieldName: "agency_id",
				Line:      route.LineNumber,
			}})
		}

		// route_id must be unique within the routes.txt file
		if usedIds[*route.Id] {
			validationResults = append(validationResults, FieldIsNotUniqueResult{SingleLineResult{
				FileName:  FileNameRoutes,
				FieldName: "route_id",
				Line:      route.LineNumber,
			}})
		} else {
			usedIds[*route.Id] = true
		}

		if agencies == nil || StringIsNilOrEmpty(route.AgencyId) {
			continue
		}

		// agency_id must be a valid agency_id from agencies.txt
		matchingAgencyFound := false
		for _, agency := range agencies {
			if agency == nil || agency.Id == nil || StringIsNilOrEmpty(route.AgencyId) {
				continue
			}

			if *route.AgencyId == *agency.Id {
				matchingAgencyFound = true
				break
			}
		}

		if !matchingAgencyFound {
			validationResults = append(validationResults, ForeignKeyViolationResult{
				ReferencingFileName:  FileNameRoutes,
				ReferencingFieldName: "agency_id",
				ReferencedFieldName:  FileNameAgency,
				ReferencedFileName:   "agency_id",
				OffendingValue:       *route.AgencyId,
				ReferencedAtRow:      route.LineNumber,
			})
		}
	}

	return validationResults
}
