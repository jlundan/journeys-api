//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"fmt"
	"testing"
)

func TestCreateStop(t *testing.T) {
	headerMap := map[string]int{"stop_id": 0, "stop_code": 1, "stop_name": 2, "tts_stop_name": 3, "stop_desc": 4,
		"stop_lat": 5, "stop_lon": 6, "zone_id": 7, "stop_url": 8, "location_type": 9, "parent_station": 10,
		"stop_timezone": 11, "wheelchair_boarding": 12, "level_id": 13, "platform_code": 14}

	tests := map[string]struct {
		headers    map[string]int
		rows       [][]string
		lineNumber int
		expected   []*Stop
	}{
		"empty-row": {
			headers: headerMap,
			rows:    [][]string{{"", "", "", "", "", "", "", "", "", "", "", "", "", "", ""}},
			expected: []*Stop{{
				Id:                 stringPtr(""),
				Code:               stringPtr(""),
				Name:               stringPtr(""),
				TTSName:            stringPtr(""),
				Desc:               stringPtr(""),
				Lat:                stringPtr(""),
				Lon:                stringPtr(""),
				ZoneId:             stringPtr(""),
				URL:                stringPtr(""),
				LocationType:       stringPtr(""),
				ParentStation:      stringPtr(""),
				Timezone:           stringPtr(""),
				WheelchairBoarding: stringPtr(""),
				LevelId:            stringPtr(""),
				PlatformCode:       stringPtr(""),
				LineNumber:         0,
			}},
		},
		"nil-values": {
			headers: headerMap,
			rows:    [][]string{nil},
			expected: []*Stop{{
				Id:                 nil,
				Code:               nil,
				Name:               nil,
				TTSName:            nil,
				Desc:               nil,
				Lat:                nil,
				Lon:                nil,
				ZoneId:             nil,
				URL:                nil,
				LocationType:       nil,
				ParentStation:      nil,
				Timezone:           nil,
				WheelchairBoarding: nil,
				LevelId:            nil,
				PlatformCode:       nil,
				LineNumber:         0,
			}},
		},
		"OK": {
			headers: headerMap,
			rows: [][]string{
				{"0001", "0001", "Place 0001", "TTS Place 0001", "Stop at place 0001", "11.1111111", "-11.1111111", "Z1", "https://acme.inc/stops/0001", "0", "4", "Europe/Helsinki", "0", "level 1", "platform_code"},
			},
			expected: []*Stop{{
				Id:                 stringPtr("0001"),
				Code:               stringPtr("0001"),
				Name:               stringPtr("Place 0001"),
				TTSName:            stringPtr("TTS Place 0001"),
				Desc:               stringPtr("Stop at place 0001"),
				Lat:                stringPtr("11.1111111"),
				Lon:                stringPtr("-11.1111111"),
				ZoneId:             stringPtr("Z1"),
				URL:                stringPtr("https://acme.inc/stops/0001"),
				LocationType:       stringPtr("0"),
				ParentStation:      stringPtr("4"),
				Timezone:           stringPtr("Europe/Helsinki"),
				WheelchairBoarding: stringPtr("0"),
				LevelId:            stringPtr("level 1"),
				PlatformCode:       stringPtr("platform_code"),
				LineNumber:         0,
			}},
		},
	}

	for name, tt := range tests {
		t.Run(fmt.Sprintf("%s", name), func(t *testing.T) {
			var actual []*Stop
			for i, row := range tt.rows {
				actual = append(actual, CreateStop(row, tt.headers, i))
			}
			handleEntityCreateResults(t, tt.expected, actual)
		})
	}
}

