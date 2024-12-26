package ggtfs

import (
	"fmt"
	"strconv"
)

type StopExtensions struct {
	MunicipalityId ID // municipality_id (optional)
}

type Stop struct {
	Id                 ID                 // stop_id 			 (required)
	Code               Text               // stop_code 			 (optional)
	Name               Text               // stop_name 			 (conditionally required)
	TTSName            Text               // tts_stop_name 		 (optional)
	Desc               Text               // stop_desc, 		 (optional)
	Lat                Latitude           // stop_lat 			 (conditionally required)
	Lon                Longitude          // stop_lon 			 (conditionally required)
	ZoneId             ID                 // zone_id 			 (optional)
	Url                URL                // stop_url 			 (optional)
	LocationType       StopLocation       // location_type 		 (optional)
	ParentStation      ID                 // parent_station 	 (conditionally required)
	Timezone           Timezone           // stop_timezone 		 (optional)
	WheelchairBoarding WheelchairBoarding // wheelchair_boarding (optional)
	PlatformCode       Text               // platform_code 		 (optional)
	LevelId            ID                 // level_id 			 (optional)
	Extensions         *StopExtensions
	LineNumber         int
}

func (s Stop) Validate() []error {
	var validationErrors []error

	requiredFields := []struct {
		fieldName string
		field     ValidAndPresentField
	}{
		{"stop_id", &s.Id},
	}
	for _, f := range requiredFields {
		if !f.field.IsValid() {
			validationErrors = append(validationErrors, createFileRowError(StopsFileName, s.LineNumber, createInvalidRequiredFieldString(f.fieldName)))
		}
	}

	optionalFields := []struct {
		field     ValidAndPresentField
		fieldName string
	}{
		{&s.Id, "stop_code"},
		{&s.Name, "stop_name"},
		{&s.TTSName, "tts_stop_name"},
		{&s.Desc, "stop_desc"},
		{&s.Lat, "stop_lat"},
		{&s.Lon, "stop_lon"},
		{&s.ZoneId, "zone_id"},
		{&s.Url, "stop_url"},
		{&s.LocationType, "location_type"},
		{&s.ParentStation, "parent_station"},
		{&s.Timezone, "stop_timezone"},
		{&s.WheelchairBoarding, "wheelchair_boarding"},
		{&s.LevelId, "level_id"},
		{&s.PlatformCode, "platform_code"},
	}

	for _, field := range optionalFields {
		if field.field != nil && field.field.IsPresent() && !field.field.IsValid() {
			validationErrors = append(validationErrors, createFileRowError(StopsFileName, s.LineNumber, createInvalidFieldString(field.fieldName)))
		}
	}

	if s.LocationType.IsValid() && s.LocationType.Int() <= 3 && !s.Name.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(StopsFileName, s.LineNumber, "stop_name must be specified for location types 0, 1, and 2"))
	}
	if s.LocationType.IsValid() && s.LocationType.Int() <= 3 && !s.Lat.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(StopsFileName, s.LineNumber, "stop_lat must be specified for location types 0, 1, and 2"))
	}
	if s.LocationType.IsValid() && s.LocationType.Int() <= 3 && !s.Lon.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(StopsFileName, s.LineNumber, "stop_lon must be specified for location types 0, 1, and 2"))
	}
	if s.LocationType.IsValid() && s.LocationType.Int() >= 2 && s.LocationType.Int() <= 4 && !s.ParentStation.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(StopsFileName, s.LineNumber, "parent_station must be specified for location types 2, 3, and 4"))
	}

	// TODO: VALIDATION: stop_desc: Should not be a duplicate of stop_name
	// TODO: VALIDATION: zone_id: If this record represents a station or station entrance, the zone_id is ignored
	// TODO: VALIDATION: parent_station: Required for locations which are entrances (location_type=2), generic nodes (location_type=3) or boarding areas (location_type=4).
	// TODO: VALIDATION: parent_station: Optional for stops/platforms (location_type=0)
	// TODO: VALIDATION: parent_station: Forbidden for stations (location_type=1) (this field must be empty)
	// TODO: VALIDATION: parent_station: Stop/platform (location_type=0): the parent_station field contains the ID of a station.
	// TODO: VALIDATION: parent_station: Entrance/exit (location_type=2) or generic node (location_type=3): the parent_station field contains the ID of a station (location_type=1)
	// TODO: VALIDATION: parent_station: Boarding Area (location_type=4): the parent_station field contains ID of a platform
	// TODO: VALIDATION: stop_timezone: If the location has a parent station, it inherits the parent station’s timezone instead of applying its own.
	// TODO: VALIDATION: stop_timezone: Stations and parentless stops with empty stop_timezone inherit the timezone specified by agency.agency_timezone.
	// TODO: VALIDATION: level_id: Foreign ID referencing levels.level_id (must exist)
	// TODO: VALIDATION: platform_code: Words like “platform” or "track" (or the feed’s language-specific equivalent) should not be included.

	return validationErrors
}

