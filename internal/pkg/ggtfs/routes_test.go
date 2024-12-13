//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"encoding/csv"
	"strings"
	"testing"
)

func TestRouteParsing(t *testing.T) {
	loadRoutesFunc := func(reader *csv.Reader) ([]interface{}, []error) {
		routes, errs := LoadRoutes(reader)
		entities := make([]interface{}, len(routes))
		for i, route := range routes {
			entities[i] = route
		}
		return entities, errs
	}

	validateRoutesFunc := func(entities []interface{}, fixtures map[string][]interface{}) ([]error, []string) {
		routes := make([]*Route, len(entities))
		for i, entity := range entities {
			if route, ok := entity.(*Route); ok {
				routes[i] = route
			}
		}

		if agenciesFixture, ok := fixtures["agencies"]; !ok || len(agenciesFixture) == 0 {
			return ValidateRoutes(routes, nil)
		}

		agencies := make([]*Agency, len(fixtures["agencies"]))
		for i, entity := range fixtures["agencies"] {
			if agency, ok := entity.(*Agency); ok {
				agencies[i] = agency
			}
		}

		return ValidateRoutes(routes, agencies)
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
	networkId := "1"

	expected1 := Route{
		Id:                NewID(stringPtr(id)),
		AgencyId:          NewOptionalID(stringPtr(agency)),
		ShortName:         NewOptionalText(stringPtr(shortName)),
		LongName:          NewOptionalText(stringPtr(longName)),
		Description:       NewOptionalText(stringPtr(desc)),
		Type:              NewRouteType(stringPtr(routeType)),
		URL:               NewOptionalURL(stringPtr(u)),
		Color:             NewOptionalColor(stringPtr(rColor)),
		TextColor:         NewOptionalColor(stringPtr(textColor)),
		SortOrder:         NewOptionalInteger(stringPtr(so)),
		ContinuousPickup:  NewOptionalContinuousPickupType(stringPtr(cpt)),
		ContinuousDropOff: NewOptionalContinuousDropOffType(stringPtr(cdt)),
		NetworkId:         NewOptionalID(stringPtr(networkId)),
	}

	cpt2 := ""
	cdt2 := ""

	expected2 := Route{
		Id:                NewID(stringPtr(id)),
		AgencyId:          NewOptionalID(stringPtr(agency)),
		ShortName:         NewOptionalText(stringPtr(shortName)),
		LongName:          NewOptionalText(stringPtr(longName)),
		Description:       NewOptionalText(stringPtr(desc)),
		Type:              NewRouteType(stringPtr(routeType)),
		URL:               NewOptionalURL(stringPtr(u)),
		Color:             NewOptionalColor(stringPtr(rColor)),
		TextColor:         NewOptionalColor(stringPtr(textColor)),
		SortOrder:         NewOptionalInteger(stringPtr(so)),
		ContinuousPickup:  NewOptionalContinuousPickupType(stringPtr(cpt2)),
		ContinuousDropOff: NewOptionalContinuousDropOffType(stringPtr(cdt2)),
		NetworkId:         NewOptionalID(stringPtr(networkId)),
	}
	testCases := make(map[string]ggtfsTestCase)
	testCases["1"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "agency_id", "route_short_name", "route_long_name", "route_desc", "route_type", "route_url", "route_color", "route_text_color", "route_sort_order", "continuous_pickup", "continuous_drop_off", "network_id"},
			{"1", "ACME", "1", "route1", "ACME route 1", "3", "https://acme.inc/1", "FFFFFF", "000000", "1", "2", "3", "1"},
		},
		expectedStructs: []interface{}{&expected1},
	}

	testCases["2"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "agency_id", "route_short_name", "route_long_name", "route_desc", "route_type", "route_url", "route_color", "route_text_color", "route_sort_order", "continuous_pickup", "continuous_drop_off", "network_id"},
			{"1", "ACME", "1", "route1", "ACME route 1", "3", "https://acme.inc/1", "FFFFFF", "000000", "1", "", "", "1"},
		},
		expectedStructs: []interface{}{&expected2},
	}

	testCases["3"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_type", "route_short_name", "agency_id"},
			{"1", "1", "1", "0"},
		},
		expectedErrors: []string{},
		fixtures: map[string][]interface{}{
			"agencies": {
				&Agency{Id: NewID(stringPtr("0"))},
			},
		},
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
			"routes.txt:1: invalid field: route_id",
			"routes.txt:1: missing mandatory field: route_type",
			"routes.txt:2: either route_short_name or route_long_name must be specified",
			"routes.txt:2: missing mandatory field: route_type",
		},
	}
	testCases["2"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "agency_id"},
			{"1", ""},
		},
		expectedErrors: []string{
			"routes.txt:0: either route_short_name or route_long_name must be specified",
			"routes.txt:0: missing mandatory field: route_type",
		},
	}
	testCases["3"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_short_name"},
			{"1", ""},
		},
		expectedErrors: []string{
			"routes.txt:0: missing mandatory field: route_type",
		},
	}
	testCases["4"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_long_name"},
			{"1", ""},
		},
		expectedErrors: []string{
			"routes.txt:0: missing mandatory field: route_type",
		},
	}
	testCases["5"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_desc"},
			{"1", ""},
		},
		expectedErrors: []string{
			"routes.txt:0: either route_short_name or route_long_name must be specified",
			"routes.txt:0: invalid field: route_desc",
			"routes.txt:0: missing mandatory field: route_type",
		},
	}
	testCases["6"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_type"},
			{"1", "100"},
		},
		expectedErrors: []string{
			"routes.txt:0: either route_short_name or route_long_name must be specified",
			"routes.txt:0: invalid field: route_type",
		},
	}
	testCases["7"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_type"},
			{"1", "malformed"},
		},
		expectedErrors: []string{
			"routes.txt:0: either route_short_name or route_long_name must be specified",
			"routes.txt:0: invalid field: route_type",
		},
	}
	testCases["8"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_type", "route_short_name", "route_url"},
			{"1", "1", "1", "\000malformed"},
		},
		expectedErrors: []string{
			"routes.txt:0: invalid field: route_url",
		},
	}
	testCases["9"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_type", "route_short_name", "route_color"},
			{"1", "1", "1", "malformed"},
			{"2", "1", "1", ""},
		},
		expectedErrors: []string{
			"routes.txt:0: invalid field: route_color",
			"routes.txt:1: invalid field: route_color",
		},
	}
	testCases["10"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_type", "route_short_name", "route_text_color"},
			{"1", "1", "1", "malformed"},
		},
		expectedErrors: []string{
			"routes.txt:0: invalid field: route_text_color",
		},
	}
	testCases["11"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_type", "route_short_name", "route_sort_order"},
			{"1", "1", "1", "malformed"},
		},
		expectedErrors: []string{
			"routes.txt:0: invalid field: route_sort_order",
		},
	}
	testCases["12"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_type", "route_short_name", "continuous_pickup"},
			{"1", "1", "1", "100"},
			{"2", "1", "1", "malformed"},
		},
		expectedErrors: []string{
			"routes.txt:0: invalid field: continuous_pickup",
			"routes.txt:1: invalid field: continuous_pickup",
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
			"routes.txt:0: invalid field: continuous_drop_off",
			"routes.txt:1: invalid field: continuous_drop_off",
			"routes.txt:3: route_id '3' is not unique within the file",
		},
	}

	testCases["14"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_type", "route_short_name", "network_id"},
			{"1", "1", "1", ""},
		},
		expectedErrors: []string{
			"routes.txt:0: invalid field: network_id",
		},
	}

	testCases["15"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_type", "route_short_name", "agency_id"},
			{"2", "1", "2", "1"},
		},
		expectedErrors: []string{
			"routes.txt:0: referenced agency_id '1' not found in agency.txt",
		},
		fixtures: map[string][]interface{}{
			"agencies": {
				&Agency{Id: NewID(stringPtr("0"))},
			},
		},
	}

	return testCases
}

