package validation

import (
	"github.com/jlundan/journeys-api/internal/pkg/ggtfs/types"
	"unicode/utf8"
)

func validateText(text string) []Result {
	if validUtf8 := utf8.ValidString(text); !validUtf8 {
		return []Result{&InvalidCharactersResult{}}
	}

	return nil
}

type validatedField interface {
	IsPresent() bool
	IsEmpty() bool
	Raw() string
}

func validateField(fieldName string, field validatedField, isRequired bool) []Result {
	hasValue := field.IsPresent() && !field.IsEmpty()

	if !isRequired && !hasValue {
		return []Result{}
	}

	if isRequired && !hasValue {
		return []Result{MissingRequiredFieldResult{FieldName: fieldName}}
	}

	// hasValue is true implicitly here

	switch field.(type) {
	case types.Text:
		return validateText(field.Raw())
	case types.ID:
		return validateText(field.Raw())
	case types.URL:
		return validateText(field.Raw())
	case types.Timezone:
		return validateText(field.Raw())
	case types.LanguageCode:
		return validateText(field.Raw())
	case types.PhoneNumber:
		return validateText(field.Raw())
	case types.Email:
		return validateText(field.Raw())
	default:
		return []Result{}
	}
}
