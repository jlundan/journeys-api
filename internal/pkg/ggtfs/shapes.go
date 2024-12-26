package ggtfs

import (
	"fmt"
)

type Shape struct {
	Id           ID              // shape_id 			(required)
	PtLat        Latitude        // shape_pt_lat 		(required)
	PtLon        Longitude       // shape_pt_lon 		(required)
	PtSequence   PositiveInteger // shape_pt_sequence 	(required)
	DistTraveled PositiveFloat   // shape_dist_traveled (optional)
	LineNumber   int
}

func (s Shape) Validate() []error {
	var validationErrors []error

	requiredFields := []struct {
		fieldName string
		field     ValidAndPresentField
	}{
		{"shape_id", &s.Id},
		{"shape_pt_lat", &s.PtLat},
		{"shape_pt_lon", &s.PtLon},
		{"shape_pt_sequence", &s.PtSequence},
	}
	for _, f := range requiredFields {
		if !f.field.IsValid() {
			validationErrors = append(validationErrors, createFileRowError(ShapesFileName, s.LineNumber, createInvalidRequiredFieldString(f.fieldName)))
		}
	}

	optionalFields := []struct {
		field     ValidAndPresentField
		fieldName string
	}{
		{&s.DistTraveled, "shape_dist_traveled"},
	}

	for _, field := range optionalFields {
		if field.field != nil && field.field.IsPresent() && !field.field.IsValid() {
			validationErrors = append(validationErrors, createFileRowError(ShapesFileName, s.LineNumber, createInvalidFieldString(field.fieldName)))
		}
	}

	return validationErrors
}

func CreateShape(row []string, headers map[string]int, lineNumber int) *Shape {

	shape := Shape{
		LineNumber: lineNumber,
	}

	for hName := range headers {
		v := getRowValueForHeaderName(row, headers, hName)

		switch hName {
		case "shape_id":
			shape.Id = NewID(v)
		case "shape_pt_lat":
			shape.PtLat = NewLatitude(v)
		case "shape_pt_lon":
			shape.PtLon = NewLongitude(v)
		case "shape_pt_sequence":
			shape.PtSequence = NewPositiveInteger(v)
		case "shape_dist_traveled":
			shape.DistTraveled = NewPositiveFloat(v)
		}
	}

	return &shape
}

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

	// TODO: VALIDATION: shape_pt_sequence: Sequence in which the shape points connect to form the shape. Values must increase along the trip but do not need to be consecutive.
	// Implement this if there is a way to check the contrast between two colors

	return validationErrors, []string{}
}
