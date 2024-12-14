package ggtfs

import (
	"encoding/csv"
	"fmt"
)

// Shape struct with fields as strings and optional fields as string pointers.
type Shape struct {
	Id           ID        // shape_id
	PtLat        Latitude  // shape_pt_lat
	PtLon        Longitude // shape_pt_lon
	PtSequence   Integer   // shape_pt_sequence
	DistTraveled *Float    // shape_dist_traveled (optional)
	LineNumber   int       // Line number in the CSV file for error reporting
}

func (s Shape) Validate() []error {
	var validationErrors []error

	fields := []struct {
		fieldName string
		field     ValidAndPresentField
	}{
		{"shape_id", &s.Id},
		{"shape_pt_lat", &s.PtLat},
		{"shape_pt_lon", &s.PtLon},
		{"shape_pt_sequence", &s.PtSequence},
	}
	for _, f := range fields {
		validationErrors = append(validationErrors, validateFieldIsPresentAndValid(f.field, f.fieldName, s.LineNumber, ShapesFileName)...)
	}

	if s.DistTraveled != nil && !s.DistTraveled.IsValid() {
		validationErrors = append(validationErrors, createFileRowError(ShapesFileName, s.LineNumber, createInvalidFieldString("shape_dist_traveled")))
	}

	return validationErrors
}

// ValidShapeHeaders defines the headers expected in the shapes CSV file.
var validShapeHeaders = []string{"shape_id", "shape_pt_lat", "shape_pt_lon", "shape_pt_sequence", "shape_dist_traveled"}

// LoadShapes loads shapes from a CSV reader and returns them along with any errors.
func LoadShapes(csvReader *csv.Reader) ([]*Shape, []error) {
	entities, errs := loadEntities(csvReader, validShapeHeaders, CreateShape, ShapesFileName)

	shapes := make([]*Shape, 0, len(entities))
	for _, entity := range entities {
		if shape, ok := entity.(*Shape); ok {
			shapes = append(shapes, shape)
		}
	}

	return shapes, errs
}

// CreateShape creates and validates a Shape instance from the CSV row data.
func CreateShape(row []string, headers map[string]int, lineNumber int) interface{} {

	shape := Shape{
		LineNumber: lineNumber,
	}

	for hName, hPos := range headers {
		switch hName {
		case "shape_id":
			shape.Id = NewID(getRowValue(row, hPos))
		case "shape_pt_lat":
			shape.PtLat = NewLatitude(getRowValue(row, hPos))
		case "shape_pt_lon":
			shape.PtLon = NewLongitude(getRowValue(row, hPos))
		case "shape_pt_sequence":
			shape.PtSequence = NewInteger(getRowValue(row, hPos))
		case "shape_dist_traveled":
			shape.DistTraveled = NewOptionalFloat(getRowValue(row, hPos))
		}
	}

	return &shape
}

// ValidateShapes performs additional validation for a list of Shape instances.
func ValidateShapes(shapes []*Shape) ([]error, []string) {
	var validationErrors []error

	if shapes == nil {
		return validationErrors, []string{}
	}

	shapeIdToPointCount := make(map[string]int)
	for _, shapeItem := range shapes {
		if shapeItem == nil {
			continue
		}

		vErr := shapeItem.Validate()
		if len(vErr) > 0 {
			validationErrors = append(validationErrors, vErr...)
			continue
		}

		shapeIdToPointCount[shapeItem.Id.String()]++
	}

	// Check that each shape has at least two points.
	for shapeId, pointCount := range shapeIdToPointCount {
		if pointCount < 2 {
			validationErrors = append(validationErrors, createFileError(ShapesFileName, fmt.Sprintf("shape (%v) has less than two shape points", shapeId)))
		}
	}

	return validationErrors, []string{}
}
