//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"encoding/csv"
	"strings"
	"testing"
)

func TestShouldReturnEmptyRouteArrayOnEmptyString(t *testing.T) {
	agencies, errors := LoadRoutes(csv.NewReader(strings.NewReader("")))
	if len(errors) > 0 {
		t.Error(errors)
	}
	if len(agencies) != 0 {
		t.Error("expected zero calendar items")
	}
}

func TestRouteParsing(t *testing.T) {
	loadRoutesFunc := func(reader *csv.Reader) ([]interface{}, []error) {
		routes, errs := LoadRoutes(reader)
		entities := make([]interface{}, len(routes))
		for i, route := range routes {
			entities[i] = route
		}
		return entities, errs
	}

	validateRoutesFunc := func(entities []interface{}) []error {
		routes := make([]*Route, len(entities))
		for i, entity := range entities {
			if route, ok := entity.(*Route); ok {
				routes[i] = route
			}
		}

		return ValidateRoutes(routes, nil)
	}

	runGenericGTFSParseTest(t, "RouteNOKTestcases", loadRoutesFunc, validateRoutesFunc, false, getRouteNOKTestcases())
	runGenericGTFSParseTest(t, "RouteOKTestcases", loadRoutesFunc, validateRoutesFunc, false, getRouteOKTestcases())
}

func getRouteOKTestcases() map[string]ggtfsTestCase {
	id := "1"
	agency := "ACME"
	shortName := "1"
	longName := "route1"
	desc := "ACME route 1"
	routeType := "3"
	u := "https://acme.inc/1"
	rColor := "FFFFFF"
	textColor := "000000"
	so := "1"
	cpt := "2"
	cdt := "3"

	expected1 := Route{
		Id:                id,
		AgencyId:          agency,
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

	//cpt2 := "1"
	//cdt2 := "1"

	expected2 := Route{
		Id:                id,
		AgencyId:          agency,
		ShortName:         &shortName,
		LongName:          &longName,
		Desc:              &desc,
		Type:              routeType,
		Url:               &u,
		Color:             &rColor,
		TextColor:         &textColor,
		SortOrder:         &so,
		ContinuousPickup:  nil,
		ContinuousDropOff: nil,
	}
	testCases := make(map[string]ggtfsTestCase)
	testCases["1"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "agency_id", "route_short_name", "route_long_name", "route_desc", "route_type", "route_url", "route_color", "route_text_color", "route_sort_order", "continuous_pickup", "continuous_drop_off"},
			{"1", "ACME", "1", "route1", "ACME route 1", "3", "https://acme.inc/1", "FFFFFF", "000000", "1", "2", "3"},
		},
		expectedStructs: []interface{}{&expected1},
	}

	testCases["2"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "agency_id", "route_short_name", "route_long_name", "route_desc", "route_type", "route_url", "route_color", "route_text_color", "route_sort_order", "continuous_pickup", "continuous_drop_off"},
			{"1", "ACME", "1", "route1", "ACME route 1", "3", "https://acme.inc/1", "FFFFFF", "000000", "1", "", ""},
		},
		expectedStructs: []interface{}{&expected2},
	}

	return testCases
}

