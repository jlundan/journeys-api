package validation

import "github.com/jlundan/journeys-api/internal/pkg/ggtfs/types"

type DummyCsvAgency struct {
	Id       *string // agency_id 		(conditionally required)
	Name     *string // agency_name 		(required)
	URL      *string // agency_url 		(required)
	Timezone *string // agency_timezone 	(required)
	Lang     *string // agency_lang 		(optional)
	Phone    *string // agency_phone 		(optional)
	FareURL  *string // agency_fare_url 	(optional)
	Email    *string // agency_email 		(optional)
}

func (a DummyCsvAgency) GetID() types.ID {
	return types.NewID(a.Id)
}

func (a DummyCsvAgency) GetName() types.Text {
	return types.NewText(a.Name)
}

func (a DummyCsvAgency) GetURL() types.URL {
	return types.NewURL(a.URL)
}

func (a DummyCsvAgency) GetTimezone() types.Timezone {
	return types.NewTimezone(a.Timezone)
}

func (a DummyCsvAgency) GetLang() types.LanguageCode {
	return types.NewLanguageCode(a.Lang)
}

func (a DummyCsvAgency) GetPhone() types.PhoneNumber {
	return types.NewPhoneNumber(a.Phone)
}

func (a DummyCsvAgency) GetFareURL() types.URL {
	return types.NewURL(a.FareURL)
}

func (a DummyCsvAgency) GetEmail() types.Email {
	return types.NewEmail(a.Email)
}
