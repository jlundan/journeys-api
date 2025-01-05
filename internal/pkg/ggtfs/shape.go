package ggtfs

type Shape struct {
	Id           *string // shape_id 			(required)
	PtLat        *string // shape_pt_lat 		(required)
	PtLon        *string // shape_pt_lon 		(required)
	PtSequence   *string // shape_pt_sequence 	(required)
	DistTraveled *string // shape_dist_traveled (optional)
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
			shape.Id = v
		case "shape_pt_lat":
			shape.PtLat = v
		case "shape_pt_lon":
			shape.PtLon = v
		case "shape_pt_sequence":
			shape.PtSequence = v
		case "shape_dist_traveled":
			shape.DistTraveled = v
		}
	}

	return &shape
}

func ValidateShape(s Shape) []ValidationNotice {

	var validationResults []ValidationNotice

	fields := []struct {
		fieldType FieldType
		name      string
		value     *string
		required  bool
	}{
		{FieldTypeID, "shape_id", s.Id, true},
		{FieldTypeLatitude, "shape_pt_lat", s.PtLat, true},
		{FieldTypeLongitude, "shape_pt_lon", s.PtLon, true},
		{FieldTypeInteger, "shape_pt_sequence", s.PtSequence, true},
		{FieldTypeFloat, "shape_dist_traveled", s.DistTraveled, false},
	}

	for _, field := range fields {
		validationResults = append(validationResults, validateField(field.fieldType, field.name, field.value, field.required, FileNameShapes, s.LineNumber)...)
	}

	return validationResults
}

func ValidateShapes(shapes []*Shape) []ValidationNotice {
	var validationResults []ValidationNotice

	if shapes == nil {
		return validationResults
	}

	shapeIdToPointCount := make(map[string]int)
	for _, shapeItem := range shapes {
		if shapeItem == nil {
			continue
		}

		vRes := ValidateShape(*shapeItem)
		if len(vRes) > 0 {
			validationResults = append(validationResults, vRes...)
		}

		shapeIdToPointCount[*shapeItem.Id]++
	}

	// Check that each shape has at least two points.
	for shapeId, pointCount := range shapeIdToPointCount {
		if pointCount < 2 {
			validationResults = append(validationResults, TooFewShapePointsNotice{
				FileName: FileNameShapes,
				ShapeId:  shapeId,
			})
		}
	}

	return validationResults
}
