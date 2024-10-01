package tre

import (
	"encoding/csv"
	"github.com/dimchansky/utfbom"
	"github.com/jlundan/journeys-api/internal/pkg/ggtfs"
	"os"
	"path"
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

	files := []string{ggtfs.AgenciesFileName, ggtfs.RoutesFileName, ggtfs.StopsFileName, ggtfs.TripsFileName, ggtfs.StopTimesFileName,
		ggtfs.CalendarFileName, ggtfs.CalendarDatesFileName, ggtfs.ShapesFileName}

	for _, file := range files {
		reader, err := CreateCSVReaderForFile(path.Join(gtfsPath, file))
		if err != nil {
			errs = append(errs, err)
		}

		switch file {
		case ggtfs.AgenciesFileName:
			context.Agencies, gtfsErrors = ggtfs.LoadAgencies(reader)
		case ggtfs.RoutesFileName:
			context.Routes, gtfsErrors = ggtfs.LoadRoutes(reader)
		case ggtfs.StopsFileName:
			context.Stops, gtfsErrors = ggtfs.LoadStops(reader)
		case ggtfs.TripsFileName:
			context.Trips, gtfsErrors = ggtfs.LoadTrips(reader)
		case ggtfs.StopTimesFileName:
			context.StopTimes, gtfsErrors = ggtfs.LoadStopTimes(reader)
		case ggtfs.CalendarFileName:
			context.CalendarItems, gtfsErrors = ggtfs.LoadCalendarItems(reader)
		case ggtfs.CalendarDatesFileName:
			context.CalendarDates, gtfsErrors = ggtfs.LoadCalendarDates(reader)
		case ggtfs.ShapesFileName:
			context.Shapes, gtfsErrors = ggtfs.LoadShapes(reader)
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

	return csv.NewReader(sr), nil
}
