package validation

//func ValidateAgencies(agencies []Agency) ([]error, []string) {
//	var validationErrors []error
//	var recommendations []string
//
//	var filteredAgencies []Agency
//
//	for _, a := range agencies {
//		if a != nil {
//			filteredAgencies = append(filteredAgencies, a)
//		}
//	}
//
//	aLength := len(filteredAgencies)
//
//	if aLength == 0 {
//		return validationErrors, recommendations
//	}
//
//	if aLength == 1 && !filteredAgencies[0].GetID().IsValid() {
//		validationErrors = append(validationErrors, ValidateAgency(filteredAgencies[0])...)
//		recommendations = append(recommendations, createFileRowRecommendation(AgenciesFileName, agencies[0].LineNumber, "it is recommended that agency_id is specified even when there is only one agency"))
//		return validationErrors, recommendations
//	}
//
//	if aLength == 1 {
//		validationErrors = append(validationErrors, ValidateAgency(filteredAgencies[0])...)
//		return validationErrors, recommendations
//	}
//
//	usedIds := make(map[string]bool)
//	for _, a := range filteredAgencies {
//		validationErrors = append(validationErrors, ValidateAgency(a)...)
//
//		aID := a.GetID()
//		if !aID.IsValid() {
//			validationErrors = append(validationErrors, createFileRowError(AgenciesFileName, a.LineNumber, "a valid agency_id must be specified when multiple agencies are declared"))
//			continue
//		}
//
//		r := aID.Raw()
//		if usedIds[r] {
//			validationErrors = append(validationErrors, createFileRowError(AgenciesFileName, a.LineNumber, fmt.Sprintf("agency_id is not unique within the file")))
//		} else {
//			usedIds[r] = true
//		}
//	}
//
//	return validationErrors, recommendations
//}
