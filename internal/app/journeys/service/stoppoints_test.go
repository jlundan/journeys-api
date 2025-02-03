//go:build journeys_stops_tests || journeys_tests || all_tests

package service

import (
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/app/journeys/repository"
	"github.com/jlundan/journeys-api/internal/testutil"
	"testing"
)

func TestStopPointMatchesConditions(t *testing.T) {
	testCases := []struct {
		item       *model.StopPoint
		conditions map[string]string
		expected   bool
	}{
		{nil, nil, false},
		{&model.StopPoint{Municipality: nil}, map[string]string{"municipalityName": "foo"}, false},
		{&model.StopPoint{Municipality: nil}, map[string]string{"municipalityShortName": "foo"}, false},
	}

	for _, tc := range testCases {
		matches := stopPointMatchesConditions(tc.item, tc.conditions)
		if matches != tc.expected {
			t.Error(fmt.Sprintf("expected %v, got %v", matches, tc.expected))
		}
	}
}

func TestStopPointsService_Search(t *testing.T) {
	dataStore := &repository.JourneysRepository{
		StopPoints: &repository.JourneysStopPointsRepository{
			All: []*model.StopPoint{
				{Name: "StopPoint1", ShortName: "SP1", TariffZone: "Zone1", Municipality: &model.Municipality{Name: "Municipality1", PublicCode: "M1"}},
				{Name: "StopPoint2", ShortName: "SP2", TariffZone: "Zone2", Municipality: &model.Municipality{Name: "Municipality2", PublicCode: "M2"}},
			},
		},
	}
	service := StopPointsService{Repository: dataStore}

	testCases := []struct {
		id       string
		params   map[string]string
		expected int
	}{
		{"1", map[string]string{"name": "StopPoint1"}, 1},
		{"2", map[string]string{"shortName": "SP1"}, 1},
		{"3", map[string]string{"tariffZone": "Zone1"}, 1},
		{"4", map[string]string{"municipalityName": "Municipality1"}, 1},
		{"5", map[string]string{"municipalityShortName": "M1"}, 1},
		{"6", map[string]string{"name": "NonExistent"}, 0},
	}

	for _, tc := range testCases {
		result := service.Search(tc.params)
		testutil.CompareVariablesAndPrintResults(t, tc.expected, len(result), tc.id)
	}
}

func TestStopPointsService_GetOneById(t *testing.T) {
	dataStore := &repository.JourneysRepository{
		StopPoints: &repository.JourneysStopPointsRepository{
			ById: map[string]*model.StopPoint{
				"1": {Name: "StopPoint1", ShortName: "SP1", TariffZone: "Zone1", Municipality: &model.Municipality{Name: "Municipality1", PublicCode: "M1"}},
				"2": {Name: "StopPoint2", ShortName: "SP2", TariffZone: "Zone2", Municipality: &model.Municipality{Name: "Municipality2", PublicCode: "M2"}},
			},
		},
	}
	service := StopPointsService{Repository: dataStore}

	testCases := []struct {
		id          string
		stopPointId string
		expected    *model.StopPoint
		err         error
	}{
		{"1", "1", &model.StopPoint{Name: "StopPoint1", ShortName: "SP1", TariffZone: "Zone1", Municipality: &model.Municipality{Name: "Municipality1", PublicCode: "M1"}}, nil},
		{"2", "2", &model.StopPoint{Name: "StopPoint2", ShortName: "SP2", TariffZone: "Zone2", Municipality: &model.Municipality{Name: "Municipality2", PublicCode: "M2"}}, nil},
		{"3", "3", nil, model.ErrNoSuchElement},
	}

	for _, tc := range testCases {
		result, err := service.GetOneById(tc.stopPointId)
		if err != nil && tc.err == nil {
			t.Error(err)
		} else if err != nil && tc.err != nil {
			testutil.CompareVariablesAndPrintResults(t, tc.err, err, tc.id)
		} else {
			testutil.CompareVariablesAndPrintResults(t, tc.expected, result, tc.id)
		}
	}
}

func TestStopPointMatchesConditions2(t *testing.T) {
	testCases := []struct {
		id         string
		item       *model.StopPoint
		conditions map[string]string
		expected   bool
	}{
		{"1", nil, nil, false},
		{"2", &model.StopPoint{Name: "StopPoint1"}, map[string]string{"name": "StopPoint1"}, true},
		{"3", &model.StopPoint{Name: "StopPoint1"}, map[string]string{"name": "StopPoint2"}, false},
		{"4", &model.StopPoint{ShortName: "SP1"}, map[string]string{"shortName": "SP1"}, true},
		{"5", &model.StopPoint{ShortName: "SP2"}, map[string]string{"shortName": "SP1"}, false},
		{"6", &model.StopPoint{TariffZone: "Zone1"}, map[string]string{"tariffZone": "Zone1"}, true},
		{"7", &model.StopPoint{TariffZone: "Zone2"}, map[string]string{"tariffZone": "Zone1"}, false},
		{"8", &model.StopPoint{Municipality: &model.Municipality{Name: "Municipality1"}}, map[string]string{"municipalityName": "Municipality1"}, true},
		{"9", &model.StopPoint{Municipality: &model.Municipality{Name: "Municipality2"}}, map[string]string{"municipalityName": "Municipality1"}, false},
		{"10", &model.StopPoint{Municipality: &model.Municipality{PublicCode: "M1"}}, map[string]string{"municipalityShortName": "M1"}, true},
		{"11", &model.StopPoint{Municipality: &model.Municipality{PublicCode: "M2"}}, map[string]string{"municipalityShortName": "M1"}, false},
		{"12", &model.StopPoint{Latitude: 60.0, Longitude: 24.0}, map[string]string{"location": "59.0,23.0:61.0,25.0"}, true},
		{"13", &model.StopPoint{Latitude: 60.0, Longitude: 24.0}, map[string]string{"location": "61.0,25.0:62.0,26.0"}, false},
		{"14", &model.StopPoint{Latitude: 60.0, Longitude: 24.0}, map[string]string{"location": "61.0:62.0,26.0"}, false},
		{"15", &model.StopPoint{Latitude: 60.0, Longitude: 24.0}, map[string]string{"location": "61.0,foo:62.0,26.0"}, false},
		{"16", &model.StopPoint{Latitude: 60.0, Longitude: 24.0}, map[string]string{"location": "foo,25.0:62.0,26.0"}, false},
		{"17", &model.StopPoint{Latitude: 60.0, Longitude: 24.0}, map[string]string{"location": "61.0,25.0:62.0"}, false},
		{"18", &model.StopPoint{Latitude: 60.0, Longitude: 24.0}, map[string]string{"location": "61.0,25.0:foo,26.0"}, false},
		{"19", &model.StopPoint{Latitude: 60.0, Longitude: 24.0}, map[string]string{"location": "61.0,25.0:62.0,foo"}, false},
		{"20", &model.StopPoint{Latitude: 60.0, Longitude: 24.0}, map[string]string{"location": "61.0,25.0"}, false},
	}

	for _, tc := range testCases {
		matches := stopPointMatchesConditions(tc.item, tc.conditions)
		testutil.CompareVariablesAndPrintResults(t, tc.expected, matches, tc.id)
	}
}
