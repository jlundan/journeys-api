//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"encoding/csv"
	"strings"
	"testing"
)

var validRouteHeaders = []string{"route_id", "agency_id", "route_short_name", "route_long_name", "route_desc",
	"route_type", "route_url", "route_color", "route_text_color", "route_sort_order", "continuous_pickup",
	"continuous_drop_off", "network_id"}

func TestRouteParsing(t *testing.T) {
	loadRoutesFunc := func(reader *csv.Reader) ([]interface{}, []error) {
		routes, errs := LoadEntitiesFromCSV[*Route](reader, validRouteHeaders, CreateRoute, RoutesFileName)
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
		AgencyId:          NewID(stringPtr(agency)),
		ShortName:         NewText(stringPtr(shortName)),
		LongName:          NewText(stringPtr(longName)),
		Desc:              NewText(stringPtr(desc)),
		Type:              NewRouteType(stringPtr(routeType)),
		URL:               NewURL(stringPtr(u)),
		Color:             NewColor(stringPtr(rColor)),
		TextColor:         NewColor(stringPtr(textColor)),
		SortOrder:         NewPositiveInteger(stringPtr(so)),
		ContinuousPickup:  NewContinuousPickupType(stringPtr(cpt)),
		ContinuousDropOff: NewContinuousDropOffType(stringPtr(cdt)),
		NetworkId:         NewID(stringPtr(networkId)),
		LineNumber:        2,
	}

	cpt2 := ""
	cdt2 := ""

	expected2 := Route{
		Id:                NewID(stringPtr(id)),
		AgencyId:          NewID(stringPtr(agency)),
		ShortName:         NewText(stringPtr(shortName)),
		LongName:          NewText(stringPtr(longName)),
		Desc:              NewText(stringPtr(desc)),
		Type:              NewRouteType(stringPtr(routeType)),
		URL:               NewURL(stringPtr(u)),
		Color:             NewColor(stringPtr(rColor)),
		TextColor:         NewColor(stringPtr(textColor)),
		SortOrder:         NewPositiveInteger(stringPtr(so)),
		ContinuousPickup:  NewContinuousPickupType(stringPtr(cpt2)),
		ContinuousDropOff: NewContinuousDropOffType(stringPtr(cdt2)),
		NetworkId:         NewID(stringPtr(networkId)),
		LineNumber:        2,
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
	testCases["invalid-fields-must-error-out"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "agency_id", "route_short_name", "route_long_name", "route_desc", "route_type", "route_url", "route_color", "route_text_color", "route_sort_order", "continuous_pickup", "continuous_drop_off", "network_id"},
			{"", "", "", "", "", "", "", "", "", "", "", "", ""},
			{"id3", "valid agency", "OK short n", "valid long name", "valid desc", "-1", "not an url", "not a color", "not a color", "not an integer", "999", "999", "network id"},
			{"id4", "valid agency", "OK short n", "valid long name", "valid desc", "8", "http://example.com", "FFFFFF", "FFFFFF", "-1", "-1", "-1", "valid network id"},
			{"id5", "valid agency", "OK short n", "valid long name", "valid desc", "9", "http://example.com", "FFFFFF", "FFFFFF", "1", "4", "4", "valid network id"},
			{"id6", "valid agency", "OK short n", "valid long name", "valid desc", "10", "http://example.com", "FFFFFF", "FFFFFF", "1", "1", "1", "valid network id"},
			{"id7", "valid agency", "OK short n", "valid long name", "valid desc", "13", "http://example.com", "FFFFFF", "FFFFFF", "1", "1", "1", "valid network id"},
		},
		expectedErrors: []string{
			// "routes.txt:2: invalid field: network_id", // validation not implemented yet
			"routes.txt:2: invalid mandatory field: route_id",
			"routes.txt:2: invalid mandatory field: route_type",
			"routes.txt:2: route_long_name must be specified when route_short_name is empty or not present",
			"routes.txt:2: route_short_name must be specified when route_long_name is empty or not present",
			"routes.txt:3: invalid field: continuous_drop_off",
			"routes.txt:3: invalid field: continuous_pickup",
			"routes.txt:3: invalid field: route_color",
			"routes.txt:3: invalid field: route_sort_order",
			"routes.txt:3: invalid field: route_text_color",
			"routes.txt:3: invalid field: route_url",
			"routes.txt:3: invalid mandatory field: route_type",
			"routes.txt:4: invalid field: continuous_drop_off",
			"routes.txt:4: invalid field: continuous_pickup",
			"routes.txt:4: invalid field: route_sort_order",
			"routes.txt:4: invalid mandatory field: route_type",
			"routes.txt:5: invalid field: continuous_drop_off",
			"routes.txt:5: invalid field: continuous_pickup",
			"routes.txt:5: invalid mandatory field: route_type",
			"routes.txt:6: invalid mandatory field: route_type",
			"routes.txt:7: invalid mandatory field: route_type",
		},
		expectedRecommendations: []string{
			"routes.txt:2: specifying agency_id is recommended", // when there isn't multiple agencies in agencies.txt (or agencies.txt is not present)
		},
	}

	testCases["short_name-length-should-be-less-than-12-chars"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "agency_id", "route_short_name", "route_long_name", "route_desc", "route_type", "route_url", "route_color", "route_text_color", "route_sort_order", "continuous_pickup", "continuous_drop_off", "network_id"},
			{"id", "agency", "a name longer than thirteen characters", "long name", "desc", "1", "http://example.com", "FFFFFF", "FFFFFF", "1", "1", "1", "1"},
		},
		expectedRecommendations: []string{
			"routes.txt:2: route_short_name should be less than 12 characters",
		},
	}

	// route_short_name is required if routes.route_long_name is empty (id)
	// route_long_name is required if routes.route_short_name is empty (id)
	// there should be no errors if either (id2 and id3) or both (id4) are present
	testCases["short_name-and-long_name-must-be-present-if-the-other-is-not"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "agency_id", "route_short_name", "route_long_name", "route_desc", "route_type", "route_url", "route_color", "route_text_color", "route_sort_order", "continuous_pickup", "continuous_drop_off", "network_id"},
			{"id", "agency", "", "", "desc", "1", "http://example.com", "FFFFFF", "FFFFFF", "1", "1", "1", "1"},
			{"id2", "agency", "short", "", "desc", "1", "http://example.com", "FFFFFF", "FFFFFF", "1", "1", "1", "1"},
			{"id3", "agency", "", "long", "desc", "1", "http://example.com", "FFFFFF", "FFFFFF", "1", "1", "1", "1"},
			{"id4", "agency", "short", "long", "desc", "1", "http://example.com", "FFFFFF", "FFFFFF", "1", "1", "1", "1"},
		},
		expectedErrors: []string{
			"routes.txt:2: route_long_name must be specified when route_short_name is empty or not present",
			"routes.txt:2: route_short_name must be specified when route_long_name is empty or not present",
		},
	}

	testCases["agency-id-must-match-valid-agency"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_type", "route_short_name", "agency_id"},
			{"2", "1", "2", "1"},
		},
		expectedErrors: []string{
			"routes.txt:2: referenced agency_id '1' not found in agency.txt",
		},
		fixtures: map[string][]interface{}{
			"agencies": {
				&Agency{Id: NewID(stringPtr("0"))},
			},
		},
	}

	testCases["agency-id-must-be-specified-when-multiple-agencies-in-agency.txt"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_type", "route_short_name", "agency_id"},
			{"2", "1", "2", ""},
		},
		expectedErrors: []string{
			"routes.txt:2: agency_id is required when there are multiple agencies in agencies.txt",
		},
		fixtures: map[string][]interface{}{
			"agencies": {
				&Agency{Id: NewID(stringPtr("0"))},
				&Agency{Id: NewID(stringPtr("1"))},
			},
		},
	}

	testCases["route-ids-must-be-unique"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_type", "route_short_name", "agency_id"},
			{"1", "1", "2", "1"},
			{"1", "1", "2", "1"},
		},
		expectedErrors: []string{
			"routes.txt:3: route_id '1' is not unique within the file",
		},
	}

	testCases["route-desc-should-not-match-short-or-long-name"] = ggtfsTestCase{
		csvRows: [][]string{
			{"route_id", "route_type", "agency_id", "route_short_name", "route_long_name", "route_desc"},
			{"1", "1", "1", "foo", "", "foo"},
			{"2", "1", "2", "", "foo", "foo"},
		},
		expectedRecommendations: []string{
			"routes.txt:2: route_desc should not be the same as route_short_name or route_long_name",
			"routes.txt:3: route_desc should not be the same as route_short_name or route_long_name",
		},
	}

	return testCases
}

func TestShouldReturnEmptyRouteArrayOnEmptyString(t *testing.T) {
	routes, errors := LoadEntitiesFromCSV[*Route](csv.NewReader(strings.NewReader("")), validRouteHeaders, CreateRoute, RoutesFileName)
	if len(errors) > 0 {
		t.Error(errors)
	}
	if len(routes) != 0 {
		t.Error("expected zero route items")
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
		AgencyId:   NewID(stringPtr("foo")),
		LineNumber: 0,
	}}, []*Agency{nil})

	ValidateRoutes([]*Route{{
		LineNumber: 0,
	}}, []*Agency{nil})
}

func TestStructsShouldNotBePresentOnNilInput(t *testing.T) {
	cpt := NewContinuousPickupType(nil)
	if cpt.IsPresent() {
		t.Error("expected not present field")
	}

	cdo := NewContinuousDropOffType(nil)
	if cdo.IsPresent() {
		t.Error("expected not present field")
	}
	rt := NewRouteType(nil)
	if rt.IsPresent() {
		t.Error("expected not present field")
	}
}
