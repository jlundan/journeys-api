package ggtfs

type StopExtensions struct {
	MunicipalityId *string // municipality_id (optional)
}

type Stop struct {
	Id                 *string // stop_id               (required)
	Code               *string // stop_code             (optional)
	Name               *string // stop_name             (conditionally required)
	TTSName            *string // tts_stop_name         (optional)
	Desc               *string // stop_desc             (optional)
	Lat                *string // stop_lat              (conditionally required)
	Lon                *string // stop_lon              (conditionally required)
	ZoneId             *string // zone_id               (optional)
	URL                *string // stop_url              (optional)
	LocationType       *string // location_type         (optional)
	ParentStation      *string // parent_station        (conditionally required)
	Timezone           *string // stop_timezone         (optional)
	WheelchairBoarding *string // wheelchair_boarding   (optional)
	PlatformCode       *string // platform_code         (optional)
	LevelId            *string // level_id              (optional)
	Extensions         *StopExtensions
	LineNumber         int
}

func CreateStop(row []string, headers map[string]int, lineNumber int) *Stop {
	stop := Stop{
		LineNumber: lineNumber,
	}

	for hName := range headers {
		v := getRowValueForHeaderName(row, headers, hName)
		switch hName {
		case "stop_id":
			stop.Id = v
		case "stop_code":
			stop.Code = v
		case "stop_name":
			stop.Name = v
		case "tts_stop_name":
			stop.TTSName = v
		case "stop_desc":
			stop.Desc = v
		case "stop_lat":
			stop.Lat = v
		case "stop_lon":
			stop.Lon = v
		case "zone_id":
			stop.ZoneId = v
		case "stop_url":
			stop.URL = v
		case "location_type":
			stop.LocationType = v
		case "parent_station":
			stop.ParentStation = v
		case "stop_timezone":
			stop.Timezone = v
		case "wheelchair_boarding":
			stop.WheelchairBoarding = v
		case "level_id":
			stop.LevelId = v
		case "platform_code":
			stop.PlatformCode = v
		case "municipality_id":
			stop.Extensions = &StopExtensions{
				MunicipalityId: v,
			}
		}
	}

	return &stop
}

func ValidateStop(s Stop) []Result {
	var validationResults []Result

	fields := []struct {
		fieldType FieldType
		name      string
		value     *string
		required  bool
	}{
		{FieldTypeID, "stop_id", s.Id, true},
		{FieldTypeText, "stop_code", s.Code, false},
		{FieldTypeText, "stop_name", s.Name, false},
		{FieldTypeText, "tts_stop_name", s.TTSName, false},
		{FieldTypeText, "stop_desc", s.Desc, false},
		{FieldTypeID, "zone_id", s.ZoneId, false},
		{FieldTypeURL, "stop_url", s.URL, false},
		{FieldTypeLocationType, "location_type", s.LocationType, false},
		{FieldTypeID, "parent_station", s.ParentStation, false},
		{FieldTypeTimezone, "stop_timezone", s.Timezone, false},
		{FieldTypeWheelchairBoarding, "wheelchair_boarding", s.WheelchairBoarding, false},
		{FieldTypeID, "level_id", s.LevelId, false},
		{FieldTypeText, "platform_code", s.PlatformCode, false},
	}

	for _, field := range fields {
		validationResults = append(validationResults, validateField(field.fieldType, field.name, field.value, field.required, FileNameStops, s.LineNumber)...)
	}

	requiredFields := []struct {
		fieldName  string
		fieldValue *string
	}{
		{"stop_name", s.Name},
		{"stop_lat", s.Lat},
		{"stop_lon", s.Lon},
	}

	for _, field := range requiredFields {
		if StringIsNilOrEmpty(field.fieldValue) && IsLocationTypeValid(s.LocationType) {
			locationType := *s.LocationType
			if locationType == "0" || locationType == "1" || locationType == "2" {
				validationResults = append(validationResults, FieldRequiredForLocationTypeResult{
					RequiredField: field.fieldName,
					LocationType:  locationType,
					FileName:      FileNameStops,
					Line:          s.LineNumber,
				})
			}
		}
	}

	if StringIsNilOrEmpty(s.ParentStation) && IsLocationTypeValid(s.LocationType) {
		locationType := *s.LocationType
		if locationType == "2" || locationType == "3" || locationType == "4" {
			validationResults = append(validationResults, FieldRequiredForLocationTypeResult{
				RequiredField: "parent_station",
				LocationType:  locationType,
				FileName:      FileNameStops,
				Line:          s.LineNumber,
			})
		}
	}

	return validationResults
}

func ValidateStops(stops []*Stop) []Result {
	var validationResults []Result

	if stops == nil {
		return validationResults
	}

	// Check for unique stop_id values.
	usedIds := make(map[string]bool)
	for _, stop := range stops {
		if stop == nil {
			continue
		}

		vRes := ValidateStop(*stop)
		if len(vRes) > 0 {
			validationResults = append(validationResults, vRes...)
		}

		if StringIsNilOrEmpty(stop.Id) {
			continue
		}

		if usedIds[*stop.Id] {
			validationResults = append(validationResults, FieldIsNotUniqueResult{SingleLineResult{
				FileName:  FileNameStops,
				FieldName: "stop_id",
				Line:      stop.LineNumber,
			}})
		} else {
			usedIds[*stop.Id] = true
		}
	}

	return validationResults
}
