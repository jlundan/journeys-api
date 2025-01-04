package ggtfs

import (
	"net/mail"
	"net/url"
	"regexp"
	"strconv"
	"unicode/utf8"
)

func validateURL(fieldName string, fieldValue string, fileName string, line int) []Result {
	_, err := url.ParseRequestURI(fieldValue)
	if err != nil {
		return []Result{InvalidURLResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validateColor(fieldName string, fieldValue string, fileName string, line int) []Result {
	match, _ := regexp.MatchString(`^[0-9A-Fa-f]{6}$`, fieldValue)
	if !match {
		return []Result{InvalidColorResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validateTimezone(fieldName string, fieldValue string, fileName string, line int) []Result {
	// Basic regex to validate Continent/City or Continent/City_Name format.
	match, _ := regexp.MatchString(`^[A-Za-z]+/[A-Za-z_]+$|^[A-Za-z]+/[A-Za-z]+$`, fieldValue)
	if !match {
		return []Result{InvalidTimezoneResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validateLanguageCode(fieldName string, fieldValue string, fileName string, line int) []Result {
	// Basic validation for language codes: e.g., "en", "en-US"
	match, _ := regexp.MatchString(`^[a-zA-Z]{2,3}(-[a-zA-Z]{2,3})?$`, fieldValue)
	if !match {
		return []Result{InvalidLanguageCodeResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validatePhoneNumber(fieldName string, fieldValue string, fileName string, line int) []Result {
	// Check for minimum length, only contains digits, and common phone number symbols
	match, _ := regexp.MatchString(`^[\d\s\-+()]{5,}$`, fieldValue)
	if !match {
		return []Result{InvalidPhoneNumberResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validateEmail(fieldName string, fieldValue string, fileName string, line int) []Result {
	_, err := mail.ParseAddress(fieldValue)
	if err != nil {
		return []Result{InvalidEmailResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validateLatitude(fieldName string, fieldValue string, fileName string, line int) []Result {
	value, err := strconv.ParseFloat(fieldValue, 64)
	if err != nil || value < -90.0 && value > 90.0 {
		return []Result{InvalidLatitudeResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validateLongitude(fieldName string, fieldValue string, fileName string, line int) []Result {
	value, err := strconv.ParseFloat(fieldValue, 64)
	if err != nil || value < -180.0 && value > 180.0 {
		return []Result{InvalidLongitudeResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validateDate(fieldName string, fieldValue string, fileName string, line int) []Result {
	match, _ := regexp.MatchString(`^(19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12][0-9]|3[01])$`, fieldValue)
	if !match {
		return []Result{InvalidDateResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validateFloat(fieldName string, fieldValue string, fileName string, line int) []Result {
	_, err := strconv.ParseFloat(fieldValue, 64)
	if err != nil {
		return []Result{InvalidFloatResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validateInteger(fieldName string, fieldValue string, fileName string, line int) []Result {
	_, err := strconv.Atoi(fieldValue)
	if err != nil {
		return []Result{InvalidIntegerResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validateCurrencyAmount(fieldName string, fieldValue string, fileName string, line int) []Result {
	_, err := strconv.ParseFloat(fieldValue, 64)
	if err != nil {
		return []Result{InvalidCurrencyAmountResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validateCurrencyCode(fieldName string, fieldValue string, fileName string, line int) []Result {
	match, _ := regexp.MatchString(`^[A-Z]{3}$`, fieldValue)
	if !match {
		return []Result{InvalidCurrencyAmountResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validateTime(fieldName string, fieldValue string, fileName string, line int) []Result {
	// Checks if the Time is in the valid HH:MM:SS or H:MM:SS format. The hour is between 0 and 47, since the trips on the service day might run
	// through the night. For example, 25:00:00 represents 1:00:00 AM the next day.
	match, _ := regexp.MatchString(`^(0[0-9]|1[0-9]|2[0-9]|3[0-9]|4[0-7]):([0-5][0-9]):([0-5][0-9])$|^([0-9]|1[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9])$`, fieldValue)
	if !match {
		return []Result{InvalidTimeResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validateCalendarDay(fieldName string, fieldValue string, fileName string, line int) []Result {
	i, err := strconv.Atoi(fieldValue)
	if err != nil || i < 0 || i > 1 {
		return []Result{InvalidCalendarDayResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validateCalendarException(fieldName string, fieldValue string, fileName string, line int) []Result {
	i, err := strconv.Atoi(fieldValue)
	if err != nil || i < 1 || i > 2 {
		return []Result{InvalidCalendarExceptionResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validateRouteType(fieldName string, fieldValue string, fileName string, line int) []Result {
	i, err := strconv.Atoi(fieldValue)
	if err != nil || (i != 11 && i != 12 && i < 0 && i > 7) {
		return []Result{InvalidRouteTypeResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validateContinuousPickup(fieldName string, fieldValue string, fileName string, line int) []Result {
	i, err := strconv.Atoi(fieldValue)
	if err != nil || i < 1 || i > 4 {
		return []Result{InvalidContinuousPickupResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validateContinuousDropOff(fieldName string, fieldValue string, fileName string, line int) []Result {
	i, err := strconv.Atoi(fieldValue)
	if err != nil || i < 1 || i > 4 {
		return []Result{InvalidContinuousDropOffResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validateLocationType(fieldName string, fieldValue string, fileName string, line int) []Result {
	i, err := strconv.Atoi(fieldValue)
	if err != nil || i < 0 || i > 4 {
		return []Result{InvalidLocationTypeResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validateWheelchairBoarding(fieldName string, fieldValue string, fileName string, line int) []Result {
	i, err := strconv.Atoi(fieldValue)
	if err != nil || i < 0 || i > 2 {
		return []Result{InvalidWheelchairBoardingResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validatePickupType(fieldName string, fieldValue string, fileName string, line int) []Result {
	i, err := strconv.Atoi(fieldValue)
	if err != nil || i < 0 || i > 3 {
		return []Result{InvalidPickupTypeResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validateDropOffType(fieldName string, fieldValue string, fileName string, line int) []Result {
	i, err := strconv.Atoi(fieldValue)
	if err != nil || i < 0 || i > 3 {
		return []Result{InvalidDropOffTypeResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validateTimepoint(fieldName string, fieldValue string, fileName string, line int) []Result {
	i, err := strconv.Atoi(fieldValue)
	if err != nil || i < 0 || i > 1 {
		return []Result{InvalidTimepointResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validateDirectionId(fieldName string, fieldValue string, fileName string, line int) []Result {
	i, err := strconv.Atoi(fieldValue)
	if err != nil || i < 0 || i > 1 {
		return []Result{InvalidDirectionIdResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validateWheelchairAccessible(fieldName string, fieldValue string, fileName string, line int) []Result {
	i, err := strconv.Atoi(fieldValue)
	if err != nil || i < 0 || i > 2 {
		return []Result{InvalidWheelchairAccessibleResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validateTypeBikesAllowed(fieldName string, fieldValue string, fileName string, line int) []Result {
	i, err := strconv.Atoi(fieldValue)
	if err != nil || i < 0 || i > 2 {
		return []Result{InvalidBikesAllowedResult{SingleLineResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}}
	}

	return []Result{}
}

func validateField(fieldType FieldType, fieldName string, fieldValue *string, isRequired bool, fileName string, line int) []Result {
	hasValue := fieldValue != nil && *fieldValue != ""

	if !isRequired && !hasValue {
		return []Result{}
	}

	if isRequired && !hasValue {
		return []Result{MissingRequiredFieldResult{SingleLineResult{FileName: fileName, FieldName: fieldName, Line: line}}}
	}

	// hasValue is true implicitly here

	var results []Result

	if !utf8.ValidString(*fieldValue) {
		results = append(results, &InvalidCharactersResult{SingleLineResult{FileName: fileName, FieldName: fieldName, Line: line}})
	}

	switch fieldType {
	case FieldTypeURL:
		results = append(results, validateURL(fieldName, *fieldValue, fileName, line)...)
	case FieldTypeTimezone:
		results = append(results, validateTimezone(fieldName, *fieldValue, fileName, line)...)
	case FieldTypeLanguageCode:
		results = append(results, validateLanguageCode(fieldName, *fieldValue, fileName, line)...)
	case FieldTypePhoneNumber:
		results = append(results, validatePhoneNumber(fieldName, *fieldValue, fileName, line)...)
	case FieldTypeEmail:
		results = append(results, validateEmail(fieldName, *fieldValue, fileName, line)...)
	case FieldTypeColor:
		results = append(results, validateColor(fieldName, *fieldValue, fileName, line)...)
	case FieldTypeInteger:
		results = append(results, validateInteger(fieldName, *fieldValue, fileName, line)...)
	case FieldTypeFloat:
		results = append(results, validateFloat(fieldName, *fieldValue, fileName, line)...)
	case FieldTypeTime:
		results = append(results, validateTime(fieldName, *fieldValue, fileName, line)...)
	case FieldTypeCurrencyCode:
		results = append(results, validateCurrencyCode(fieldName, *fieldValue, fileName, line)...)
	case FieldTypeCurrencyAmount:
		results = append(results, validateCurrencyAmount(fieldName, *fieldValue, fileName, line)...)
	case FieldTypeDate:
		results = append(results, validateDate(fieldName, *fieldValue, fileName, line)...)
	case FieldTypeLatitude:
		results = append(results, validateLatitude(fieldName, *fieldValue, fileName, line)...)
	case FieldTypeLongitude:
		results = append(results, validateLongitude(fieldName, *fieldValue, fileName, line)...)
	case FieldTypeCalendarDay:
		results = append(results, validateCalendarDay(fieldName, *fieldValue, fileName, line)...)
	case FieldTypeCalendarException:
		results = append(results, validateCalendarException(fieldName, *fieldValue, fileName, line)...)
	case FieldTypeRouteType:
		results = append(results, validateRouteType(fieldName, *fieldValue, fileName, line)...)
	case FieldTypeContinuousPickup:
		results = append(results, validateContinuousPickup(fieldName, *fieldValue, fileName, line)...)
	case FieldTypeContinuousDropOff:
		results = append(results, validateContinuousDropOff(fieldName, *fieldValue, fileName, line)...)
	case FieldTypeLocationType:
		results = append(results, validateLocationType(fieldName, *fieldValue, fileName, line)...)
	case FieldTypeWheelchairBoarding:
		results = append(results, validateWheelchairBoarding(fieldName, *fieldValue, fileName, line)...)
	case FieldTypePickupType:
		results = append(results, validatePickupType(fieldName, *fieldValue, fileName, line)...)
	case FieldTypeDropOffType:
		results = append(results, validateDropOffType(fieldName, *fieldValue, fileName, line)...)
	case FieldTypeTimepoint:
		results = append(results, validateTimepoint(fieldName, *fieldValue, fileName, line)...)
	case FieldTypeDirectionId:
		results = append(results, validateDirectionId(fieldName, *fieldValue, fileName, line)...)
	case FieldTypeWheelchairAccessible:
		results = append(results, validateWheelchairAccessible(fieldName, *fieldValue, fileName, line)...)
	case FieldTypeBikesAllowed:
		results = append(results, validateTypeBikesAllowed(fieldName, *fieldValue, fileName, line)...)
	}

	return results
}

func IsLocationTypeValid(locationType *string) bool {
	if locationType == nil {
		return false
	}

	i, err := strconv.Atoi(*locationType)
	if err != nil {
		return false
	}

	return i >= 0 && i <= 4
}
