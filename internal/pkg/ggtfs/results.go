package ggtfs

import "unicode/utf8"

type Result interface {
	Code() string
}

type InvalidCharactersResult struct {
	FileName  string
	FieldName string
	Line      int
}

func (e InvalidCharactersResult) Code() string {
	return "invalid_characters"
}

type MissingRequiredFieldResult struct {
	FileName  string
	FieldName string
	Line      int
}

func (e MissingRequiredFieldResult) Code() string {
	return "missing_required_field"
}

type SingleAgencyRecommendedResult struct {
	FileName string
}

func (e SingleAgencyRecommendedResult) Code() string {
	return "single_agency_recommended"
}

type ValidAgencyIdRequiredWhenMultipleAgenciesResult struct {
	FileName string
	Line     int
}

func (e ValidAgencyIdRequiredWhenMultipleAgenciesResult) Code() string {
	return "valid_agency_id_required_when_multiple_agencies"
}

type FieldIsNotUniqueResult struct {
	FileName  string
	FieldName string
	Line      int
}

func (e FieldIsNotUniqueResult) Code() string {
	return "field_is_not_unique"
}

func validateField(fieldName string, fieldValue *string, isRequired bool, fileName string, line int) []Result {
	hasValue := fieldValue != nil && *fieldValue != ""

	if !isRequired && !hasValue {
		return []Result{}
	}

	if isRequired && !hasValue {
		return []Result{MissingRequiredFieldResult{FileName: fileName, FieldName: fieldName, Line: line}}
	}

	// hasValue is true implicitly here

	var results []Result

	if !utf8.ValidString(*fieldValue) {
		results = append(results, &InvalidCharactersResult{FileName: fileName, FieldName: fieldName, Line: line})
	}

	//switch field.(type) {
	//case Text:
	//	return validateText(field.Raw())
	//case ID:
	//	return validateText(field.Raw())
	//case URL:
	//	return validateText(field.Raw())
	//case Timezone:
	//	return validateText(field.Raw())
	//case LanguageCode:
	//	return validateText(field.Raw())
	//case PhoneNumber:
	//	return validateText(field.Raw())
	//case Email:
	//	return validateText(field.Raw())
	//default:
	//	return []Result{}
	//}

	return results
}
