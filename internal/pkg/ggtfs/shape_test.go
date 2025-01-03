//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"fmt"
	"testing"
)

func TestCreateShape(t *testing.T) {
	headerMap := map[string]int{"shape_id": 0, "shape_pt_lat": 1, "shape_pt_lon": 2, "shape_pt_sequence": 3, "shape_dist_traveled": 4}

	tests := map[string]struct {
		headers    map[string]int
		rows       [][]string
		lineNumber int
		expected   []*Shape
	}{
		"empty-row": {
			headers: headerMap,
			rows:    [][]string{{"", "", "", "", ""}},
			expected: []*Shape{{
				Id:           stringPtr(""),
				PtLat:        stringPtr(""),
				PtLon:        stringPtr(""),
				PtSequence:   stringPtr(""),
				DistTraveled: stringPtr(""),
				LineNumber:   0,
			}},
		},
		"nil-values": {
			headers: headerMap,
			rows:    [][]string{nil},
			expected: []*Shape{{
				Id:           nil,
				PtLat:        nil,
				PtLon:        nil,
				PtSequence:   nil,
				DistTraveled: nil,
				LineNumber:   0,
			}},
		},
		"OK": {
			headers: headerMap,
			rows: [][]string{
				{"1", "1.111", "1.111", "1", "100"},
				{"1", "2.111", "2.111", "2", "100"},
			},
			expected: []*Shape{{
				Id:           stringPtr("1"),
				PtLat:        stringPtr("1.111"),
				PtLon:        stringPtr("1.111"),
				PtSequence:   stringPtr("1"),
				DistTraveled: stringPtr("100"),
			}, {
				Id:           stringPtr("1"),
				PtLat:        stringPtr("2.111"),
				PtLon:        stringPtr("2.111"),
				PtSequence:   stringPtr("2"),
				DistTraveled: stringPtr("100"),
				LineNumber:   1,
			}},
		},
	}

	for name, tt := range tests {
		t.Run(fmt.Sprintf("%s", name), func(t *testing.T) {
			var actual []*Shape
			for i, row := range tt.rows {
				actual = append(actual, CreateShape(row, tt.headers, i))
			}
			handleEntityCreateResults(t, tt.expected, actual)
		})
	}
}

func TestValidateShapes(t *testing.T) {
	tests := map[string]struct {
		actualEntities  []*Shape
		expectedResults []Result
	}{
		"nil-slice": {
			actualEntities:  nil,
			expectedResults: []Result{},
		},
		"nil-slice-items": {
			actualEntities:  []*Shape{nil},
			expectedResults: []Result{},
		},
		"invalid-fields": {
			actualEntities: []*Shape{
				{
					Id:           stringPtr("1"),
					PtLat:        stringPtr("Not a latitude"),
					PtLon:        stringPtr("Not a longitude"),
					PtSequence:   stringPtr("Not a sequence"),
					DistTraveled: stringPtr("Not a distance"),
				},
			},
			expectedResults: []Result{
				InvalidLatitudeResult{SingleLineResult{FileName: "shapes.txt", FieldName: "shape_pt_lat"}},
				InvalidLongitudeResult{SingleLineResult{FileName: "shapes.txt", FieldName: "shape_pt_lon"}},
				InvalidIntegerResult{SingleLineResult{FileName: "shapes.txt", FieldName: "shape_pt_sequence"}},
				InvalidFloatResult{SingleLineResult{FileName: "shapes.txt", FieldName: "shape_dist_traveled"}},
				TooFewShapePointsResult{
					FileName: "shapes.txt",
					ShapeId:  "1",
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(fmt.Sprintf("%s", name), func(t *testing.T) {
			handleValidationResults(t, ValidateShapes(tt.actualEntities), tt.expectedResults)
		})
	}
}
