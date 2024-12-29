//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"encoding/csv"
	"strings"
	"testing"
)

var validShapeHeaders = []string{"shape_id", "shape_pt_lat", "shape_pt_lon", "shape_pt_sequence", "shape_dist_traveled"}

func TestShouldReturnEmptyShapeArrayOnEmptyString(t *testing.T) {
	agencies, errors := LoadEntitiesFromCSV[*Shape](csv.NewReader(strings.NewReader("")), validShapeHeaders, CreateShape, ShapesFileName)
	if len(errors) > 0 {
		t.Error(errors)
	}
	if len(agencies) != 0 {
		t.Error("expected zero calendar items")
	}
}

func TestShapeParsing(t *testing.T) {
	loadShapesFunc := func(reader *csv.Reader) ([]interface{}, []error) {
		shapes, errs := LoadEntitiesFromCSV[*Shape](reader, validShapeHeaders, CreateShape, ShapesFileName)
		entities := make([]interface{}, len(shapes))
		for i, shape := range shapes {
			entities[i] = shape
		}
		return entities, errs
	}

	validateShapesFunc := func(entities []interface{}, _ map[string][]interface{}) ([]error, []string) {
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
		Id:           NewID(stringPtr("1")),
		PtLat:        NewLatitude(stringPtr("1.111")),
		PtLon:        NewLongitude(stringPtr("1.111")),
		PtSequence:   NewPositiveInteger(stringPtr("1")),
		DistTraveled: NewPositiveFloat(stringPtr("100")),
		LineNumber:   2,
	}

	expected2 := Shape{
		Id:           NewID(stringPtr("1")),
		PtLat:        NewLatitude(stringPtr("1.211")),
		PtLon:        NewLongitude(stringPtr("1.211")),
		PtSequence:   NewPositiveInteger(stringPtr("2")),
		DistTraveled: NewPositiveFloat(stringPtr("100")),
		LineNumber:   3,
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
			"shapes.txt: record on line 2: wrong number of fields",
			"shapes.txt:3: invalid mandatory field: shape_id",
			"shapes.txt:3: invalid mandatory field: shape_pt_lat",
			"shapes.txt:3: invalid mandatory field: shape_pt_lon",
			"shapes.txt:3: invalid mandatory field: shape_pt_sequence",
		},
	}

	testCases["2"] = ggtfsTestCase{
		csvRows: [][]string{
			{"shape_id", "shape_pt_lat", "shape_pt_lon", "shape_pt_sequence", "shape_dist_traveled"},
			{"1", "11.1", "11.1", "1", ""},
			{"2", "11.1", "11.1", "1", "invalid"},
			{"3", "11.1", "11.1", "1", "100"},
		},
		expectedErrors: []string{
			"shapes.txt: shape (1) has less than two shape points",
			"shapes.txt: shape (3) has less than two shape points",
			"shapes.txt:3: invalid field: shape_dist_traveled",
		},
	}

	return testCases
}

func TestShouldNotFailValidationOnNilShapes(t *testing.T) {
	ValidateShapes(nil)
}

func TestShouldNotFailValidationOnNilShapeItem(t *testing.T) {
	ValidateShapes([]*Shape{nil})
}
