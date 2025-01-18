//go:build journeys_stops_tests || journeys_tests || all_tests

package service

import (
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
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
