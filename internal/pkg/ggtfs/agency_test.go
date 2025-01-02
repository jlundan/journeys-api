//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"fmt"
	"testing"
)

func TestValidateAgency(t *testing.T) {
	tests := map[string]struct {
		agencies        []*Agency
		expectedResults []Result
	}{
		"nil-agency-slice": {
			agencies:        nil,
			expectedResults: []Result{},
		},
		"nil-agency-slice-items": {
			agencies:        []*Agency{nil},
			expectedResults: []Result{},
		},
		"recommend-agency-id": {
			agencies: []*Agency{
				{Name: stringPtr("acme"), URL: stringPtr("http://acme.inc"), Timezone: stringPtr("Europe/Helsinki")},
			},
			expectedResults: []Result{
				SingleAgencyRecommendedResult{FileName: "agency.txt"},
			},
		},
		"unique-agencies": {
			agencies: []*Agency{
				{Id: stringPtr("1"), Name: stringPtr("acme"), URL: stringPtr("http://acme.inc"), Timezone: stringPtr("Europe/Helsinki")},
				{Id: stringPtr("1"), Name: stringPtr("acme"), URL: stringPtr("http://acme.inc"), Timezone: stringPtr("Europe/Helsinki")},
			},
			expectedResults: []Result{
				FieldIsNotUniqueResult{
					FileName:  "agency.txt",
					FieldName: "agency_id",
				},
			},
		},
		"valid-id-when-multiple-agencies": {
			agencies: []*Agency{
				{Name: stringPtr("acme"), URL: stringPtr("http://acme.inc"), Timezone: stringPtr("Europe/Helsinki"), LineNumber: 0},
				{Name: stringPtr("acme"), URL: stringPtr("http://acme.inc"), Timezone: stringPtr("Europe/Helsinki"), LineNumber: 1},
			},
			expectedResults: []Result{
				ValidAgencyIdRequiredWhenMultipleAgenciesResult{FileName: "agency.txt", Line: 0},
				ValidAgencyIdRequiredWhenMultipleAgenciesResult{FileName: "agency.txt", Line: 1},
			},
		},
	}

	for name, tt := range tests {
		t.Run(fmt.Sprintf("%s", name), func(t *testing.T) {
			handleValidationResults(t, ValidateAgencies(tt.agencies), tt.expectedResults)
		})
	}
}

