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

	//tripWarnings, tripRecommendations := ggtfs.ValidateTrips(ctx.Trips, ctx.Routes, ctx.CalendarItems, ctx.Shapes)
	//warnings = append(warnings, tripWarnings...)
	//recommendations = append(recommendations, tripRecommendations...)
	//
	//shapeWarnings, shapeRecommendations := ggtfs.ValidateShapes(ctx.Shapes)
	//warnings = append(warnings, shapeWarnings...)
	//recommendations = append(recommendations, shapeRecommendations...)
	//
	//calendarDateWarnings, calendarDateRecommendations := ggtfs.ValidateCalendarDates(ctx.CalendarDates, ctx.CalendarItems)
	//warnings = append(warnings, calendarDateWarnings...)
	//recommendations = append(recommendations, calendarDateRecommendations...)
	//
	//routeWarnings, routeRecommendations := ggtfs.ValidateRoutes(ctx.Routes, ctx.Agencies)
	//warnings = append(warnings, routeWarnings...)
	//recommendations = append(recommendations, routeRecommendations...)
	//
	//stopTimeWarnings, stopTimeRecommendations := ggtfs.ValidateStopTimes(ctx.StopTimes, ctx.Stops)
	//warnings = append(warnings, stopTimeWarnings...)
	//recommendations = append(recommendations, stopTimeRecommendations...)

	return warnings, recommendations
}

func NewGTFSContextForDirectory(gtfsPath string) (*GTFSContext, []error) {
	errs := make([]error, 0)
	context := GTFSContext{}
	var gtfsErrors []error

	files := []string{ggtfs.FileNameAgency, ggtfs.FileNameRoutes, ggtfs.FileNameStops, ggtfs.FileNameTrips, ggtfs.FileNameStopTimes,
		ggtfs.FileNameCalendar, ggtfs.FileNameCalendarDate, ggtfs.FileNameShapes}

	for _, file := range files {
		reader, err := CreateCSVReaderForFile(path.Join(gtfsPath, file))
		if err != nil {
			errs = append(errs, err)
		}

		switch file {
		case ggtfs.FileNameAgency:
			context.Agencies, gtfsErrors = ggtfs.LoadAgencies(ggtfs.NewReader(reader))
		case ggtfs.FileNameRoutes:
			context.Routes, gtfsErrors = ggtfs.LoadRoutes(ggtfs.NewReader(reader))
		case ggtfs.FileNameStops:
			context.Stops, gtfsErrors = ggtfs.LoadStops(ggtfs.NewReader(reader))
		case ggtfs.FileNameTrips:
			context.Trips, gtfsErrors = ggtfs.LoadTrips(ggtfs.NewReader(reader))
		case ggtfs.FileNameStopTimes:
			context.StopTimes, gtfsErrors = ggtfs.LoadStopTimes(ggtfs.NewReader(reader))
		case ggtfs.FileNameCalendar:
			context.CalendarItems, gtfsErrors = ggtfs.LoadCalendar(ggtfs.NewReader(reader))
		case ggtfs.FileNameCalendarDate:
			context.CalendarDates, gtfsErrors = ggtfs.LoadCalendarDates(ggtfs.NewReader(reader))
		case ggtfs.FileNameShapes:
			context.Shapes, gtfsErrors = ggtfs.LoadShapes(ggtfs.NewReader(reader))
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
