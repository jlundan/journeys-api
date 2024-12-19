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
	stops, errors := LoadEntities[*Stop](csv.NewReader(strings.NewReader("")), validStopHeaders, CreateStop, StopsFileName)
	if len(errors) > 0 {
		t.Error(errors)
	}
	if len(stops) != 0 {
		t.Error("expected zero calendar items")
	}
}

func TestStopParsing(t *testing.T) {
	loadStopsFunc := func(reader *csv.Reader) ([]interface{}, []error) {
		stops, errs := LoadEntities[*Stop](reader, validStopHeaders, CreateStop, StopsFileName)
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
		Code:               NewOptionalText(stringPtr("0001")),
		Name:               NewOptionalText(stringPtr("Place 0001")),
		Desc:               NewOptionalText(stringPtr("Stop at place 0001")),
		Lat:                NewOptionalLatitude(stringPtr("11.1111111")),
		Lon:                NewOptionalLongitude(stringPtr("-11.1111111")),
		ZoneId:             NewOptionalID(stringPtr("Z1")),
		Url:                NewOptionalURL(stringPtr("https://acme.inc/stops/0001")),
		LocationType:       NewOptionalStopLocation(stringPtr("0")),
		ParentStation:      NewOptionalID(stringPtr("4")),
		Timezone:           NewOptionalTimezone(stringPtr("Europe/Helsinki")),
		WheelchairBoarding: NewOptionalWheelchairBoarding(stringPtr("0")),
		PlatformCode:       NewOptionalText(stringPtr("0001")),
		LevelId:            NewOptionalID(stringPtr("1")),
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
	testCases["1"] = ggtfsTestCase{
		csvRows: [][]string{
			{"stop_id", "stop_name", "stop_lat", "stop_lon", "parent_station"},
			{","},
			{"", "foo", "11.11", "22.22", "0002"},
			{"0001", "foo", "11.11", "22.22", "0002"},
			{"0001", "foo", "invalid", "22.22", "0002"},
		},
		expectedErrors: []string{
			"stops.txt: record on line 2: wrong number of fields",
			"stops.txt:1: stop_id must be specified",
			"stops.txt:3: stop_id '0001' is not unique within the file",
			//"stops.txt: record on line 2: wrong number of fields",
			//"stops.txt:1: stop_id: empty value not allowed",
			//"stops.txt:1: stop_id must be specified",
			//"stops.txt:3: non-unique id: stop_id",
			//"stops.txt:3: stop_lat: strconv.ParseFloat: parsing \"invalid\": invalid syntax",
		},
	}
	testCases["2"] = ggtfsTestCase{
		csvRows: [][]string{
			{"stop_id", "location_type"},
			{"0001", "2"},
		},
		expectedErrors: []string{
			"stops.txt:0: parent_station must be specified for location types 2, 3, and 4",
			"stops.txt:0: stop_lat must be specified for location types 0, 1, and 2",
			"stops.txt:0: stop_lon must be specified for location types 0, 1, and 2",
			"stops.txt:0: stop_name must be specified for location types 0, 1, and 2",
		},
	}
	testCases["3"] = ggtfsTestCase{
		csvRows: [][]string{
			{"stop_id", "stop_name", "stop_lat", "stop_lon", "parent_station", "location_type", "wheelchair_boarding"},
			{"0000", "foo", "11.11", "22.22", "0002", "10", "4"},
			{"0001", "foo", "11.11", "22.22", "0002", "4", "10"},
			{"0002", "foo", "11.11", "22.22", "0002", "invalid", "4"},
			{"0003", "foo", "11.11", "22.22", "0002", "4", "invalid"},
		},
		expectedErrors: []string{
			//"stops.txt:0: location_type: invalid value",
			//"stops.txt:1: wheelchair_boarding: invalid value",
			//"stops.txt:2: location_type: strconv.ParseInt: parsing \"invalid\": invalid syntax",
			//"stops.txt:3: wheelchair_boarding: strconv.ParseInt: parsing \"invalid\": invalid syntax",
		},
	}

	return testCases
}
