package ggtfs

import (
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
)

type Stop struct {
	Id                 string
	Code               *string
	Name               *string
	Desc               *string
	Lat                *float64
	Lon                *float64
	ZoneId             *string
	Url                *string
	LocationType       *int
	ParentStation      *string
	Timezone           *string
	WheelchairBoarding *int
	PlatformCode       *string
	LevelId            *string
	MunicipalityId     *string
	lineNumber         int
}

func LoadStops(csvReader *csv.Reader) ([]*Stop, []error) {
	stops := make([]*Stop, 0)
	errs := make([]error, 0)

	headers, err := ReadHeaderRow(csvReader)
	if err != nil {
		errs = append(errs, createFileError(StopsFileName, fmt.Sprintf("read error: %v", err.Error())))
		return stops, errs
	}
	if headers == nil {
		return stops, errs
	}

	usedIds := make([]string, 0)
	index := 0
	for {
		row, err := ReadDataRow(csvReader)
		if err != nil {
			errs = append(errs, createFileError(StopsFileName, fmt.Sprintf("%v", err.Error())))
			index++
			continue
		}
		if row == nil {
			break
		}

		rowErrs := make([]error, 0)
		stop := Stop{
			lineNumber: index,
		}

		var stopId *string
		for name, column := range headers {
			switch name {
			case "stop_id":
				stopId = handleIDField(row[column], StopsFileName, name, index, &rowErrs)
			case "stop_code":
				stop.Code = handleTextField(row[column], StopsFileName, name, index, &rowErrs)
			case "stop_name":
				stop.Name = handleTextField(row[column], StopsFileName, name, index, &rowErrs)
			case "stop_desc":
				stop.Desc = handleTextField(row[column], StopsFileName, name, index, &rowErrs)
			case "stop_lat":
				stop.Lat = handleFloat64Field(row[column], StopsFileName, name, index, &rowErrs)
			case "stop_lon":
				stop.Lon = handleFloat64Field(row[column], StopsFileName, name, index, &rowErrs)
			case "zone_id":
				stop.ZoneId = &row[column]
				//stop.ZoneId = handleIDField(row[column], StopsFileName, name, index, &rowErrs)
			case "stop_url":
				stop.Url = handleURLField(row[column], StopsFileName, name, index, &rowErrs)
			case "location_type":
				stop.LocationType = handleLocationTypeField(row[column], StopsFileName, name, index, &rowErrs)
			case "parent_station":
				stop.ParentStation = handleIDField(row[column], StopsFileName, name, index, &rowErrs)
			case "stop_timezone":
				stop.Timezone = handleTimeZoneField(row[column], StopsFileName, name, index, &rowErrs)
			case "wheelchair_boarding":
				stop.WheelchairBoarding = handleWheelchairBoardingField(row[column], StopsFileName, name, index, &rowErrs)
			case "level_id":
				stop.LevelId = handleIDField(row[column], StopsFileName, name, index, &rowErrs)
			case "platform_code":
				stop.PlatformCode = handleTextField(row[column], StopsFileName, name, index, &rowErrs)
			case "municipality_id":
				stop.MunicipalityId = handleIDField(row[column], StopsFileName, name, index, &rowErrs)
			}
		}

		if stopId == nil {
			rowErrs = append(rowErrs, createFileRowError(StopsFileName, stop.lineNumber, "stop_id must be specified"))
		} else {
			stop.Id = *stopId
			if StringArrayContainsItem(usedIds, *stopId) {
				errs = append(errs, createFileRowError(StopsFileName, index, fmt.Sprintf("%s: stop_id", nonUniqueId)))
			} else {
				usedIds = append(usedIds, *stopId)
			}

		}

		if stop.Name == nil && stop.LocationType != nil && *stop.LocationType <= 3 {
			rowErrs = append(rowErrs, createFileRowError(StopsFileName, stop.lineNumber, "stop_name must be specified for location types 0,1 and 2"))
		}

		if stop.Lat == nil && stop.LocationType != nil && *stop.LocationType <= 3 {
			rowErrs = append(rowErrs, createFileRowError(StopsFileName, stop.lineNumber, "stop_lat must be specified for location types 0,1 and 2"))
		}

		if stop.Lon == nil && stop.LocationType != nil && *stop.LocationType <= 3 {
			rowErrs = append(rowErrs, createFileRowError(StopsFileName, stop.lineNumber, "stop_lon must be specified for location types 0,1 and 2"))
		}

		if stop.ParentStation == nil && stop.LocationType != nil && *stop.LocationType >= 2 && *stop.LocationType <= 4 {
			rowErrs = append(rowErrs, createFileRowError(StopsFileName, stop.lineNumber, "parent_station must be specified for location types 2,3 and 4"))
		}

		if len(rowErrs) > 0 {
			errs = append(errs, rowErrs...)
		} else {
			stops = append(stops, &stop)
		}

		index++
	}

	return stops, errs
}

func handleLocationTypeField(str string, fileName string, fieldName string, index int, errs *[]error) *int {
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, err))
		return nil
	}

	if lt := int(n); lt <= 4 {
		return &lt
	} else {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, errors.New(invalidValue)))
		return nil
	}
}

func handleWheelchairBoardingField(str string, fileName string, fieldName string, index int, errs *[]error) *int {
	if str == "" {
		vcb := 0
		return &vcb
	}

	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, err))
		return nil
	}

	if wcb := int(n); wcb <= 4 {
		return &wcb
	} else {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, errors.New(invalidValue)))
		return nil
	}
}
