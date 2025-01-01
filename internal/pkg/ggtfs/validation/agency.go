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
	}{
		{"agency_name", a.GetName()},         // required
		{"agency_url", a.GetURL()},           // required
		{"agency_timezone", a.GetTimezone()}, // required
		{"agency_id", a.GetID()},
		{"agency_lang", a.GetLang()},
		{"agency_phone", a.GetPhone()},
		{"agency_fare_url", a.GetURL()},
		{"agency_email", a.GetEmail()},
	}

	for _, field := range fields {
		validationResults = append(validationResults, validateField("agency", field.fieldName, field.field)...)
	}

	return validationResults
}
