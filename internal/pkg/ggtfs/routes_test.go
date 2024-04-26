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

func TestRouteCSVParsing(t *testing.T) {
	items, errors := LoadRoutes(csv.NewReader(strings.NewReader("")))
	if len(errors) > 0 {
		t.Error(errors)
	}
	if len(items) != 0 {
		t.Error("expected zero items")
	}

	reader := csv.NewReader(strings.NewReader("foo,bar\n1,2"))
	reader.Comma = ','
	reader.Comment = ','
	_, errors = LoadRoutes(reader)
	if len(errors) == 0 {
		t.Error("expected to throw error")
	}
}

func TestRouteParsingOK(t *testing.T) {
	id := "1"
	agency := "ACME"
	shortName := "1"
	longName := "route1"
	desc := "ACME route 1"
	routeType := 3
	u := "https://acme.inc/1"
	rColor := "FFFFFF"
	textColor := "000000"
	so := 1
	cpt := 2
	cdt := 3

	expected1 := Route{
		Id:                id,
		AgencyId:          &agency,
		ShortName:         &shortName,
		LongName:          &longName,
		Desc:              &desc,
		Type:              routeType,
		Url:               &u,
		Color:             &rColor,
		TextColor:         &textColor,
		SortOrder:         &so,
		ContinuousPickup:  &cpt,
		ContinuousDropOff: &cdt,
	}

	cpt2 := 1
	cdt2 := 1

	expected2 := Route{
		Id:                id,
		AgencyId:          &agency,
		ShortName:         &shortName,
		LongName:          &longName,
		Desc:              &desc,
		Type:              routeType,
		Url:               &u,
		Color:             &rColor,
		TextColor:         &textColor,
		SortOrder:         &so,
		ContinuousPickup:  &cpt2,
		ContinuousDropOff: &cdt2,
	}

	testCases := []struct {
		rows     [][]string
		expected Route
	}{
		{
			rows: [][]string{
				{"route_id", "agency_id", "route_short_name", "route_long_name", "route_desc", "route_type", "route_url", "route_color", "route_text_color", "route_sort_order", "continuous_pickup", "continuous_drop_off"},
				{"1", "ACME", "1", "route1", "ACME route 1", "3", "https://acme.inc/1", "FFFFFF", "000000", "1", "2", "3"},
			},
			expected: expected1,
		},
		{
			rows: [][]string{
				{"route_id", "agency_id", "route_short_name", "route_long_name", "route_desc", "route_type", "route_url", "route_color", "route_text_color", "route_sort_order", "continuous_pickup", "continuous_drop_off"},
				{"1", "ACME", "1", "route1", "ACME route 1", "3", "https://acme.inc/1", "FFFFFF", "000000", "1", "", ""},
			},
			expected: expected2,
		},
	}

	for _, tc := range testCases {
		routes, err := LoadRoutes(csv.NewReader(strings.NewReader(tableToString(tc.rows))))
		if err != nil && len(err) > 0 {
			t.Error(err)
			continue
		}

		if len(routes) != 1 {
			t.Error("expected one row")
			continue
		}

		if !routesMatch(tc.expected, *routes[0]) {
			r1, err := json.Marshal(tc.expected)
			if err != nil {
				t.Error(err)
			}
			r2, err := json.Marshal(*routes[0])
			if err != nil {
				t.Error(err)
			}
			t.Error(fmt.Sprintf("expected %v, got %v", string(r1), string(r2)))
		}
	}
}

