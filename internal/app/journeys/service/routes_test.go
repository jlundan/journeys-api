package service

import (
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
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
