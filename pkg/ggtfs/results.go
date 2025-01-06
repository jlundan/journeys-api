package ggtfs

type ValidationNotice interface {
	Code() string
	Severity() ValidationNoticeSeverity
}

type SingleLineNotice struct {
	FileName  string
	FieldName string
	Line      int
}

type InvalidCharacterNotice struct {
	SingleLineNotice
}

func (r InvalidCharacterNotice) Code() string {
	return "invalid_characters"
}
func (r InvalidCharacterNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type InvalidURLNotice struct {
	SingleLineNotice
}

func (r InvalidURLNotice) Code() string {
	return "invalid_url"
}
func (r InvalidURLNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type InvalidColorNotice struct {
	SingleLineNotice
}

func (r InvalidColorNotice) Code() string {
	return "invalid_color"
}
func (r InvalidColorNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type InvalidIntegerNotice struct {
	SingleLineNotice
}

func (r InvalidIntegerNotice) Code() string {
	return "invalid_integer"
}
func (r InvalidIntegerNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type InvalidFloatNotice struct {
	SingleLineNotice
}

func (r InvalidFloatNotice) Code() string {
	return "invalid_float"
}
func (r InvalidFloatNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type InvalidTimeNotice struct {
	SingleLineNotice
}

func (r InvalidTimeNotice) Code() string {
	return "invalid_time"
}
func (r InvalidTimeNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type InvalidCurrencyCodeNotice struct {
	SingleLineNotice
}

func (r InvalidCurrencyCodeNotice) Code() string {
	return "invalid_currency_code"
}
func (r InvalidCurrencyCodeNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type InvalidCurrencyAmountNotice struct {
	SingleLineNotice
}

func (r InvalidCurrencyAmountNotice) Code() string {
	return "invalid_currency_amount"
}
func (r InvalidCurrencyAmountNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type InvalidDateNotice struct {
	SingleLineNotice
}

func (r InvalidDateNotice) Code() string {
	return "invalid_date"
}
func (r InvalidDateNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type InvalidLatitudeNotice struct {
	SingleLineNotice
}

func (r InvalidLatitudeNotice) Code() string {
	return "invalid_latitude"
}
func (r InvalidLatitudeNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type InvalidLongitudeNotice struct {
	SingleLineNotice
}

func (r InvalidLongitudeNotice) Code() string {
	return "invalid_longitude"
}
func (r InvalidLongitudeNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type InvalidLanguageCodeNotice struct {
	SingleLineNotice
}

func (r InvalidLanguageCodeNotice) Code() string {
	return "invalid_language_code"
}
func (r InvalidLanguageCodeNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type InvalidPhoneNumberNotice struct {
	SingleLineNotice
}

func (r InvalidPhoneNumberNotice) Code() string {
	return "invalid_phone_number"
}
func (r InvalidPhoneNumberNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type InvalidEmailNotice struct {
	SingleLineNotice
}

func (r InvalidEmailNotice) Code() string {
	return "invalid_email"
}
func (r InvalidEmailNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type InvalidTimezoneNotice struct {
	SingleLineNotice
}

func (r InvalidTimezoneNotice) Code() string {
	return "invalid_timezone"
}
func (r InvalidTimezoneNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type InvalidCalendarDayNotice struct {
	SingleLineNotice
}

func (r InvalidCalendarDayNotice) Code() string {
	return "invalid_calendar_day"
}
func (r InvalidCalendarDayNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type InvalidCalendarExceptionNotice struct {
	SingleLineNotice
}

func (r InvalidCalendarExceptionNotice) Code() string {
	return "invalid_calendar_exception"
}
func (r InvalidCalendarExceptionNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type MissingRequiredFieldNotice struct {
	SingleLineNotice
}

func (r MissingRequiredFieldNotice) Code() string {
	return "missing_required_field"
}
func (r MissingRequiredFieldNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type SingleAgencyRecommendedNotice struct {
	FileName string
}

func (r SingleAgencyRecommendedNotice) Code() string {
	return "single_agency_recommended"
}
func (r SingleAgencyRecommendedNotice) Severity() ValidationNoticeSeverity {
	return SeverityRecommendation
}

type ValidAgencyIdRequiredWhenMultipleAgenciesNotice struct {
	FileName string
	Line     int
}

func (r ValidAgencyIdRequiredWhenMultipleAgenciesNotice) Code() string {
	return "valid_agency_id_required_when_multiple_agencies"
}
func (r ValidAgencyIdRequiredWhenMultipleAgenciesNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type FieldIsNotUniqueNotice struct {
	SingleLineNotice
}

func (r FieldIsNotUniqueNotice) Code() string {
	return "field_is_not_unique"
}
func (r FieldIsNotUniqueNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type InvalidRouteTypeNotice struct {
	SingleLineNotice
}

func (r InvalidRouteTypeNotice) Code() string {
	return "invalid_route_type"
}
func (r InvalidRouteTypeNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type InvalidContinuousPickupNotice struct {
	SingleLineNotice
}

func (r InvalidContinuousPickupNotice) Code() string {
	return "invalid_continuous_pickup"
}
func (r InvalidContinuousPickupNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type InvalidContinuousDropOffNotice struct {
	SingleLineNotice
}

func (r InvalidContinuousDropOffNotice) Code() string {
	return "invalid_continuous_drop_off"
}
func (r InvalidContinuousDropOffNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type InvalidPickupTypeNotice struct {
	SingleLineNotice
}

func (r InvalidPickupTypeNotice) Code() string {
	return "invalid_pickup_type"
}
func (r InvalidPickupTypeNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type InvalidDropOffTypeNotice struct {
	SingleLineNotice
}

func (r InvalidDropOffTypeNotice) Code() string {
	return "invalid_drop_off_type"
}
func (r InvalidDropOffTypeNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type InvalidTimepointNotice struct {
	SingleLineNotice
}

func (r InvalidTimepointNotice) Code() string {
	return "invalid_timepoint"
}
func (r InvalidTimepointNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type InvalidDirectionIdNotice struct {
	SingleLineNotice
}

func (r InvalidDirectionIdNotice) Code() string {
	return "invalid_direction_id"
}
func (r InvalidDirectionIdNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type InvalidWheelchairAccessibleNotice struct {
	SingleLineNotice
}

func (r InvalidWheelchairAccessibleNotice) Code() string {
	return "invalid_wheelchair_accessible"
}
func (r InvalidWheelchairAccessibleNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type InvalidBikesAllowedNotice struct {
	SingleLineNotice
}

func (r InvalidBikesAllowedNotice) Code() string {
	return "invalid_bikes_allowed"
}
func (r InvalidBikesAllowedNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type MissingRouteShortNameWhenLongNameIsNotPresentNotice struct {
	SingleLineNotice
}

func (r MissingRouteShortNameWhenLongNameIsNotPresentNotice) Code() string {
	return "missing_route_short_name_when_long_name_is_not_present"
}
func (r MissingRouteShortNameWhenLongNameIsNotPresentNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type MissingRouteLongNameWhenShortNameIsNotPresentNotice struct {
	SingleLineNotice
}

func (r MissingRouteLongNameWhenShortNameIsNotPresentNotice) Code() string {
	return "missing_route_long_name_when_short_name_is_not_present"
}
func (r MissingRouteLongNameWhenShortNameIsNotPresentNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type TooLongRouteShortNameNotice struct {
	SingleLineNotice
}

func (r TooLongRouteShortNameNotice) Code() string {
	return "too_long_route_short_name"
}
func (r TooLongRouteShortNameNotice) Severity() ValidationNoticeSeverity {
	return SeverityRecommendation
}

type RouteDescriptionDuplicatesNameNotice struct {
	SingleLineNotice
	DuplicatingField string
}

func (r RouteDescriptionDuplicatesNameNotice) Code() string {
	return "description_duplicates_route_name"
}
func (r RouteDescriptionDuplicatesNameNotice) Severity() ValidationNoticeSeverity {
	return SeverityRecommendation
}

type AgencyIdRequiredForRouteWhenMultipleAgenciesNotice struct {
	SingleLineNotice
}

func (r AgencyIdRequiredForRouteWhenMultipleAgenciesNotice) Code() string {
	return "agency_id_required_for_route_when_multiple_agencies"
}
func (r AgencyIdRequiredForRouteWhenMultipleAgenciesNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type AgencyIdRecommendedForRouteNotice struct {
	SingleLineNotice
}

func (r AgencyIdRecommendedForRouteNotice) Code() string {
	return "agency_id_recommended_for_route"
}
func (r AgencyIdRecommendedForRouteNotice) Severity() ValidationNoticeSeverity {
	return SeverityRecommendation
}

type InvalidLocationTypeNotice struct {
	SingleLineNotice
}

func (r InvalidLocationTypeNotice) Code() string {
	return "invalid_location_type"
}
func (r InvalidLocationTypeNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type InvalidWheelchairBoardingValueNotice struct {
	SingleLineNotice
}

func (r InvalidWheelchairBoardingValueNotice) Code() string {
	return "invalid_wheelchair_boarding_value"
}
func (r InvalidWheelchairBoardingValueNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type TooFewShapePointsNotice struct {
	FileName string
	ShapeId  string
}

func (r TooFewShapePointsNotice) Code() string {
	return "too_few_shape_points"
}
func (r TooFewShapePointsNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type FieldRequiredForStopLocationTypeNotice struct {
	RequiredField string
	LocationType  string
	FileName      string
	Line          int
}

func (r FieldRequiredForStopLocationTypeNotice) Code() string {
	return "field_required_for_location_type"
}
func (r FieldRequiredForStopLocationTypeNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}

type ForeignKeyViolationNotice struct {
	ReferencingFileName  string
	ReferencingFieldName string
	ReferencedFieldName  string
	ReferencedFileName   string
	OffendingValue       string
	ReferencedAtRow      int
}

func (r ForeignKeyViolationNotice) Code() string {
	return "foreign_key_violation"
}
func (r ForeignKeyViolationNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
