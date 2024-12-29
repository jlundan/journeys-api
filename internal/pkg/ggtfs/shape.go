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

func ValidateShape(s Shape) []error {
	var validationErrors []error

	requiredFields := map[string]FieldTobeValidated{
		"shape_id":          &s.Id,
		"shape_pt_lat":      &s.PtLat,
		"shape_pt_lon":      &s.PtLon,
		"shape_pt_sequence": &s.PtSequence,
	}
	validateRequiredFields(requiredFields, &validationErrors, s.LineNumber, ShapesFileName)

	optionalFields := map[string]FieldTobeValidated{
		"shape_dist_traveled": &s.DistTraveled,
	}
	validateOptionalFields(optionalFields, &validationErrors, s.LineNumber, ShapesFileName)

	return validationErrors
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

		vErr := ValidateShape(*shapeItem)
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
