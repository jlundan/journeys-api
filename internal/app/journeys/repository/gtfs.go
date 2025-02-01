package repository

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

func newGTFSBundle(gtfsPath string, skipValidation bool) *GTFSBundle {
	bundle := GTFSBundle{}

	files := []string{ggtfs.FileNameAgency, ggtfs.FileNameRoutes, ggtfs.FileNameStops, ggtfs.FileNameTrips, ggtfs.FileNameStopTimes,
		ggtfs.FileNameCalendar, ggtfs.FileNameCalendarDate, ggtfs.FileNameShapes, "municipalities.txt"}

	for _, file := range files {
		reader, err := createCSVReaderForFile(path.Join(gtfsPath, file))
		if err != nil {
			bundle.Errors = append(bundle.Errors, err)
		}

		var gtfsErrors []error
		var municipalityError error

		switch file {
		case ggtfs.FileNameAgency:
			bundle.Agencies, gtfsErrors = ggtfs.LoadAgencies(ggtfs.NewReader(reader))
		case ggtfs.FileNameRoutes:
			bundle.Routes, gtfsErrors = ggtfs.LoadRoutes(ggtfs.NewReader(reader))
		case ggtfs.FileNameStops:
			bundle.Stops, gtfsErrors = ggtfs.LoadStops(ggtfs.NewReader(reader))
		case ggtfs.FileNameTrips:
			bundle.Trips, gtfsErrors = ggtfs.LoadTrips(ggtfs.NewReader(reader))
		case ggtfs.FileNameStopTimes:
			bundle.StopTimes, gtfsErrors = ggtfs.LoadStopTimes(ggtfs.NewReader(reader))
		case ggtfs.FileNameCalendar:
			bundle.CalendarItems, gtfsErrors = ggtfs.LoadCalendar(ggtfs.NewReader(reader))
		case ggtfs.FileNameCalendarDate:
			bundle.CalendarDates, gtfsErrors = ggtfs.LoadCalendarDates(ggtfs.NewReader(reader))
		case ggtfs.FileNameShapes:
			bundle.Shapes, gtfsErrors = ggtfs.LoadShapes(ggtfs.NewReader(reader))
		case "municipalities.txt":
			bundle.Municipalities, municipalityError = readMunicipalities(gtfsPath)
		}

		bundle.Errors = append(bundle.Errors, gtfsErrors...)
		if municipalityError != nil {
			bundle.Errors = append(bundle.Errors, municipalityError)
		}
	}

	if !skipValidation {
		bundle.ValidationNotices = append(bundle.ValidationNotices, ggtfs.ValidateTrips(bundle.Trips, bundle.Routes, bundle.CalendarItems, bundle.Shapes)...)
		bundle.ValidationNotices = append(bundle.ValidationNotices, ggtfs.ValidateShapes(bundle.Shapes)...)
		bundle.ValidationNotices = append(bundle.ValidationNotices, ggtfs.ValidateCalendarDates(bundle.CalendarDates, bundle.CalendarItems)...)
		bundle.ValidationNotices = append(bundle.ValidationNotices, ggtfs.ValidateRoutes(bundle.Routes, bundle.Agencies)...)
		bundle.ValidationNotices = append(bundle.ValidationNotices, ggtfs.ValidateStopTimes(bundle.StopTimes, bundle.Stops)...)
	}

	return &bundle
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

type GTFSBundle struct {
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

const MunicipalityFileName = "municipalities.txt"

type municipalityData struct {
	municipalityHeaders map[string]uint8
	municipalityRows    [][]string
}

func readMunicipalities(gtfsPath string) (*municipalityData, error) {
	var err error
	m := &municipalityData{}
	m.municipalityHeaders, m.municipalityRows, err = parseFile(fmt.Sprintf("%v/%v", gtfsPath, MunicipalityFileName), true)
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
