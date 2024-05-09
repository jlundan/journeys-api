package ggtfs

import (
	"encoding/csv"
	"fmt"
)

type Shape struct {
	Id           string
	PtLat        float64
	PtLon        float64
	PtSequence   int
	DistTraveled *float64
	lineNumber   int
}

var validShapeHeaders = []string{"shape_id", "shape_pt_lat", "shape_pt_lon", "shape_pt_sequence",
	"shape_dist_traveled"}

func LoadShapes(csvReader *csv.Reader) ([]*Shape, []error) {
	shapes := make([]*Shape, 0)
	errs := make([]error, 0)

	headers, err := ReadHeaderRow(csvReader, validShapeHeaders)
	if err != nil {
		errs = append(errs, createFileError(ShapesFileName, fmt.Sprintf("read error: %v", err.Error())))
		return shapes, errs
	}
	if headers == nil {
		return shapes, errs
	}

	index := 0
	for {
		row, err := ReadDataRow(csvReader)
		if err != nil {
			errs = append(errs, createFileError(ShapesFileName, fmt.Sprintf("%v", err.Error())))
			index++
			continue
		}
		if row == nil {
			break
		}

		rowErrs := make([]error, 0)
		shape := Shape{
			lineNumber: index,
		}

		var (
			shapeId  *string
			lat, lon *float64
			sequence *int
		)
		for name, column := range headers {
			switch name {
			case "shape_id":
				shapeId = handleIDField(row[column], ShapesFileName, name, index, &rowErrs)
			case "shape_pt_lat":
				lat = handleFloat64Field(row[column], ShapesFileName, name, index, &rowErrs)
			case "shape_pt_lon":
				lon = handleFloat64Field(row[column], ShapesFileName, name, index, &rowErrs)
			case "shape_pt_sequence":
				sequence = handleIntField(row[column], ShapesFileName, name, index, &rowErrs)
			case "shape_dist_traveled":
				shape.DistTraveled = handleFloat64Field(row[column], ShapesFileName, name, index, &rowErrs)
			}
		}

		if shapeId == nil {
			rowErrs = append(rowErrs, createFileRowError(ShapesFileName, shape.lineNumber, "shape_id must be specified"))
		} else {
			shape.Id = *shapeId
		}

		if lat == nil {
			rowErrs = append(rowErrs, createFileRowError(ShapesFileName, shape.lineNumber, "shape_pt_lat must be specified"))
		} else {
			shape.PtLat = *lat
		}

		if lon == nil {
			rowErrs = append(rowErrs, createFileRowError(ShapesFileName, shape.lineNumber, "shape_pt_lon must be specified"))
		} else {
			shape.PtLon = *lon
		}

		if sequence == nil {
			rowErrs = append(rowErrs, createFileRowError(ShapesFileName, shape.lineNumber, "shape_pt_sequence must be specified"))
		} else {
			shape.PtSequence = *sequence
		}

		if len(rowErrs) > 0 {
			errs = append(errs, rowErrs...)
		} else {
			shapes = append(shapes, &shape)
		}

		index++
	}

	return shapes, errs
}

func ValidateShapes(shapes []*Shape) []error {
	var validationErrors []error

	if shapes == nil {
		return validationErrors
	}

	shapeIdToPointCount := make(map[string]uint64)
	for _, shapeItem := range shapes {
		if shapeItem == nil {
			continue
		}
		shapeIdToPointCount[shapeItem.Id]++
	}

	for shapeId, pointCount := range shapeIdToPointCount {
		if pointCount < 2 {
			validationErrors = append(validationErrors, createFileError(ShapesFileName, fmt.Sprintf("shape (%v) has less than two shape points", shapeId)))
		}
	}

	return validationErrors
}
