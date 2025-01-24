package service

import (
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"testing"
)

func TestJourneyPatternMatchesConditions(t *testing.T) {
	testCases := []struct {
		item       *model.JourneyPattern
		conditions map[string]string
		expected   bool
	}{
		{nil, nil, false},
		{&model.JourneyPattern{Route: nil}, map[string]string{"lineId": "1"}, false},
		{&model.JourneyPattern{Route: &model.Route{Line: nil}}, map[string]string{"lineId": "1"}, false},
		{&model.JourneyPattern{StopPoints: nil}, map[string]string{"firstStopPointId": "1"}, false},
		{&model.JourneyPattern{StopPoints: nil}, map[string]string{"lastStopPointId": "1"}, false},
		{&model.JourneyPattern{StopPoints: nil}, map[string]string{"stopPointId": "1"}, false},
	}

	for _, tc := range testCases {
		matches := journeyPatternMatchesConditions(tc.item, tc.conditions)
		if matches != tc.expected {
			t.Error(fmt.Sprintf("expected %v, got %v", matches, tc.expected))
		}
	}
}
