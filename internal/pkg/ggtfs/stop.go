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

func ValidateStop(s Stop) []error {
	var validationErrors []error

	requiredFields := map[string]FieldTobeValidated{
		"stop_id": &s.Id,
	}
	validateRequiredFields(requiredFields, &validationErrors, s.LineNumber, StopsFileName)

	optionalFields := map[string]FieldTobeValidated{
		"stop_code":           &s.Id,
		"stop_name":           &s.Name,
		"tts_stop_name":       &s.TTSName,
		"stop_desc":           &s.Desc,
		"stop_lat":            &s.Lat,
		"stop_lon":            &s.Lon,
		"zone_id":             &s.ZoneId,
		"stop_url":            &s.Url,
		"location_type":       &s.LocationType,
		"parent_station":      &s.ParentStation,
		"stop_timezone":       &s.Timezone,
		"wheelchair_boarding": &s.WheelchairBoarding,
		"level_id":            &s.LevelId,
		"platform_code":       &s.PlatformCode,
	}
	validateOptionalFields(optionalFields, &validationErrors, s.LineNumber, StopsFileName)

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

	return validationErrors
}

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

		vErr := ValidateStop(*stop)
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