//var validAgencyHeaders = []string{"agency_id", "agency_name", "agency_url", "agency_timezone",
//	"agency_lang", "agency_phone", "agency_fare_url", "agency_email"}
//
//func TestShouldReturnEmptyAgencyArrayOnEmptyString(t *testing.T) {
//	agencies, errors := LoadEntitiesFromCSV[*Agency](csv.NewReader(strings.NewReader("")), validAgencyHeaders, CreateAgency, AgenciesFileName)
//	if len(errors) > 0 {
//		t.Error(errors)
//	}
//	if len(agencies) != 0 {
//		t.Error("expected zero agencies")
//	}
//}
//
//func TestShouldValidateOnNilAgencies(t *testing.T) {
//	errors, recommendations := ValidateAgencies([]*Agency{nil, nil, nil})
//	if len(errors) > 0 {
//		t.Error(errors)
//	}
//	if len(recommendations) != 0 {
//		t.Error("expected zero recommendations")
//	}
//}
//
//func TestAgencyParsing(t *testing.T) {
//	loadAgenciesFunc := func(reader *csv.Reader) ([]interface{}, []error) {
//		agencies, errs := LoadEntitiesFromCSV[*Agency](reader, validAgencyHeaders, CreateAgency, AgenciesFileName)
//		entities := make([]interface{}, len(agencies))
//		for i, agency := range agencies {
//			entities[i] = agency
//		}
//		return entities, errs
//	}
//
//	validateAgenciesFunc := func(entities []interface{}, _fixtures map[string][]interface{}) ([]error, []string) {
//		agencies := make([]*Agency, len(entities))
//		for i, entity := range entities {
//			if agency, ok := entity.(*Agency); ok {
//				agencies[i] = agency
//			}
//		}
//		return ValidateAgencies(agencies)
//	}
//
//	runGenericGTFSParseTest(t, "AgencyNOKTestcases", loadAgenciesFunc, validateAgenciesFunc, false, getAgencyNOKTestcases())
//	runGenericGTFSParseTest(t, "AgencyOKTestcases", loadAgenciesFunc, validateAgenciesFunc, false, getAgencyOKTestcases())
//}
//
//func getAgencyNOKTestcases() map[string]ggtfsTestCase {
//	testCases := make(map[string]ggtfsTestCase)
//
//	testCases["invalid-fields"] = ggtfsTestCase{
//		csvRows: [][]string{
//			{"agency_name", "agency_url", "agency_timezone", "agency_lang", "agency_phone", "agency_fare_url", "agency_email"},
//			{",,,,,,"},
//			{"", "", "", "", "", "", ""},
//			{" ", " ", " ", " ", " ", " ", " "},
//			{"ACME", "Not an URL", "Europe/Helsinki2", "Not a language code", "Not a phone", "Not an URL", "Not a phone number"},
//		},
//		expectedErrors: []string{
//			"agency.txt:2: a valid agency_id must be specified when multiple agencies are declared",
//			"agency.txt:2: invalid mandatory field: agency_name",
//			"agency.txt:2: invalid mandatory field: agency_timezone",
//			"agency.txt:2: invalid mandatory field: agency_url",
//			"agency.txt:3: a valid agency_id must be specified when multiple agencies are declared",
//			"agency.txt:3: invalid mandatory field: agency_name",
//			"agency.txt:3: invalid mandatory field: agency_timezone",
//			"agency.txt:3: invalid mandatory field: agency_url",
//			"agency.txt:4: a valid agency_id must be specified when multiple agencies are declared",
//			"agency.txt:4: invalid mandatory field: agency_name",
//			"agency.txt:4: invalid mandatory field: agency_timezone",
//			"agency.txt:4: invalid mandatory field: agency_url",
//			"agency.txt:5: a valid agency_id must be specified when multiple agencies are declared",
//			"agency.txt:5: invalid field: agency_email",
//			"agency.txt:5: invalid field: agency_fare_url",
//			"agency.txt:5: invalid field: agency_lang",
//			"agency.txt:5: invalid field: agency_phone",
//			"agency.txt:5: invalid mandatory field: agency_timezone",
//			"agency.txt:5: invalid mandatory field: agency_url",
//		},
//	}
//	testCases["unique-agencies"] = ggtfsTestCase{
//		csvRows: [][]string{
//			{"agency_id", "agency_name", "agency_url", "agency_timezone"},
//			{"ACME", "ACME", "http://acme.inc", "Europe/Helsinki"},
//			{"ACME", "ACME", "http://acme.inc", "Europe/Helsinki"},
//		},
//		expectedErrors: []string{
//			"agency.txt:3: agency_id is not unique within the file",
//		},
//	}
//	testCases["recommend-agency-id"] = ggtfsTestCase{
//		csvRows: [][]string{
//			{"agency_id", "agency_name", "agency_url", "agency_timezone"},
//			{"", "ACME", "http://acme2.inc", "Europe/Helsinki"},
//		},
//		expectedRecommendations: []string{
//			"agency.txt:2: it is recommended that agency_id is specified even when there is only one agency",
//		},
//	}
//
//	return testCases
//}
//
//func getAgencyOKTestcases() map[string]ggtfsTestCase {
//	expected1 := Agency{
//		Id:         stringPtr("1"),
//		Name:       stringPtr("ACME"),
//		URL:        stringPtr("https://acme.inc"),
//		Timezone:   stringPtr("Europe/Helsinki"),
//		Lang:       stringPtr("fi"),
//		Phone:      stringPtr("+358123456"),
//		FareURL:    stringPtr("https://acme.inc/fares"),
//		Email:      stringPtr("acme@acme.inc"),
//		LineNumber: 2,
//	}
//
//	expected2 := Agency{
//		Id:         stringPtr("2"),
//		Name:       stringPtr("FOO"),
//		URL:        stringPtr("https://foo.com"),
//		Timezone:   stringPtr("Europe/Helsinki"),
//		LineNumber: 2,
//	}
//
//	testCases := make(map[string]ggtfsTestCase)
//	testCases["1"] = ggtfsTestCase{
//		csvRows: [][]string{
//			{"agency_id", "agency_name", "agency_url", "agency_timezone", "agency_lang", "agency_phone", "agency_fare_url", "agency_email"},
//			{"1", "ACME", "https://acme.inc", "Europe/Helsinki", "fi", "+358123456", "https://acme.inc/fares", "acme@acme.inc"},
//		},
//		expectedStructs: []interface{}{&expected1},
//	}
//
//	testCases["2"] = ggtfsTestCase{
//		csvRows: [][]string{
//			{"agency_id", "agency_name", "agency_url", "agency_timezone"},
//			{"2", "FOO", "https://foo.com", "Europe/Helsinki"},
//		},
//		expectedStructs: []interface{}{&expected2},
//	}
//
//	return testCases
//}