func TestRouteParsingNOK(t *testing.T) {
	testCases := []struct {
		rows     [][]string
		expected []string
	}{
		{
			rows: [][]string{
				{"route_id"},
				{","},
				{" "},
				{"1"},
			},
			expected: []string{
				"routes.txt: record on line 2: wrong number of fields",
				"routes.txt:1: route_id: empty value not allowed",
				"routes.txt:1: route_id must be specified",
				"routes.txt:1: either route_short_name or route_long_name must be specified",
				"routes.txt:2: either route_short_name or route_long_name must be specified",
				"routes.txt:2: route_type must be specified",
				"routes.txt:1: route_type must be specified",
			},
		},
		{
			rows: [][]string{
				{"route_id", "agency_id"},
				{"1", ""},
			},
			expected: []string{
				"routes.txt:0: agency_id: empty value not allowed",
				"routes.txt:0: either route_short_name or route_long_name must be specified",
				"routes.txt:0: route_type must be specified",
			},
		},
		{
			rows: [][]string{
				{"route_id", "route_short_name"},
				{"1", ""},
			},
			expected: []string{
				"routes.txt:0: route_short_name: empty value not allowed",
				"routes.txt:0: either route_short_name or route_long_name must be specified",
				"routes.txt:0: route_type must be specified",
			},
		},
		{
			rows: [][]string{
				{"route_id", "route_long_name"},
				{"1", ""},
			},
			expected: []string{
				//"routes.txt:0: route_long_name: empty value not allowed",
				//"routes.txt:0: either route_short_name or route_long_name must be specified",
				"routes.txt:0: route_type must be specified",
			},
		},
		{
			rows: [][]string{
				{"route_id", "route_desc"},
				{"1", ""},
			},
			expected: []string{
				"routes.txt:0: route_desc: empty value not allowed",
				"routes.txt:0: either route_short_name or route_long_name must be specified",
				"routes.txt:0: route_type must be specified",
			},
		},
		{
			rows: [][]string{
				{"route_id", "route_type"},
				{"1", "100"},
			},
			expected: []string{
				"routes.txt:0: route_type: invalid value",
				"routes.txt:0: either route_short_name or route_long_name must be specified",
				"routes.txt:0: route_type must be specified",
			},
		},
		{
			rows: [][]string{
				{"route_id", "route_type"},
				{"1", "malformed"},
			},
			expected: []string{
				"routes.txt:0: route_type: strconv.ParseInt: parsing \"malformed\": invalid syntax",
				"routes.txt:0: either route_short_name or route_long_name must be specified",
			},
		},
		{
			rows: [][]string{
				{"route_id", "route_type", "route_short_name", "route_url"},
				{"1", "1", "1", "\000malformed"},
			},
			expected: []string{"routes.txt:0: route_url: parse \"\\x00malformed\": net/url: invalid control character in URL"},
		},
		{
			rows: [][]string{
				{"route_id", "route_type", "route_short_name", "route_color"},
				{"1", "1", "1", "malformed"},
				{"2", "1", "1", ""},
			},
			expected: []string{
				"routes.txt:0: route_color: encoding/hex: invalid byte: U+006D 'm'",
				"routes.txt:1: route_color: empty value not allowed",
			},
		},
		{
			rows: [][]string{
				{"route_id", "route_type", "route_short_name", "route_text_color"},
				{"1", "1", "1", "malformed"},
			},
			expected: []string{"routes.txt:0: route_text_color: encoding/hex: invalid byte: U+006D 'm'"},
		},
		{
			rows: [][]string{
				{"route_id", "route_type", "route_short_name", "route_sort_order"},
				{"1", "1", "1", "malformed"},
			},
			expected: []string{"routes.txt:0: route_sort_order: strconv.ParseInt: parsing \"malformed\": invalid syntax"},
		},
		{
			rows: [][]string{
				{"route_id", "route_type", "route_short_name", "continuous_pickup"},
				{"1", "1", "1", "100"},
				{"2", "1", "1", "malformed"},
			},
			expected: []string{
				"routes.txt:0: continuous_pickup: invalid value",
				"routes.txt:1: continuous_pickup: strconv.ParseInt: parsing \"malformed\": invalid syntax",
			},
		},
		{
			rows: [][]string{
				{"route_id", "route_type", "route_short_name", "continuous_drop_off"},
				{"1", "1", "1", "100"},
				{"2", "1", "1", "malformed"},
				{"3", "1", "1", "1"},
				{"3", "2", "2", "1"},
			},
			expected: []string{
				"routes.txt:0: continuous_drop_off: invalid value",
				"routes.txt:1: continuous_drop_off: strconv.ParseInt: parsing \"malformed\": invalid syntax",
				"routes.txt:3: non-unique id: route_id",
			},
		},
	}

	for _, tc := range testCases {
		_, err := LoadRoutes(csv.NewReader(strings.NewReader(tableToString(tc.rows))))

		if len(err) == 0 {
			t.Error("expected to throw an error")
			continue
		}

		sort.Slice(err, func(x, y int) bool {
			return err[x].Error() < err[y].Error()
		})

		sort.Slice(tc.expected, func(x, y int) bool {
			return tc.expected[x] < tc.expected[y]
		})

		if len(err) != len(tc.expected) {
			t.Error(fmt.Sprintf("expected %v errors, got %v", len(tc.expected), len(err)))
			for _, e := range err {
				fmt.Println(e)
			}
			continue
		}

		for i, e := range err {
			if e.Error() != tc.expected[i] {
				t.Error(fmt.Sprintf("expected error %s, got %s", tc.expected[i], e.Error()))
			}
		}
	}
}

