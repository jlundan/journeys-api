package ggtfs

import (
	"encoding/csv"
	"fmt"
)

// Stop struct with fields as strings and optional fields as string pointers.
type Stop struct {
	Id                 string  // stop_id
	Code               *string // stop_code (optional)
	Name               *string // stop_name (optional)
	Desc               *string // stop_desc (optional)
	Lat                *string // stop_lat (optional)
	Lon                *string // stop_lon (optional)
	ZoneId             *string // zone_id (optional)
	Url                *string // stop_url (optional)
	LocationType       *string // location_type (optional)
	ParentStation      *string // parent_station (optional)
	Timezone           *string // stop_timezone (optional)
	WheelchairBoarding *string // wheelchair_boarding (optional)
	PlatformCode       *string // platform_code (optional)
	LevelId            *string // level_id (optional)
	MunicipalityId     *string // municipality_id (optional)
	LineNumber         int     // Line number in the CSV file for error reporting
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
func CreateStop(row []string, headers map[string]uint8, lineNumber int) (interface{}, []error) {
	var parseErrors []error

	stop := Stop{
		LineNumber: lineNumber,
	}

	for hName, hPos := range headers {
		switch hName {
		case "stop_id":
			stop.Id = getField(row, hName, hPos, &parseErrors, lineNumber, StopsFileName)
		case "stop_code":
			stop.Code = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, StopsFileName)
		case "stop_name":
			stop.Name = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, StopsFileName)
		case "stop_desc":
			stop.Desc = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, StopsFileName)
		case "stop_lat":
			stop.Lat = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, StopsFileName)
		case "stop_lon":
			stop.Lon = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, StopsFileName)
		case "zone_id":
			stop.ZoneId = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, StopsFileName)
		case "stop_url":
			stop.Url = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, StopsFileName)
		case "location_type":
			stop.LocationType = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, StopsFileName)
		case "parent_station":
			stop.ParentStation = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, StopsFileName)
		case "stop_timezone":
			stop.Timezone = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, StopsFileName)
		case "wheelchair_boarding":
			stop.WheelchairBoarding = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, StopsFileName)
		case "level_id":
			stop.LevelId = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, StopsFileName)
		case "platform_code":
			stop.PlatformCode = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, StopsFileName)
		case "municipality_id":
			stop.MunicipalityId = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, StopsFileName)
		}
	}

	if len(parseErrors) > 0 {
		return &stop, parseErrors
	}
	return &stop, nil
}

