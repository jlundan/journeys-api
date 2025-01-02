package ggtfs

import (
	"net/mail"
	"net/url"
	"regexp"
	"unicode/utf8"
)

func validateURL(fieldName string, fieldValue string, fileName string, line int) []Result {
	_, err := url.ParseRequestURI(fieldValue)
	if err != nil {
		return []Result{InvalidURLResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}
	}

	return []Result{}
}

func validateTimezone(fieldName string, fieldValue string, fileName string, line int) []Result {
	// Basic regex to validate Continent/City or Continent/City_Name format.
	match, _ := regexp.MatchString(`^[A-Za-z]+/[A-Za-z_]+$|^[A-Za-z]+/[A-Za-z]+$`, fieldValue)
	if !match {
		return []Result{InvalidTimezoneResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}
	}

	return []Result{}
}

func validateLanguageCode(fieldName string, fieldValue string, fileName string, line int) []Result {
	// Basic validation for language codes: e.g., "en", "en-US"
	match, _ := regexp.MatchString(`^[a-zA-Z]{2,3}(-[a-zA-Z]{2,3})?$`, fieldValue)
	if !match {
		return []Result{InvalidLanguageCodeResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}
	}

	return []Result{}
}

func validatePhoneNumber(fieldName string, fieldValue string, fileName string, line int) []Result {
	// Check for minimum length, only contains digits, and common phone number symbols
	match, _ := regexp.MatchString(`^[\d\s\-+()]{5,}$`, fieldValue)
	if !match {
		return []Result{InvalidPhoneNumberResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}
	}

	return []Result{}
}

func validateEmail(fieldName string, fieldValue string, fileName string, line int) []Result {
	_, err := mail.ParseAddress(fieldValue)
	if err != nil {
		return []Result{InvalidEmailResult{
			FileName:  fileName,
			FieldName: fieldName,
			Line:      line,
		}}
	}

	return []Result{}
}

func validateField(fieldType FieldType, fieldName string, fieldValue *string, isRequired bool, fileName string, line int) []Result {
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
