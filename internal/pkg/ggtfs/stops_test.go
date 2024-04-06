package ggtfs

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"testing"
)

func TestStopCSVParsing(t *testing.T) {
	items, errors := LoadStops(csv.NewReader(strings.NewReader("")))
	if len(errors) > 0 {
		t.Error(errors)
	}
	if len(items) != 0 {
		t.Error("expected zero items")
	}

	reader := csv.NewReader(strings.NewReader("foo,bar\n1,2"))
	reader.Comma = ','
	reader.Comment = ','
	_, errors = LoadStops(reader)
	if len(errors) == 0 {
		t.Error("expected to throw error")
	}
}

func TestStopParsingOK(t *testing.T) {
	code := "0001"
	name := "Place 0001"
	desc := "Stop at place 0001"
	lat := 11.1111111
	lon := -11.1111111
	zone := "Z1"
	stopUrl := "https://acme.inc/stops/0001"
	location := 0
	parentStation := "4"
	timeZone := "Europe/Helsinki"
	wheelchairBoarding := 0
	level := "1"
	platformCode := "0001"

	expected1 := Stop{
		Id:                 "0001",
		Code:               &code,
		Name:               &name,
		Desc:               &desc,
		Lat:                &lat,
		Lon:                &lon,
		ZoneId:             &zone,
		Url:                &stopUrl,
		LocationType:       &location,
		ParentStation:      &parentStation,
		Timezone:           &timeZone,
		WheelchairBoarding: &wheelchairBoarding,
		PlatformCode:       &platformCode,
		LevelId:            &level,
	}

	testCases := []struct {
		headers  map[string]uint8
		rows     [][]string
		expected Stop
	}{
		{
			rows: [][]string{
				{"stop_id", "stop_code", "stop_name", "stop_desc", "stop_lat", "stop_lon",
					"zone_id", "stop_url", "location_type", "parent_station", "stop_timezone", "wheelchair_boarding", "level_id", "platform_code"},
				{"0001", "0001", "Place 0001", "Stop at place 0001", "11.1111111", "-11.1111111", "Z1",
					"https://acme.inc/stops/0001", "0", "4", "Europe/Helsinki", "0", "1", "0001"},
			},
			expected: expected1,
		},
	}

	for _, tc := range testCases {
		stops, err := LoadStops(csv.NewReader(strings.NewReader(tableToString(tc.rows))))
		if err != nil && len(err) > 0 {
			t.Error(err)
			continue
		}

		if len(stops) != 1 {
			t.Error("expected one row")
			continue
		}

		if !stopsMatch(tc.expected, *stops[0]) {
			s1, err := json.Marshal(tc.expected)
			if err != nil {
				t.Error(err)
			}
			s2, err := json.Marshal(*stops[0])
			if err != nil {
				t.Error(err)
			}
			t.Error(fmt.Sprintf("expected %v, got %v", string(s1), string(s2)))
		}
	}
}

func TestStopParsingNOK(t *testing.T) {
	testCases := []struct {
		rows     [][]string
		expected []string
	}{
		{
			rows: [][]string{
				{"stop_id", "stop_name", "stop_lat", "stop_lon", "parent_station"},
				{","},
				{"", "foo", "11.11", "22.22", "0002"},
				{"0001", "foo", "11.11", "22.22", "0002"},
				{"0001", "foo", "invalid", "22.22", "0002"},
			},
			expected: []string{
				"stops.txt: record on line 2: wrong number of fields",
				"stops.txt:1: stop_id: empty value not allowed",
				"stops.txt:1: stop_id must be specified",
				"stops.txt:3: non-unique id: stop_id",
				"stops.txt:3: stop_lat: strconv.ParseFloat: parsing \"invalid\": invalid syntax",
			},
		},
		{
			rows: [][]string{
				{"stop_id", "location_type"},
				{"0001", "2"},
			},
			expected: []string{
				"stops.txt:0: stop_name must be specified for location types 0,1 and 2",
				"stops.txt:0: stop_lat must be specified for location types 0,1 and 2",
				"stops.txt:0: stop_lon must be specified for location types 0,1 and 2",
				"stops.txt:0: parent_station must be specified for location types 2,3 and 4",
			},
		},
		{
			rows: [][]string{
				{"stop_id", "stop_name", "stop_lat", "stop_lon", "parent_station", "location_type", "wheelchair_boarding"},
				{"0000", "foo", "11.11", "22.22", "0002", "10", "4"},
				{"0001", "foo", "11.11", "22.22", "0002", "4", "10"},
				{"0002", "foo", "11.11", "22.22", "0002", "invalid", "4"},
				{"0003", "foo", "11.11", "22.22", "0002", "4", "invalid"},
			},
			expected: []string{
				"stops.txt:0: location_type: invalid value",
				"stops.txt:1: wheelchair_boarding: invalid value",
				"stops.txt:2: location_type: strconv.ParseInt: parsing \"invalid\": invalid syntax",
				"stops.txt:3: wheelchair_boarding: strconv.ParseInt: parsing \"invalid\": invalid syntax",
			},
		},
	}

	for _, tc := range testCases {
		_, err := LoadStops(csv.NewReader(strings.NewReader(tableToString(tc.rows))))

		sort.Slice(err, func(x, y int) bool {
			return err[x].Error() < err[y].Error()
		})

		sort.Slice(tc.expected, func(x, y int) bool {
			return tc.expected[x] < tc.expected[y]
		})

		if len(err) == 0 {
			t.Error("expected to throw an error")
			continue
		}

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

func stopsMatch(a Stop, b Stop) bool {
	return a.Id == b.Id && *a.Code == *b.Code && *a.Name == *b.Name && *a.Desc == *b.Desc && *a.Lat == *b.Lat &&
		*a.Lon == *b.Lon && *a.ZoneId == *b.ZoneId && *a.Url == *b.Url && *a.LocationType == *b.LocationType &&
		*a.ParentStation == *b.ParentStation && *a.Timezone == *b.Timezone && *a.WheelchairBoarding == *b.WheelchairBoarding &&
		*a.PlatformCode == *b.PlatformCode && *a.LevelId == *b.LevelId
}
