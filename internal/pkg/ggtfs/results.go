package ggtfs

type Result interface {
	Code() string
}

type SingleLineResult struct {
	FileName  string
	FieldName string
	Line      int
}

type InvalidCharactersResult struct {
	SingleLineResult
}

func (e InvalidCharactersResult) Code() string {
	return "invalid_characters"
}

type InvalidURLResult struct {
	SingleLineResult
}

func (e InvalidURLResult) Code() string {
	return "invalid_url"
}

type InvalidColorResult struct {
	SingleLineResult
}

func (e InvalidColorResult) Code() string {
	return "invalid_color"
}

type InvalidIntegerResult struct {
	SingleLineResult
}

func (e InvalidIntegerResult) Code() string {
	return "invalid_integer"
}

type InvalidFloatResult struct {
	SingleLineResult
}

func (e InvalidFloatResult) Code() string {
	return "invalid_float"
}

type InvalidTimeResult struct {
	SingleLineResult
}

func (e InvalidTimeResult) Code() string {
	return "invalid_time"
}

type InvalidCurrencyCodeResult struct {
	SingleLineResult
}

func (e InvalidCurrencyCodeResult) Code() string {
	return "invalid_currency_code"
}

type InvalidCurrencyAmountResult struct {
	SingleLineResult
}

func (e InvalidCurrencyAmountResult) Code() string {
	return "invalid_currency_amount"
}

type InvalidDateResult struct {
	SingleLineResult
}

func (e InvalidDateResult) Code() string {
	return "invalid_date"
}

type InvalidLatitudeResult struct {
	SingleLineResult
}

func (e InvalidLatitudeResult) Code() string {
	return "invalid_latitude"
}

type InvalidLongitudeResult struct {
	SingleLineResult
}

func (e InvalidLongitudeResult) Code() string {
	return "invalid_longitude"
}

type InvalidLanguageCodeResult struct {
	SingleLineResult
}

func (e InvalidLanguageCodeResult) Code() string {
	return "invalid_language_code"
}

type InvalidPhoneNumberResult struct {
	SingleLineResult
}

func (e InvalidPhoneNumberResult) Code() string {
	return "invalid_phone_number"
}

type InvalidEmailResult struct {
	SingleLineResult
}

func (e InvalidEmailResult) Code() string {
	return "invalid_email"
}

type InvalidTimezoneResult struct {
	SingleLineResult
}

func (e InvalidTimezoneResult) Code() string {
	return "invalid_timezone"
}

type InvalidCalendarDayResult struct {
	SingleLineResult
}

func (e InvalidCalendarDayResult) Code() string {
	return "invalid_calendar_day"
}

type MissingRequiredFieldResult struct {
	SingleLineResult
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
	SingleLineResult
}

func (e FieldIsNotUniqueResult) Code() string {
	return "field_is_not_unique"
}
