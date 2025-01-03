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

type InvalidCalendarExceptionResult struct {
	SingleLineResult
}

func (e InvalidCalendarExceptionResult) Code() string {
	return "invalid_calendar_exception"
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

type InvalidRouteTypeResult struct {
	SingleLineResult
}

func (e InvalidRouteTypeResult) Code() string {
	return "invalid_route_type"
}

type InvalidContinuousPickupResult struct {
	SingleLineResult
}

func (e InvalidContinuousPickupResult) Code() string {
	return "invalid_continuous_pickup"
}

type InvalidContinuousDropOffResult struct {
	SingleLineResult
}

func (e InvalidContinuousDropOffResult) Code() string {
	return "invalid_continuous_drop_off"
}

type MissingRouteShortNameWhenLongNameIsNotPresentResult struct {
	SingleLineResult
}

func (e MissingRouteShortNameWhenLongNameIsNotPresentResult) Code() string {
	return "missing_route_short_name_when_long_name_is_not_present"
}

type MissingRouteLongNameWhenShortNameIsNotPresentResult struct {
	SingleLineResult
}

func (e MissingRouteLongNameWhenShortNameIsNotPresentResult) Code() string {
	return "missing_route_long_name_when_short_name_is_not_present"
}

type TooLongRouteShortNameResult struct {
	SingleLineResult
}

func (e TooLongRouteShortNameResult) Code() string {
	return "too_long_route_short_name"
}

type DescriptionDuplicatesRouteNameResult struct {
	SingleLineResult
	DuplicatingField string
}

func (e DescriptionDuplicatesRouteNameResult) Code() string {
	return "description_duplicates_route_name"
}

type AgencyIdRequiredForRouteWhenMultipleAgenciesResult struct {
	SingleLineResult
}

func (e AgencyIdRequiredForRouteWhenMultipleAgenciesResult) Code() string {
	return "agency_id_required_for_route_when_multiple_agencies"
}

type AgencyIdRecommendedForRouteResult struct {
	SingleLineResult
}

func (e AgencyIdRecommendedForRouteResult) Code() string {
	return "agency_id_recommended_for_route"
}

type ForeignKeyViolationResult struct {
	ReferencingFileName  string
	ReferencingFieldName string
	ReferencedFieldName  string
	ReferencedFileName   string
	OffendingValue       string
	ReferencedAtRow      int
}

func (e ForeignKeyViolationResult) Code() string {
	return "foreign_key_violation"
}
