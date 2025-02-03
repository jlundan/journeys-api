package service

import (
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/app/journeys/repository"
	"github.com/jlundan/journeys-api/internal/testutil"
	"testing"
)

func TestMunicipalitiesMatchesConditions(t *testing.T) {
	testCases := []struct {
		item       *model.Municipality
		conditions map[string]string
		expected   bool
	}{
		{nil, nil, false},
	}

	for _, tc := range testCases {
		matches := municipalityMatchesConditions(tc.item, tc.conditions)
		if matches != tc.expected {
			t.Error(fmt.Sprintf("expected %v, got %v", matches, tc.expected))
		}
	}
}

func TestMunicipalitiesService_Search(t *testing.T) {
	dataStore := &repository.JourneysRepository{
		Municipalities: &repository.JourneysMunicipalitiesRepository{
			All: []*model.Municipality{
				{Name: "Municipality1", PublicCode: "M1"},
				{Name: "Municipality2", PublicCode: "M2"},
			},
		},
	}
	service := MunicipalitiesService{Repository: dataStore}

	testCases := []struct {
		id       string
		params   map[string]string
		expected int
	}{
		{"1", map[string]string{"name": "Municipality1"}, 1},
		{"2", map[string]string{"shortName": "M1"}, 1},
		{"3", map[string]string{"name": "NonExistent"}, 0},
	}

	for _, tc := range testCases {
		result := service.Search(tc.params)
		testutil.CompareVariablesAndPrintResults(t, tc.expected, len(result), tc.id)
	}
}

func TestMunicipalitiesService_GetOneById(t *testing.T) {
	dataStore := &repository.JourneysRepository{
		Municipalities: &repository.JourneysMunicipalitiesRepository{
			ById: map[string]*model.Municipality{
				"1": {Name: "Municipality1", PublicCode: "M1"},
				"2": {Name: "Municipality2", PublicCode: "M2"},
			},
		},
	}
	service := MunicipalitiesService{Repository: dataStore}

	testCases := []struct {
		id             string
		municipalityId string
		expected       *model.Municipality
		err            error
	}{
		{"1", "1", &model.Municipality{Name: "Municipality1", PublicCode: "M1"}, nil},
		{"2", "2", &model.Municipality{Name: "Municipality2", PublicCode: "M2"}, nil},
		{"3", "3", nil, model.ErrNoSuchElement},
	}

	for _, tc := range testCases {
		result, err := service.GetOneById(tc.municipalityId)
		if err != nil && tc.err == nil {
			t.Error(err)
		} else if err != nil && tc.err != nil {
			testutil.CompareVariablesAndPrintResults(t, tc.err, err, tc.id)
		} else {
			testutil.CompareVariablesAndPrintResults(t, tc.expected, result, tc.id)
		}
	}
}

func TestMunicipalityMatchesConditions(t *testing.T) {
	testCases := []struct {
		id         string
		item       *model.Municipality
		conditions map[string]string
		expected   bool
	}{
		{"1", nil, nil, false},
		{"2", &model.Municipality{Name: "Municipality1"}, map[string]string{"name": "Municipality1"}, true},
		{"3", &model.Municipality{Name: "Municipality1"}, map[string]string{"name": "Municipality2"}, false},
		{"4", &model.Municipality{PublicCode: "M1"}, map[string]string{"shortName": "M1"}, true},
		{"5", &model.Municipality{PublicCode: "M2"}, map[string]string{"shortName": "M1"}, false},
	}

	for _, tc := range testCases {
		matches := municipalityMatchesConditions(tc.item, tc.conditions)
		testutil.CompareVariablesAndPrintResults(t, tc.expected, matches, tc.id)
	}
}
