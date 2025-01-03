//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"fmt"
	"testing"
)

func TestCreateRoute(t *testing.T) {
	headerMap := map[string]int{"route_id": 0, "agency_id": 1, "route_short_name": 2, "route_long_name": 3,
		"route_desc": 4, "route_type": 5, "route_url": 6, "route_color": 7, "route_text_color": 8, "route_sort_order": 9,
		"continuous_pickup": 10, "continuous_drop_off": 11, "network_id": 12}

	tests := map[string]struct {
		headers    map[string]int
		rows       [][]string
		lineNumber int
		expected   []*Route
	}{
		"empty-row": {
			headers: headerMap,
			rows:    [][]string{{"", "", "", "", "", "", "", "", "", "", "", "", ""}},
			expected: []*Route{{
				Id:                stringPtr(""),
				AgencyId:          stringPtr(""),
				ShortName:         stringPtr(""),
				LongName:          stringPtr(""),
				Desc:              stringPtr(""),
				Type:              stringPtr(""),
				URL:               stringPtr(""),
				Color:             stringPtr(""),
				TextColor:         stringPtr(""),
				SortOrder:         stringPtr(""),
				ContinuousPickup:  stringPtr(""),
				ContinuousDropOff: stringPtr(""),
				NetworkId:         stringPtr(""),
				LineNumber:        0,
			}},
		},
		"nil-values": {
			headers: headerMap,
			rows:    [][]string{nil},
			expected: []*Route{{
				Id:                nil,
				AgencyId:          nil,
				ShortName:         nil,
				LongName:          nil,
				Desc:              nil,
				Type:              nil,
				URL:               nil,
				Color:             nil,
				TextColor:         nil,
				SortOrder:         nil,
				ContinuousPickup:  nil,
				ContinuousDropOff: nil,
				NetworkId:         nil,
				LineNumber:        0,
			}},
		},
		"OK": {
			headers: headerMap,
			rows: [][]string{
				{"1", "Agency", "route1", "route 1", "route description", "3", "https://acme.inc/1", "FFFFFF", "FFF000", "1", "2", "3", "network"},
			},
			expected: []*Route{{
				Id:                stringPtr("1"),
				AgencyId:          stringPtr("Agency"),
				ShortName:         stringPtr("route1"),
				LongName:          stringPtr("route 1"),
				Desc:              stringPtr("route description"),
				Type:              stringPtr("3"),
				URL:               stringPtr("https://acme.inc/1"),
				Color:             stringPtr("FFFFFF"),
				TextColor:         stringPtr("FFF000"),
				SortOrder:         stringPtr("1"),
				ContinuousPickup:  stringPtr("2"),
				ContinuousDropOff: stringPtr("3"),
				NetworkId:         stringPtr("network"),
				LineNumber:        0,
			}},
		},
	}

	for name, tt := range tests {
		t.Run(fmt.Sprintf("%s", name), func(t *testing.T) {
			var actual []*Route
			for i, row := range tt.rows {
				actual = append(actual, CreateRoute(row, tt.headers, i))
			}
			handleEntityCreateResults(t, tt.expected, actual)
		})
	}
}

