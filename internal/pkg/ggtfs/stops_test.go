//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"encoding/csv"
	"strings"
	"testing"
)

var validStopHeaders = []string{"stop_id", "stop_code", "stop_name", "stop_desc", "stop_lat", "stop_lon", "zone_id",
	"stop_url", "location_type", "parent_station", "stop_timezone", "wheelchair_boarding", "level_id", "platform_code", "municipality_id"}

func TestShouldReturnEmptyStopArrayOnEmptyString(t *testing.T) {
	stops, errors := LoadEntitiesFromCSV[*Stop](csv.NewReader(strings.NewReader("")), validStopHeaders, CreateStop, StopsFileName)
	if len(errors) > 0 {
		t.Error(errors)
	}
	if len(stops) != 0 {
		t.Error("expected zero calendar items")
	}
}

func TestStopParsing(t *testing.T) {
	loadStopsFunc := func(reader *csv.Reader) ([]interface{}, []error) {
		stops, errs := LoadEntitiesFromCSV[*Stop](reader, validStopHeaders, CreateStop, StopsFileName)
		entities := make([]interface{}, len(stops))
		for i, stop := range stops {
			entities[i] = stop
		}
		return entities, errs
	}

	validateStopsFunc := func(entities []interface{}, _ map[string][]interface{}) ([]error, []string) {
		stops := make([]*Stop, len(entities))
		for i, entity := range entities {
			if stop, ok := entity.(*Stop); ok {
				stops[i] = stop
			}
		}

		return ValidateStops(stops)
	}

	runGenericGTFSParseTest(t, "StopOKTestcases", loadStopsFunc, validateStopsFunc, false, getStopOKTestcases())
	runGenericGTFSParseTest(t, "StopNOKTestcases", loadStopsFunc, validateStopsFunc, false, getStopNOKTestcases())
}

func getStopOKTestcases() map[string]ggtfsTestCase {
	expected1 := Stop{
		Id:                 NewID(stringPtr("0001")),
		Code:               NewText(stringPtr("0001")),
		Name:               NewText(stringPtr("Place 0001")),
		Desc:               NewText(stringPtr("Stop at place 0001")),
		Lat:                NewLatitude(stringPtr("11.1111111")),
		Lon:                NewLongitude(stringPtr("-11.1111111")),
		ZoneId:             NewID(stringPtr("Z1")),
		Url:                NewURL(stringPtr("https://acme.inc/stops/0001")),
		LocationType:       NewStopLocation(stringPtr("0")),
		ParentStation:      NewID(stringPtr("4")),
		Timezone:           NewTimezone(stringPtr("Europe/Helsinki")),
		WheelchairBoarding: NewWheelchairBoarding(stringPtr("0")),
		PlatformCode:       NewText(stringPtr("0001")),
		LevelId:            NewID(stringPtr("1")),
		LineNumber:         2,
	}

	testCases := make(map[string]ggtfsTestCase)
	testCases["1"] = ggtfsTestCase{
		csvRows: [][]string{
			{"stop_id", "stop_code", "stop_name", "stop_desc", "stop_lat", "stop_lon",
				"zone_id", "stop_url", "location_type", "parent_station", "stop_timezone", "wheelchair_boarding", "level_id", "platform_code"},
			{"0001", "0001", "Place 0001", "Stop at place 0001", "11.1111111", "-11.1111111", "Z1",
				"https://acme.inc/stops/0001", "0", "4", "Europe/Helsinki", "0", "1", "0001"},
		},
		expectedStructs: []interface{}{&expected1},
	}

	return testCases
}

func getStopNOKTestcases() map[string]ggtfsTestCase {
	testCases := make(map[string]ggtfsTestCase)

	testCases["invalid-fields-must-error-out"] = ggtfsTestCase{
		csvRows: [][]string{
			{"stop_id", "stop_code", "stop_name", "stop_desc", "stop_lat", "stop_lon", "zone_id",
				"stop_url", "location_type", "parent_station", "stop_timezone", "wheelchair_boarding", "level_id", "platform_code", "municipality_id"},
			{"", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		},
		expectedErrors: []string{
			"stops.txt:2: invalid mandatory field: stop_id",
		},
	}

	return testCases
}
