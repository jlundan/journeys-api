package service

import (
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/app/journeys/repository"
	"github.com/jlundan/journeys-api/internal/testutil"
	"testing"
)

func TestRoutesMatchesConditions(t *testing.T) {
	testCases := []struct {
		item       *model.Route
		conditions map[string]string
		expected   bool
	}{
		{nil, nil, false},
		{&model.Route{Line: nil}, map[string]string{"lineId": "1"}, false},
	}

	for _, tc := range testCases {
		matches := routeMatchesConditions(tc.item, tc.conditions)
		if matches != tc.expected {
			t.Error(fmt.Sprintf("expected %v, got %v", matches, tc.expected))
		}
	}
}

func TestRoutesService_Search(t *testing.T) {
	dataStore := &repository.JourneysRepository{
		Routes: &repository.JourneysRoutesRepository{
			All: []*model.Route{
				{Id: "1", Name: "Route1", Line: &model.Line{Name: "Line1"}},
				{Id: "2", Name: "Route2", Line: &model.Line{Name: "Line2"}},
			},
		},
	}
	service := RoutesService{Repository: dataStore}

	testCases := []struct {
		id       string
		params   map[string]string
		expected int
	}{
		{"1", map[string]string{"name": "Route1"}, 1},
		{"2", map[string]string{"lineId": "Line1"}, 1},
		{"3", map[string]string{"name": "NonExistent"}, 0},
	}

	for _, tc := range testCases {
		result := service.Search(tc.params)
		testutil.CompareVariablesAndPrintResults(t, tc.expected, len(result), tc.id)
	}
}

func TestRoutesService_GetOneById(t *testing.T) {
	dataStore := &repository.JourneysRepository{
		Routes: &repository.JourneysRoutesRepository{
			ById: map[string]*model.Route{
				"1": {Id: "1", Name: "Route1", Line: &model.Line{Name: "Line1"}},
				"2": {Id: "2", Name: "Route2", Line: &model.Line{Name: "Line2"}},
			},
		},
	}
	service := RoutesService{Repository: dataStore}

	testCases := []struct {
		id       string
		routeId  string
		expected *model.Route
		err      error
	}{
		{"1", "1", &model.Route{Id: "1", Name: "Route1", Line: &model.Line{Name: "Line1"}}, nil},
		{"2", "2", &model.Route{Id: "2", Name: "Route2", Line: &model.Line{Name: "Line2"}}, nil},
		{"3", "3", nil, model.ErrNoSuchElement},
	}

	for _, tc := range testCases {
		result, err := service.GetOneById(tc.routeId)
		if err != nil && tc.err == nil {
			t.Error(err)
		} else if err != nil && tc.err != nil {
			testutil.CompareVariablesAndPrintResults(t, tc.err, err, tc.id)
		} else {
			testutil.CompareVariablesAndPrintResults(t, tc.expected, result, tc.id)
		}
	}
}

func TestRouteMatchesConditions(t *testing.T) {
	testCases := []struct {
		id         string
		route      *model.Route
		conditions map[string]string
		expected   bool
	}{
		{"1", nil, nil, false},
		{"2", &model.Route{Name: "Route1"}, map[string]string{"name": "Route1"}, true},
		{"3", &model.Route{Name: "Route1"}, map[string]string{"name": "Route2"}, false},
		{"4", &model.Route{Line: &model.Line{Name: "Line1"}}, map[string]string{"lineId": "Line1"}, true},
		{"5", &model.Route{Line: &model.Line{Name: "Line2"}}, map[string]string{"lineId": "Line1"}, false},
	}

	for _, tc := range testCases {
		matches := routeMatchesConditions(tc.route, tc.conditions)
		testutil.CompareVariablesAndPrintResults(t, tc.expected, matches, tc.id)
	}
}
