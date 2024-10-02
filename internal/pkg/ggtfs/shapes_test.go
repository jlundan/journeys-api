//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"encoding/csv"
	"strings"
	"testing"
)

func TestShouldReturnEmptyShapeArrayOnEmptyString(t *testing.T) {
	agencies, errors := LoadShapes(csv.NewReader(strings.NewReader("")))
	if len(errors) > 0 {
		t.Error(errors)
	}
	if len(agencies) != 0 {
		t.Error("expected zero calendar items")
	}
}

func TestShapeParsing(t *testing.T) {
	loadShapesFunc := func(reader *csv.Reader) ([]interface{}, []error) {
		shapes, errs := LoadShapes(reader)
		entities := make([]interface{}, len(shapes))
		for i, shape := range shapes {
			entities[i] = shape
		}
		return entities, errs
	}

	validateShapesFunc := func(entities []interface{}) ([]error, []string) {
		shapes := make([]*Shape, len(entities))
		for i, entity := range entities {
			if shape, ok := entity.(*Shape); ok {
				shapes[i] = shape
			}
		}

		return ValidateShapes(shapes)
	}

	runGenericGTFSParseTest(t, "ShapeOKTestcases", loadShapesFunc, validateShapesFunc, false, getShapeOKTestcases())
	runGenericGTFSParseTest(t, "ShapeNOKTestcases", loadShapesFunc, validateShapesFunc, false, getShapeNOKTestcases())
}

func getShapeOKTestcases() map[string]ggtfsTestCase {
	expected1 := Shape{
		Id:           "1",
		PtLat:        "1.111",
		PtLon:        "1.111",
		PtSequence:   "1",
		DistTraveled: stringPtr("100"),
		LineNumber:   0,
	}

	expected2 := Shape{
		Id:           "1",
		PtLat:        "1.211",
		PtLon:        "1.211",
		PtSequence:   "2",
		DistTraveled: stringPtr("100"),
		LineNumber:   1,
	}

	testCases := make(map[string]ggtfsTestCase)
	testCases["1"] = ggtfsTestCase{
		csvRows: [][]string{
			{"shape_id", "shape_pt_lat", "shape_pt_lon", "shape_pt_sequence", "shape_dist_traveled"},
			{"1", "1.111", "1.111", "1", "100"},
			{"1", "1.211", "1.211", "2", "100"},
		},
		expectedStructs: []interface{}{&expected1, &expected2},
	}

	return testCases
}

func getShapeNOKTestcases() map[string]ggtfsTestCase {
	testCases := make(map[string]ggtfsTestCase)
	testCases["1"] = ggtfsTestCase{
		csvRows: [][]string{
			{"shape_id"},
			{","},
			{" "},
		},
		expectedErrors: []string{
			// TODO: Fix these
			"shapes.txt: record on line 2: wrong number of fields",
			"shapes.txt: shape ( ) has less than two shape points",
			//"shapes.txt:1: shape_id must be specified",
			"shapes.txt:1: shape_pt_lat must be specified and non-empty",
			"shapes.txt:1: shape_pt_lon must be specified and non-empty",
			"shapes.txt:1: shape_pt_sequence must be specified",
			//"shapes.txt: record on line 2: wrong number of fields",
			//"shapes.txt:1: shape_id must be specified",
			//"shapes.txt:1: shape_id: empty value not allowed",
			//"shapes.txt:1: shape_pt_lat must be specified",
			//"shapes.txt:1: shape_pt_lon must be specified",
			//"shapes.txt:1: shape_pt_sequence must be specified",
		},
	}

	return testCases
}

// TODO: Integrate these into tests run by runGenericGTFSParseTest
//
//	func TestValidateShapes(t *testing.T) {
//		testCases := []struct {
//			shapes         []*Shape
//			expectedErrors []string
//		}{
//			{
//				shapes: []*Shape{
//					{Id: "1000", LineNumber: 0},
//					{Id: "1000", LineNumber: 1},
//					{Id: "1000", LineNumber: 2},
//				},
//				expectedErrors: []string{},
//			},
//			{
//				shapes: []*Shape{
//					{Id: "1000", LineNumber: 0},
//					{Id: "1000", LineNumber: 1},
//				},
//				expectedErrors: []string{},
//			},
//			{
//				shapes:         nil,
//				expectedErrors: []string{},
//			},
//			{
//				shapes: []*Shape{
//					{Id: "1000", LineNumber: 0},
//					nil,
//				},
//				expectedErrors: []string{
//					"shapes.txt: shape (1000) has less than two shape points",
//				},
//			},
//			{
//				shapes: []*Shape{
//					{Id: "1000", LineNumber: 0},
//				},
//				expectedErrors: []string{
//					"shapes.txt: shape (1000) has less than two shape points",
//				},
//			},
//		}
//
//		for _, tc := range testCases {
//			err := ValidateShapes(tc.shapes)
//
//			checkErrors(tc.expectedErrors, err, t)
//		}
//	}