// ValidateStops performs additional validation for a list of Stop instances.
func ValidateStops(stops []*Stop) []error {
	var validationErrors []error

	if stops == nil {
		return validationErrors
	}

	for _, stop := range stops {
		// Additional required field checks for individual Stop.
		if stop.Id == "" {
			validationErrors = append(validationErrors, createFileRowError(StopsFileName, stop.LineNumber, "stop_id must be specified"))
		}
		if stop.Name == nil && stop.LocationType != nil && *stop.LocationType <= "3" {
			validationErrors = append(validationErrors, createFileRowError(StopsFileName, stop.LineNumber, "stop_name must be specified for location types 0, 1, and 2"))
		}
		if stop.Lat == nil && stop.LocationType != nil && *stop.LocationType <= "3" {
			validationErrors = append(validationErrors, createFileRowError(StopsFileName, stop.LineNumber, "stop_lat must be specified for location types 0, 1, and 2"))
		}
		if stop.Lon == nil && stop.LocationType != nil && *stop.LocationType <= "3" {
			validationErrors = append(validationErrors, createFileRowError(StopsFileName, stop.LineNumber, "stop_lon must be specified for location types 0, 1, and 2"))
		}
		if stop.ParentStation == nil && stop.LocationType != nil && *stop.LocationType >= "2" && *stop.LocationType <= "4" {
			validationErrors = append(validationErrors, createFileRowError(StopsFileName, stop.LineNumber, "parent_station must be specified for location types 2, 3, and 4"))
		}
	}

	// Check for unique stop_id values.
	usedIds := make(map[string]bool)
	for _, stop := range stops {
		if stop == nil {
			continue
		}
		if usedIds[stop.Id] {
			validationErrors = append(validationErrors, createFileRowError(StopsFileName, stop.LineNumber, fmt.Sprintf("stop_id '%s' is not unique within the file", stop.Id)))
		} else {
			usedIds[stop.Id] = true
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
//type Stop struct {
//	Id                 string
//	Code               *string
//	Name               *string
//	Desc               *string
//	Lat                *float64
//	Lon                *float64
//	ZoneId             *string
//	Url                *string
//	LocationType       *int
//	ParentStation      *string
//	Timezone           *string
//	WheelchairBoarding *int
//	PlatformCode       *string
//	LevelId            *string
//	MunicipalityId     *string
//	lineNumber         int
//}
//
//var validStopHeaders = []string{"stop_id", "stop_code", "stop_name", "tts_stop_name", "stop_desc", "stop_lat",
//	"stop_lon", "zone_id", "stop_url", "location_type", "parent_station", "stop_timezone", "wheelchair_boarding",
//	"level_id", "platform_code", "municipality_id"}
//
//func LoadStops(csvReader *csv.Reader) ([]*Stop, []error) {
//	stops := make([]*Stop, 0)
//	errs := make([]error, 0)
//
//	headers, err := ReadHeaderRow(csvReader, validStopHeaders)
//	if err != nil {
//		errs = append(errs, createFileError(StopsFileName, fmt.Sprintf("read error: %v", err.Error())))
//		return stops, errs
//	}
//	if headers == nil {
//		return stops, errs
//	}
//
//	usedIds := make([]string, 0)
//	index := 0
//	for {
//		row, err := ReadDataRow(csvReader)
//		if err != nil {
//			errs = append(errs, createFileError(StopsFileName, fmt.Sprintf("%v", err.Error())))
//			index++
//			continue
//		}
//		if row == nil {
//			break
//		}
//
//		rowErrs := make([]error, 0)
//		stop := Stop{
//			lineNumber: index,
//		}
//
//		var stopId *string
//		for name, column := range headers {
//			switch name {
//			case "stop_id":
//				stopId = handleIDField(row[column], StopsFileName, name, index, &rowErrs)
//			case "stop_code":
//				stop.Code = handleTextField(row[column], StopsFileName, name, index, &rowErrs)
//			case "stop_name":
//				stop.Name = handleTextField(row[column], StopsFileName, name, index, &rowErrs)
//			case "stop_desc":
//				stop.Desc = handleTextField(row[column], StopsFileName, name, index, &rowErrs)
//			case "stop_lat":
//				stop.Lat = handleFloat64Field(row[column], StopsFileName, name, index, &rowErrs)
//			case "stop_lon":
//				stop.Lon = handleFloat64Field(row[column], StopsFileName, name, index, &rowErrs)
//			case "zone_id":
//				stop.ZoneId = handleIDField(row[column], StopsFileName, name, index, &rowErrs)
//			case "stop_url":
//				stop.Url = handleURLField(row[column], StopsFileName, name, index, &rowErrs)
//			case "location_type":
//				stop.LocationType = handleLocationTypeField(row[column], StopsFileName, name, index, &rowErrs)
//			case "parent_station":
//				stop.ParentStation = handleIDField(row[column], StopsFileName, name, index, &rowErrs)
//			case "stop_timezone":
//				stop.Timezone = handleTimeZoneField(row[column], StopsFileName, name, index, &rowErrs)
//			case "wheelchair_boarding":
//				stop.WheelchairBoarding = handleWheelchairBoardingField(row[column], StopsFileName, name, index, &rowErrs)
//			case "level_id":
//				stop.LevelId = handleIDField(row[column], StopsFileName, name, index, &rowErrs)
//			case "platform_code":
//				stop.PlatformCode = handleTextField(row[column], StopsFileName, name, index, &rowErrs)
//			case "municipality_id":
//				stop.MunicipalityId = handleIDField(row[column], StopsFileName, name, index, &rowErrs)
//			}
//		}
//
//		if stopId == nil {
//			rowErrs = append(rowErrs, createFileRowError(StopsFileName, stop.lineNumber, "stop_id must be specified"))
//		} else {
//			stop.Id = *stopId
//			if StringArrayContainsItem(usedIds, *stopId) {
//				errs = append(errs, createFileRowError(StopsFileName, index, fmt.Sprintf("%s: stop_id", nonUniqueId)))
//			} else {
//				usedIds = append(usedIds, *stopId)
//			}
//
//		}
//
//		if stop.Name == nil && stop.LocationType != nil && *stop.LocationType <= 3 {
//			rowErrs = append(rowErrs, createFileRowError(StopsFileName, stop.lineNumber, "stop_name must be specified for location types 0,1 and 2"))
//		}
//
//		if stop.Lat == nil && stop.LocationType != nil && *stop.LocationType <= 3 {
//			rowErrs = append(rowErrs, createFileRowError(StopsFileName, stop.lineNumber, "stop_lat must be specified for location types 0,1 and 2"))
//		}
//
//		if stop.Lon == nil && stop.LocationType != nil && *stop.LocationType <= 3 {
//			rowErrs = append(rowErrs, createFileRowError(StopsFileName, stop.lineNumber, "stop_lon must be specified for location types 0,1 and 2"))
//		}
//
//		if stop.ParentStation == nil && stop.LocationType != nil && *stop.LocationType >= 2 && *stop.LocationType <= 4 {
//			rowErrs = append(rowErrs, createFileRowError(StopsFileName, stop.lineNumber, "parent_station must be specified for location types 2,3 and 4"))
//		}
//
//		if len(rowErrs) > 0 {
//			errs = append(errs, rowErrs...)
//		} else {
//			stops = append(stops, &stop)
//		}
//
//		index++
//	}
//
//	return stops, errs
//}
//
//func handleLocationTypeField(str string, fileName string, fieldName string, index int, errs *[]error) *int {
//	n, err := strconv.ParseInt(str, 10, 64)
//	if err != nil {
//		*errs = append(*errs, createFieldError(fileName, fieldName, index, err))
//		return nil
//	}
//
//	if lt := int(n); lt <= 4 {
//		return &lt
//	} else {
//		*errs = append(*errs, createFieldError(fileName, fieldName, index, errors.New(invalidValue)))
//		return nil
//	}
//}
//
//func handleWheelchairBoardingField(str string, fileName string, fieldName string, index int, errs *[]error) *int {
//	if str == "" {
//		vcb := 0
//		return &vcb
//	}
//
//	n, err := strconv.ParseInt(str, 10, 64)
//	if err != nil {
//		*errs = append(*errs, createFieldError(fileName, fieldName, index, err))
//		return nil
//	}
//
//	if wcb := int(n); wcb <= 4 {
//		return &wcb
//	} else {
//		*errs = append(*errs, createFieldError(fileName, fieldName, index, errors.New(invalidValue)))
//		return nil
//	}
//}
