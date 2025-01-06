package tre

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/dimchansky/utfbom"
	"github.com/jlundan/journeys-api/pkg/ggtfs"
	"golang.org/x/text/encoding"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
)

type GTFSContext struct {
	Agencies          []*ggtfs.Agency
	Routes            []*ggtfs.Route
	Stops             []*ggtfs.Stop
	Trips             []*ggtfs.Trip
	StopTimes         []*ggtfs.StopTime
	CalendarItems     []*ggtfs.CalendarItem
	CalendarDates     []*ggtfs.CalendarDate
	Shapes            []*ggtfs.Shape
	Municipalities    *municipalityData
	ValidationNotices []ggtfs.ValidationNotice
	Errors            []error
}

func NewGTFSContextForDirectory(gtfsPath string) *GTFSContext {
	context := GTFSContext{}

	files := []string{ggtfs.FileNameAgency, ggtfs.FileNameRoutes, ggtfs.FileNameStops, ggtfs.FileNameTrips, ggtfs.FileNameStopTimes,
		ggtfs.FileNameCalendar, ggtfs.FileNameCalendarDate, ggtfs.FileNameShapes, "municipalities.txt"}

	for _, file := range files {
		reader, err := createCSVReaderForFile(path.Join(gtfsPath, file))
		if err != nil {
			context.Errors = append(context.Errors, err)
		}

		var gtfsErrors []error
		var municipalityError error

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
		case "municipalities.txt":
			context.Municipalities, municipalityError = readMunicipalities()
		}

		context.Errors = append(context.Errors, gtfsErrors...)
		if municipalityError != nil {
			context.Errors = append(context.Errors, municipalityError)
		}
	}

	context.ValidationNotices = append(context.ValidationNotices, ggtfs.ValidateTrips(context.Trips, context.Routes, context.CalendarItems, context.Shapes)...)
	context.ValidationNotices = append(context.ValidationNotices, ggtfs.ValidateShapes(context.Shapes)...)
	context.ValidationNotices = append(context.ValidationNotices, ggtfs.ValidateCalendarDates(context.CalendarDates, context.CalendarItems)...)
	context.ValidationNotices = append(context.ValidationNotices, ggtfs.ValidateRoutes(context.Routes, context.Agencies)...)
	context.ValidationNotices = append(context.ValidationNotices, ggtfs.ValidateStopTimes(context.StopTimes, context.Stops)...)

	return &context
}

func createCSVReaderForFile(path string) (*csv.Reader, error) {
	csvFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	sr, _ := utfbom.Skip(csvFile)

	filteredReader := newSkippingReader(sr)

	return csv.NewReader(filteredReader), nil
}

func newSkippingReader(r io.Reader) io.Reader {
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

const MunicipalityFileName = "municipalities.txt"

type municipalityData struct {
	municipalityHeaders map[string]uint8
	municipalityRows    [][]string
}

func readMunicipalities() (*municipalityData, error) {
	var err error
	m := &municipalityData{}
	m.municipalityHeaders, m.municipalityRows, err = parseFile(fmt.Sprintf("%v/%v", os.Getenv(GtfsEnvKey), MunicipalityFileName), true)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return m, nil
}

func parseFile(path string, firstLineAsHeaders bool) (map[string]uint8, [][]string, error) {
	return parseFileWithDecoderAndDelimiter(path, firstLineAsHeaders, nil, ',')
}

func parseFileWithDecoderAndDelimiter(path string, firstLineAsHeaders bool, decoder *encoding.Decoder, delimiter rune) (map[string]uint8, [][]string, error) {
	csvFile, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	var r *csv.Reader
	if decoder != nil {
		r = csv.NewReader(decoder.Reader(trimReader{csvFile}))
	} else {
		r = csv.NewReader(trimReader{csvFile})
	}

	if delimiter != ',' {
		r.Comma = delimiter
	}

	var headers = map[string]uint8{}
	var data = make([][]string, 0)

	headersRead := false
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, errors.New(fmt.Sprintf("%v: %v", path, err.Error()))
		}

		if firstLineAsHeaders && !headersRead {
			for i, v := range record {
				headers[v] = uint8(i)
			}
			headersRead = true
		} else {
			var row []string
			for _, v := range record {
				row = append(row, strings.TrimSpace(v))
			}
			data = append(data, row)
		}

	}
	return headers, data, nil
}

var trailingWs = regexp.MustCompile(`\s\n`)

type trimReader struct{ io.Reader }

func (tr trimReader) Read(bs []byte) (int, error) {
	n, err := tr.Reader.Read(bs)
	if err != nil {
		return n, err
	}

	lines := string(bs[:n])
	trimmed := []byte(trailingWs.ReplaceAllString(lines, "\n"))
	copy(bs, trimmed)
	return len(trimmed), nil
}
