package service

import (
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/app/journeys/repository"
	"github.com/jlundan/journeys-api/internal/testutil"
	"testing"
)

func TestJourneyPatternsService_Search(t *testing.T) {
	dataStore := &repository.JourneysRepository{
		JourneyPatterns: &repository.JourneysJourneyPatternRepository{
			All: []*model.JourneyPattern{
				{Name: "Pattern1", Route: &model.Route{Line: &model.Line{Name: "1"}}, StopPoints: []*model.StopPoint{{ShortName: "SP1"}, {ShortName: "SP2"}}},
				{Name: "Pattern2", Route: &model.Route{Line: &model.Line{Name: "2"}}, StopPoints: []*model.StopPoint{{ShortName: "SP3"}, {ShortName: "SP4"}}},
			},
		},
	}
	service := JourneyPatternsService{Repository: dataStore}

	testCases := []struct {
		id       string
		params   map[string]string
		expected int
	}{
		{"1", map[string]string{"name": "Pattern1"}, 1},
		{"2", map[string]string{"lineId": "1"}, 1},
		{"3", map[string]string{"firstStopPointId": "SP1"}, 1},
		{"4", map[string]string{"lastStopPointId": "SP2"}, 1},
		{"5", map[string]string{"stopPointId": "SP3"}, 1},
		{"6", map[string]string{"name": "NonExistent"}, 0},
	}

	for _, tc := range testCases {
		result := service.Search(tc.params)
		testutil.CompareVariablesAndPrintResults(t, tc.expected, len(result), tc.id)
	}
}

func TestJourneyPatternsService_GetOneById(t *testing.T) {
	repo := &repository.JourneysRepository{
		JourneyPatterns: &repository.JourneysJourneyPatternRepository{
			ById: map[string]*model.JourneyPattern{
				"1": {Name: "Pattern1"},
				"2": {Name: "Pattern2"},
			},
		},
	}
	service := JourneyPatternsService{Repository: repo}

	testCases := []struct {
		id       string
		expected *model.JourneyPattern
		err      error
	}{
		{"1", &model.JourneyPattern{Name: "Pattern1"}, nil},
		{"2", &model.JourneyPattern{Name: "Pattern2"}, nil},
		{"3", nil, model.ErrNoSuchElement},
	}

	for _, tc := range testCases {
		result, err := service.GetOneById(tc.id)
		if err != nil && tc.err == nil {
			t.Error(err)
		} else if err != nil && tc.err == nil {
			testutil.CompareVariablesAndPrintResults(t, tc.err, err, tc.id)
		} else {
			testutil.CompareVariablesAndPrintResults(t, tc.expected, result, tc.id)
		}
	}
}

func TestJourneyPatternMatchesConditions(t *testing.T) {
	testCases := []struct {
		id         string
		item       *model.JourneyPattern
		conditions map[string]string
		expected   bool
	}{
		{"1", nil, nil, false},
		{"2", &model.JourneyPattern{Name: "Pattern1"}, map[string]string{"name": "Pattern1"}, true},
		{"3", &model.JourneyPattern{Name: "Pattern1"}, map[string]string{"name": "Pattern2"}, false},
		{"4", &model.JourneyPattern{Route: &model.Route{Line: &model.Line{Name: "1"}}}, map[string]string{"lineId": "1"}, true},
		{"5", &model.JourneyPattern{Route: &model.Route{Line: &model.Line{Name: "2"}}}, map[string]string{"lineId": "1"}, false},
		{"6", &model.JourneyPattern{StopPoints: []*model.StopPoint{{ShortName: "SP1"}}}, map[string]string{"firstStopPointId": "SP1"}, true},
		{"7", &model.JourneyPattern{StopPoints: []*model.StopPoint{{ShortName: "SP2"}}}, map[string]string{"firstStopPointId": "SP1"}, false},
		{"8", &model.JourneyPattern{StopPoints: []*model.StopPoint{{ShortName: "SP1"}, {ShortName: "SP2"}}}, map[string]string{"lastStopPointId": "SP2"}, true},
		{"9", &model.JourneyPattern{StopPoints: []*model.StopPoint{{ShortName: "SP1"}, {ShortName: "SP3"}}}, map[string]string{"lastStopPointId": "SP2"}, false},
		{"10", &model.JourneyPattern{StopPoints: []*model.StopPoint{{ShortName: "SP1"}, {ShortName: "SP2"}}}, map[string]string{"stopPointId": "SP2"}, true},
		{"11", &model.JourneyPattern{StopPoints: []*model.StopPoint{{ShortName: "SP1"}, {ShortName: "SP3"}}}, map[string]string{"stopPointId": "SP2"}, false},
	}

	for _, tc := range testCases {
		matches := journeyPatternMatchesConditions(tc.item, tc.conditions)
		testutil.CompareVariablesAndPrintResults(t, tc.expected, matches, tc.id)
	}
}
