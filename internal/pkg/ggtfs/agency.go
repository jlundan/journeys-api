package ggtfs

type Agency struct {
	Id         *string // agency_id 		(conditionally required)
	Name       *string // agency_name 		(required)
	URL        *string // agency_url 		(required)
	Timezone   *string // agency_timezone 	(required)
	Lang       *string // agency_lang 		(optional)
	Phone      *string // agency_phone 		(optional)
	FareURL    *string // agency_fare_url 	(optional)
	Email      *string // agency_email 		(optional)
	LineNumber int
}

func CreateAgency(row []string, headers map[string]int, lineNumber int) *Agency {
	agency := Agency{
		LineNumber: lineNumber,
	}

	for hName := range headers {
		v := getRowValueForHeaderName(row, headers, hName)

		switch hName {
		case "agency_id":
			agency.Id = v
		case "agency_name":
			agency.Name = v
		case "agency_url":
			agency.URL = v
		case "agency_timezone":
			agency.Timezone = v
		case "agency_lang":
			agency.Lang = v
		case "agency_phone":
			agency.Phone = v
		case "agency_fare_url":
			agency.FareURL = v
		case "agency_email":
			agency.Email = v
		}
	}

	return &agency
}

func ValidateAgency(a Agency) []Result {
	var validationResults []Result

	fields := []struct {
		fieldName  string
		fieldValue *string
		required   bool
	}{
		{"agency_name", a.Name, true},
		{"agency_url", a.URL, true},
		{"agency_timezone", a.Timezone, true},
		{"agency_id", a.Id, false},
		{"agency_lang", a.Lang, false},
		{"agency_phone", a.Phone, false},
		{"agency_fare_url", a.URL, false},
		{"agency_email", a.Email, false},
	}

	for _, field := range fields {
		validationResults = append(validationResults, validateField(field.fieldName, field.fieldValue, field.required, "agency.txt", a.LineNumber)...)
	}

	return validationResults
}

func ValidateAgencies(agencies []*Agency) []Result {
	var results []Result

	var filteredAgencies []*Agency

	for _, a := range agencies {
		if a != nil {
			filteredAgencies = append(filteredAgencies, a)
		}
	}

	aLength := len(filteredAgencies)

	if aLength == 0 {
		return []Result{}
	}

	if aLength == 1 && StringIsNilOrEmpty(filteredAgencies[0].Id) {
		return []Result{SingleAgencyRecommendedResult{
			FileName: "agency.txt",
		}}
	}

	if aLength == 1 {
		return ValidateAgency(*filteredAgencies[0])
	}

	usedIds := make(map[string]bool)
	for _, a := range filteredAgencies {
		results = append(results, ValidateAgency(*a)...)

		if StringIsNilOrEmpty(a.Id) {
			results = append(results, ValidAgencyIdRequiredWhenMultipleAgenciesResult{
				FileName: "agency.txt",
				Line:     a.LineNumber,
			})
			continue
		}

		if usedIds[*a.Id] {
			results = append(results, FieldIsNotUniqueResult{
				FileName:  "agency.txt",
				FieldName: "agency_id",
				Line:      a.LineNumber,
			})
		} else {
			usedIds[*a.Id] = true
		}
	}

	return results
}
