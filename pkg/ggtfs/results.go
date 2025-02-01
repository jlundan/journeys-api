package ggtfs

import (
	"fmt"
)

type ValidationNotice interface {
	Code() string
	Severity() ValidationNoticeSeverity
	AsText() string
}

type SingleLineNotice struct {
	FileName  string
	FieldName string
	Line      int
}

type InvalidCharacterNotice struct {
	SingleLineNotice
}

func (n InvalidCharacterNotice) Code() string {
	return "invalid_characters"
}
func (n InvalidCharacterNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidCharacterNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidURLNotice struct {
	SingleLineNotice
}

func (n InvalidURLNotice) Code() string {
	return "invalid_url"
}
func (n InvalidURLNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidURLNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidColorNotice struct {
	SingleLineNotice
}

func (n InvalidColorNotice) Code() string {
	return "invalid_color"
}
func (n InvalidColorNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidColorNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidIntegerNotice struct {
	SingleLineNotice
}

func (n InvalidIntegerNotice) Code() string {
	return "invalid_integer"
}
func (n InvalidIntegerNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidIntegerNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidFloatNotice struct {
	SingleLineNotice
}

func (n InvalidFloatNotice) Code() string {
	return "invalid_float"
}
func (n InvalidFloatNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidFloatNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidTimeNotice struct {
	SingleLineNotice
}

func (n InvalidTimeNotice) Code() string {
	return "invalid_time"
}
func (n InvalidTimeNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidTimeNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidCurrencyCodeNotice struct {
	SingleLineNotice
}

func (n InvalidCurrencyCodeNotice) Code() string {
	return "invalid_currency_code"
}
func (n InvalidCurrencyCodeNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidCurrencyCodeNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidCurrencyAmountNotice struct {
	SingleLineNotice
}

func (n InvalidCurrencyAmountNotice) Code() string {
	return "invalid_currency_amount"
}
func (n InvalidCurrencyAmountNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidCurrencyAmountNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidDateNotice struct {
	SingleLineNotice
}

func (n InvalidDateNotice) Code() string {
	return "invalid_date"
}
func (n InvalidDateNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidDateNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidLatitudeNotice struct {
	SingleLineNotice
}

func (n InvalidLatitudeNotice) Code() string {
	return "invalid_latitude"
}
func (n InvalidLatitudeNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidLatitudeNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidLongitudeNotice struct {
	SingleLineNotice
}

func (n InvalidLongitudeNotice) Code() string {
	return "invalid_longitude"
}
func (n InvalidLongitudeNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidLongitudeNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidLanguageCodeNotice struct {
	SingleLineNotice
}

func (n InvalidLanguageCodeNotice) Code() string {
	return "invalid_language_code"
}
func (n InvalidLanguageCodeNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidLanguageCodeNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidPhoneNumberNotice struct {
	SingleLineNotice
}

func (n InvalidPhoneNumberNotice) Code() string {
	return "invalid_phone_number"
}
func (n InvalidPhoneNumberNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidPhoneNumberNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidEmailNotice struct {
	SingleLineNotice
}

func (n InvalidEmailNotice) Code() string {
	return "invalid_email"
}
func (n InvalidEmailNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidEmailNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidTimezoneNotice struct {
	SingleLineNotice
}

func (n InvalidTimezoneNotice) Code() string {
	return "invalid_timezone"
}
func (n InvalidTimezoneNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidTimezoneNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidCalendarDayNotice struct {
	SingleLineNotice
}

func (n InvalidCalendarDayNotice) Code() string {
	return "invalid_calendar_day"
}
func (n InvalidCalendarDayNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidCalendarDayNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidCalendarExceptionNotice struct {
	SingleLineNotice
}

func (n InvalidCalendarExceptionNotice) Code() string {
	return "invalid_calendar_exception"
}
func (n InvalidCalendarExceptionNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidCalendarExceptionNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type MissingRequiredFieldNotice struct {
	SingleLineNotice
}

func (n MissingRequiredFieldNotice) Code() string {
	return "missing_required_field"
}
func (n MissingRequiredFieldNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n MissingRequiredFieldNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type SingleAgencyRecommendedNotice struct {
	FileName string
}

func (n SingleAgencyRecommendedNotice) Code() string {
	return "single_agency_recommended"
}
func (n SingleAgencyRecommendedNotice) Severity() ValidationNoticeSeverity {
	return SeverityRecommendation
}
func (n SingleAgencyRecommendedNotice) AsText() string {
	return fmt.Sprintf("%s in %v", n.Code(), n.FileName)
}

type ValidAgencyIdRequiredWhenMultipleAgenciesNotice struct {
	FileName string
	Line     int
}

func (n ValidAgencyIdRequiredWhenMultipleAgenciesNotice) Code() string {
	return "valid_agency_id_required_when_multiple_agencies"
}
func (n ValidAgencyIdRequiredWhenMultipleAgenciesNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n ValidAgencyIdRequiredWhenMultipleAgenciesNotice) AsText() string {
	return fmt.Sprintf("%s in %v -> agency_id (line %v)", n.Code(), n.FileName, n.Line)
}

type FieldIsNotUniqueNotice struct {
	SingleLineNotice
}

func (n FieldIsNotUniqueNotice) Code() string {
	return "field_is_not_unique"
}
func (n FieldIsNotUniqueNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n FieldIsNotUniqueNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidRouteTypeNotice struct {
	SingleLineNotice
}

func (n InvalidRouteTypeNotice) Code() string {
	return "invalid_route_type"
}
func (n InvalidRouteTypeNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidRouteTypeNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidContinuousPickupNotice struct {
	SingleLineNotice
}

func (n InvalidContinuousPickupNotice) Code() string {
	return "invalid_continuous_pickup"
}
func (n InvalidContinuousPickupNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidContinuousPickupNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidContinuousDropOffNotice struct {
	SingleLineNotice
}

func (n InvalidContinuousDropOffNotice) Code() string {
	return "invalid_continuous_drop_off"
}
func (n InvalidContinuousDropOffNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidContinuousDropOffNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidPickupTypeNotice struct {
	SingleLineNotice
}

func (n InvalidPickupTypeNotice) Code() string {
	return "invalid_pickup_type"
}
func (n InvalidPickupTypeNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidPickupTypeNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidDropOffTypeNotice struct {
	SingleLineNotice
}

func (n InvalidDropOffTypeNotice) Code() string {
	return "invalid_drop_off_type"
}
func (n InvalidDropOffTypeNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidDropOffTypeNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidTimepointNotice struct {
	SingleLineNotice
}

func (n InvalidTimepointNotice) Code() string {
	return "invalid_timepoint"
}
func (n InvalidTimepointNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidTimepointNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidDirectionIdNotice struct {
	SingleLineNotice
}

func (n InvalidDirectionIdNotice) Code() string {
	return "invalid_direction_id"
}
func (n InvalidDirectionIdNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidDirectionIdNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidWheelchairAccessibleNotice struct {
	SingleLineNotice
}

func (n InvalidWheelchairAccessibleNotice) Code() string {
	return "invalid_wheelchair_accessible"
}
func (n InvalidWheelchairAccessibleNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidWheelchairAccessibleNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidBikesAllowedNotice struct {
	SingleLineNotice
}

func (n InvalidBikesAllowedNotice) Code() string {
	return "invalid_bikes_allowed"
}
func (n InvalidBikesAllowedNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidBikesAllowedNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type MissingRouteShortNameWhenLongNameIsNotPresentNotice struct {
	SingleLineNotice
}

func (n MissingRouteShortNameWhenLongNameIsNotPresentNotice) Code() string {
	return "missing_route_short_name_when_long_name_is_not_present"
}
func (n MissingRouteShortNameWhenLongNameIsNotPresentNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n MissingRouteShortNameWhenLongNameIsNotPresentNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type MissingRouteLongNameWhenShortNameIsNotPresentNotice struct {
	SingleLineNotice
}

func (n MissingRouteLongNameWhenShortNameIsNotPresentNotice) Code() string {
	return "missing_route_long_name_when_short_name_is_not_present"
}
func (n MissingRouteLongNameWhenShortNameIsNotPresentNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n MissingRouteLongNameWhenShortNameIsNotPresentNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type TooLongRouteShortNameNotice struct {
	SingleLineNotice
}

func (n TooLongRouteShortNameNotice) Code() string {
	return "too_long_route_short_name"
}
func (n TooLongRouteShortNameNotice) Severity() ValidationNoticeSeverity {
	return SeverityRecommendation
}
func (n TooLongRouteShortNameNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type RouteDescriptionDuplicatesNameNotice struct {
	SingleLineNotice
	DuplicatingField string
}

func (n RouteDescriptionDuplicatesNameNotice) Code() string {
	return "description_duplicates_route_name"
}
func (n RouteDescriptionDuplicatesNameNotice) Severity() ValidationNoticeSeverity {
	return SeverityRecommendation
}
func (n RouteDescriptionDuplicatesNameNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type AgencyIdRequiredForRouteWhenMultipleAgenciesNotice struct {
	SingleLineNotice
}

func (n AgencyIdRequiredForRouteWhenMultipleAgenciesNotice) Code() string {
	return "agency_id_required_for_route_when_multiple_agencies"
}
func (n AgencyIdRequiredForRouteWhenMultipleAgenciesNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n AgencyIdRequiredForRouteWhenMultipleAgenciesNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type AgencyIdRecommendedForRouteNotice struct {
	SingleLineNotice
}

func (n AgencyIdRecommendedForRouteNotice) Code() string {
	return "agency_id_recommended_for_route"
}
func (n AgencyIdRecommendedForRouteNotice) Severity() ValidationNoticeSeverity {
	return SeverityRecommendation
}
func (n AgencyIdRecommendedForRouteNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidLocationTypeNotice struct {
	SingleLineNotice
}

func (n InvalidLocationTypeNotice) Code() string {
	return "invalid_location_type"
}
func (n InvalidLocationTypeNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidLocationTypeNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type InvalidWheelchairBoardingValueNotice struct {
	SingleLineNotice
}

func (n InvalidWheelchairBoardingValueNotice) Code() string {
	return "invalid_wheelchair_boarding_value"
}
func (n InvalidWheelchairBoardingValueNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n InvalidWheelchairBoardingValueNotice) AsText() string {
	return convertSingleLineNotice(n.Code(), n.FileName, n.FieldName, n.Line)
}

type TooFewShapePointsNotice struct {
	FileName string
	ShapeId  string
}

func (n TooFewShapePointsNotice) Code() string {
	return "too_few_shape_points"
}
func (n TooFewShapePointsNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n TooFewShapePointsNotice) AsText() string {
	return fmt.Sprintf("%s in %v (shape %v)", n.Code(), n.FileName, n.ShapeId)
}

type FieldRequiredForStopLocationTypeNotice struct {
	RequiredField string
	LocationType  string
	FileName      string
	Line          int
}

func (n FieldRequiredForStopLocationTypeNotice) Code() string {
	return "field_required_for_location_type"
}
func (n FieldRequiredForStopLocationTypeNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n FieldRequiredForStopLocationTypeNotice) AsText() string {
	return fmt.Sprintf("%s %v in %v (line %v)", n.Code(), n.RequiredField, n.FileName, n.Line)
}

type ForeignKeyViolationNotice struct {
	ReferencingFileName  string
	ReferencingFieldName string
	ReferencedFieldName  string
	ReferencedFileName   string
	OffendingValue       string
	ReferencedAtRow      int
}

func (n ForeignKeyViolationNotice) Code() string {
	return "foreign_key_violation"
}
func (n ForeignKeyViolationNotice) Severity() ValidationNoticeSeverity {
	return SeverityViolation
}
func (n ForeignKeyViolationNotice) AsText() string {
	return fmt.Sprintf("%s from %v:%v->%v(value: %v) to %v->%v", n.Code(), n.ReferencingFileName, n.ReferencedAtRow, n.ReferencingFieldName, n.OffendingValue, n.ReferencedFileName, n.ReferencedFieldName)
}

func convertSingleLineNotice(code string, fileName string, fieldName string, line int) string {
	return fmt.Sprintf("%s in %v->%v (line %v)", code, fileName, fieldName, line)
}