func TestShouldReturnEmptyRouteArrayOnEmptyString(t *testing.T) {
	agencies, errors := LoadRoutes(csv.NewReader(strings.NewReader("")))
	if len(errors) > 0 {
		t.Error(errors)
	}
	if len(agencies) != 0 {
		t.Error("expected zero calendar items")
	}
}

func TestShouldNotFailValidationOnNilRoutes(t *testing.T) {
	ValidateRoutes(nil, nil)
}

func TestShouldNotFailValidationOnNilRouteItem(t *testing.T) {
	ValidateRoutes([]*Route{nil}, nil)
}

func TestShouldNotFailValidationOnNilAgencyItem(t *testing.T) {
	ValidateRoutes([]*Route{{
		Id: ID{
			base: base{
				raw:       "",
				isPresent: false,
			},
		},
		AgencyId:    NewOptionalID(stringPtr("foo")),
		ShortName:   nil,
		LongName:    nil,
		Description: nil,
		Type: RouteType{
			Integer: Integer{
				base: base{
					raw:       "",
					isPresent: false,
				},
			},
		},
		URL:               nil,
		Color:             nil,
		TextColor:         nil,
		SortOrder:         nil,
		ContinuousPickup:  nil,
		ContinuousDropOff: nil,
		NetworkId:         nil,
		LineNumber:        0,
	}}, []*Agency{nil})

	ValidateRoutes([]*Route{{
		Id: ID{
			base: base{
				raw:       "",
				isPresent: false,
			},
		},
		AgencyId:    nil,
		ShortName:   nil,
		LongName:    nil,
		Description: nil,
		Type: RouteType{
			Integer: Integer{
				base: base{
					raw:       "",
					isPresent: false,
				},
			},
		},
		URL:               nil,
		Color:             nil,
		TextColor:         nil,
		SortOrder:         nil,
		ContinuousPickup:  nil,
		ContinuousDropOff: nil,
		NetworkId:         nil,
		LineNumber:        0,
	}}, []*Agency{nil})
}

func TestStructsShouldNotBePresentOnNilInput(t *testing.T) {
	cpt := NewOptionalContinuousPickupType(nil)
	if cpt.IsPresent() {
		t.Error("expected not present field")
	}

	cdo := NewOptionalContinuousDropOffType(nil)
	if cdo.IsPresent() {
		t.Error("expected not present field")
	}
	rt := NewRouteType(nil)
	if rt.IsPresent() {
		t.Error("expected not present field")
	}
}