func TestValidateRoutes(t *testing.T) {
	agencyId1 := "1000"
	agencyId2 := "1001"
	agencyId3 := "1002"

	testCases := []struct {
		routes         []*Route
		agencies       []*Agency
		expectedErrors []string
	}{
		{
			routes: []*Route{
				{AgencyId: &agencyId1, lineNumber: 0},
			},
			agencies: []*Agency{
				{Id: &agencyId1, lineNumber: 0},
				{Id: &agencyId2, lineNumber: 1},
			},
			expectedErrors: []string{},
		},
		{
			routes:         nil,
			expectedErrors: []string{},
		},
		{
			routes: []*Route{
				nil,
			},
			agencies: []*Agency{
				{Id: &agencyId2, lineNumber: 0},
				{Id: &agencyId3, lineNumber: 1},
			},
			expectedErrors: []string{},
		},
		{
			routes: []*Route{
				{AgencyId: &agencyId1, lineNumber: 0},
			},
			agencies: []*Agency{
				{Id: &agencyId2, lineNumber: 0},
				{Id: &agencyId3, lineNumber: 1},
			},
			expectedErrors: []string{
				"routes.txt:0: referenced agency_id not found in agency.txt",
			},
		},
		{
			routes: []*Route{
				{AgencyId: &agencyId1, lineNumber: 0},
			},
			agencies: []*Agency{
				nil,
			},
			expectedErrors: []string{
				"routes.txt:0: referenced agency_id not found in agency.txt",
			},
		},
	}

	for _, tc := range testCases {
		err := ValidateRoutes(tc.routes, tc.agencies)
		checkErrors(tc.expectedErrors, err, t)
	}
}

func routesMatch(r1 Route, r2 Route) bool {
	// GTFS spec says that Id and type are the only absolutely mandatory field, therefore it is not a pointer (and cannot be nil)
	return r1.Id == r2.Id &&
		((r1.AgencyId == nil && r2.AgencyId == nil) || *r1.AgencyId == *r2.AgencyId) &&
		((r1.ShortName != nil && r2.ShortName != nil) || (r1.LongName != nil && r2.LongName != nil)) && // either short_name or long_name must not be nil
		((r1.ShortName == nil && r2.ShortName == nil) || *r1.ShortName == *r2.ShortName) &&
		((r1.LongName == nil && r2.LongName == nil) || *r1.LongName == *r2.LongName) &&
		((r1.Desc == nil && r2.Desc == nil) || *r1.Desc == *r2.Desc) &&
		r1.Type == r2.Type &&
		((r1.Url == nil && r2.Url == nil) || *r1.Url == *r2.Url) &&
		((r1.Color == nil && r2.Color == nil) || *r1.Color == *r2.Color) &&
		((r1.TextColor == nil && r2.TextColor == nil) || *r1.TextColor == *r2.TextColor) &&
		((r1.SortOrder == nil && r2.SortOrder == nil) || *r1.SortOrder == *r2.SortOrder) &&
		((r1.ContinuousPickup == nil && r2.ContinuousPickup == nil) || *r1.ContinuousPickup == *r2.ContinuousPickup) &&
		((r1.ContinuousDropOff == nil && r2.ContinuousDropOff == nil) || *r1.ContinuousDropOff == *r2.ContinuousDropOff)
}
