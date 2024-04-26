//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"testing"
)

func TestTripsCSVParsing(t *testing.T) {
	items, errors := LoadTrips(csv.NewReader(strings.NewReader("")))
	if len(errors) > 0 {
		t.Error(errors)
	}
	if len(items) != 0 {
		t.Error("expected zero items")
	}

	reader := csv.NewReader(strings.NewReader("foo,bar\n1,2"))
	reader.Comma = ','
	reader.Comment = ','
	_, errors = LoadTrips(reader)
	if len(errors) == 0 {
		t.Error("expected to throw error")
	}
}

func TestTripParsingOK(t *testing.T) {
	headSign := "Trip 1"
	shortName := "trip1"
	direction := 0
	block := "4"
	shape := "5"
	wca := 1
	ba := 2

	expected1 := Trip{
		Id:                   "1",
		RouteId:              "2",
		ServiceId:            "3",
		HeadSign:             &headSign,
		ShortName:            &shortName,
		DirectionId:          &direction,
		BlockId:              &block,
		ShapeId:              &shape,
		WheelchairAccessible: &wca,
		BikesAllowed:         &ba,
	}

	testCases := []struct {
		rows     [][]string
		expected Trip
	}{
		{
			rows: [][]string{
				{"route_id", "service_id", "trip_id", "trip_headsign", "trip_short_name", "direction_id", "block_id", "shape_id", "wheelchair_accessible", "bikes_allowed"},
				{"2", "3", "1", "Trip 1", "trip1", "0", "4", "5", "1", "2"},
			},
			expected: expected1,
		},
	}

	for _, tc := range testCases {
		stops, err := LoadTrips(csv.NewReader(strings.NewReader(tableToString(tc.rows))))
		if err != nil && len(err) > 0 {
			t.Error(err)
			continue
		}

		if len(stops) != 1 {
			t.Error("expected one row")
			continue
		}

		if !tipsMatch(tc.expected, *stops[0]) {
			s1, err := json.Marshal(tc.expected)
			if err != nil {
				t.Error(err)
			}
			s2, err := json.Marshal(*stops[0])
			if err != nil {
				t.Error(err)
			}
			t.Error(fmt.Sprintf("expected %v, got %v", string(s1), string(s2)))
		}
	}
}

func TestTripParsingNOK(t *testing.T) {
	testCases := []struct {
		rows     [][]string
		expected []string
	}{
		{
			rows: [][]string{
				{"route_id"},
				{","},
				{" "},
			},
			expected: []string{
				"trips.txt: record on line 2: wrong number of fields",
				"trips.txt:1: route_id must be specified",
				"trips.txt:1: route_id: empty value not allowed",
				"trips.txt:1: trip_id must be specified",
				"trips.txt:1: service_id must be specified",
			},
		},
		{
			rows: [][]string{
				{"trip_id", "service_id", "route_id", "direction_id", "wheelchair_accessible", "bikes_allowed"},
				{"001", "002", "001", "invalid", "invalid", "invalid"},
				{"002", "002", "001", "5", "5", "5"},
				{"002", "002", "001", "0", "0", "0"},
			},
			expected: []string{
				"trips.txt:0: bikes_allowed: strconv.ParseInt: parsing \"invalid\": invalid syntax",
				"trips.txt:0: direction_id: strconv.ParseInt: parsing \"invalid\": invalid syntax",
				"trips.txt:0: wheelchair_accessible: strconv.ParseInt: parsing \"invalid\": invalid syntax",
				"trips.txt:1: bikes_allowed: invalid value",
				"trips.txt:1: direction_id: invalid value",
				"trips.txt:1: wheelchair_accessible: invalid value",
				"trips.txt:2: non-unique id: trip_id",
			},
		},
	}

	for _, tc := range testCases {
		_, err := LoadTrips(csv.NewReader(strings.NewReader(tableToString(tc.rows))))

		sort.Slice(err, func(x, y int) bool {
			return err[x].Error() < err[y].Error()
		})

		sort.Slice(tc.expected, func(x, y int) bool {
			return tc.expected[x] < tc.expected[y]
		})

		if len(err) == 0 {
			t.Error("expected to throw an error")
			continue
		}

		if len(err) != len(tc.expected) {
			t.Error(fmt.Sprintf("expected %v errors, got %v", len(tc.expected), len(err)))
			for _, e := range err {
				fmt.Println(e)
			}
			continue
		}

		for i, e := range err {
			if e.Error() != tc.expected[i] {
				t.Error(fmt.Sprintf("expected error %s, got %s", tc.expected[i], e.Error()))
			}
		}
	}
}

