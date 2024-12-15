package ggtfs

import (
	"encoding/csv"
	"fmt"
	"strconv"
)

type StopExtensions struct {
	MunicipalityId *ID // municipality_id (optional)
}

// Stop struct with fields as strings and optional fields as string pointers.
type Stop struct {
	Id                 ID                  // stop_id
	Code               *Text               // stop_code (optional)
	Name               *Text               // stop_name (optional)
	TTSName            *Text               // tts_stop_name (optional)
	Desc               *Text               // stop_desc (optional)
	Lat                *Latitude           // stop_lat (optional)
	Lon                *Longitude          // stop_lon (optional)
	ZoneId             *ID                 // zone_id (optional)
	Url                *URL                // stop_url (optional)
	LocationType       *StopLocation       // location_type (optional)
	ParentStation      *ID                 // parent_station (optional)
	Timezone           *Timezone           // stop_timezone (optional)
	WheelchairBoarding *WheelchairBoarding // wheelchair_boarding (optional)
	PlatformCode       *Text               // platform_code (optional)
	LevelId            *ID                 // level_id (optional)
	Extensions         *StopExtensions
	LineNumber         int
}

func (s Stop) Validate() []error {
	var validationErrors []error

	// stop_name is handled in the ValidateStops function since it is conditionally required
	// stop_lat is handled in the ValidateStops function since it is conditionally required
	// stop_lon is handled in the ValidateStops function since it is conditionally required
	// parent_station is handled in the ValidateStops function since it is conditionally required

	fields := []struct {
		fieldName string
		field     ValidAndPresentField
	}{
		{"stop_id", &s.Id},
	}
	for _, f := range fields {
		validationErrors = append(validationErrors, validateFieldIsPresentAndValid(f.field, f.fieldName, s.LineNumber, StopsFileName)...)
	}

	if s.Code != nil && !s.Code.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(StopsFileName, s.LineNumber, createInvalidFieldString("stop_code")))
	}
	if s.TTSName != nil && !s.TTSName.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(StopsFileName, s.LineNumber, createInvalidFieldString("tts_stop_name")))
	}
	if s.Desc != nil && !s.Desc.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(StopsFileName, s.LineNumber, createInvalidFieldString("stop_desc")))
	}
	if s.ZoneId != nil && !s.ZoneId.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(StopsFileName, s.LineNumber, createInvalidFieldString("zone_id")))
	}
	if s.Url != nil && !s.Url.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(StopsFileName, s.LineNumber, createInvalidFieldString("stop_url")))
	}
	if s.LocationType != nil && !s.LocationType.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(StopsFileName, s.LineNumber, createInvalidFieldString("location_type")))
	}
	if s.Timezone != nil && !s.Timezone.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(StopsFileName, s.LineNumber, createInvalidFieldString("stop_timezone")))
	}
	if s.WheelchairBoarding != nil && !s.WheelchairBoarding.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(StopsFileName, s.LineNumber, createInvalidFieldString("wheelchair_boarding")))
	}
	if s.LevelId != nil && !s.LevelId.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(StopsFileName, s.LineNumber, createInvalidFieldString("level_id")))
	}
	if s.PlatformCode != nil && !s.PlatformCode.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(StopsFileName, s.LineNumber, createInvalidFieldString("platform_code")))
	}
	if s.Extensions.MunicipalityId != nil && !s.Extensions.MunicipalityId.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(StopsFileName, s.LineNumber, createInvalidFieldString("municipality_id")))
	}

	return validationErrors
}

// ValidStopHeaders defines the headers expected in the stops CSV file.
var validStopHeaders = []string{"stop_id", "stop_code", "stop_name", "stop_desc", "stop_lat", "stop_lon", "zone_id",
	"stop_url", "location_type", "parent_station", "stop_timezone", "wheelchair_boarding", "level_id", "platform_code", "municipality_id"}

// LoadStops loads stops from a CSV reader and returns them along with any errors.
func LoadStops(csvReader *csv.Reader) ([]*Stop, []error) {
	entities, errs := loadEntities(csvReader, validStopHeaders, CreateStop, StopsFileName)

	stops := make([]*Stop, 0, len(entities))
	for _, entity := range entities {
		if stop, ok := entity.(*Stop); ok {
			stops = append(stops, stop)
		}
	}

	return stops, errs
}