func TestValidateRoutes(t *testing.T) {
	tests := map[string]struct {
		actualEntities  []*Route
		expectedResults []Result
		agencies        []*Agency
	}{
		"nil-slice": {
			actualEntities:  nil,
			expectedResults: []Result{},
		},
		"nil-slice-items": {
			actualEntities:  []*Route{nil},
			expectedResults: []Result{},
		},
		"invalid-fields": {
			actualEntities: []*Route{
				{
					Id:                stringPtr("1"),
					AgencyId:          stringPtr("Agency"),
					ShortName:         stringPtr("a way too long short name"),
					LongName:          stringPtr("route 1"),
					Desc:              stringPtr("route description"),
					Type:              stringPtr("8"),
					URL:               stringPtr("Not an URL"),
					Color:             stringPtr("not a color"),
					TextColor:         stringPtr("not a color"),
					SortOrder:         stringPtr("not an integer"),
					ContinuousPickup:  stringPtr("4"),
					ContinuousDropOff: stringPtr("4"),
					NetworkId:         stringPtr("network"),
				},
			},
			expectedResults: []Result{
				InvalidURLResult{SingleLineResult{FileName: "routes.txt", FieldName: "route_url", Line: 0}},
				InvalidColorResult{SingleLineResult{FileName: "routes.txt", FieldName: "route_color", Line: 0}},
				InvalidColorResult{SingleLineResult{FileName: "routes.txt", FieldName: "route_text_color", Line: 0}},
				InvalidIntegerResult{SingleLineResult{FileName: "routes.txt", FieldName: "route_sort_order", Line: 0}},
				TooLongRouteShortNameResult{SingleLineResult{FileName: "routes.txt", FieldName: "route_short_name", Line: 0}},
			},
		},
		"empty-short-and-long-name": {
			actualEntities: []*Route{
				{
					Id:                stringPtr("1"),
					AgencyId:          stringPtr("Agency"),
					ShortName:         stringPtr(""),
					LongName:          stringPtr(""),
					Desc:              stringPtr("route description"),
					Type:              stringPtr("3"),
					URL:               stringPtr("https://acme.inc/1"),
					Color:             stringPtr("FFFFFF"),
					TextColor:         stringPtr("FFF000"),
					SortOrder:         stringPtr("1"),
					ContinuousPickup:  stringPtr("2"),
					ContinuousDropOff: stringPtr("3"),
					NetworkId:         stringPtr("network"),
				},
			},
			expectedResults: []Result{
				MissingRouteShortNameWhenLongNameIsNotPresentResult{SingleLineResult{FileName: "routes.txt", FieldName: "route_short_name", Line: 0}},
				MissingRouteLongNameWhenShortNameIsNotPresentResult{SingleLineResult{FileName: "routes.txt", FieldName: "route_long_name", Line: 0}},
			},
		},
		"desc-duplicates-route-names": {
			actualEntities: []*Route{
				{
					Id:                stringPtr("1"),
					AgencyId:          stringPtr("Agency"),
					ShortName:         stringPtr("route1"),
					LongName:          stringPtr("route1"),
					Desc:              stringPtr("route1"),
					Type:              stringPtr("3"),
					URL:               stringPtr("https://acme.inc/1"),
					Color:             stringPtr("FFFFFF"),
					TextColor:         stringPtr("FFF000"),
					SortOrder:         stringPtr("1"),
					ContinuousPickup:  stringPtr("2"),
					ContinuousDropOff: stringPtr("3"),
					NetworkId:         stringPtr("network"),
				},
			},
			expectedResults: []Result{
				DescriptionDuplicatesRouteNameResult{SingleLineResult{FileName: "routes.txt", FieldName: "route_desc", Line: 0}, "route_short_name"},
				DescriptionDuplicatesRouteNameResult{SingleLineResult{FileName: "routes.txt", FieldName: "route_desc", Line: 0}, "route_long_name"},
			},
		},
		"agency-id-required-when-multiple-agencies": {
			actualEntities: []*Route{
				{
					Id:        stringPtr("1"),
					AgencyId:  stringPtr(""),
					ShortName: stringPtr("route1"),
					LongName:  stringPtr("route 1"),
					Type:      stringPtr("3"),
				},
			},
			agencies: []*Agency{
				{Id: stringPtr("111")},
				{Id: stringPtr("112")},
				{Id: stringPtr("")}, // Do not crash on empty agency ID
				{Id: nil},           // Do not crash on nil agency ID
			},
			expectedResults: []Result{
				AgencyIdRequiredForRouteWhenMultipleAgenciesResult{SingleLineResult{FileName: "routes.txt", FieldName: "agency_id", Line: 0}},
			},
		},
		"recommend-agency-id": {
			actualEntities: []*Route{
				{
					Id:        stringPtr("1"),
					AgencyId:  stringPtr(""),
					ShortName: stringPtr("route1"),
					LongName:  stringPtr("route 1"),
					Type:      stringPtr("3"),
				},
			},
			expectedResults: []Result{
				AgencyIdRecommendedForRouteResult{SingleLineResult{FileName: "routes.txt", FieldName: "agency_id", Line: 0}},
			},
		},
		"unique-route-id": {
			actualEntities: []*Route{
				{
					Id:        stringPtr("1"),
					AgencyId:  stringPtr("agency"),
					ShortName: stringPtr("route1"),
					LongName:  stringPtr("route 1"),
					Type:      stringPtr("3"),
				},
				{
					Id:        stringPtr("1"),
					AgencyId:  stringPtr("agency"),
					ShortName: stringPtr("route1"),
					LongName:  stringPtr("route 1"),
					Type:      stringPtr("3"),
				},
			},
			expectedResults: []Result{
				FieldIsNotUniqueResult{SingleLineResult{FileName: "routes.txt", FieldName: "route_id", Line: 0}},
			},
		},
		"foreign-key-failure": {
			actualEntities: []*Route{
				{
					Id:        stringPtr("1"),
					AgencyId:  stringPtr("113"),
					ShortName: stringPtr("route1"),
					LongName:  stringPtr("route 1"),
					Type:      stringPtr("3"),
				},
			},
			agencies: []*Agency{
				{Id: stringPtr("111")},
				{Id: stringPtr("112")},
				{Id: stringPtr("")}, // Do not crash on empty agency ID
				{Id: nil},           // Do not crash on nil agency ID
			},
			expectedResults: []Result{
				ForeignKeyViolationResult{
					ReferencingFileName:  "routes.txt",
					ReferencingFieldName: "agency_id",
					ReferencedFieldName:  "agency.txt",
					ReferencedFileName:   "agency_id",
					OffendingValue:       "113",
					ReferencedAtRow:      0,
				},
			},
		},
		"foreign-key-OK": {
			actualEntities: []*Route{
				{
					Id:        stringPtr("1"),
					AgencyId:  stringPtr("113"),
					ShortName: stringPtr("route1"),
					LongName:  stringPtr("route 1"),
					Type:      stringPtr("3"),
				},
			},
			agencies: []*Agency{
				{Id: stringPtr("113")},
				{Id: stringPtr("112")},
			},
			expectedResults: []Result{},
		},
	}

	for name, tt := range tests {
		t.Run(fmt.Sprintf("%s", name), func(t *testing.T) {
			handleValidationResults(t, ValidateRoutes(tt.actualEntities, tt.agencies), tt.expectedResults)
		})
	}
}