func TestValidateStops(t *testing.T) {
	tests := map[string]struct {
		actualEntities  []*Stop
		expectedResults []ValidationNotice
	}{
		"nil-slice": {
			actualEntities:  nil,
			expectedResults: []ValidationNotice{},
		},
		"nil-slice-items": {
			actualEntities:  []*Stop{nil},
			expectedResults: []ValidationNotice{},
		},
		"nil-stop-id": {
			actualEntities: []*Stop{{Id: nil}},
			expectedResults: []ValidationNotice{
				MissingRequiredFieldNotice{SingleLineNotice{FileName: "stops.txt", FieldName: "stop_id", Line: 0}},
			},
		},
		"invalid-fields": {
			actualEntities: []*Stop{
				{
					Id:                 stringPtr("0001"),
					Lat:                stringPtr("Not a latitude"),
					Lon:                stringPtr("Not a longitude"),
					URL:                stringPtr("Not an URL"),
					LocationType:       stringPtr("5"),
					ParentStation:      stringPtr("0004"),
					Timezone:           stringPtr("Not a timezone"),
					WheelchairBoarding: stringPtr("3"),
				},
			},
			expectedResults: []ValidationNotice{
				InvalidURLNotice{SingleLineNotice{FileName: "stops.txt", FieldName: "stop_url"}},
				InvalidLocationTypeNotice{SingleLineNotice{FileName: "stops.txt", FieldName: "location_type"}},
				InvalidTimezoneNotice{SingleLineNotice{FileName: "stops.txt", FieldName: "stop_timezone"}},
				InvalidWheelchairBoardingValueNotice{SingleLineNotice{FileName: "stops.txt", FieldName: "wheelchair_boarding"}},
			},
		},
		"empty-fields-for-location-types": {
			actualEntities: []*Stop{
				{Id: stringPtr("0001"), Name: stringPtr(""), LocationType: stringPtr("0"), Lat: stringPtr("11.1111111"), Lon: stringPtr("11.1111111"), ParentStation: stringPtr("0004"), LineNumber: 0},
				{Id: stringPtr("0002"), Name: stringPtr(""), LocationType: stringPtr("1"), Lat: stringPtr("11.1111111"), Lon: stringPtr("11.1111111"), ParentStation: stringPtr("0004"), LineNumber: 1},
				{Id: stringPtr("0003"), Name: stringPtr(""), LocationType: stringPtr("2"), Lat: stringPtr("11.1111111"), Lon: stringPtr("11.1111111"), ParentStation: stringPtr("0004"), LineNumber: 2},

				{Id: stringPtr("0004"), Name: stringPtr("Place 0001"), Lat: stringPtr(""), LocationType: stringPtr("0"), Lon: stringPtr("11.1111111"), ParentStation: stringPtr("0004"), LineNumber: 3},
				{Id: stringPtr("0005"), Name: stringPtr("Place 0001"), Lat: stringPtr(""), LocationType: stringPtr("1"), Lon: stringPtr("11.1111111"), ParentStation: stringPtr("0004"), LineNumber: 4},
				{Id: stringPtr("0006"), Name: stringPtr("Place 0001"), Lat: stringPtr(""), LocationType: stringPtr("2"), Lon: stringPtr("11.1111111"), ParentStation: stringPtr("0004"), LineNumber: 5},

				{Id: stringPtr("0007"), Name: stringPtr("Place 0001"), Lon: stringPtr(""), LocationType: stringPtr("0"), Lat: stringPtr("11.1111111"), ParentStation: stringPtr("0004"), LineNumber: 6},
				{Id: stringPtr("0008"), Name: stringPtr("Place 0001"), Lon: stringPtr(""), LocationType: stringPtr("1"), Lat: stringPtr("11.1111111"), ParentStation: stringPtr("0004"), LineNumber: 7},
				{Id: stringPtr("0009"), Name: stringPtr("Place 0001"), Lon: stringPtr(""), LocationType: stringPtr("2"), Lat: stringPtr("11.1111111"), ParentStation: stringPtr("0004"), LineNumber: 8},

				{Id: stringPtr("0010"), Name: stringPtr("Place 0001"), ParentStation: stringPtr(""), LocationType: stringPtr("2"), Lat: stringPtr("11.1111111"), Lon: stringPtr("11.1111111"), LineNumber: 9},
				{Id: stringPtr("0011"), Name: stringPtr("Place 0001"), ParentStation: stringPtr(""), LocationType: stringPtr("3"), Lat: stringPtr("11.1111111"), Lon: stringPtr("11.1111111"), LineNumber: 10},
				{Id: stringPtr("0012"), Name: stringPtr("Place 0001"), ParentStation: stringPtr(""), LocationType: stringPtr("4"), Lat: stringPtr("11.1111111"), Lon: stringPtr("11.1111111"), LineNumber: 11},
			},
			expectedResults: []ValidationNotice{
				FieldRequiredForStopLocationTypeNotice{RequiredField: "stop_name", LocationType: "0", FileName: "stops.txt", Line: 0},
				FieldRequiredForStopLocationTypeNotice{RequiredField: "stop_name", LocationType: "1", FileName: "stops.txt", Line: 1},
				FieldRequiredForStopLocationTypeNotice{RequiredField: "stop_name", LocationType: "2", FileName: "stops.txt", Line: 2},

				FieldRequiredForStopLocationTypeNotice{RequiredField: "stop_lat", LocationType: "0", FileName: "stops.txt", Line: 3},
				FieldRequiredForStopLocationTypeNotice{RequiredField: "stop_lat", LocationType: "1", FileName: "stops.txt", Line: 4},
				FieldRequiredForStopLocationTypeNotice{RequiredField: "stop_lat", LocationType: "2", FileName: "stops.txt", Line: 5},

				FieldRequiredForStopLocationTypeNotice{RequiredField: "stop_lon", LocationType: "0", FileName: "stops.txt", Line: 6},
				FieldRequiredForStopLocationTypeNotice{RequiredField: "stop_lon", LocationType: "1", FileName: "stops.txt", Line: 7},
				FieldRequiredForStopLocationTypeNotice{RequiredField: "stop_lon", LocationType: "2", FileName: "stops.txt", Line: 8},

				FieldRequiredForStopLocationTypeNotice{RequiredField: "parent_station", LocationType: "2", FileName: "stops.txt", Line: 9},
				FieldRequiredForStopLocationTypeNotice{RequiredField: "parent_station", LocationType: "3", FileName: "stops.txt", Line: 10},
				FieldRequiredForStopLocationTypeNotice{RequiredField: "parent_station", LocationType: "4", FileName: "stops.txt", Line: 11},
			},
		},
		"unique-stop-id": {
			actualEntities: []*Stop{
				{Id: stringPtr("1"), Lat: stringPtr("11.1111111"), Lon: stringPtr("11.1111111")},
				{Id: stringPtr("1"), Lat: stringPtr("11.1111111"), Lon: stringPtr("11.1111111")},
			},
			expectedResults: []ValidationNotice{
				FieldIsNotUniqueNotice{SingleLineNotice{FileName: "stops.txt", FieldName: "stop_id", Line: 0}},
			},
		},
	}

	for name, tt := range tests {
		t.Run(fmt.Sprintf("%s", name), func(t *testing.T) {
			handleValidationResults(t, ValidateStops(tt.actualEntities), tt.expectedResults)
		})
	}
}
