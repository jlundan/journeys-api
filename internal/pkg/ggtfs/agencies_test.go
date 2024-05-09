//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"testing"
)

func TestAgencyCSVParsing(t *testing.T) {
	agencies, errors := LoadAgencies(csv.NewReader(strings.NewReader("")))
	if len(errors) > 0 {
		t.Error(errors)
	}
	if len(agencies) != 0 {
		t.Error("expected zero agencies")
	}

	reader := csv.NewReader(strings.NewReader("foo,bar\n1,2"))
	reader.Comma = ','
	reader.Comment = ','
	_, errors = LoadAgencies(reader)
	if len(errors) == 0 {
		t.Error("expected to throw error")
	}
}

func TestAgencyParsingOK(t *testing.T) {
	expected1 := Agency{
		Id:       "1",
		Name:     "ACME",
		Url:      "https://acme.inc",
		Timezone: "Europe/Helsinki",
		Lang:     "fi",
		Phone:    "+358123456",
		FareURL:  "https://acme.inc/fares",
		Email:    "acme@acme.inc",
	}

	expected2 := Agency{
		Id:       "2",
		Name:     "FOO",
		Url:      "https://foo.com",
		Timezone: "Europe/Helsinki",
		Lang:     "",
		Phone:    "",
		FareURL:  "",
		Email:    "",
	}

	testCases := []struct {
		headers  map[string]uint8
		rows     [][]string
		expected Agency
	}{
		{
			rows: [][]string{
				{"agency_id", "agency_name", "agency_url", "agency_timezone", "agency_lang", "agency_phone", "agency_fare_url", "agency_email"},
				{"1", "ACME", "https://acme.inc", "Europe/Helsinki", "fi", "+358123456", "https://acme.inc/fares", "acme@acme.inc"},
			},
			expected: expected1,
		},
		{
			rows: [][]string{
				{"agency_id", "agency_name", "agency_url", "agency_timezone"},
				{"2", "FOO", "https://foo.com", "Europe/Helsinki"},
			},
			expected: expected2,
		},
	}

	for _, tc := range testCases {
		agencies, err := LoadAgencies(csv.NewReader(strings.NewReader(tableToString(tc.rows))))
		if err != nil && len(err) > 0 {
			t.Error(err)
			continue
		}

		if len(agencies) != 1 {
			t.Error("expected one row")
			continue
		}

		if !agenciesMatch(tc.expected, *agencies[0]) {
			a1, err := json.Marshal(tc.expected)
			if err != nil {
				t.Error(err)
			}
			a2, err := json.Marshal(*agencies[0])
			if err != nil {
				t.Error(err)
			}
			t.Error(fmt.Sprintf("expected %v, got %v", string(a1), string(a2)))
		}
	}
}

func TestAgencyParsingNOK(t *testing.T) {
	testCases := []struct {
		rows     [][]string
		expected []string
	}{
		{
			rows: [][]string{
				{"agency_name"},
				{","},
			},
			expected: []string{
				"agency.txt: record on line 2: wrong number of fields",
			},
		},
		{
			rows: [][]string{
				{"agency_name"},
				{" "},
			},
			expected: []string{
				//"agency.txt:0: agency_name must be specified",
				//"agency.txt:0: agency_name: empty value not allowed",
				"agency.txt:0: agency_timezone must be specified",
				"agency.txt:0: agency_url must be specified",
			},
		},
		{
			rows: [][]string{
				{"agency_name"},
				{"1"},
			},
			expected: []string{
				"agency.txt:0: agency_timezone must be specified",
				"agency.txt:0: agency_url must be specified",
			},
		},
		{
			rows: [][]string{
				{"agency_name", "agency_url"},
				{"ACME", ""},
			},
			expected: []string{
				"agency.txt:0: agency_timezone must be specified",
				"agency.txt:0: agency_url must be specified",
			},
		},
		{
			rows: [][]string{
				{"agency_name", "agency_url", "agency_timezone"},
				{"ACME", "http://acme.inc", ""},
			},
			expected: []string{
				"agency.txt:0: agency_timezone must be specified",
				//"agency.txt:0: agency_timezone: empty value not allowed",
			},
		},
		//{
		//	rows: [][]string{
		//		{"agency_name", "agency_url", "agency_timezone", "agency_lang", "agency_phone", "agency_fare_url", "agency_email"},
		//		{"ACME", "http://acme.inc", "Europe/Helsinki", "", "", "", ""},
		//	},
		//	expected: []string{
		//		"agency.txt:0: agency_email: empty value not allowed",
		//		"agency.txt:0: agency_lang: empty value not allowed",
		//		"agency.txt:0: agency_phone: empty value not allowed",
		//	},
		//},
		{
			rows: [][]string{
				{"agency_name", "agency_url", "agency_timezone"},
				{"ACME", "http://acme.inc", "Europe/Helsinki"},
				{"ACME", "http://acme.inc", "Europe/Helsinki"},
			},
			expected: []string{
				"agency.txt:0: agency id must be specified when multiple agencies are declared",
				"agency.txt:1: agency id must be specified when multiple agencies are declared",
			},
		},
		//{
		//	rows: [][]string{
		//		{"agency_name", "agency_url", "agency_timezone", "agency_id"},
		//		{"ACME", "http://acme.inc", "Europe/Helsinki", ""},
		//	},
		//	expected: []string{
		//		"agency.txt:0: agency_id: empty value not allowed",
		//	},
		//},
		{
			rows: [][]string{
				{"agency_name", "agency_url", "agency_timezone", "agency_id"},
				{"ACME", "http://acme.inc", "Europe/Helsinki", "ACME"},
				{"ACME2", "http://acme2.inc", "Europe/Helsinki2", "ACME"},
			},
			expected: []string{
				"agency.txt:2: non-unique id: agency_id",
			},
		},
	}

	for tcIndex, tc := range testCases {
		_, err := LoadAgencies(csv.NewReader(strings.NewReader(tableToString(tc.rows))))

		sort.Slice(err, func(x, y int) bool {
			return err[x].Error() < err[y].Error()
		})

		sort.Slice(tc.expected, func(x, y int) bool {
			return tc.expected[x] < tc.expected[y]
		})

		if len(err) == 0 {
			t.Error(fmt.Sprintf("%v: expected to throw an error", tcIndex))
			continue
		}

		if len(err) != len(tc.expected) {
			t.Error(fmt.Sprintf("%v: expected %v errors, got %v", tcIndex, len(tc.expected), len(err)))
			for _, e := range err {
				fmt.Println(e)
			}
			continue
		}

		for i, e := range err {
			if e.Error() != tc.expected[i] {
				t.Error(fmt.Sprintf("%v: expected error %s, got %s", tcIndex, tc.expected[i], e.Error()))
			}
		}
	}
}

func agenciesMatch(a Agency, b Agency) bool {
	// Name, Url and Timezone are mandatory according to GTFS spec -> test should fail if those are nil
	return a.Name == b.Name &&
		a.Url == b.Url &&
		a.Timezone == b.Timezone &&
		((a.Id == "" && b.Id == "") || a.Id == b.Id) &&
		((a.Email == "" && b.Email == "") || a.Email == b.Email) &&
		((a.FareURL == "" && b.FareURL == "") || a.FareURL == b.FareURL) &&
		((a.Lang == "" && b.Lang == "") || a.Lang == b.Lang) &&
		((a.Phone == "" && b.Phone == "") || a.Phone == b.Phone)
}
