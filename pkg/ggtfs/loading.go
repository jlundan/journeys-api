package ggtfs

import (
	"encoding/csv"
	"fmt"
	"io"
)

func NewReader(r *csv.Reader) *GtfsCsvReader {
	return &GtfsCsvReader{
		csvReader:          r,
		FailOnHeaderErrors: true,
		SkipRowsWithErrors: true,
	}
}

func LoadAgencies(reader *GtfsCsvReader) ([]*Agency, []error) {
	return loadCsvEntities[*Agency](defaultAgencyHeaders, reader, CreateAgency)
}

func LoadRoutes(reader *GtfsCsvReader) ([]*Route, []error) {
	return loadCsvEntities[*Route](defaultRouteHeaders, reader, CreateRoute)
}

func LoadStops(reader *GtfsCsvReader) ([]*Stop, []error) {
	return loadCsvEntities[*Stop](defaultStopHeaders, reader, CreateStop)
}

func LoadTrips(reader *GtfsCsvReader) ([]*Trip, []error) {
	return loadCsvEntities[*Trip](defaultTripHeaders, reader, CreateTrip)
}

func LoadStopTimes(reader *GtfsCsvReader) ([]*StopTime, []error) {
	return loadCsvEntities[*StopTime](defaultStopTimeHeaders, reader, CreateStopTime)
}

func LoadCalendar(reader *GtfsCsvReader) ([]*CalendarItem, []error) {
	return loadCsvEntities[*CalendarItem](defaultCalendarHeaders, reader, CreateCalendarItem)
}

func LoadCalendarDates(reader *GtfsCsvReader) ([]*CalendarDate, []error) {
	return loadCsvEntities[*CalendarDate](defaultCalendarDateHeaders, reader, CreateCalendarDate)
}

func LoadShapes(reader *GtfsCsvReader) ([]*Shape, []error) {
	return loadCsvEntities[*Shape](defaultShapeHeaders, reader, CreateShape)
}

func loadCsvEntities[T CsvEntity](headerNames []string, reader *GtfsCsvReader, entityCreator csvEntityCreator[T]) ([]T, []error) {
	var errs []error

	headers, indexingErrors := getHeaderIndex(reader.csvReader, headerNames)

	if len(indexingErrors) > 0 {
		for _, err := range indexingErrors {
			errs = append(errs, fmt.Errorf("line 1: %v", err.Error()))
		}

		if reader.FailOnHeaderErrors {
			return []T{}, errs
		}
	}

	if len(headers) == 0 {
		return []T{}, errs
	}

	var entities []T

	lineNumber := 2
	for {
		row, rErr := reader.csvReader.Read()
		if rErr == io.EOF {
			break
		}

		if rErr != nil {
			errs = append(errs, fmt.Errorf("line %d: %v", lineNumber, rErr.Error()))

			if reader.SkipRowsWithErrors {
				lineNumber++
				continue
			}
		}

		entities = append(entities, entityCreator(row, headers, lineNumber))

		lineNumber++
	}

	return entities, errs
}

type GtfsCsvReader struct {
	csvReader          *csv.Reader
	FailOnHeaderErrors bool
	SkipRowsWithErrors bool
}

var defaultAgencyHeaders = []string{"agency_id", "agency_name", "agency_url", "agency_timezone",
	"agency_lang", "agency_phone", "agency_fare_url", "agency_email"}
var defaultShapeHeaders = []string{"shape_id", "shape_pt_lat", "shape_pt_lon", "shape_pt_sequence", "shape_dist_traveled"}
var defaultStopHeaders = []string{"stop_id", "stop_code", "stop_name", "stop_desc", "stop_lat", "stop_lon", "zone_id",
	"stop_url", "location_type", "parent_station", "stop_timezone", "wheelchair_boarding", "level_id", "platform_code", "municipality_id"}
var defaultCalendarHeaders = []string{
	"service_id", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday", "start_date", "end_date"}
var defaultCalendarDateHeaders = []string{"service_id", "date", "exception_type"}
var defaultRouteHeaders = []string{"route_id", "agency_id", "route_short_name", "route_long_name", "route_desc",
	"route_type", "route_url", "route_color", "route_text_color", "route_sort_order", "continuous_pickup",
	"continuous_drop_off", "network_id"}
var defaultStopTimeHeaders = []string{"trip_id", "arrival_time", "departure_time", "stop_id", "stop_sequence",
	"stop_headsign", "pickup_type", "drop_off_type", "continuous_pickup", "continuous_drop_off",
	"shape_dist_traveled", "timepoint"}
var defaultTripHeaders = []string{"route_id", "service_id", "trip_id", "trip_headsign", "trip_short_name",
	"direction_id", "block_id", "shape_id", "wheelchair_accessible", "bikes_allowed"}

type csvEntityCreator[T CsvEntity] func(row []string, headers map[string]int, lineNumber int) T

type CsvEntity interface {
	*Agency | any
}