// CreateStop creates and validates a Stop instance from the CSV row data.
func CreateStop(row []string, headers map[string]int, lineNumber int) *Stop {
	var parseErrors []error

	stop := Stop{
		LineNumber: lineNumber,
	}

	for hName := range headers {
		v := getRowValueForHeaderName(row, headers, hName)
		switch hName {
		case "stop_id":
			stop.Id = NewID(v)
		case "stop_code":
			stop.Code = NewText(v)
		case "stop_name":
			stop.Name = NewText(v)
		case "tts_stop_name":
			stop.TTSName = NewText(v)
		case "stop_desc":
			stop.Desc = NewText(v)
		case "stop_lat":
			stop.Lat = NewLatitude(v)
		case "stop_lon":
			stop.Lon = NewLongitude(v)
		case "zone_id":
			stop.ZoneId = NewID(v)
		case "stop_url":
			stop.Url = NewURL(v)
		case "location_type":
			stop.LocationType = NewStopLocation(v)
		case "parent_station":
			stop.ParentStation = NewID(v)
		case "stop_timezone":
			stop.Timezone = NewTimezone(v)
		case "wheelchair_boarding":
			stop.WheelchairBoarding = NewWheelchairBoarding(v)
		case "level_id":
			stop.LevelId = NewID(v)
		case "platform_code":
			stop.PlatformCode = NewText(v)
		case "municipality_id":
			stop.Extensions = &StopExtensions{
				MunicipalityId: NewID(v),
			}
		}
	}

	if len(parseErrors) > 0 {
		return &stop
	}
	return &stop
}

// ValidateStops performs additional validation for a list of Stop instances.
func ValidateStops(stops []*Stop) ([]error, []string) {
	var validationErrors []error
	var recommendations []string

	if stops == nil {
		return validationErrors, recommendations
	}

	// Check for unique stop_id values.
	usedIds := make(map[string]bool)
	for _, stop := range stops {
		if stop == nil {
			continue
		}

		vErr := stop.Validate()
		if len(vErr) > 0 {
			validationErrors = append(validationErrors, vErr...)
			continue
		}

		if usedIds[stop.Id.String()] {
			validationErrors = append(validationErrors, createFileRowError(StopsFileName, stop.LineNumber, fmt.Sprintf("stop_id '%s' is not unique within the file", stop.Id.String())))
		} else {
			usedIds[stop.Id.String()] = true
		}
	}

	// TODO: VALIDATION: stop_id: ID must be unique across all stops.stop_id, locations.geojson id, and location_groups.location_group_id values.
	// TODO: VALIDATION: stop_url: This should be different from the agency.agency_url and the routes.route_url field values.

	return validationErrors, recommendations
}

//goland:noinspection GoUnusedConst
const (
	StopLocationTypeStop           = 0
	StopLocationTypeStation        = 1
	StopLocationTypeEntranceOrExit = 2
	StopLocationTypeGenericNode    = 3
	StopLocationTypeBoardingArea   = 4

	ParentlessStopWheelchairBoardingNoInformation = 0
	ParentlessStopWheelchairBoardingSomeVehicles  = 1
	ParentlessStopWheelchairBoardingNotPossible   = 2

	ChildStopWheelchairBoardingInheritFromParent = 0
	ChildStopWheelchairBoardingSomePath          = 1
	ChildStopWheelchairBoardingNoPath            = 2

	StationStopWheelchairBoardingInheritFromParent  = 0
	StationStopWheelchairBoardingEntranceAccessible = 1
	StationStopWheelchairBoardingNoAccessiblePath   = 2
)

type StopLocation struct {
	Integer
}

func (s StopLocation) IsValid() bool {
	val, err := strconv.Atoi(s.Integer.base.raw)
	if err != nil {
		return false
	}

	return val == StopLocationTypeStop || val == StopLocationTypeStation || val == StopLocationTypeEntranceOrExit ||
		val == StopLocationTypeGenericNode || val == StopLocationTypeBoardingArea
}

func NewStopLocation(raw *string) StopLocation {
	if raw == nil {
		return StopLocation{
			Integer{base: base{raw: ""}}}
	}
	return StopLocation{Integer{base: base{raw: *raw, isPresent: true}}}
}

type WheelchairBoarding struct {
	Integer
}

func (wcb WheelchairBoarding) IsValid() bool {
	val, err := strconv.Atoi(wcb.Integer.base.raw)
	if err != nil {
		return false
	}

	return val >= 0 && val <= 2
}

func NewWheelchairBoarding(raw *string) WheelchairBoarding {
	if raw == nil {
		return WheelchairBoarding{
			Integer{base: base{raw: ""}}}
	}
	return WheelchairBoarding{Integer{base: base{raw: *raw, isPresent: true}}}
}
