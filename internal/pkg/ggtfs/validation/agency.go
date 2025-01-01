package validation

import (
	"github.com/jlundan/journeys-api/internal/pkg/ggtfs/types"
)

type Agency interface {
	GetID() types.ID
	GetName() types.Text
	GetURL() types.URL
	GetTimezone() types.Timezone
	GetLang() types.LanguageCode
	GetPhone() types.PhoneNumber
	GetFareURL() types.URL
	GetEmail() types.Email
}

func ValidateAgency(a Agency) []Result {
	var validationResults []Result

	fields := []struct {
		fieldName string
		field     validatedField
		required  bool
	}{
		{"agency_name", a.GetName(), true},
		{"agency_url", a.GetURL(), true},
		{"agency_timezone", a.GetTimezone(), true},
		{"agency_id", a.GetID(), false},
		{"agency_lang", a.GetLang(), false},
		{"agency_phone", a.GetPhone(), false},
		{"agency_fare_url", a.GetURL(), false},
		{"agency_email", a.GetEmail(), false},
	}

	for _, field := range fields {
		validationResults = append(validationResults, validateField(field.fieldName, field.field, field.required)...)
	}

	return validationResults
}