func getRouteNOKTestcases() map[string]ggtfsTestCase {
	testCases := make(map[string]ggtfsTestCase)
	testCases["1"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id"},
			{","},
			{" "},
			{"1"},
		},
		expectedErrors: []string{
			"routes.txt: record on line 2: wrong number of fields",
			"routes.txt:1: either route_short_name or route_long_name must be specified",
			"routes.txt:1: route_id must be specified",
			"routes.txt:1: route_type must be specified",
			"routes.txt:2: either route_short_name or route_long_name must be specified",
			"routes.txt:2: route_type must be specified",
		},
	}
	testCases["2"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "agency_id"},
			{"1", ""},
		},
		expectedErrors: []string{
			"routes.txt:0: either route_short_name or route_long_name must be specified",
			"routes.txt:0: route_type must be specified",
		},
	}
	testCases["3"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_short_name"},
			{"1", ""},
		},
		expectedErrors: []string{
			//"routes.txt:0: route_short_name: empty value not allowed",
			"routes.txt:0: either route_short_name or route_long_name must be specified",
			"routes.txt:0: route_type must be specified",
		},
	}
	testCases["4"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_long_name"},
			{"1", ""},
		},
		expectedErrors: []string{
			"routes.txt:0: either route_short_name or route_long_name must be specified",
			"routes.txt:0: route_type must be specified",
		},
	}
	testCases["5"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_desc"},
			{"1", ""},
		},
		expectedErrors: []string{
			//"routes.txt:0: route_desc: empty value not allowed",
			"routes.txt:0: either route_short_name or route_long_name must be specified",
			"routes.txt:0: route_type must be specified",
		},
	}
	testCases["6"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_type"},
			{"1", "100"},
		},
		expectedErrors: []string{
			//"routes.txt:0: route_type: invalid value",
			"routes.txt:0: either route_short_name or route_long_name must be specified",
			//"routes.txt:0: route_type must be specified",
		},
	}
	testCases["7"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_type"},
			{"1", "malformed"},
		},
		expectedErrors: []string{
			//"routes.txt:0: route_type: strconv.ParseInt: parsing \"malformed\": invalid syntax",
			"routes.txt:0: either route_short_name or route_long_name must be specified",
		},
	}
	testCases["8"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_type", "route_short_name", "route_url"},
			{"1", "1", "1", "\000malformed"},
		},
		expectedErrors: []string{
			//"routes.txt:0: route_url: parse \"\\x00malformed\": net/url: invalid control character in URL",
		},
	}
	testCases["9"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_type", "route_short_name", "route_color"},
			{"1", "1", "1", "malformed"},
			{"2", "1", "1", ""},
		},
		expectedErrors: []string{
			//"routes.txt:0: route_color: encoding/hex: invalid byte: U+006D 'm'",
			//"routes.txt:1: route_color: empty value not allowed",
		},
	}
	testCases["10"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_type", "route_short_name", "route_text_color"},
			{"1", "1", "1", "malformed"},
		},
		expectedErrors: []string{
			//"routes.txt:0: route_text_color: encoding/hex: invalid byte: U+006D 'm'",
		},
	}
	testCases["11"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_type", "route_short_name", "route_sort_order"},
			{"1", "1", "1", "malformed"},
		},
		expectedErrors: []string{
			//"routes.txt:0: route_sort_order: strconv.ParseInt: parsing \"malformed\": invalid syntax",
		},
	}
	testCases["12"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_type", "route_short_name", "continuous_pickup"},
			{"1", "1", "1", "100"},
			{"2", "1", "1", "malformed"},
		},
		expectedErrors: []string{
			//"routes.txt:0: continuous_pickup: invalid value",
			//"routes.txt:1: continuous_pickup: strconv.ParseInt: parsing \"malformed\": invalid syntax",
		},
	}
	testCases["13"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_type", "route_short_name", "continuous_drop_off"},
			{"1", "1", "1", "100"},
			{"2", "1", "1", "malformed"},
			{"3", "1", "1", "1"},
			{"3", "2", "2", "1"},
		},
		expectedErrors: []string{
			//"routes.txt:0: continuous_drop_off: invalid value",
			//"routes.txt:1: continuous_drop_off: strconv.ParseInt: parsing \"malformed\": invalid syntax",
			//"routes.txt:3: non-unique id: route_id",
		},
	}

	return testCases
}

//TODO: Integrate these to the test cases run by runGenericGTFSParseTest
//func TestValidateRoutes(t *testing.T) {
//	agencyId1 := "1000"
//	agencyId2 := "1001"
//	agencyId3 := "1002"
//
//	testCases := []struct {
//		routes         []*Route
//		agencies       []*Agency
//		expectedErrors []string
//	}{
//		{
//			routes: []*Route{
//				{AgencyId: agencyId1, LineNumber: 0},
//			},
//			agencies: []*Agency{
//				{Id: agencyId1, LineNumber: 0},
//				{Id: agencyId2, LineNumber: 1},
//			},
//			expectedErrors: []string{},
//		},
//		{
//			routes:         nil,
//			expectedErrors: []string{},
//		},
//		{
//			routes: []*Route{
//				nil,
//			},
//			agencies: []*Agency{
//				{Id: agencyId2, LineNumber: 0},
//				{Id: agencyId3, LineNumber: 1},
//			},
//			expectedErrors: []string{},
//		},
//		{
//			routes: []*Route{
//				{AgencyId: agencyId1, LineNumber: 0},
//			},
//			agencies: []*Agency{
//				{Id: agencyId2, LineNumber: 0},
//				{Id: agencyId3, LineNumber: 1},
//			},
//			expectedErrors: []string{
//				"routes.txt:0: referenced agency_id not found in agency.txt",
//			},
//		},
//		{
//			routes: []*Route{
//				{AgencyId: agencyId1, LineNumber: 0},
//			},
//			agencies: []*Agency{
//				nil,
//			},
//			expectedErrors: []string{
//				"routes.txt:0: referenced agency_id not found in agency.txt",
//			},
//		},
//	}
//
//	for _, tc := range testCases {
//		err := ValidateRoutes(tc.routes, tc.agencies)
//		checkErrors(tc.expectedErrors, err, t)
//	}
//}
