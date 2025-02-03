package service

import (
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/app/journeys/repository"
	"github.com/jlundan/journeys-api/internal/testutil"
	"testing"
)

func TestLineMatchesConditions(t *testing.T) {
	testCases := []struct {
		item       *model.Line
		conditions map[string]string
		expected   bool
	}{
		{
			nil,
			nil,
			false,
		},
		{
			&model.Line{Name: "1", Description: "Vatiala - Pirkkala"},
			nil,
			true,
		},
	}

	for _, tc := range testCases {
		matches := lineMatchesConditions(tc.item, tc.conditions)
		if matches != tc.expected {
			t.Error(fmt.Sprintf("expected %v, got %v", matches, tc.expected))
		}
	}
}

func TestLinesService_Search(t *testing.T) {
	dataStore := &repository.JourneysRepository{
		Lines: &repository.JourneysLinesRepository{
			All: []*model.Line{
				{Name: "Line1", Description: "Description1"},
				{Name: "Line2", Description: "Description2"},
			},
		},
	}
	service := LinesService{Repository: dataStore}

	testCases := []struct {
		id       string
		params   map[string]string
		expected int
	}{
		{"1", map[string]string{"name": "Line1"}, 1},
		{"2", map[string]string{"description": "Description1"}, 1},
		{"3", map[string]string{"name": "NonExistent"}, 0},
	}

	for _, tc := range testCases {
		result := service.Search(tc.params)
		testutil.CompareVariablesAndPrintResults(t, tc.expected, len(result), tc.id)
	}
}

func TestLinesService_GetOneById(t *testing.T) {
	dataStore := &repository.JourneysRepository{
		Lines: &repository.JourneysLinesRepository{
			ById: map[string]*model.Line{
				"1": {Name: "Line1", Description: "Description1"},
				"2": {Name: "Line2", Description: "Description2"},
			},
		},
	}
	service := LinesService{Repository: dataStore}

	testCases := []struct {
		id       string
		lineId   string
		expected *model.Line
		err      error
	}{
		{"1", "1", &model.Line{Name: "Line1", Description: "Description1"}, nil},
		{"2", "2", &model.Line{Name: "Line2", Description: "Description2"}, nil},
		{"3", "3", nil, model.ErrNoSuchElement},
	}

	for _, tc := range testCases {
		result, err := service.GetOneById(tc.lineId)
		if err != nil && tc.err == nil {
			t.Error(err)
		} else if err != nil && tc.err != nil {
			testutil.CompareVariablesAndPrintResults(t, tc.err, err, tc.id)
		} else {
			testutil.CompareVariablesAndPrintResults(t, tc.expected, result, tc.id)
		}
	}
}

func TestLineMatchesConditions2(t *testing.T) {
	testCases := []struct {
		id         string
		line       *model.Line
		conditions map[string]string
		expected   bool
	}{
		{"1", nil, nil, false},
		{"2", &model.Line{Name: "Line1"}, map[string]string{"name": "Line1"}, true},
		{"3", &model.Line{Name: "Line1"}, map[string]string{"name": "Line2"}, false},
		{"4", &model.Line{Description: "Description1"}, map[string]string{"description": "Description1"}, true},
		{"5", &model.Line{Description: "Description2"}, map[string]string{"description": "Description1"}, false},
	}

	for _, tc := range testCases {
		matches := lineMatchesConditions(tc.line, tc.conditions)
		testutil.CompareVariablesAndPrintResults(t, tc.expected, matches, tc.id)
	}
}
