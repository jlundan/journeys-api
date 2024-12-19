package tre

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"github.com/dimchansky/utfbom"
	"github.com/jlundan/journeys-api/internal/pkg/ggtfs"
	"io"
	"os"
	"path"
	"strings"
)

type GTFSContext struct {
	Agencies      []*ggtfs.Agency
	Routes        []*ggtfs.Route
	Stops         []*ggtfs.Stop
	Trips         []*ggtfs.Trip
	StopTimes     []*ggtfs.StopTime
	CalendarItems []*ggtfs.CalendarItem
	CalendarDates []*ggtfs.CalendarDate
	Shapes        []*ggtfs.Shape
}

func Validate(ctx *GTFSContext) ([]error, []string) {
	var warnings []error
	var recommendations []string

	tripWarnings, tripRecommendations := ggtfs.ValidateTrips(ctx.Trips, ctx.Routes, ctx.CalendarItems, ctx.Shapes)
	warnings = append(warnings, tripWarnings...)
	recommendations = append(recommendations, tripRecommendations...)

	shapeWarnings, shapeRecommendations := ggtfs.ValidateShapes(ctx.Shapes)
	warnings = append(warnings, shapeWarnings...)
	recommendations = append(recommendations, shapeRecommendations...)

	calendarDateWarnings, calendarDateRecommendations := ggtfs.ValidateCalendarDates(ctx.CalendarDates, ctx.CalendarItems)
	warnings = append(warnings, calendarDateWarnings...)
	recommendations = append(recommendations, calendarDateRecommendations...)

	routeWarnings, routeRecommendations := ggtfs.ValidateRoutes(ctx.Routes, ctx.Agencies)
	warnings = append(warnings, routeWarnings...)
	recommendations = append(recommendations, routeRecommendations...)

	stopTimeWarnings, stopTimeRecommendations := ggtfs.ValidateStopTimes(ctx.StopTimes, ctx.Stops)
	warnings = append(warnings, stopTimeWarnings...)
	recommendations = append(recommendations, stopTimeRecommendations...)

	return warnings, recommendations
}

func NewGTFSContextForDirectory(gtfsPath string) (*GTFSContext, []error) {
	errs := make([]error, 0)
	context := GTFSContext{}
	var gtfsErrors []error

	var validShapeHeaders = []string{"shape_id", "shape_pt_lat", "shape_pt_lon", "shape_pt_sequence", "shape_dist_traveled"}
	var validStopHeaders = []string{"stop_id", "stop_code", "stop_name", "stop_desc", "stop_lat", "stop_lon", "zone_id",
		"stop_url", "location_type", "parent_station", "stop_timezone", "wheelchair_boarding", "level_id", "platform_code", "municipality_id"}
	var validAgencyHeaders = []string{"agency_id", "agency_name", "agency_url", "agency_timezone",
		"agency_lang", "agency_phone", "agency_fare_url", "agency_email"}
	var validCalendarHeaders = []string{
		"service_id", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday", "start_date", "end_date",
	}
	var validCalendarDateHeaders = []string{"service_id", "date", "exception_type"}
	var validRouteHeaders = []string{"route_id", "agency_id", "route_short_name", "route_long_name", "route_desc",
		"route_type", "route_url", "route_color", "route_text_color", "route_sort_order", "continuous_pickup",
		"continuous_drop_off", "network_id"}
	var validStopTimeHeaders = []string{"trip_id", "arrival_time", "departure_time", "stop_id", "stop_sequence",
		"stop_headsign", "pickup_type", "drop_off_type", "continuous_pickup", "continuous_drop_off",
		"shape_dist_traveled", "timepoint"}
	var validTripHeaders = []string{"route_id", "service_id", "trip_id", "trip_headsign", "trip_short_name",
		"direction_id", "block_id", "shape_id", "wheelchair_accessible", "bikes_allowed"}

	files := []string{ggtfs.AgenciesFileName, ggtfs.RoutesFileName, ggtfs.StopsFileName, ggtfs.TripsFileName, ggtfs.StopTimesFileName,
		ggtfs.CalendarFileName, ggtfs.CalendarDatesFileName, ggtfs.ShapesFileName}

	for _, file := range files {
		reader, err := CreateCSVReaderForFile(path.Join(gtfsPath, file))
		if err != nil {
			errs = append(errs, err)
		}

		switch file {
		case ggtfs.AgenciesFileName:
			context.Agencies, gtfsErrors = ggtfs.LoadEntities[*ggtfs.Agency](reader, validAgencyHeaders, ggtfs.CreateAgency, ggtfs.AgenciesFileName)
		case ggtfs.RoutesFileName:
			context.Routes, gtfsErrors = ggtfs.LoadEntities[*ggtfs.Route](reader, validRouteHeaders, ggtfs.CreateRoute, ggtfs.RoutesFileName)
		case ggtfs.StopsFileName:
			context.Stops, gtfsErrors = ggtfs.LoadEntities[*ggtfs.Stop](reader, validStopHeaders, ggtfs.CreateStop, ggtfs.StopsFileName)
		case ggtfs.TripsFileName:
			context.Trips, gtfsErrors = ggtfs.LoadEntities[*ggtfs.Trip](reader, validTripHeaders, ggtfs.CreateTrip, ggtfs.TripsFileName)
		case ggtfs.StopTimesFileName:
			context.StopTimes, gtfsErrors = ggtfs.LoadEntities[*ggtfs.StopTime](reader, validStopTimeHeaders, ggtfs.CreateStopTime, ggtfs.StopTimesFileName)
		case ggtfs.CalendarFileName:
			context.CalendarItems, gtfsErrors = ggtfs.LoadEntities[*ggtfs.CalendarItem](reader, validCalendarHeaders, ggtfs.CreateCalendarItem, ggtfs.CalendarFileName)
		case ggtfs.CalendarDatesFileName:
			context.CalendarDates, gtfsErrors = ggtfs.LoadEntities[*ggtfs.CalendarDate](reader, validCalendarDateHeaders, ggtfs.CreateCalendarDate, ggtfs.CalendarDatesFileName)
		case ggtfs.ShapesFileName:
			context.Shapes, gtfsErrors = ggtfs.LoadEntities[*ggtfs.Shape](reader, validShapeHeaders, ggtfs.CreateShape, ggtfs.ShapesFileName)
		}

		if len(gtfsErrors) > 0 {
			errs = append(errs, gtfsErrors...)
		}
	}

	return &context, errs
}

func CreateCSVReaderForFile(path string) (*csv.Reader, error) {
	csvFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	sr, _ := utfbom.Skip(csvFile)

	filteredReader := NewSkippingReader(sr)

	return csv.NewReader(filteredReader), nil
}

func NewSkippingReader(r io.Reader) io.Reader {
	var buf bytes.Buffer

	// Use a bufio.Scanner to read through the input line by line.
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip empty lines and lines that contain only whitespace.
		if strings.TrimSpace(line) == "" {
			continue
		}
		// Write non-empty lines to the buffer.
		buf.WriteString(line + "\n")
	}

	// Return a new reader that reads from the buffer.
	return &buf
}