func TestValidateTrips(t *testing.T) {
	shapeId1 := "S001"
	shapeId2 := "S002"

	testCases := []struct {
		trips          []*Trip
		routes         []*Route
		services       []*CalendarItem
		shapes         []*Shape
		expectedErrors []string
	}{
		{
			trips: []*Trip{
				{RouteId: "R001", ServiceId: "C001", ShapeId: &shapeId1, lineNumber: 0},
				{RouteId: "R002", ServiceId: "C002", ShapeId: &shapeId2, lineNumber: 1},
			},
			routes: []*Route{{
				Id: "R001",
			}},
			services: []*CalendarItem{{
				ServiceId: "C001",
			}},
			shapes: []*Shape{{
				Id: "S001",
			}},
			expectedErrors: []string{
				"trips.txt:1: referenced route_id not found in routes.txt",
				"trips.txt:1: referenced service_id not found in calendar.txt",
				"trips.txt:1: referenced shape_id not found in shapes.txt",
			},
		},
		{
			trips:          nil,
			expectedErrors: []string{},
		},
		{
			trips: []*Trip{
				{RouteId: "R001", ServiceId: "C001", ShapeId: &shapeId1, lineNumber: 0},
				{RouteId: "R002", ServiceId: "C002", ShapeId: &shapeId2, lineNumber: 1},
			},
			expectedErrors: []string{},
		},
		{
			trips:          []*Trip{nil},
			expectedErrors: []string{},
		},
		{
			trips: []*Trip{
				{RouteId: "R001", ServiceId: "C001", ShapeId: &shapeId1, lineNumber: 0},
				{RouteId: "R002", ServiceId: "C002", ShapeId: &shapeId2, lineNumber: 1},
			},
			routes: []*Route{nil},
			expectedErrors: []string{
				"trips.txt:0: referenced route_id not found in routes.txt",
				"trips.txt:1: referenced route_id not found in routes.txt",
			},
		},
		{
			trips: []*Trip{
				{RouteId: "R001", ServiceId: "C001", ShapeId: &shapeId1, lineNumber: 0},
				{RouteId: "R002", ServiceId: "C002", ShapeId: &shapeId2, lineNumber: 1},
			},
			services: []*CalendarItem{nil},
			expectedErrors: []string{
				"trips.txt:0: referenced service_id not found in calendar.txt",
				"trips.txt:1: referenced service_id not found in calendar.txt",
			},
		},
		{
			trips: []*Trip{
				{RouteId: "R001", ServiceId: "C001", ShapeId: &shapeId1, lineNumber: 0},
				{RouteId: "R002", ServiceId: "C002", ShapeId: &shapeId2, lineNumber: 1},
			},
			shapes: []*Shape{nil},
			expectedErrors: []string{
				"trips.txt:0: referenced shape_id not found in shapes.txt",
				"trips.txt:1: referenced shape_id not found in shapes.txt",
			},
		},
	}

	for _, tc := range testCases {
		err := ValidateTrips(tc.trips, tc.routes, tc.services, tc.shapes)
		checkErrors(tc.expectedErrors, err, t)
	}
}

func tipsMatch(a Trip, b Trip) bool {
	return a.Id == b.Id && a.ServiceId == b.ServiceId && a.RouteId == b.RouteId && *a.HeadSign == *b.HeadSign && *a.ShortName == *b.ShortName &&
		*a.DirectionId == *b.DirectionId && *a.BlockId == *b.BlockId && *a.ShapeId == *b.ShapeId && *a.WheelchairAccessible == *b.WheelchairAccessible &&
		*a.BikesAllowed == *b.BikesAllowed
}