// CreateStop creates and validates a Stop instance from the CSV row data.
func CreateStop(row []string, headers map[string]int, lineNumber int) interface{} {
	var parseErrors []error

	stop := Stop{
		LineNumber: lineNumber,
	}

	for hName, hPos := range headers {
		switch hName {
		case "stop_id":
			stop.Id = NewID(getRowValue(row, hPos))
		case "stop_code":
			stop.Code = NewOptionalText(getRowValue(row, hPos))
		case "stop_name":
			stop.Name = NewOptionalText(getRowValue(row, hPos))
		case "tts_stop_name":
			stop.TTSName = NewOptionalText(getRowValue(row, hPos))
		case "stop_desc":
			stop.Desc = NewOptionalText(getRowValue(row, hPos))
		case "stop_lat":
			stop.Lat = NewOptionalLatitude(getRowValue(row, hPos))
		case "stop_lon":
			stop.Lon = NewOptionalLongitude(getRowValue(row, hPos))
		case "zone_id":
			stop.ZoneId = NewOptionalID(getRowValue(row, hPos))
		case "stop_url":
			stop.Url = NewOptionalURL(getRowValue(row, hPos))
		case "location_type":
			stop.LocationType = NewOptionalStopLocation(getRowValue(row, hPos))
		case "parent_station":
			stop.ParentStation = NewOptionalID(getRowValue(row, hPos))
		case "stop_timezone":
			stop.Timezone = NewOptionalTimezone(getRowValue(row, hPos))
		case "wheelchair_boarding":
			stop.WheelchairBoarding = NewOptionalWheelchairBoarding(getRowValue(row, hPos))
		case "level_id":
			stop.LevelId = NewOptionalID(getRowValue(row, hPos))
		case "platform_code":
			stop.PlatformCode = NewOptionalText(getRowValue(row, hPos))
		case "municipality_id":
			stop.Extensions = &StopExtensions{
				MunicipalityId: NewOptionalID(getRowValue(row, hPos)),
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

	for _, stop := range stops {
		// Additional required field checks for individual Stop.
		if stop.Id.String() == "" {
			validationErrors = append(validationErrors, createFileRowError(StopsFileName, stop.LineNumber, "stop_id must be specified"))
		}
		if stop.Name == nil && stop.LocationType != nil && stop.LocationType.Int() <= 3 {
			validationErrors = append(validationErrors, createFileRowError(StopsFileName, stop.LineNumber, "stop_name must be specified for location types 0, 1, and 2"))
		}
		if stop.Lat == nil && stop.LocationType != nil && stop.LocationType.Int() <= 3 {
			validationErrors = append(validationErrors, createFileRowError(StopsFileName, stop.LineNumber, "stop_lat must be specified for location types 0, 1, and 2"))
		}
		if stop.Lon == nil && stop.LocationType != nil && stop.LocationType.Int() <= 3 {
			validationErrors = append(validationErrors, createFileRowError(StopsFileName, stop.LineNumber, "stop_lon must be specified for location types 0, 1, and 2"))
		}
		if stop.ParentStation == nil && stop.LocationType != nil && stop.LocationType.Int() >= 2 && stop.LocationType.Int() <= 4 {
			validationErrors = append(validationErrors, createFileRowError(StopsFileName, stop.LineNumber, "parent_station must be specified for location types 2, 3, and 4"))
		}
	}

	// Check for unique stop_id values.
	usedIds := make(map[string]bool)
	for _, stop := range stops {
		if stop == nil {
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

func NewOptionalStopLocation(raw *string) *StopLocation {
	if raw == nil {
		return &StopLocation{
			Integer{base: base{raw: ""}}}
	}
	return &StopLocation{Integer{base: base{raw: *raw}}}
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

func NewOptionalWheelchairBoarding(raw *string) *WheelchairBoarding {
	if raw == nil {
		return &WheelchairBoarding{
			Integer{base: base{raw: ""}}}
	}
	return &WheelchairBoarding{Integer{base: base{raw: *raw}}}
}
