//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"encoding/csv"
	"strings"
	"testing"
)

func TestShouldReturnEmptyAgencyArrayOnEmptyString(t *testing.T) {
	agencies, errors := LoadAgencies(csv.NewReader(strings.NewReader("")))
	if len(errors) > 0 {
		t.Error(errors)
	}
	if len(agencies) != 0 {
		t.Error("expected zero agencies")
	}
}

func TestAgencyParsing(t *testing.T) {
	loadAgenciesFunc := func(reader *csv.Reader) ([]interface{}, []error) {
		agencies, errs := LoadAgencies(reader)
		entities := make([]interface{}, len(agencies))
		for i, agency := range agencies {
			entities[i] = agency
		}
		return entities, errs
	}

	validateAgenciesFunc := func(entities []interface{}, _ map[string][]interface{}) ([]error, []string) {
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
	testCases["1"] = ggtfsTestCase{
		csvRows: [][]string{
			{"agency_name"},
			{","},
		},
		expectedErrors: []string{
			"agency.txt: record on line 2: wrong number of fields",
		},
	}
	testCases["2"] = ggtfsTestCase{
		csvRows: [][]string{
			{"agency_id", "agency_name"},
			{"Foo", " "},
		},
		expectedErrors: []string{
			"agency.txt:0: invalid field: agency_name",
			"agency.txt:0: missing mandatory field: agency_timezone",
			"agency.txt:0: missing mandatory field: agency_url",
		},
	}

	testCases["3"] = ggtfsTestCase{
		csvRows: [][]string{
			{"agency_name", "agency_url", "agency_timezone"},
			{"", "", ""},
		},
		expectedErrors: []string{
			"agency.txt:0: invalid field: agency_name",
			"agency.txt:0: invalid field: agency_url",
			"agency.txt:0: invalid field: agency_timezone",
		},
	}

	testCases["4"] = ggtfsTestCase{
		csvRows: [][]string{
			{"agency_name", "agency_url", "agency_timezone", "agency_lang", "agency_phone", "agency_fare_url", "agency_email"},
			{"ACME", "http://acme.inc", "Europe/Helsinki", "", "", "", ""},
		},
		expectedErrors: []string{
			"agency.txt:0: invalid field: agency_lang",
			"agency.txt:0: invalid field: agency_phone",
			"agency.txt:0: invalid field: agency_fare_url",
			"agency.txt:0: invalid field: agency_email",
		},
	}

	testCases["5"] = ggtfsTestCase{
		csvRows: [][]string{
			{"agency_name", "agency_url", "agency_timezone"},
			{"ACME", "http://acme.inc", "Europe/Helsinki"},
			{"ACME", "http://acme.inc", "Europe/Helsinki"},
		},
		expectedErrors: []string{
			"agency.txt:0: a valid agency_id must be specified when multiple agencies are declared",
			"agency.txt:1: a valid agency_id must be specified when multiple agencies are declared",
		},
	}

	testCases["6"] = ggtfsTestCase{
		csvRows: [][]string{
			{"agency_name", "agency_url", "agency_timezone", "agency_id"},
			{"ACME", "http://acme.inc", "Europe/Helsinki", "ACME"},
			{"ACME2", "http://acme2.inc", "Europe/Helsinki2", "ACME"},
		},
		expectedErrors: []string{
			"agency.txt:1: agency_id is not unique within the file",
			"agency.txt:1: invalid field: agency_timezone",
		},
	}

	testCases["7"] = ggtfsTestCase{
		csvRows: [][]string{
			{"agency_url", "agency_timezone", "agency_id"},
			{"http://acme.inc", "Europe/Helsinki", "ACME"},
		},
		expectedErrors: []string{
			"agency.txt:0: missing mandatory field: agency_name",
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
