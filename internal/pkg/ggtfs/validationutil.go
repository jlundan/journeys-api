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
		return []Result{InvalidCurrencyAmountResult{SingleLineResult{
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

	}

	return results
}

type FieldTobeValidated interface {
	IsValid() bool
	IsPresent() bool
	IsEmpty() bool
}

func validateRequiredFields(fields map[string]FieldTobeValidated, validationErrors *[]error, lineNumber int, fileName string) {
	for name, value := range fields {
		if !value.IsValid() {
			*validationErrors = append(*validationErrors, createFileRowError(fileName, lineNumber, createInvalidRequiredFieldString(name)))
		}
	}
}

func validateOptionalFields(fields map[string]FieldTobeValidated, validationErrors *[]error, lineNumber int, fileName string) {
	for name, value := range fields {
		if value.IsPresent() && !value.IsEmpty() && !value.IsValid() {
			*validationErrors = append(*validationErrors, createFileRowError(fileName, lineNumber, createInvalidFieldString(name)))
		}
	}
}
