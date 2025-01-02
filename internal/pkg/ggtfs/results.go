package ggtfs

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

type InvalidURLResult struct {
	FileName  string
	FieldName string
	Line      int
}

func (e InvalidURLResult) Code() string {
	return "invalid_url"
}

type InvalidLanguageCodeResult struct {
	FileName  string
	FieldName string
	Line      int
}

func (e InvalidLanguageCodeResult) Code() string {
	return "invalid_language_code"
}

type InvalidPhoneNumberResult struct {
	FileName  string
	FieldName string
	Line      int
}

func (e InvalidPhoneNumberResult) Code() string {
	return "invalid_phone_number"
}

type InvalidEmailResult struct {
	FileName  string
	FieldName string
	Line      int
}

func (e InvalidEmailResult) Code() string {
	return "invalid_email"
}

type InvalidTimezoneResult struct {
	FileName  string
	FieldName string
	Line      int
}

func (e InvalidTimezoneResult) Code() string {
	return "invalid_timezone"
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
