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

	validateAgenciesFunc := func(entities []interface{}) []error {
		agencies := make([]*Agency, len(entities))
		for i, entity := range entities {
			if agency, ok := entity.(*Agency); ok {
				agencies[i] = agency
			}
		}
		return ValidateAgencies(agencies)
	}

	runGenericGTFSParseTest(t, "NOK", loadAgenciesFunc, validateAgenciesFunc, false, getAgencyNOKTestcases())
	runGenericGTFSParseTest(t, "OK", loadAgenciesFunc, validateAgenciesFunc, false, getAgencyOKTestcases())
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
			{"agency_name"},
			{" "},
		},
		expectedErrors: []string{
			"agency.txt:0: agency_timezone must be specified",
			"agency.txt:0: agency_url must be specified",
		},
	}

	testCases["3"] = ggtfsTestCase{
		csvRows: [][]string{
			{"agency_name"},
			{"1"},
		},
		expectedErrors: []string{
			"agency.txt:0: agency_timezone must be specified",
			"agency.txt:0: agency_url must be specified",
		},
	}

	testCases["4"] = ggtfsTestCase{
		csvRows: [][]string{
			{"agency_name", "agency_url"},
			{"ACME", ""},
		},
		expectedErrors: []string{
			"agency.txt:0: agency_timezone must be specified",
			"agency.txt:0: agency_url must be specified",
		},
	}

	testCases["5"] = ggtfsTestCase{
		csvRows: [][]string{
			{"agency_name", "agency_url", "agency_timezone"},
			{"ACME", "http://acme.inc", ""},
		},
		expectedErrors: []string{
			"agency.txt:0: agency_timezone must be specified",
		},
	}

	testCases["6"] = ggtfsTestCase{
		csvRows: [][]string{
			{"agency_name", "agency_url", "agency_timezone"},
			{"ACME", "http://acme.inc", "Europe/Helsinki"},
			{"ACME", "http://acme.inc", "Europe/Helsinki"},
		},
		expectedErrors: []string{
			"agency.txt:0: agency_id must be specified when multiple agencies are declared",
			"agency.txt:1: agency_id must be specified when multiple agencies are declared",
		},
	}

	testCases["7"] = ggtfsTestCase{
		csvRows: [][]string{
			{"agency_name", "agency_url", "agency_timezone", "agency_id"},
			{"ACME", "http://acme.inc", "Europe/Helsinki", "ACME"},
			{"ACME2", "http://acme2.inc", "Europe/Helsinki2", "ACME"},
		},
		expectedErrors: []string{
			"agency.txt:1: agency_id is not unique within the file",
		},
	}

	return testCases
}

func getAgencyOKTestcases() map[string]ggtfsTestCase {
	expected1 := Agency{
		Id:       "1",
		Name:     "ACME",
		Url:      "https://acme.inc",
		Timezone: "Europe/Helsinki",
		Lang:     stringPtr("fi"),
		Phone:    stringPtr("+358123456"),
		FareURL:  stringPtr("https://acme.inc/fares"),
		Email:    stringPtr("acme@acme.inc"),
	}

	expected2 := Agency{
		Id:       "2",
		Name:     "FOO",
		Url:      "https://foo.com",
		Timezone: "Europe/Helsinki",
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
