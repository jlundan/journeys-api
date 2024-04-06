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

func Validate(ctx *GTFSContext) []error {
	var warnings []error

	tripWarnings := ggtfs.ValidateTrips(ctx.Trips, ctx.Routes, ctx.CalendarItems, ctx.Shapes)
	if len(tripWarnings) > 0 {
		warnings = append(warnings, tripWarnings...)
	}

	shapeWarnings := ggtfs.ValidateShapes(ctx.Shapes)
	if len(shapeWarnings) > 0 {
		warnings = append(warnings, shapeWarnings...)
	}

	calendarDateWarnings := ggtfs.ValidateCalendarDates(ctx.CalendarDates, ctx.CalendarItems)
	if len(calendarDateWarnings) > 0 {
		warnings = append(warnings, calendarDateWarnings...)
	}

	routeWarnings := ggtfs.ValidateRoutes(ctx.Routes, ctx.Agencies)
	if len(routeWarnings) > 0 {
		warnings = append(warnings, routeWarnings...)
	}

	stopTimeWarnings := ggtfs.ValidateStoptimes(ctx.StopTimes, ctx.Stops)
	if len(stopTimeWarnings) > 0 {
		warnings = append(warnings, stopTimeWarnings...)
	}

	return warnings
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
