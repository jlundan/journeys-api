//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"encoding/csv"
	"strings"
	"testing"
)

var validAgencyHeaders = []string{"agency_id", "agency_name", "agency_url", "agency_timezone",
	"agency_lang", "agency_phone", "agency_fare_url", "agency_email"}

func TestShouldReturnEmptyAgencyArrayOnEmptyString(t *testing.T) {
	agencies, errors := LoadEntities[*Agency](csv.NewReader(strings.NewReader("")), validAgencyHeaders, CreateAgency, AgenciesFileName)
	if len(errors) > 0 {
		t.Error(errors)
	}
	if len(agencies) != 0 {
		t.Error("expected zero agencies")
	}
}

func TestAgencyParsing(t *testing.T) {
	loadAgenciesFunc := func(reader *csv.Reader) ([]interface{}, []error) {
		agencies, errs := LoadEntities[*Agency](reader, validAgencyHeaders, CreateAgency, AgenciesFileName)
		entities := make([]interface{}, len(agencies))
		for i, agency := range agencies {
			entities[i] = agency
		}
		return entities, errs
	}

	validateAgenciesFunc := func(entities []interface{}, _fixtures map[string][]interface{}) ([]error, []string) {
		agencies := make([]*Agency, len(entities))
		for i, entity := range entities {
			if agency, ok := entity.(*Agency); ok {
				agencies[i] = agency
			}
		}
		return ValidateAgencies(agencies)
	}

	runGenericGTFSParseTest(t, "AgencyNOKTestcases", loadAgenciesFunc, validateAgenciesFunc, false, getAgencyNOKTestcases())
	runGenericGTFSParseTest(t, "AgencyOKTestcases", loadAgenciesFunc, validateAgenciesFunc, false, getAgencyOKTestcases())
}

func getAgencyNOKTestcases() map[string]ggtfsTestCase {
	testCases := make(map[string]ggtfsTestCase)
	testCases["parse-failures"] = ggtfsTestCase{
		csvRows: [][]string{
			{"agency_name", "agency_url", "agency_timezone"},
			{" "},
			{","},
			{"", ""},
			{" ", " "},
		},
		expectedErrors: []string{
			"agency.txt: record on line 2: wrong number of fields",
			"agency.txt: record on line 3: wrong number of fields",
			"agency.txt: record on line 4: wrong number of fields",
			"agency.txt: record on line 5: wrong number of fields",
		},
	}
	testCases["invalid-fields"] = ggtfsTestCase{
		csvRows: [][]string{
			{"agency_name", "agency_url", "agency_timezone", "agency_lang", "agency_phone", "agency_fare_url", "agency_email"},
			{",,,,,,"},
			{"", "", "", "", "", "", ""},
			{" ", " ", " ", " ", " ", " ", " "},
			{"ACME", "Not an URL", "Europe/Helsinki2", "Not a language code", "Not a phone", "Not an URL", "Not a phone number"},
		},
		expectedErrors: []string{
			"agency.txt:0: a valid agency_id must be specified when multiple agencies are declared",
			"agency.txt:0: invalid field: agency_email",
			"agency.txt:0: invalid field: agency_fare_url",
			"agency.txt:0: invalid field: agency_lang",
			"agency.txt:0: invalid field: agency_phone",
			"agency.txt:0: missing mandatory field: agency_name",
			"agency.txt:0: missing mandatory field: agency_timezone",
			"agency.txt:0: missing mandatory field: agency_url",
			"agency.txt:1: a valid agency_id must be specified when multiple agencies are declared",
			"agency.txt:1: invalid field: agency_email",
			"agency.txt:1: invalid field: agency_fare_url",
			"agency.txt:1: invalid field: agency_lang",
			"agency.txt:1: invalid field: agency_phone",
			"agency.txt:1: missing mandatory field: agency_name",
			"agency.txt:1: missing mandatory field: agency_timezone",
			"agency.txt:1: missing mandatory field: agency_url",
			"agency.txt:2: a valid agency_id must be specified when multiple agencies are declared",
			"agency.txt:2: invalid field: agency_email",
			"agency.txt:2: invalid field: agency_fare_url",
			"agency.txt:2: invalid field: agency_lang",
			"agency.txt:2: invalid field: agency_phone",
			"agency.txt:2: missing mandatory field: agency_name",
			"agency.txt:2: missing mandatory field: agency_timezone",
			"agency.txt:2: missing mandatory field: agency_url",
			"agency.txt:3: a valid agency_id must be specified when multiple agencies are declared",
			"agency.txt:3: invalid field: agency_email",
			"agency.txt:3: invalid field: agency_fare_url",
			"agency.txt:3: invalid field: agency_lang",
			"agency.txt:3: invalid field: agency_phone",
			"agency.txt:3: invalid field: agency_timezone",
			"agency.txt:3: invalid field: agency_url",
		},
	}
	testCases["unique-agencies"] = ggtfsTestCase{
		csvRows: [][]string{
			{"agency_id", "agency_name", "agency_url", "agency_timezone"},
			{"ACME", "ACME", "http://acme.inc", "Europe/Helsinki"},
			{"ACME", "ACME", "http://acme.inc", "Europe/Helsinki"},
		},
		expectedErrors: []string{
			"agency.txt:1: agency_id is not unique within the file",
		},
	}

	return testCases
}

func getAgencyOKTestcases() map[string]ggtfsTestCase {
	expected1 := Agency{
		Id:       NewID(stringPtr("1")),
		Name:     NewText(stringPtr("ACME")),
		URL:      NewURL(stringPtr("https://acme.inc")),
		Timezone: NewTimezone(stringPtr("Europe/Helsinki")),
		Lang:     NewOptionalLanguageCode(stringPtr("fi")),
		Phone:    NewOptionalPhoneNumber(stringPtr("+358123456")),
		FareURL:  NewOptionalURL(stringPtr("https://acme.inc/fares")),
		Email:    NewOptionalEmail(stringPtr("acme@acme.inc")),
	}

	expected2 := Agency{
		Id:       NewID(stringPtr("2")),
		Name:     NewText(stringPtr("FOO")),
		URL:      NewURL(stringPtr("https://foo.com")),
		Timezone: NewTimezone(stringPtr("Europe/Helsinki")),
		Lang:     nil,
		Phone:    nil,
		FareURL:  nil,
		Email:    nil,
	}

	testCases := make(map[string]ggtfsTestCase)
	testCases["1"] = ggtfsTestCase{
		csvRows: [][]string{
			{"agency_id", "agency_name", "agency_url", "agency_timezone", "agency_lang", "agency_phone", "agency_fare_url", "agency_email"},
			{"1", "ACME", "https://acme.inc", "Europe/Helsinki", "fi", "+358123456", "https://acme.inc/fares", "acme@acme.inc"},
		},
		expectedStructs: []interface{}{&expected1},
	}

	testCases["2"] = ggtfsTestCase{
		csvRows: [][]string{
			{"agency_id", "agency_name", "agency_url", "agency_timezone"},
			{"2", "FOO", "https://foo.com", "Europe/Helsinki"},
		},
		expectedStructs: []interface{}{&expected2},
	}

	return testCases
}
