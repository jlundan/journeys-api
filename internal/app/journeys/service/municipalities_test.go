package service

import (
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
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
