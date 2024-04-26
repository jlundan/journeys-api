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

func TestShapeCSVParsing(t *testing.T) {
	items, errors := LoadShapes(csv.NewReader(strings.NewReader("")))
	if len(errors) > 0 {
		t.Error(errors)
	}
	if len(items) != 0 {
		t.Error("expected zero items")
	}

	reader := csv.NewReader(strings.NewReader("foo,bar\n1,2"))
	reader.Comma = ','
	reader.Comment = ','
	_, errors = LoadShapes(reader)
	if len(errors) == 0 {
		t.Error("expected to throw error")
	}
}

func TestShapesParsingOK(t *testing.T) {
	dist := float64(100)

	expected1 := Shape{
		Id:           "1",
		PtLat:        1.111,
		PtLon:        1.111,
		PtSequence:   1,
		DistTraveled: &dist,
	}

	testCases := []struct {
		headers  map[string]uint8
		rows     [][]string
		expected Shape
	}{
		{
			rows: [][]string{
				{"shape_id", "shape_pt_lat", "shape_pt_lon", "shape_pt_sequence", "shape_dist_traveled"},
				{"1", "1.111", "1.111", "1", "100"},
			},
			expected: expected1,
		},
	}

	for _, tc := range testCases {
		stops, err := LoadShapes(csv.NewReader(strings.NewReader(tableToString(tc.rows))))
		if err != nil && len(err) > 0 {
			t.Error(err)
			continue
		}

		if len(stops) != 1 {
			t.Error("expected one row")
			continue
		}

		if !shapesMatch(tc.expected, *stops[0]) {
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

func TestShapesParsingNOK(t *testing.T) {
	testCases := []struct {
		headers  map[string]uint8
		rows     [][]string
		expected []string
	}{
		{
			headers: map[string]uint8{"shape_id": 0},
			rows: [][]string{
				{"shape_id"},
				{","},
				{" "},
			},
			expected: []string{
				"shapes.txt: record on line 2: wrong number of fields",
				"shapes.txt:1: shape_id must be specified",
				"shapes.txt:1: shape_id: empty value not allowed",
				"shapes.txt:1: shape_pt_lat must be specified",
				"shapes.txt:1: shape_pt_lon must be specified",
				"shapes.txt:1: shape_pt_sequence must be specified",
			},
		},
	}

	for _, tc := range testCases {
		_, err := LoadShapes(csv.NewReader(strings.NewReader(tableToString(tc.rows))))

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

func TestValidateShapes(t *testing.T) {
	testCases := []struct {
		shapes         []*Shape
		expectedErrors []string
	}{
		{
			shapes: []*Shape{
				{Id: "1000", lineNumber: 0},
				{Id: "1000", lineNumber: 1},
				{Id: "1000", lineNumber: 2},
			},
			expectedErrors: []string{},
		},
		{
			shapes: []*Shape{
				{Id: "1000", lineNumber: 0},
				{Id: "1000", lineNumber: 1},
			},
			expectedErrors: []string{},
		},
		{
			shapes:         nil,
			expectedErrors: []string{},
		},
		{
			shapes: []*Shape{
				{Id: "1000", lineNumber: 0},
				nil,
			},
			expectedErrors: []string{
				"shapes.txt: shape (1000) has less than two shape points",
			},
		},
		{
			shapes: []*Shape{
				{Id: "1000", lineNumber: 0},
			},
			expectedErrors: []string{
				"shapes.txt: shape (1000) has less than two shape points",
			},
		},
	}

	for _, tc := range testCases {
		err := ValidateShapes(tc.shapes)

		checkErrors(tc.expectedErrors, err, t)
	}
}

func shapesMatch(a Shape, b Shape) bool {
	return a.Id == b.Id && a.PtLat == b.PtLat && a.PtLon == b.PtLon &&
		a.PtSequence == b.PtSequence && *a.DistTraveled == *b.DistTraveled
}
