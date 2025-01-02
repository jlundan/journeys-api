//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"fmt"
	"testing"
)

func TestCreateAgency(t *testing.T) {
	headerMap := map[string]int{"agency_id": 0, "agency_name": 1, "agency_url": 2, "agency_timezone": 3,
		"agency_lang": 4, "agency_phone": 5, "agency_fare_url": 6, "agency_email": 7,
	}

	tests := map[string]struct {
		headers          map[string]int
		rows             [][]string
		lineNumber       int
		expectedAgencies []*Agency
	}{
		"empty-row": {
			headers: headerMap,
			rows:    [][]string{{"", "", "", "", "", "", "", ""}},
			expectedAgencies: []*Agency{{
				Id:         stringPtr(""),
				Name:       stringPtr(""),
				URL:        stringPtr(""),
				Timezone:   stringPtr(""),
				Lang:       stringPtr(""),
				Phone:      stringPtr(""),
				FareURL:    stringPtr(""),
				Email:      stringPtr(""),
				LineNumber: 0,
			}},
		},
		"nil-values": {
			headers: headerMap,
			rows:    [][]string{nil},
			expectedAgencies: []*Agency{{
				Id:         nil,
				Name:       nil,
				URL:        nil,
				Timezone:   nil,
				Lang:       nil,
				Phone:      nil,
				FareURL:    nil,
				Email:      nil,
				LineNumber: 0,
			}},
		},
		"OK": {
			headers: headerMap,
			rows: [][]string{
				{"1", "ACME", "https://acme.inc", "Europe/Helsinki", "fi", "+358123456", "https://acme.inc/fares", "acme@acme.inc"},
				{"2", "ACME2", "https://acme2.inc", "Europe/Helsinki", "fi", "+3589876543", "https://acme2.inc/fares", "acme@acme2.inc"},
			},
			expectedAgencies: []*Agency{{
				Id:         stringPtr("1"),
				Name:       stringPtr("ACME"),
				URL:        stringPtr("https://acme.inc"),
				Timezone:   stringPtr("Europe/Helsinki"),
				Lang:       stringPtr("fi"),
				Phone:      stringPtr("+358123456"),
				FareURL:    stringPtr("https://acme.inc/fares"),
				Email:      stringPtr("acme@acme.inc"),
				LineNumber: 0,
			}, {
				Id:         stringPtr("2"),
				Name:       stringPtr("ACME2"),
				URL:        stringPtr("https://acme2.inc"),
				Timezone:   stringPtr("Europe/Helsinki"),
				Lang:       stringPtr("fi"),
				Phone:      stringPtr("+3589876543"),
				FareURL:    stringPtr("https://acme2.inc/fares"),
				Email:      stringPtr("acme@acme2.inc"),
				LineNumber: 1,
			}},
		},
	}

	for name, tt := range tests {
		t.Run(fmt.Sprintf("%s", name), func(t *testing.T) {
			var agencies []*Agency
			for i, row := range tt.rows {
				agencies = append(agencies, CreateAgency(row, tt.headers, i))
			}
			handleEntityCreateResults(t, tt.expectedAgencies, agencies)
		})
	}
}

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
		"invalid-fields": {
			agencies: []*Agency{
				{
					Id:       stringPtr("valid id to avoid recommend agency result"),
					Name:     stringPtr("valid name to avoid required name result"),
					URL:      stringPtr("Not an url"),
					Timezone: stringPtr("Not a city"),
					Lang:     stringPtr("Not a language"),
					Phone:    stringPtr("Not a phone number"),
					FareURL:  stringPtr("Not an url"),
					Email:    stringPtr("Not an email"),
				},
			},
			expectedResults: []Result{
				InvalidURLResult{FileName: "agency.txt", FieldName: "agency_url"},
				InvalidTimezoneResult{FileName: "agency.txt", FieldName: "agency_timezone"},
				InvalidLanguageCodeResult{FileName: "agency.txt", FieldName: "agency_lang"},
				InvalidPhoneNumberResult{FileName: "agency.txt", FieldName: "agency_phone"},
				InvalidURLResult{FileName: "agency.txt", FieldName: "agency_fare_url"},
				InvalidEmailResult{FileName: "agency.txt", FieldName: "agency_email"},
			},
		},
	}

	for name, tt := range tests {
		t.Run(fmt.Sprintf("%s", name), func(t *testing.T) {
			handleValidationResults(t, ValidateAgencies(tt.agencies), tt.expectedResults)
		})
	}
}
