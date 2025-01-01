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

	fields := map[string]validatedField{
		"agency_name":     a.GetName(),     // agency_name 		(required)
		"agency_url":      a.GetURL(),      // agency_url 		(required)
		"agency_timezone": a.GetTimezone(), // agency_timezone 	(required)
		"agency_id":       a.GetID(),
		"agency_lang":     a.GetLang(),
		"agency_phone":    a.GetPhone(),
		"agency_fare_url": a.GetURL(),
		"agency_email":    a.GetEmail(),
	}

	for fieldName, field := range fields {
		validationResults = append(validationResults, validateField("agency", fieldName, field)...)
	}

	return validationResults
}

//validationResults = append(validationResults, validateField("agency", "agency_name", a.GetName())...)
//validationResults = append(validationResults, validateField("agency", "agency_url", a.GetURL())...)
//validationResults = append(validationResults, validateField("agency", "agency_timezone", a.GetTimezone())...)
//validationResults = append(validationResults, validateField("agency", "agency_id", a.GetID())...)
//validationResults = append(validationResults, validateField("agency", "agency_lang", a.GetLang())...)
//validationResults = append(validationResults, validateField("agency", "agency_phone", a.GetPhone())...)
//validationResults = append(validationResults, validateField("agency", "agency_fare_url", a.GetURL())...)
//validationResults = append(validationResults, validateField("agency", "agency_email", a.GetEmail())...)
